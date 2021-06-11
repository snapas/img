package iccjpeg

import (
	"log"
	"os"
	"testing"
)

func TestGetICCRaw(t *testing.T) {
	f, err := os.Open("../testdata/porto-1.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	icc, err := GetICCRaw(f)
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", icc)
}
