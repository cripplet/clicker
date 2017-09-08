package firebase_db

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
)

//  See https://firebase.google.com/docs/reference/rest/database/.

func do(c *http.Client, req *http.Request) ([]byte, int, error) {
	resp, err := c.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	return b, resp.StatusCode, err
}

func paramToURL(p map[string]string) string {
	s := ""
	if len(p) > 0 {
		s += "?"
	}
	for k, v := range p {
		s += fmt.Sprintf("%s=%s&", k, v)
	}
	var last_char int = 0
	if len(s) > 0 {
		last_char = len(s) - 1
	}
	return s[:last_char]
}

func Get(
	c *http.Client,
	path string,
	x_firebase_etag bool,
	query_parameters map[string]string) ([]byte, int, error) {

	path += paramToURL(query_parameters)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	if x_firebase_etag {
		req.Header.Set("X-Firebase-ETag", "true")
	}

	return do(c, req)
}

func Put(
	c *http.Client,
	path string,
	data []byte,
	x_firebase_etag bool,
	if_match string,
	query_parameters map[string]string) ([]byte, int, error) {

	path += paramToURL(query_parameters)
	req, err := http.NewRequest(http.MethodPut, path, bytes.NewReader(data))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	if x_firebase_etag {
		req.Header.Set("X-Firebase-ETag", "true")
	}
	if if_match != "" {
		req.Header.Set("if-match", if_match)
	}

	return do(c, req)
}

func Post(
	c *http.Client,
	path string,
	data []byte,
	x_firebase_etag bool,
	query_parameters map[string]string) ([]byte, int, error) {

	path += paramToURL(query_parameters)
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(data))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	if x_firebase_etag {
		req.Header.Set("X-Firebase-ETag", "true")
	}

	return do(c, req)
}

func Patch(
	c *http.Client,
	path string,
	data []byte,
	query_parameters map[string]string) ([]byte, int, error) {

	path += paramToURL(query_parameters)
	req, err := http.NewRequest(http.MethodPatch, path, bytes.NewReader(data))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	return do(c, req)
}

func Delete(
	c *http.Client,
	path string,
	x_firebase_etag bool,
	if_match string,
	query_parameters map[string]string) ([]byte, int, error) {

	path += paramToURL(query_parameters)
	req, err := http.NewRequest(http.MethodDelete, path, nil)
	if err != nil {
		return nil, 0, err
	}

	if x_firebase_etag {
		req.Header.Set("X-Firebase-ETag", "true")
	}
	if if_match != "" {
		req.Header.Set("if-match", if_match)
	}

	_, status_code, err := do(c, req)
	return nil, status_code, err
}
