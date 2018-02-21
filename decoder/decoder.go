package decoder

import (
	"compress/gzip"
	"compress/zlib"
	"fmt"
	"net/http"

	"io/ioutil"
	"strings"

	"github.com/pkg/errors"

	"github.com/elastic/apm-server/utility"
	"github.com/elastic/apm-server/models"
	"bytes"
)

type Decoder func(req *http.Request) (*models.TransactionsPayload, error)

func DecodeLimitJSONData(maxSize int64) Decoder {
	return func(req *http.Request) (*models.TransactionsPayload, error) {
		contentType := req.Header.Get("Content-Type")
		if !strings.Contains(contentType, "application/json") {
			return nil, fmt.Errorf("invalid content type: %s", req.Header.Get("Content-Type"))
		}

		reader := req.Body
		if reader == nil {
			return nil, errors.New("no content")
		}

		switch req.Header.Get("Content-Encoding") {
		case "deflate":
			var err error
			reader, err = zlib.NewReader(reader)
			if err != nil {
				return nil, err
			}

		case "gzip":
			var err error
			reader, err = gzip.NewReader(reader)
			if err != nil {
				return nil, err
			}
		}

		buf := new(bytes.Buffer)

		if _, err := buf.ReadFrom(reader); err != nil {
			return nil, errors.Wrap(err, "data read error")
		}
		by := buf.Bytes()
		//fmt.Println(string(by))
		return DecodeTransactionsPayload(by)
	}
}


func DecodeTransactionsPayload(bytes []byte) (*models.TransactionsPayload, error) {
	payload := models.TransactionsPayload{}
	if err := payload.UnmarshalJSON(bytes); err != nil {
		return nil, errors.Wrap(err, "data read error")
	}
	return &payload, nil
}


func DecodeSourcemapFormData(req *http.Request) (map[string]interface{}, error) {
	contentType := req.Header.Get("Content-Type")
	if !strings.Contains(contentType, "multipart/form-data") {
		return nil, fmt.Errorf("invalid content type: %s", req.Header.Get("Content-Type"))
	}

	file, _, err := req.FormFile("sourcemap")
	if err != nil {
		return nil, err
	}
	defer file.Close()

	sourcemapBytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	payload := map[string]interface{}{
		"sourcemap":       string(sourcemapBytes),
		"service_name":    req.FormValue("service_name"),
		"service_version": req.FormValue("service_version"),
		"bundle_filepath": utility.CleanUrlPath(req.FormValue("bundle_filepath")),
	}

	return payload, nil
}
//
//func DecodeUserData(decoder Decoder) Decoder {
//	augment := func(req *http.Request) map[string]interface{} {
//		return map[string]interface{}{
//			"ip":         utility.ExtractIP(req),
//			"user_agent": req.Header.Get("User-Agent"),
//		}
//	}
//	return augmentData(decoder, "user", augment)
//}
//
//func DecodeSystemData(decoder Decoder) Decoder {
//	augment := func(req *http.Request) map[string]interface{} {
//		return map[string]interface{}{"ip": utility.ExtractIP(req)}
//	}
//	return augmentData(decoder, "system", augment)
//}
//
//func augmentData(decoder Decoder, key string, augment func(req *http.Request) map[string]interface{}) Decoder {
//	return func(req *http.Request) (map[string]interface{}, error) {
//		v, err := decoder(req)
//		if err != nil {
//			return v, err
//		}
//		utility.InsertInMap(v, key, augment(req))
//		return v, nil
//	}
//}
