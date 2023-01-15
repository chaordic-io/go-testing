package testhelpers

import (
	"fmt"
	"io"
	"os"
)

type chunkedFile struct {
	name         string
	chunkSize    int
	fileCounter  int
	bytesWritten int
	currentFile  *os.File
}

func newChunkedFile(name string, chunkSize int) *chunkedFile {
	c := &chunkedFile{name: name, chunkSize: chunkSize}
	return c
}

// Writer is the interface that wraps the basic Write method.
//
// Write writes len(p) bytes from p to the underlying data stream.
// It returns the number of bytes written from p (0 <= n <= len(p))
// and any error encountered that caused the write to stop early.
// Write must return a non-nil error if it returns n < len(p).
// Write must not modify the slice data, even temporarily.
//
// Implementations must not retain p.
func (c *chunkedFile) Write(p []byte) (int, error) {
	if c.currentFile == nil {
		err := c.newFile()
		if err != nil {
			return 0, err
		}
	}
	i := 0
	for {
		size := c.chunkSize + i
		if c.bytesWritten+c.chunkSize > c.chunkSize {
			size = c.chunkSize - c.bytesWritten
		}
		if size > len(p) {
			size = len(p)
		}
		b := p[i:size]
		n, err := c.currentFile.Write(b)
		c.bytesWritten += n
		if err != nil {
			return n, err
		}
		if c.bytesWritten >= c.chunkSize {
			err = c.newFile()
			if err != nil {
				return 0, err
			}
		}
		if len(p) == size {
			return size, nil
		}
		i += n
		if len(b) == 0 {
			return len(p), nil
		}
	}
}

// Reader is the interface that wraps the basic Read method.
//
// Read reads up to len(p) bytes into p. It returns the number of bytes
// read (0 <= n <= len(p)) and any error encountered. Even if Read
// returns n < len(p), it may use all of p as scratch space during the call.
// If some data is available but not len(p) bytes, Read conventionally
// returns what is available instead of waiting for more.
//
// When Read encounters an error or end-of-file condition after
// successfully reading n > 0 bytes, it returns the number of
// bytes read. It may return the (non-nil) error from the same call
// or return the error (and n == 0) from a subsequent call.
// An instance of this general case is that a Reader returning
// a non-zero number of bytes at the end of the input stream may
// return either err == EOF or err == nil. The next Read should
// return 0, EOF.
//
// Callers should always process the n > 0 bytes returned before
// considering the error err. Doing so correctly handles I/O errors
// that happen after reading some bytes and also both of the
// allowed EOF behaviors.
//
// Implementations of Read are discouraged from returning a
// zero byte count with a nil error, except when len(p) == 0.
// Callers should treat a return of 0 and nil as indicating that
// nothing happened; in particular it does not indicate EOF.
//
// Implementations must not retain p.
func (c *chunkedFile) Read(p []byte) (int, error) {
	if c.currentFile == nil {
		err := c.openFile()
		if err != nil {
			return 0, io.EOF
		}
	}
	n, err := c.currentFile.Read(p)
	if n == c.chunkSize {
		e := c.currentFile.Close()
		if e != nil {
			return 0, e
		}
		c.currentFile = nil
	}
	if err == io.EOF {
		err = nil
		e := c.currentFile.Close()
		if e != nil {
			return 0, e
		}
		c.currentFile = nil
	}
	return n, err
}

func (c *chunkedFile) Close() error {
	if c.currentFile != nil {
		return c.currentFile.Close()
	}
	return nil
}

func (c *chunkedFile) newFile() error {
	err := c.Close()
	if err != nil {
		return err
	}
	f, err := os.Create(fmt.Sprintf("%s.%d", c.name, c.fileCounter))
	if err != nil {
		return err
	}
	c.currentFile = f
	c.bytesWritten = 0
	c.fileCounter++
	if c.fileCounter > 100 {
		panic("arrgghhh ")
	}
	return nil
}

func (c *chunkedFile) openFile() error {
	err := c.Close()
	if err != nil {
		return err
	}
	f, err := os.Open(fmt.Sprintf("%s.%d", c.name, c.fileCounter))
	c.fileCounter++
	if err != nil {
		return err
	}
	c.currentFile = f
	if c.fileCounter > 100 {
		panic("arrgghhh")
	}
	return nil
}
