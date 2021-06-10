# iccjpeg #

**Description**

A small utility package to extract ICC profiles from JPEG buffers. Forked from: https://github.com/vimeo/go-iccjpeg

**Installing**

```
go get github.com/snapas/img/iccjpeg
```

**API**

The API is a single function call:

```go
import "github.com/snapas/img/iccjpeg"

iccjpeg.GetICCBuf(input io.Reader) ([]byte, error)
```

It takes an `io.Reader` with a JPEG, and returns a buffer with the embedded ICC profile from that JPEG, if there is one. If there is not one, it returns an empty buffer.