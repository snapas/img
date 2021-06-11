// Copyright 2014, Vimeo, LLC. All rights reserved.
// Copyright 2021 A Bunch Tell LLC. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found here: https://github.com/vimeo/go-iccjpeg/blob/master/LICENSE

// Package iccjpeg implements ICC profile extraction from JPEG files.
package iccjpeg

import (
	"bufio"
	"errors"
	"io"
	"io/ioutil"
)

const (
	// JPEG Markers
	soiMarker  = 0xD8
	eoiMarker  = 0xD9
	app2Marker = 0xE2
	rst0Marker = 0xD0
	rst7Marker = 0xD7

	// Others
	iccHeaderLen = 14
)

type (
	Parser struct {
		count int
		in    *bufio.Reader
	}

	Segment struct {
		MarkerID   byte
		MarkerName string
		Size       int
		Offset     int
		Data       []byte
	}
)

func NewParser(input io.Reader) *Parser {
	return &Parser{
		in: bufio.NewReader(input),
	}
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
			var s *Segment
			return s, nil
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

// GetICCRaw reads a JPEG from input and returns a buffer containing the raw ICC profile data.
// If no ICC profile is present, then the buffer may be of length 0.
func GetICCRaw(input io.Reader) ([]byte, error) {
	var err error

	/*
		var iccData [][]byte
		iccLength := 0
	*/
	readProfs := 0
	numMarkers := -1
	p := NewParser(input)
	p.ReadSOI()
	seg, err := p.GetSegment(app2Marker)
	if err != nil {
		return nil, err
	}
	if seg.Size < iccHeaderLen {
		return nil, errors.New("ICC segment invalid")
	}

	i := 11
	if string(seg.Data[:i]) != "ICC_PROFILE" || seg.Data[i] != 0 {
		return nil, errors.New("ICC segment invalid")
	}
	i++

	seqN := seg.Data[i]
	i++
	if seqN == 0 {
		return nil, errors.New("invalid sequence number")
	}

	num := seg.Data[i]
	i++
	if numMarkers == -1 {
		numMarkers = int(num)
		//iccData = make([][]byte, numMarkers)
	} else if int(num) != numMarkers {
		return nil, errors.New("invalid ICC segment (numMarkers != cur_num_markers)")
	}

	if int(seqN) > numMarkers {
		return nil, errors.New("invalid ICC segment (seqN > numMarkers)")
	}

	/*
		// Non-raw data
		iccData[seqN-1] = make([]byte, seg.Size-iccHeaderLen)
		copy(iccData[seqN-1], seg.Data[i:])
		iccLength += seg.Size - iccHeaderLen
	*/

	readProfs++

	if readProfs == numMarkers {
		return seg.Data, nil
	}

	return nil, nil
}
