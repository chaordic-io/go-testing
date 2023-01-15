package testhelpers

import (
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"strings"
	"testing"

	fuzz "github.com/google/gofuzz"
	"github.com/stretchr/testify/assert"
)

func TestWriteAndReadChunkedFile(t *testing.T) {
	for i := 2; i < 30; i++ {
		dd := fuzzBytes()
		data := bytes.NewBuffer(dd)
		f := newChunkedFile("test_", i*100)
		_, err := io.Copy(f, data)
		assert.NoError(t, err)

		f = newChunkedFile("test_", i*100)
		d, err := io.ReadAll(f)
		assert.NoError(t, err)
		assert.Equal(t, dd, d)
		numberOfFiles := 0

		dirs, err := os.ReadDir("./")
		assert.NoError(t, err)
		for _, d := range dirs {
			if strings.Contains(d.Name(), "test_.") {
				assert.NoError(t, os.Remove(d.Name()))
				numberOfFiles++
			}
		}
		assert.Equal(t, len(dd)/(i*100)+1, numberOfFiles)
	}
}

func TestOneMore(t *testing.T) {
	data := bytes.NewBufferString("1234")
	f := newChunkedFile("test_", 9)
	writer := gzip.NewWriter(f)
	_, err := io.Copy(writer, data)
	assert.NoError(t, err)
	assert.NoError(t, writer.Close())

	f = newChunkedFile("test_", 3)
	reader, err := gzip.NewReader(f)
	assert.NoError(t, err)
	d, err := io.ReadAll(reader)
	assert.NoError(t, err)
	assert.Equal(t, "1234", string(d))

	numberOfFiles := 0
	dirs, err := os.ReadDir("./")
	assert.NoError(t, err)
	for _, d := range dirs {
		if strings.Contains(d.Name(), "test_.") {
			assert.NoError(t, os.Remove(d.Name()))
			numberOfFiles++
		}
	}
	assert.Equal(t, 4, numberOfFiles)
}

func fuzzBytes() []byte {
	f := fuzz.New().NilChance(0).NumElements(1, 3000)
	strings := make([]byte, 0)
	f.Fuzz(&strings)
	return strings
}
