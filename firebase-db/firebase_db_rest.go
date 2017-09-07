package firebase_db

import (
	"net/http"
)

//  See https://firebase.google.com/docs/reference/rest/database/.

func Get(
	c *http.Client,
	path string,
	x_firebase_etag bool,
	access_token string,
	query_parameters map[string]string,
	resp *[]byte) (int, error) {
	return http.StatusOK, nil
}

func Put(
	c *http.Client,
	path string,
	req *[]byte,
	x_firebase_etag bool,
	if_match string,
	query_parameters map[string]string,
	resp *[]byte) (int, error) {
	return http.StatusOK, nil
}

func Post(
	c *http.Client,
	path string,
	req *[]byte,
	x_firebase_etag bool,
	query_parameters map[string]string,
	resp *[]byte) (int, error) {
	return http.StatusOK, nil
}

func Patch(
	c *http.Client,
	path string,
	req *[]byte,
	query_parameters map[string]string,
	resp *[]byte) (int, error) {
	return http.StatusOK, nil
}

func Delete(
	c *http.Client,
	path string,
	x_firebase_etag bool,
	if_match string,
	query_parameters map[string]string,
	resp *[]byte) (int, error) {
	return http.StatusOK, nil
}
