package cc_firebase

type Environment int

const (
	ENVIRONMENT_DEV Environment = iota
	ENVIRONMENT_PROD
)

var ENVIRONMENT_LOOKUP map[Environment]string = map[Environment]string{
	ENVIRONMENT_DEV:  "dev",
	ENVIRONMENT_PROD: "prod",
}
