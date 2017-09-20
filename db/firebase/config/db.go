package cc_fb_config

import (
	"flag"
	"fmt"
	"github.com/cripplet/clicker/firebase-db"
)

var CC_FIREBASE_CONFIG CCFirebaseConfig

func init() {
	var credentials string
	var project string
	var env_string string

	flag.StringVar(&credentials, "credentials", "", "Google JSON credentials path")
	flag.StringVar(&project, "project", "", "Firebase project ID")
	flag.StringVar(&env_string, "environment", "dev", "Firebase DB environment")
	flag.Parse()

	c, _ := firebase_db.NewGoogleClient(credentials)
	CC_FIREBASE_CONFIG.Client = c
	CC_FIREBASE_CONFIG.Environment = ENVIRONMENT_TYPE_REVERSE_LOOKUP[env_string]
	CC_FIREBASE_CONFIG.ProjectPath = fmt.Sprintf("https://%s.firebaseio.com/%s", project, ENVIRONMENT_TYPE_LOOKUP[CC_FIREBASE_CONFIG.Environment])
}
