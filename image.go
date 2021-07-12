// Package img enables basic image operations with automatic features like orientation correction
// and EXIF / ICC profile preservation, so you can transform photos without losing important
// metadata.
package img

import (
	"bytes"
	"fmt"
	"github.com/snapas/imageorient"
	"github.com/snapas/img/iccjpeg"
	"image"
	"io"
)

// Image contains an image.Image plus any metadata we want to preserve through future image transformations.
type Image struct {
	buf   *bytes.Buffer
	Image image.Image
	App2  []byte
}

// Decode decodes an image and changes its orientation according to the EXIF orientation tag (if present), while also
// preserving any ICC profile (APP2 data) in the returned Image.
func Decode(r io.Reader) (Image, string, error) {
	i := Image{
		buf: &bytes.Buffer{},
	}
	buf := &bytes.Buffer{}
	var err error

	// Parse out needed metadata we need to retain
	tr := io.TeeReader(r, buf)
	i.App2, err = iccjpeg.GetICCRaw(tr)
	if err != nil {
		return i, "", fmt.Errorf("GetICCRaw: %s", err)
	}

	// Fix orientation
	ri, s, err := imageorient.Decode(io.MultiReader(buf, r))
	if err != nil {
		return i, "", fmt.Errorf("imageorient.Decode: %s", err)
	}

	i.Image = ri
	return i, s, nil
}

// Len returns the number of bytes of the unread portion of the Image's buffer.
func (i Image) Len() int {
	return i.buf.Len()
}

// Bytes returns a slice of length i.Len() holding the unread portion of the Image's buffer.
func (i Image) Bytes() []byte {
	return i.buf.Bytes()
}

// Write appends the contents of p to the Image's buffer.
func (i Image) Write(p []byte) (n int, err error) {
	n, err = i.buf.Write(p)
	return
}
