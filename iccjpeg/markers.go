package iccjpeg

const (
	// JPEG Markers
	soiMarker  = 0xD8
	eoiMarker  = 0xD9
	app0Marker = 0xE0
	app1Marker = 0xE1
	app2Marker = 0xE2
	rst0Marker = 0xD0
	rst7Marker = 0xD7
)

var markerNames = map[byte]string{
	soiMarker:  "SOI",
	eoiMarker:  "EOI",
	app0Marker: "APP0",
	app1Marker: "APP1",
	app2Marker: "APP2",
}
