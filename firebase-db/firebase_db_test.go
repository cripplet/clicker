package firebase_db

import (
	"flag"
	"fmt"
)

var credentials string
var project_root string

func init() {
	var test_root string
	var project string

	flag.StringVar(&credentials, "credentials", "", "")
	flag.StringVar(&test_root, "test_root", "", "")
	flag.StringVar(&project, "project", "", "")
	flag.Parse()

	project_root = fmt.Sprintf("https://%s.firebaseio.com/%s", project, test_root)
}
