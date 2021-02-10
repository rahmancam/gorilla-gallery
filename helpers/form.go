package helpers

import (
	"net/http"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

// ParseForm form
func ParseForm(r *http.Request, dst interface{}) error {
	if err := r.ParseForm(); err != nil {
		return err
	}

	if err := decoder.Decode(dst, r.PostForm); err != nil {
		return err
	}

	return nil
}
