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

type Image struct {
	buf   *bytes.Buffer
	Image image.Image
	App2  []byte
}

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

func (i Image) Len() int {
	return i.buf.Len()
}

func (i Image) Bytes() []byte {
	return i.buf.Bytes()
}

func (i Image) Write(p []byte) (n int, err error) {
	n, err = i.buf.Write(p)
	return
}
