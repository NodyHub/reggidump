package registry

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type Server struct {
	Address         string
	Images          []Image
	Protocol        string
	Port            string
	ProcessedHashes map[string]struct{}
	pinged          bool
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
		Address:         address,
		Protocol:        protocol,
		Port:            port,
		ProcessedHashes: make(map[string]struct{}),
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

// Ping performs a http request and checks if the server is a registry server
func (s *Server) Ping() error {

	// Check if the server has already been pinged
	if s.pinged {
		return nil
	}

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
func (s *Server) Dump(logger *slog.Logger, path string, manifestOnly bool, failCount int) error {
	logger.Debug("dump layers from server", "proto", s.Protocol, "address", s.Address, "port", s.Port)
	// Download the layers
	if s.Images == nil {
		s.FetchImagesAndTags(logger)
	}

	// setup counter for failed downloads
	if failCount <= 0 {
		failCount = 3
	}

	for _, img := range s.Images {
		for _, tag := range img.Tags {

			// check if we have reached the fail count
			if failCount == 0 {
				return fmt.Errorf("failed too many times, aborting")
			}

			// check if the manifest has already been fetched
			if tag.Manifest == nil {
				if err := s.fetchManifest(&img, &tag); err != nil {
					logger.Error("failed to fetch manifest", "image", img.Name, "tag", tag.Name, "err", err)
					failCount--
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
					return fmt.Errorf("failed to create directory: %s", imagePath)
				}
			} else if err != nil {
				return fmt.Errorf("failed to check directory: %s", imagePath)
			}

			// store manifest
			manifestPath := filepath.Join(imagePath, "manifest.json")
			if _, err := os.Stat(manifestPath); os.IsNotExist(err) {
				manifest, err := json.MarshalIndent(tag.Manifest, "", "  ")
				if err != nil {
					return fmt.Errorf("failed to marshal manifest: %s", err)
				}
				if err := os.WriteFile(manifestPath, manifest, 0644); err != nil {
					return fmt.Errorf("failed to write manifest: %s", err)
				}
				logger.Debug("manifest stored", "image", img.Name, "tag", tag.Name, "path", manifestPath)
			} else if err != nil {
				return fmt.Errorf("failed to check manifest: %s", err)
			}

			// check if we only want the manifest
			if manifestOnly {
				continue
			}

			// download the layers
			if err := s.downloadLayers(logger, path, &img, &tag); err != nil {
				logger.Error("failed to download layers", "image", img.Name, "tag", tag.Name, "err", err)
				failCount--
				continue
			}
		}
	}

	return nil
}

func (s *Server) downloadLayers(logger *slog.Logger, path string, img *Image, tag *Tag) error {

	logger.Debug("download layers", "image", img.Name, "tag", tag.Name, "layers", len(tag.Manifest.FsLayers))

	// Download the layers
	for _, layer := range tag.Manifest.FsLayers {

		requestURL := fmt.Sprintf("%s/v2/%s/blobs/%s", s.GetUrl(), img.Name, layer.BlobSum)
		layerPath := filepath.Join(path, "layer", layer.BlobSum)

		// check if layer exists, if not download it
		if _, err := os.Stat(layerPath); err != nil {
			logger.Debug("downloading layer", "address", s.Address, "image", img.Name, "tag", tag.Name, "digest", layer.BlobSum[7:19])
			if err := download(requestURL, layerPath); err != nil {
				return fmt.Errorf("failed to download layer: %s", err)
			}
		}

		// link layer to image
		imagePath := filepath.Join(path, s.Address, img.Name, tag.Name)
		linkedLayerPath := filepath.Join(imagePath, layer.BlobSum)
		if _, err := os.Lstat(linkedLayerPath); os.IsNotExist(err) {
			imgNameParts := strings.Split(img.Name, "/")
			traversal := strings.Repeat("../", len(imgNameParts))
			if err := os.Symlink(filepath.Join("..", traversal, "..", "layer", layer.BlobSum), linkedLayerPath); err != nil {
				return fmt.Errorf("failed to create symlink: %s", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to check symlink: %s", err)
		}

		// store the blob sum
		s.ProcessedHashes[layer.BlobSum] = struct{}{}

	}
	return nil
}
