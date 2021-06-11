package iccjpeg

import (
	"log"
	"os"
	"testing"
)

func TestGetCommonAppSegments(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{"orientation 1, ICC profile", "porto-1.jpg"},
		{"orientation 3, no ICC profile", "holden-3-noicc.jpg"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := os.Open("../testdata/" + test.filename)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			p := NewParser(f)
			err = p.ReadSOI()
			if err != nil {
				t.Fatal("ReadSOI failed:", err)
			}
			seg, err := p.GetCommonAppSegments()
			if err != nil {
				t.Fatal("GetCommonAppSegments failed:", err)
			}
			if seg == nil {
				t.Fatal("EOF")
			}
			for _, s := range seg {
				t.Logf("ID %x Name %s Size %d Offset %d", s.MarkerID, s.MarkerName, s.Size, s.Offset)
			}
		})
	}
}

func TestGetSegment(t *testing.T) {
	f, err := os.Open("../testdata/porto-1.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	p := NewParser(f)
	err = p.ReadSOI()
	if err != nil {
		t.Fatal("ReadSOI failed:", err)
	}
	seg, err := p.GetSegment(app2Marker)
	if err != nil {
		t.Fatal("GetSegment failed:", err)
	}
	if seg == nil {
		t.Fatal("EOF")
	}
	t.Logf("ID %x Name %s Size %d Offset %d", seg.MarkerID, seg.MarkerName, seg.Size, seg.Offset)
}
