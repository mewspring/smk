// Package smk implements decoding of Smacker video files.
//
// File format reference:
//    https://wiki.multimedia.cx/index.php?title=Smacker
package smk

import (
	"bufio"
	"io"
	"os"

	"github.com/pkg/errors"
)

// A File is a container of Smacker video and audio tracks.
type File struct {
	// File header.
	FileHeader

	// Underlying io.Reader.
	r io.Reader
	// Underlying io.Closer of reader if present, and nil otherwise.
	c io.Closer
}

// Parse returns a new File for accessing the video and audio tracks of r.
//
// It reads and parses the Smacker file header, the frame size and type
// information, and the Huffman decoding tables, but skips all frame data.
func Parse(r io.Reader) (*File, error) {
	// Parse file header.
	f := &File{
		r: bufio.NewReader(r),
	}
	if c, ok := r.(io.Closer); ok {
		f.c = c
	}
	if err := f.parseFileHeader(); err != nil {
		return nil, err
	}
	// TODO: Parse Huffman decoding tables.
	return f, nil
}

// ParseFile returns a new File for accessing the video and audio tracks of
// path.
//
// It reads and parses the Smacker file header, the frame size and type
// information, and the Huffman decoding tables, but skips all frame data.
func ParseFile(path string) (*File, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return Parse(f)
}

// Close closes the underlying reader if it implements io.Closer, and performs
// no operation otherwise.
func (f *File) Close() error {
	if f.c != nil {
		return f.c.Close()
	}
	return nil
}
