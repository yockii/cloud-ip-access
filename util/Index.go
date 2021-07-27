package util

import (
	flate "compress/flate"
	"compress/gzip"
	"io"
	"net/http"
)

func init() {
	initConfig()

	initLog()
}

// 检测返回的body是否经过压缩，并返回解压的内容
func SwitchContentEncoding(res *http.Response) (bodyReader io.Reader, err error) {
	switch res.Header.Get("Content-Encoding") {
	case "gzip":
		bodyReader, err = gzip.NewReader(res.Body)
	case "deflate":
		bodyReader = flate.NewReader(res.Body)
	default:
		bodyReader = res.Body
	}
	return
}
