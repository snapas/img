package img

import (
	"log"
	"os"
	"testing"
)

func TestDecode(t *testing.T) {
	f, err := os.Open("testdata/porto-1.jpg")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	_, _, err = Decode(f)
	if err != nil {
		t.Fatal("Decode failed:", err)
	}
}
