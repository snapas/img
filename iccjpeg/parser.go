package iccjpeg

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
)

type (
	// Parser reads through a JPEG and parses out information about it.
	Parser struct {
		count int
		in    *bufio.Reader
	}

	// Segment represents one segment in a JPEG file.
	Segment struct {
		MarkerID   byte
		MarkerName string
		Size       int
		Offset     int
		Data       []byte
	}
)

// NewParser creates a new Parser with the given io.Reader.
func NewParser(input io.Reader) *Parser {
	return &Parser{
		in: bufio.NewReader(input),
	}
}

// ReadSOI reads the Start of Image marker at the beginning of the JPEG. Always call this before parsing anything else.
func (p *Parser) ReadSOI() error {
	var buf [1024]byte
	_, err := io.ReadFull(p.in, buf[0:2])
	if err != nil {
		return err
	}
	if buf[0] != 0xFF && buf[1] != soiMarker {
		return errors.New("no SOI Marker")
	}
	return nil
}

// GetCommonAppSegments parses the JPEG and returns all APP0, APP1, and APP2 segments.
func (p *Parser) GetCommonAppSegments() ([]Segment, error) {
	var buf [1024]byte
	var err error
	var n int
	var marker byte
	segs := []Segment{}

	for {
		n, err = io.ReadFull(p.in, buf[0:2])
		if err != nil {
			return nil, err
		}
		p.count += n

		// Handle broken jpegs
		for buf[0] != 0xFF {
			buf[0] = buf[1]
			buf[1], err = p.in.ReadByte()
			if err != nil {
				return nil, err
			}
		}

		// Skip 00 markers
		if buf[1] == 0 {
			continue
		}

		// Skip stuffing
		for buf[1] == 0xFF {
			buf[1], err = p.in.ReadByte()
			if err != nil {
				return nil, err
			}
		}

		// We reached the end of the image
		if buf[1] == eoiMarker {
			break
		}

		if buf[1] == app0Marker || buf[1] == app1Marker || buf[1] == app2Marker {
			marker = buf[1]

			// Found the marker we're looking for
			seg := &Segment{
				MarkerID:   marker,
				MarkerName: markerNames[marker],
				Offset:     p.count,
			}
			seg.Size, n, err = getSize(p.in)
			if err != nil {
				return nil, err
			}
			p.count += n

			seg.Data = make([]byte, seg.Size)
			n, err = io.ReadFull(p.in, seg.Data)
			p.count += n

			segs = append(segs, *seg)
		} else {
			// Skip RST if need be
			if buf[1] >= rst0Marker && buf[1] <= rst7Marker {
				continue
			}

			size, n, err := getSize(p.in)
			if err != nil {
				return nil, err
			}
			p.count += n

			// Skip sections we're not looking for
			n64, err := io.CopyN(ioutil.Discard, p.in, int64(size))
			if err != nil {
				return nil, err
			}
			p.count += int(n64)
		}
	}

	return segs, nil
}

// GetSegment searches for the given marker and returns the first instance it encounters. Important: This does NOT find
// multiple instances of segments that might be split up, e.g. APP1.
func (p *Parser) GetSegment(marker uint8) (*Segment, error) {
	var buf [1024]byte
	var err error
	var n int

	for {
		n, err = io.ReadFull(p.in, buf[0:2])
		if err != nil {
			return nil, err
		}
		p.count += n

		// Handle broken jpegs
		for buf[0] != 0xFF {
			buf[0] = buf[1]
			buf[1], err = p.in.ReadByte()
			if err != nil {
				return nil, err
			}
		}

		// Skip 00 markers
		if buf[1] == 0 {
			continue
		}

		// Skip stuffing
		for buf[1] == 0xFF {
			buf[1], err = p.in.ReadByte()
			if err != nil {
				return nil, err
			}
		}

		// We reached the end of the image
		if buf[1] == eoiMarker {
			return nil, nil
		}

		if buf[1] == marker {
			// Found the marker we're looking for
			break
		} else {
			// Skip RST if need be
			if buf[1] >= rst0Marker && buf[1] <= rst7Marker {
				continue
			}

			size, n, err := getSize(p.in)
			if err != nil {
				return nil, err
			}
			p.count += n

			// Skip sections we're not looking for
			n64, err := io.CopyN(ioutil.Discard, p.in, int64(size))
			if err != nil {
				return nil, err
			}
			p.count += int(n64)
		}
	}

	seg := &Segment{
		MarkerID:   marker,
		MarkerName: markerNames[marker],
		Offset:     p.count,
	}
	seg.Size, n, err = getSize(p.in)
	if err != nil {
		return nil, err
	}
	p.count += n

	seg.Data = make([]byte, seg.Size)
	n, err = io.ReadFull(p.in, seg.Data)
	p.count += n
	return seg, nil
}

// getSize returns the segment length, the number of bytes read, and any error.
func getSize(input io.Reader) (int, int, error) {
	var buf [2]byte
	_, err := io.ReadFull(input, buf[0:2])
	if err != nil {
		return 0, 0, err
	}

	ret := int(buf[0])<<8 + int(buf[1]) - 2
	if ret < 0 {
		return ret, 2, errors.New("invalid segment length")
	}

	return ret, 2, nil
}
