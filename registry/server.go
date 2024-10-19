package registry

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
)

type Server struct {
	Address  string
	Images   []Image
	Protocol string
	Port     string
	pinged   bool
}

func NewServer(address string) *Server {

	// set default protocol and port
	var protocol, port string

	// check if the address has a protocol
	if strings.HasPrefix(address, "https://") {
		address = strings.TrimPrefix(address, "https://")
		protocol = "https"
		port = "443"
	} else if strings.HasPrefix(address, "http://") {
		address = strings.TrimPrefix(address, "http://")
		protocol = "http"
		port = "80"
	}

	// split address and port
	if parts := strings.Split(address, ":"); len(parts) > 1 {
		address = parts[0]
		port = parts[1]

		// remove the port from the port if it contains a slash
		if portParts := strings.Split(port, "/"); len(portParts) > 1 {
			port = portParts[0]
		}
	}

	return &Server{
		Address:  address,
		Protocol: protocol,
		Port:     port,
	}
}

const (
	PING_ENDPOINT        = "v2/"
	CATALOG_ENDPOINT     = "v2/_catalog"
	TAG_LISTING_ENDPOINT = "v2/%s/tags/list"
	MANIFEST_ENDPOINT    = "v2/%s/manifests/%s"
	PATCH_SUFFIX         = "patched"
)

type protoPair struct {
	Protocol string
	Port     string
}

// Ping performs a http request and checks if the server is a registry server.
// The check is executed only once
func (s *Server) Ping() error {

	// Check if the server has already been pinged
	if s.pinged {
		return nil
	}
	s.pinged = true

	// Try the protocol and port specified in the address
	var protocols []protoPair
	if s.Protocol != "" {
		protocols = []protoPair{{s.Protocol, s.Port}}
	} else if s.Port != "" {
		protocols = []protoPair{{"http", s.Port}, {"https", s.Port}}
	} else {
		protocols = []protoPair{{"http", "80"}, {"https", "443"}}
	}

	// Try each protocol
	for _, proto := range protocols {
		requestURL := fmt.Sprintf("%s://%s:%s/%s", proto.Protocol, s.Address, proto.Port, PING_ENDPOINT)
		res, err := retryClient.Get(requestURL)
		if err != nil {
			continue // Try the next protocol
		}
		defer res.Body.Close()

		if res.Header.Get("Docker-Distribution-API-Version") == "registry/2.0" {
			s.Protocol = proto.Protocol
			s.Port = proto.Port
			return nil // Success
		}
	}

	return fmt.Errorf("server is not a docker registry")
}

// FetchImagesAndTags returns a list of images from the registry server
func (s *Server) FetchImagesAndTags(logger *slog.Logger) error {

	if err := s.Ping(); err != nil {
		return fmt.Errorf("cannot fetch images and tags: %s", err)
	}

	// Check if the images have already been fetched
	if s.Images != nil {
		return nil
	}

	// Get the list of images from the registry
	requestURL := fmt.Sprintf("%s/%s", s.GetUrl(), CATALOG_ENDPOINT)
	res, err := retryClient.Get(requestURL)
	if err != nil {
		return fmt.Errorf("request failed: %s", err)
	}
	defer res.Body.Close()

	// Parse the response
	var imageList CatalogEndpointResponse
	if err := json.NewDecoder(res.Body).Decode(&imageList); err != nil {
		return fmt.Errorf("failed to decode response: %s", err)
	}

	// store the list of images
	s.Images = make([]Image, len(imageList.Repositories))
	for i, imageName := range imageList.Repositories {
		s.Images[i] = Image{Name: imageName}
	}

	// fetch tags for each image
	for i := range s.Images {
		logger.Debug("fetch tags", "image", s.Images[i].Name)
		if err := s.FetchTags(&s.Images[i]); err != nil {
			return fmt.Errorf("failed to fetch tags: %s", err)
		}
	}

	return nil
}

