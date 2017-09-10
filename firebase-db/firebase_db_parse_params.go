package firebase_db

import (
	"fmt"
)

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
