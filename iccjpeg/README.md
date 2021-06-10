# iccjpeg #

**Description**

A small utility package to extract a raw ICC profile from a JPEG, meant for later writing directly to another JPEG. This is useful for retaining this metadata after transforming the image in some way, as the standard `image/jpeg` library normally discards this information.

This library is based on: https://github.com/vimeo/go-iccjpeg

**Installing**

```
go get github.com/snapas/img/iccjpeg
```

**API**

The API is a single function call:

```go
import "github.com/snapas/img/iccjpeg"

iccjpeg.GetICCRaw(input io.Reader) ([]byte, error)
```

It takes an `io.Reader` with a JPEG, and returns the embedded ICC profile from that JPEG, including header and segment size information. If there is no profile, it returns an empty slice.