func (s *Server) FetchTags(image *Image) error {
	requestURL := fmt.Sprintf("%s/%s", s.GetUrl(), fmt.Sprintf(TAG_LISTING_ENDPOINT, image.Name))
	res, err := retryClient.Get(requestURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Parse the response
	var tagList TagListingEndpointResponse
	if err := json.NewDecoder(res.Body).Decode(&tagList); err != nil {
		return err
	}

	// store the list of tags
	for _, tagName := range tagList.Tags {
		t := Tag{Name: tagName}
		image.Tags = append(image.Tags, t)
	}
	return nil
}

func (s *Server) GetUrl() string {
	addr := s.Address
	if len(s.Port) > 0 {
		addr += ":" + s.Port
	}
	if len(s.Protocol) > 0 {
		addr = s.Protocol + "://" + addr
	} else {
		addr = "http://" + addr
	}
	return addr
}

func (s *Server) GetService() string {
	if s.Port == "80" || len(s.Port) == 0 {
		return s.Address
	}
	return fmt.Sprintf("%s:%s", s.Address, s.Port)
}

func (s *Server) fetchManifest(img *Image, tag *Tag) error {

	// Fetch the manifest
	requestURL := fmt.Sprintf("%s/v2/%s/manifests/%s", s.GetUrl(), img.Name, tag.Name)
	res, err := retryClient.Get(requestURL)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Parse the response
	m := &Manifest{}
	if err := json.NewDecoder(res.Body).Decode(m); err != nil {
		return err
	}
	tag.Manifest = m

	// Parse the history entry
	historyEntries, err := m.ParseHistoryEntries()
	if err != nil {
		return err
	}
	for i := range historyEntries {
		m.History[i].HistoryEntry = historyEntries[i]
	}

	return nil
}

// Dump downloads the image from the registry to the local machine by fetching the manifest
// for each image and tag combination and then downloading the layers
// This is a simplified version of the Docker pull command
func (s *Server) Dump(logger *slog.Logger, path string, manifestOnly bool, failCount int32, parallel int) error {
	logger.Debug("dump layers from server", "proto", s.Protocol, "address", s.Address, "port", s.Port)

	// Fetch list of images and tags
	if err := s.FetchImagesAndTags(logger); err != nil {
		return fmt.Errorf("failed to dump layer: %w", err)
	}

	// setup counter for failed downloads
	if failCount <= 0 {
		failCount = 3
	}

	// setup 5 worker
	jobs := make(chan dlJob, 100)
	var wg sync.WaitGroup
	for w := 0; w < parallel; w++ {
		wg.Add(1)
		go func(worker int) {
			defer func() {
				logger.Debug("worker done", "worker", worker)
				wg.Done()

			}()
			for job := range jobs {

				// download the layers
				err := downloadAndLink(logger, worker, path, s, job.img, job.tag, job.blobSum)
				if err != nil {
					logger.Error("failed to download layer", "worker", worker, "image", job.img.Name, "tag", job.tag.Name, "blobSum", job.blobSum, "err", err)
					atomic.AddInt32(&failCount, -1)
				}

				// check if we have reached the fail count
				if atomic.LoadInt32(&failCount) <= 0 {
					logger.Error("failed too many times, aborting", "worker", worker)
					return
				}
			}
		}(w)
	}

	// iterate over images and tags and download the layers
	for _, img := range s.Images {
		for _, tag := range img.Tags {
			if atomic.LoadInt32(&failCount) <= 0 {
				close(jobs)
				wg.Wait()
				return fmt.Errorf("failed too many times, aborting")
			}
			if tag.Manifest == nil {
				if err := s.fetchManifest(&img, &tag); err != nil {
					logger.Error("failed to fetch manifest", "image", img.Name, "tag", tag.Name, "err", err)
					atomic.AddInt32(&failCount, -1)
					continue
				}
			}
			if tag.Manifest == nil {
				logger.Warn("no manifest found, skip", "image", img.Name, "tag", tag.Name)
				continue
			}

			// check if the image directory exists, if not create it
			imagePath := filepath.Join(path, s.Address, img.Name, tag.Name)
			if _, err := os.Stat(imagePath); os.IsNotExist(err) {
				if err := os.MkdirAll(imagePath, 0755); err != nil {
					logger.Error("failed to create directory", "path", imagePath, "err", err)
					continue
				}
			} else if err != nil {
				logger.Error("failed to check directory", "path", imagePath, "err", err)
				continue
			}

			manifestPath := filepath.Join(imagePath, "manifest.json")
			if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
				manifest, err := json.MarshalIndent(tag.Manifest, "", "  ")
				if err != nil {
					logger.Error("failed to marshal manifest", "image", img.Name, "tag", tag.Name, "err", err)
					continue
				}
				if err := os.WriteFile(manifestPath, manifest, 0644); err != nil {
					logger.Error("failed to write manifest", "image", img.Name, "tag", tag.Name, "err", err)
					continue
				}
				logger.Debug("manifest stored", "image", img.Name, "tag", tag.Name, "path", manifestPath)
			} else if err != nil {
				logger.Error("failed to check manifest", "image", img.Name, "tag", tag.Name, "err", err)
				continue
			}

			// check if we only want the manifest
			if manifestOnly {
				continue
			}

			// download the layers
			logger.Debug("download layers", "image", img.Name, "tag", tag.Name, "layers", len(tag.Manifest.FsLayers))

			for _, layer := range tag.Manifest.FsLayers {
				j := dlJob{
					img:     &img,
					tag:     &tag,
					blobSum: layer.BlobSum,
				}

				jobs <- j
			}
		}
	}

	close(jobs)
	wg.Wait()

	return nil
}

type dlJob struct {
	img     *Image
	tag     *Tag
	blobSum string
}

func downloadAndLink(logger *slog.Logger, worker int, dstDir string, s *Server, img *Image, tag *Tag, blobSum string) error {

	// assamble src & dst for download
	requestURL := fmt.Sprintf("%s/v2/%s/blobs/%s", s.GetUrl(), img.Name, blobSum)
	downloadDst := filepath.Join(dstDir, "layer", blobSum)

	// link layer to image
	imagePath := filepath.Join(dstDir, s.Address, img.Name, tag.Name)
	fsLayer := filepath.Join(imagePath, blobSum)

	// get relative path from fsLayer to downloadDst
	// this is needed to create the symlink
	fsLayerTarget, err := filepath.Rel(filepath.Dir(fsLayer), downloadDst)
	if err != nil {
		return fmt.Errorf("failed to get relative path: %s", err)
	}

	// check if layer exists, if not download it
	if _, err := os.Stat(downloadDst); err != nil {
		logger.Debug("downloading layer", "worker", worker, "address", s.Address, "image", img.Name, "tag", tag.Name, "digest", blobSum[7:19])
		if err := download(requestURL, downloadDst); err != nil {
			return fmt.Errorf("failed to download layer: %s", err)
		}
	}

	// check if symlink exists, if not create it
	if _, err := os.Lstat(fsLayer); os.IsNotExist(err) {
		if err := os.Symlink(fsLayerTarget, fsLayer); err != nil {
			return fmt.Errorf("failed to create symlink: %s", err)
		}
	} else if err != nil {
		return fmt.Errorf("failed to check symlink: %s", err)
	}

	return nil
}
