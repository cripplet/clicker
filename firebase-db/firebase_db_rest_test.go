package firebase_db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func ResetEnvironment(t *testing.T) {
	c, _ := NewGoogleClient(credentials)

	request, _ := http.NewRequest(
		http.MethodDelete,
		fmt.Sprintf("%s.json", project_root),
		nil,
	)
	_, err := c.Do(request)
	if err != nil {
		t.Errorf("Could not delete test DB, failed with error: %v", err)
	}
}

func TestParamToURLEmpty(t *testing.T) {
	p := paramToURL(map[string]string{})
	if p != "" {
		t.Errorf("Generated non-empty string %s with empty parameters", p)
	}
}

func TestParamToURLSingle(t *testing.T) {
	p := paramToURL(map[string]string{
		"key1": "value1",
	})
	if p != "?key1=value1" {
		t.Errorf("Generated URL query parameter mismatch: %s != %s", p, "?key1=value1")
	}
}

func TestResetEnvironment(t *testing.T) {
	ResetEnvironment(t)
}

func TestPut(t *testing.T) {
	ResetEnvironment(t)

	c, _ := NewGoogleClient(credentials)
	expected, _ := json.Marshal("test-data")

	b, statusCode, _, err := Put(
		c,
		fmt.Sprintf("%s/test_put.json", project_root),
		expected,
		false,
		"",
		map[string]string{},
		nil,
	)

	if err != nil {
		t.Errorf("Unexpected PUT error: %v", err)
	}

	if statusCode != http.StatusOK {
		t.Errorf("Unexpected HTTP error: %d", statusCode)
	}

	if !bytes.Equal(b, expected) {
		t.Errorf("Returned PUT data is not expected: %s != %s", string(b), string(expected))
	}
}

func TestReturnJSON(t *testing.T) {
	type TestData struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}

	ResetEnvironment(t)

	c, _ := NewGoogleClient(credentials)

	expected := TestData{
		Name: "Michael Scott",
		ID:   "dunder-mifflin-2005",
	}
	data, _ := json.Marshal(expected)
	actual := TestData{}

	Put(
		c,
		fmt.Sprintf("%s/test_return_json.json", project_root),
		data,
		false,
		"",
		map[string]string{},
		&actual,
	)

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("Expected return did not match actual returned JSON: %v != %v", expected, actual)
	}
}

func TestGet(t *testing.T) {
	ResetEnvironment(t)

	c, _ := NewGoogleClient(credentials)
	expected, _ := json.Marshal("test-data")

	Put(
		c,
		fmt.Sprintf("%s/test_get.json", project_root),
		expected,
		false,
		"",
		map[string]string{},
		nil,
	)

	b, statusCode, _, err := Get(
		c,
		fmt.Sprintf("%s/test_get.json", project_root),
		false,
		map[string]string{},
		nil,
	)

	if err != nil {
		t.Errorf("Unexpected GET error: %v", err)
	}

	if statusCode != http.StatusOK {
		t.Errorf("Unexpected HTTP error: %d", statusCode)
	}

	if !bytes.Equal(b, expected) {
		t.Errorf("Returned GET data is not expected: %s != %s", string(b), string(expected))
	}
}

func TestGetETag(t *testing.T) {
	ResetEnvironment(t)
	c, _ := NewGoogleClient(credentials)

	data, _ := json.Marshal("test-data")
	Put(
		c,
		fmt.Sprintf("%s/test_get_etag.json", project_root),
		data,
		false,
		"",
		map[string]string{},
		nil,
	)

	_, statusCode, eTag, err := Get(
		c,
		fmt.Sprintf("%s/test_get_etag.json", project_root),
		true,
		map[string]string{},
		nil,
	)

	if err != nil {
		t.Errorf("Unexpected GET error: %v", err)
	}

	if statusCode != http.StatusOK {
		t.Errorf("Unexpected HTTP error: %d", statusCode)
	}

	if eTag == "null_etag" {
		t.Error("Unexpected null ETag")
	}
}

func TestPutETagMismatch(t *testing.T) {
	ResetEnvironment(t)
	c, _ := NewGoogleClient(credentials)

	data, _ := json.Marshal("test-data")
	_, statusCode, _, err := Put(
		c,
		fmt.Sprintf("%s/test_put_etag_mismatch.json", project_root),
		data,
		false,
		"some-invalid-etag",
		map[string]string{},
		nil,
	)

	if err != nil {
		t.Errorf("Unexpected GET error: %v", err)
	}

	if statusCode != http.StatusPreconditionFailed {
		t.Errorf("Unexpected HTTP response code: %d", statusCode)
	}
}

