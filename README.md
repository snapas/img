# img

[![Go Reference](https://pkg.go.dev/badge/github.com/snapas.svg)](https://pkg.go.dev/github.com/snapas)

Drop-in functionality for decoding and encoding images, with some sugar added, including:

* Metadata preservation
* JPEG auto-rotation

This is useful for building consumer-facing services and tools that manipulate images without losing important data along the way. Above all, this aims to fill a void left by the standard Go library, where manipulating images also means losing their metadata.

## Goals

This aims to support basic functions around image manipulation in a pure Go implementation. This library was built for [Snap.as](https://snap.as), and will eventually become a part of [WriteFreely](https://github.com/writefreely/writefreely).

## HELP!

Frankly, I don't really know what I'm doing. I read Wikipedia, skim the spec, fork some code, and bang on it until it works. I could use some professional help developing this into a rock-solid library. So if you know image formats, please [reach out](mailto:hello@write.as), and I will pay you to help make that happen. 

## License

BSD 3-Clause License. Copyright 2021 A Bunch Tell LLC, and respective authors of libraries included herein.