package firebase_db

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

//  See https://firebase.google.com/docs/reference/rest/database/.

func do(c *http.Client, req *http.Request) ([]byte, int, error) {
	if c == nil {
		return nil, 0, errors.New("Unexpected nil HTTP client provided.")
	}
	resp, err := c.Do(req)
	if err != nil {
		return nil, 0, err
	}

	defer resp.Body.Close()
	b, err := ioutil.ReadAll(resp.Body)
	return b, resp.StatusCode, err
}

func Get(
	c *http.Client,
	path string,
	x_firebase_etag bool,
	query_parameters map[string]string,
	v interface{}) ([]byte, int, error) {

	path += paramToURL(query_parameters)
	req, err := http.NewRequest(http.MethodGet, path, nil)
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	if x_firebase_etag {
		req.Header.Set("X-Firebase-ETag", "true")
	}

	b, statusCode, err := do(c, req)
	if err == nil && v != nil {
		err = json.Unmarshal(b, v)
	}

	return b, statusCode, err
}

func Put(
	c *http.Client,
	path string,
	data []byte,
	x_firebase_etag bool,
	if_match string,
	query_parameters map[string]string,
	v interface{}) ([]byte, int, error) {

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

	b, statusCode, err := do(c, req)
	if err == nil && v != nil {
		err = json.Unmarshal(b, v)
	}

	return b, statusCode, err
}

func Post(
	c *http.Client,
	path string,
	data []byte,
	x_firebase_etag bool,
	query_parameters map[string]string,
	v interface{}) ([]byte, int, error) {

	path += paramToURL(query_parameters)
	req, err := http.NewRequest(http.MethodPost, path, bytes.NewReader(data))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	if x_firebase_etag {
		req.Header.Set("X-Firebase-ETag", "true")
	}

	b, statusCode, err := do(c, req)
	if err == nil && v != nil {
		err = json.Unmarshal(b, v)
	}

	return b, statusCode, err
}

func Patch(
	c *http.Client,
	path string,
	data []byte,
	query_parameters map[string]string,
	v interface{}) ([]byte, int, error) {

	path += paramToURL(query_parameters)
	req, err := http.NewRequest(http.MethodPatch, path, bytes.NewReader(data))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")

	b, statusCode, err := do(c, req)
	if err == nil && v != nil {
		err = json.Unmarshal(b, v)
	}

	return b, statusCode, err
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

	_, statusCode, err := do(c, req)
	return nil, statusCode, err
}