func TestPutETag(t *testing.T) {
	ResetEnvironment(t)
	c, _ := NewGoogleClient(credentials)

	data, _ := json.Marshal("test-data")
	_, statusCode, _, err := Put(
		c,
		fmt.Sprintf("%s/test_put_etag_mismatch.json", project_root),
		data,
		false,
		"null_etag",
		map[string]string{},
		nil,
	)

	if err != nil {
		t.Errorf("Unexpected GET error: %v", err)
	}

	if statusCode != http.StatusOK {
		t.Errorf("Unexpected HTTP response code: %d", statusCode)
	}
}

func TestPost(t *testing.T) {
	type PostResponseStruct struct {
		Name string `json:"name"`
	}

	ResetEnvironment(t)

	c, _ := NewGoogleClient(credentials)
	expected, _ := json.Marshal("test-element")

	b, statusCode, _, err := Post(
		c,
		fmt.Sprintf("%s/test_post.json", project_root),
		expected,
		false,
		map[string]string{},
		nil,
	)

	if err != nil {
		t.Errorf("Unexpected POST error: %v", err)
	}

	if statusCode != http.StatusOK {
		t.Errorf("Unexpected HTTP error: %d", statusCode)
	}

	r := PostResponseStruct{}
	json.Unmarshal(b, &r)

	b, _, _, _ = Get(
		c,
		fmt.Sprintf("%s/test_post.json", project_root),
		false,
		map[string]string{},
		nil,
	)

	resp := map[string]string{}
	json.Unmarshal(b, &resp)

	data, _ := json.Marshal(resp[r.Name])

	if !bytes.Equal(expected, data) {
		t.Errorf("Returned GET data is not expected: %s != %s", string(expected), string(data))
	}
}

func TestPatch(t *testing.T) {
	type PatchDataStruct struct {
		First string `json:"first,omitempty"`
		Last  string `json:"last,omitempty"`
	}

	ResetEnvironment(t)

	c, _ := NewGoogleClient(credentials)

	initial, _ := json.Marshal(PatchDataStruct{
		First: "John",
		Last:  "Smith",
	})

	expected, _ := json.Marshal(PatchDataStruct{
		First: "John",
		Last:  "Dillinger",
	})

	Put(
		c,
		fmt.Sprintf("%s/test_patch.json", project_root),
		initial,
		false,
		"",
		map[string]string{},
		nil,
	)

	data, _ := json.Marshal(PatchDataStruct{
		Last: "Dillinger",
	})

	b, statusCode, _, err := Patch(
		c,
		fmt.Sprintf("%s/test_patch.json", project_root),
		data,
		map[string]string{},
		nil,
	)

	if err != nil {
		t.Errorf("Unexpected DELETE error: %v", err)
	}

	if statusCode != http.StatusOK {
		t.Errorf("Unexpected HTTP error: %d", statusCode)
	}

	if !bytes.Equal(b, data) {
		t.Error("Returned PATCH data is not expected: %s != %s", string(b), string(data))
	}

	b, _, _, _ = Get(
		c,
		fmt.Sprintf("%s/test_patch.json", project_root),
		false,
		map[string]string{},
		nil,
	)

	if !bytes.Equal(b, expected) {
		t.Errorf("Returned GET data is not expected: %s != %s", string(b), string(expected))
	}
}

func TestDelete(t *testing.T) {
	ResetEnvironment(t)

	c, _ := NewGoogleClient(credentials)
	data, _ := json.Marshal("test-data")

	Put(
		c,
		fmt.Sprintf("%s/test_delete.json", project_root),
		data,
		false,
		"",
		map[string]string{},
		nil,
	)

	_, statusCode, _, err := Delete(
		c,
		fmt.Sprintf("%s/test_delete.json", project_root),
		false,
		"",
		map[string]string{},
	)

	if err != nil {
		t.Errorf("Unexpected DELETE error: %v", err)
	}

	if statusCode != http.StatusOK {
		t.Errorf("Unexpected HTTP error: %d", statusCode)
	}

	b, _, _, _ := Get(
		c,
		fmt.Sprintf("%s/test_delete.json", project_root),
		false,
		map[string]string{},
		nil,
	)

	expected, _ := json.Marshal(nil)

	if !bytes.Equal(b, expected) {
		t.Errorf("Returned GET data is not expected: %s != %s", string(b), string(expected))
	}
}
