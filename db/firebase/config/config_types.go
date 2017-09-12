package cc_fb_config

import (
	"net/http"
)

type EnvironmentType int

const (
	DEV EnvironmentType = iota
	PROD
)

var ENVIRONMENT_TYPE_LOOKUP map[EnvironmentType]string = map[EnvironmentType]string{
	DEV:  "dev",
	PROD: "prod",
}

var ENVIRONMENT_TYPE_REVERSE_LOOKUP map[string]EnvironmentType = map[string]EnvironmentType{
	"dev":  DEV,
	"prod": PROD,
}

type CCFirebaseConfig struct {
	Client      *http.Client
	ProjectPath string
	Environment EnvironmentType
}
