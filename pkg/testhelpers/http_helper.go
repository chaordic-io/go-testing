package testhelpers

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/dnaeon/go-vcr.v3/cassette"
	"gopkg.in/dnaeon/go-vcr.v3/recorder"
)

const fiftyMB = 50 * 1000 * 1000

// MakeRecorder makes an http interaction recorder
func MakeRecorder(fileToRecord string, t *testing.T) (http.RoundTripper, func(), error) {
	unzip(t, fileToRecord)
	r, err := recorder.New(fmt.Sprintf("testdata/%s", fileToRecord))
	if err != nil {
		return nil, nil, err
	}
	r.AddHook(func(i *cassette.Interaction) error {
		delete(i.Request.Headers, "Authorization")
		for k := range i.Request.Headers {
			if strings.HasPrefix(k, "X-Redacted-Header") {
				delete(i.Request.Headers, k)
			}
		}
		return nil
	}, recorder.AfterCaptureHook)

	r.AddHook(func(i *cassette.Interaction) error {
		if strings.Contains(i.Request.URL, "/oauth/token") {
			i.Response.Body = `{"access_token": "[REDACTED]"}`
		}
		return nil
	}, recorder.BeforeSaveHook)

	return r, func() {
		err = r.Stop()
		assert.NoError(t, err)
		zip(t, fileToRecord)
	}, nil
}

// TestFileExists ...
func TestFileExists(fileToRecord string) bool {
	_, err := os.Stat(fmt.Sprintf("testdata/%s.yaml.gz.0", fileToRecord))
	return errors.Is(err, os.ErrNotExist)
}

func zip(t *testing.T, fileToRecord string) {
	if _, err := os.Stat(fmt.Sprintf("testdata/%s.yaml", fileToRecord)); err == nil {
		if _, err := os.Stat(fmt.Sprintf("testdata/%s.yaml.gz.0", fileToRecord)); errors.Is(err, os.ErrNotExist) {
			gzipFile := newChunkedFile(fmt.Sprintf("testdata/%s.yaml.gz", fileToRecord), fiftyMB)
			file, err := os.Open(fmt.Sprintf("testdata/%s.yaml", fileToRecord))
			assert.NoError(t, err)
			defer func() {
				assert.NoError(t, file.Close())
			}()
			writer := gzip.NewWriter(gzipFile)
			_, err = io.Copy(writer, file)
			defer func() {
				assert.NoError(t, writer.Close())
			}()
			assert.NoError(t, err)
		}
		assert.NoError(t, os.Remove(fmt.Sprintf("testdata/%s.yaml", fileToRecord)))
	}
}

func unzip(t *testing.T, fileToRecord string) {
	gzFilename := fmt.Sprintf("testdata/%s.yaml.gz", fileToRecord)
	if _, err := os.Stat(gzFilename); err == nil {
		assert.NoError(t, os.Rename(gzFilename, gzFilename+".0"))
	}
	if _, err := os.Stat(fmt.Sprintf("testdata/%s.yaml.gz.0", fileToRecord)); err == nil {
		if _, err := os.Stat(fmt.Sprintf("testdata/%s.yaml", fileToRecord)); errors.Is(err, os.ErrNotExist) {
			gzipFile := newChunkedFile(fmt.Sprintf("testdata/%s.yaml.gz", fileToRecord), fiftyMB)
			defer func() {
				assert.NoError(t, gzipFile.Close())
			}()

			file, err := os.Create(fmt.Sprintf("testdata/%s.yaml", fileToRecord))
			assert.NoError(t, err)
			defer func() {
				assert.NoError(t, file.Close())
			}()
			reader, err := gzip.NewReader(gzipFile)
			assert.NoError(t, err)
			_, err = io.Copy(file, reader) //nolint
			assert.NoError(t, err)
		}
	}
}
