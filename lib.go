package goqueryja

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/transform"
)

func NewDocument(url_ string) (*goquery.Document, error) {
	res, err := http.Get(url_)
	if err != nil {
		return nil, err
	}
	encoding := GetResponseEncoding(res)
	reader, err := NewUTF8Reader(res.Body, encoding)
	if err != nil {
		return nil, err
	}
	return goquery.NewDocumentFromReader(reader)
}

func NewUTF8Reader(r io.Reader, encoding string) (io.Reader, error) {
	switch strings.ToLower(encoding) {
	case "utf-8":
		return r, nil
	case "euc-jp":
		return transform.NewReader(r, japanese.EUCJP.NewDecoder()), nil
	case "shift_jis":
		return transform.NewReader(r, japanese.ShiftJIS.NewDecoder()), nil
	case "iso-2022-jp":
		return transform.NewReader(r, japanese.ISO2022JP.NewDecoder()), nil
	default:
		return nil, fmt.Errorf("unsupported encoding: %s", encoding)
	}
}

func GetResponseEncoding(res *http.Response) string {
	contentType := res.Header.Get("Content-Type")
	contentTypeLower := strings.ToLower(contentType)
	if index := strings.Index(contentTypeLower, "charset="); index != -1 {
		return contentType[index+len("charset="):]
	} else {
		return ""
	}
}
