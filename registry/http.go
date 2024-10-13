package registry

import (
	"crypto/tls"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/hashicorp/go-retryablehttp"
)

// retryClient is a retryable http client
var retryClient *retryablehttp.Client

// init initializes the http client to skip certificate verification
func init() {
	retryClient = retryablehttp.NewClient()
	retryClient.HTTPClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	retryClient.Logger = nil
}

// SetHttpTimeout sets the http timeout
func SetHttpTimeout(timeout int) {
	retryClient.HTTPClient.Timeout = time.Duration(timeout) * time.Second
}

// SetHttpRetry sets the http retry
func SetHttpRetry(retry int) {
	retryClient.RetryMax = retry
}

// SetAuthToken sets the auth token
func SetAuthToken(token string) {
	retryClient.HTTPClient.Transport = NewAddAuth(retryClient.HTTPClient.Transport, token)
}

// SetUserAgent sets the user agent
func SetUserAgent(ua string) {
	retryClient.HTTPClient.Transport = NewUserAgentTransport(retryClient.HTTPClient.Transport, ua)
}

func download(url, path string) error {

	// create the request
	res, err := retryClient.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download layer: %s", err)
	}
	defer res.Body.Close()

	// Check if the request was successful
	if res.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download layer: %s", res.Status)
	}

	// Create the file
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create directory: %s", dir)
	}
	out, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("failed to create file: %s", path)
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, res.Body)
	if err != nil {
		return fmt.Errorf("failed to write to file (%s): %v", path, err)
	}

	return nil

}

// AddAuthTransport is a transport that adds an auth token to the request
type AddAuthTransport struct {
	T    http.RoundTripper
	auth string
}

func (t *AddAuthTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("Authorization", t.auth)
	return t.T.RoundTrip(req)
}

func NewAddAuth(T http.RoundTripper, auth string) *AddAuthTransport {
	if T == nil {
		T = http.DefaultTransport
	}
	return &AddAuthTransport{T, auth}
}

// AddUserAgentTransport is a transport that adds a user agent to the request
type AddUserAgentTransport struct {
	T  http.RoundTripper
	ua string
}

func (t *AddUserAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Add("User-Agent", t.ua)
	return t.T.RoundTrip(req)
}

func NewUserAgentTransport(T http.RoundTripper, ua string) *AddUserAgentTransport {
	if T == nil {
		T = http.DefaultTransport
	}
	return &AddUserAgentTransport{T, ua}
}
