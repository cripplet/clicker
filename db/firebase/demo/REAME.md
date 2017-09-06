Firebase Bindings Demo
====

**N.B.**: This basically recreates [previous work](https://www.captaincodeman.com/2017/06/09/access-firebase-database-and-authenticated-user-accounts-using-go).

Authenticating with JSON credentials:

```golang
import (
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"net/http"
)       

def main() {
	config, _ := google.JWTConfigFromJSON(
		jsonKey []byte, // raw file content, read using ioutil.ReadFile
		scope ...string)
	)
	c := config.Client(oauth2.NoContext) // HTTP client with wrapped auth token headers.
	c.Do(req *http.Request)
}       
```

It looks like the auth token does not have an expiration date; `config.Expires`
has a value of 0s, and issuing two authenticated requests more than 3600s apart
was successful.

At the very least, it looks like `config.Client` auto-refreshes the auth token.

**N.B.**: Service account needs to be granted Project.Editor IAM permission.        
