package firebase_db

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
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

	b, status_code, err := Put(
		c,
		fmt.Sprintf("%s/test_put.json", project_root),
		expected,
		false,
		"",
		map[string]string{},
	)

	if err != nil {
		t.Errorf("Unexpected PUT error: %v", err)
	}

	if status_code != http.StatusOK {
		t.Errorf("Unexpected HTTP error: %d", status_code)
	}

	if !bytes.Equal(b, expected) {
		t.Errorf("Returned PUT data is not expected: %s != %s", string(b), string(expected))
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
	)

	b, status_code, err := Get(
		c,
		fmt.Sprintf("%s/test_get.json", project_root),
		false,
		map[string]string{},
	)

	if err != nil {
		t.Errorf("Unexpected GET error: %v", err)
	}

	if status_code != http.StatusOK {
		t.Errorf("Unexpected HTTP error: %d", status_code)
	}

	if !bytes.Equal(b, expected) {
		t.Errorf("Returned GET data is not expected: %s != %s", string(b), string(expected))
	}
}

func TestPost(t *testing.T) {
	type PostResponseStruct struct {
		Name string `json:"name"`
	}

	ResetEnvironment(t)

	c, _ := NewGoogleClient(credentials)
	expected, _ := json.Marshal("test-element")

	b, status_code, err := Post(
		c,
		fmt.Sprintf("%s/test_post.json", project_root),
		expected,
		false,
		map[string]string{},
	)

	if err != nil {
		t.Errorf("Unexpected POST error: %v", err)
	}

	if status_code != http.StatusOK {
		t.Errorf("Unexpected HTTP error: %d", status_code)
	}

	r := PostResponseStruct{}
	json.Unmarshal(b, &r)

	b, _, _ = Get(
		c,
		fmt.Sprintf("%s/test_post.json", project_root),
		false,
		map[string]string{},
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
	)

	data, _ := json.Marshal(PatchDataStruct{
		Last: "Dillinger",
	})

	b, status_code, err := Patch(
		c,
		fmt.Sprintf("%s/test_patch.json", project_root),
		data,
		map[string]string{},
	)

	if err != nil {
		t.Errorf("Unexpected DELETE error: %v", err)
	}

	if status_code != http.StatusOK {
		t.Errorf("Unexpected HTTP error: %d", status_code)
	}

	if !bytes.Equal(b, data) {
		t.Error("Returned PATCH data is not expected: %s != %s", string(b), string(data))
	}

	b, _, _ = Get(
		c,
		fmt.Sprintf("%s/test_patch.json", project_root),
		false,
		map[string]string{},
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
	)

	_, status_code, err := Delete(
		c,
		fmt.Sprintf("%s/test_delete.json", project_root),
		false,
		"",
		map[string]string{},
	)

	if err != nil {
		t.Errorf("Unexpected DELETE error: %v", err)
	}

	if status_code != http.StatusOK {
		t.Errorf("Unexpected HTTP error: %d", status_code)
	}

	b, _, _ := Get(
		c,
		fmt.Sprintf("%s/test_delete.json", project_root),
		false,
		map[string]string{},
	)

	expected, _ := json.Marshal(nil)

	if !bytes.Equal(b, expected) {
		t.Errorf("Returned GET data is not expected: %s != %s", string(b), string(expected))
	}
}