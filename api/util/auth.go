package util

import (
	"errors"
	"net/http"

	"github.com/fortytw2/abdi"
	"github.com/fortytw2/embercrest/datastore"
	"github.com/fortytw2/embercrest/user"
)

// Authenticate checks for correct credentials via HTTP basic auth
func Authenticate(r *http.Request, users datastore.UserService) (u *user.User, err error) {
	username, password, ok := r.BasicAuth()
	if !ok {
		err = errors.New("basic auth not OK")
		return
	}
	u, err = users.GetUser(username)
	if err != nil {
		return
	}
	if err = abdi.Check(password, u.PasswordHash); err != nil {
		return
	}

	return
}
