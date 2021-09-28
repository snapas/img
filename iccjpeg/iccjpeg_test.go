package iccjpeg

import (
	"log"
	"os"
	"testing"
)

func TestGetICCRaw(t *testing.T) {
	tests := []struct {
		name     string
		filename string
	}{
		{"ICC profile", "porto-1.jpg"},
		{"no ICC profile", "holden-3-noicc.jpg"},
		{"unsure ICC profile", "gopro.jpg"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			f, err := os.Open("../testdata/" + test.filename)
			if err != nil {
				log.Fatal(err)
			}
			defer f.Close()

			_, err = GetICCRaw(f)
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
