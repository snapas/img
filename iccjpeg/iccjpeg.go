// Copyright 2014, Vimeo, LLC. All rights reserved.
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

func getSize(input io.Reader) (int, error) {
	var buf [2]byte
	_, err := io.ReadFull(input, buf[0:2])
	if err != nil {
		return 0, err
	}

	ret := int(buf[0])<<8 + int(buf[1]) - 2
	if ret < 0 {
		return ret, errors.New("Invalid segment length")
	}

	return ret, nil
}

// GetICCBuf reads a JPEG from input and returns a buffer containing the ICC profile.
// If no ICC profile is present, then the buffer may be of length 0.
func GetICCBuf(input io.Reader) ([]byte, error) {
	var buf [2048]byte
	var err error
	in := bufio.NewReader(input)

	_, err = io.ReadFull(in, buf[0:2])
	if err != nil {
		return nil, err
	} else if buf[0] != 0xFF && buf[1] != soiMarker {
		return nil, errors.New("no SOI Marker")
	}

	var iccData [][]byte
	iccLength := 0
	readProfs := 0
	numMarkers := -1
	for {
		_, err = io.ReadFull(in, buf[0:2])
		if err != nil {
			return nil, err
		}

		// Handle broken jpegs
		for buf[0] != 0xFF {
			buf[0] = buf[1]
			buf[1], err = in.ReadByte()
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
			buf[1], err = in.ReadByte()
			if err != nil {
				return nil, err
			}
		}

		// We reached the end of the image
		if buf[1] == eoiMarker {
			break
		}

		// Are we at an APP2 marker?
		if buf[1] != app2Marker {
			// Skip RST if need be
			if buf[1] >= rst0Marker && buf[1] <= rst7Marker {
				continue
			}

			size, err := getSize(in)
			if err != nil {
				return nil, err
			}

			// Skip non-APP2
			_, err = io.CopyN(ioutil.Discard, in, int64(size))
			if err != nil {
				return nil, err
			}
			continue
		}

		size, err := getSize(in)
		if err != nil {
			return nil, err
		} else if size < iccHeaderLen {
			return nil, errors.New("ICC segment invalid")
		}

		_, err = io.ReadFull(in, buf[0:12])
		if err != nil {
			return nil, err
		}

		if string(buf[0:11]) != "ICC_PROFILE" || buf[11] != 0 {
			return nil, errors.New("ICC segment invalid")
		}

		seqN, err := in.ReadByte()
		if err != nil {
			return nil, err
		} else if seqN == 0 {
			return nil, errors.New("invalid sequence number")
		}

		num, err := in.ReadByte()
		if err != nil {
			return nil, err
		} else if numMarkers == -1 {
			numMarkers = int(num)
			iccData = make([][]byte, numMarkers)
		} else if int(num) != numMarkers {
			return nil, errors.New("invalid ICC segment (numMarkers != cur_num_markers)")
		}

		if int(seqN) > numMarkers {
			return nil, errors.New("invalid ICC segment (seqN > numMarkers)")
		}

		iccData[seqN-1] = make([]byte, size-iccHeaderLen)
		_, err = io.ReadFull(in, iccData[seqN-1])
		if err != nil {
			return nil, err
		}

		iccLength += size - iccHeaderLen
		readProfs++

		if readProfs == numMarkers {
			ret := make([]byte, 0, iccLength)
			for _, data := range iccData {
				ret = append(ret, data...)
			}
			return ret, nil
		}
	}

	return nil, nil
}
