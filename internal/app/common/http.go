package common

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"
)

var ErrorStatusCode = errors.New("ups, something has not working, please wait a moment and re-try")

func DoRequestJSON(method string, url string, body interface{}) (*http.Response, error) {
	var buff io.Reader

	switch body := body.(type) {
	case *bytes.Buffer:
		buff = body
	default:
		b, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		buff = bytes.NewBuffer(b)
	}

	req, err := http.NewRequest(method, url, buff)
	if err != nil {
		return nil, err
	}
	if body != nil {
		req.Header.Add("Content-Type", "application/json")
	}
	resp, err := http.DefaultClient.Do(req)

	return resp, err
}

func Concat(strs ...string) string {
	return strings.Join(strs, "")
}
