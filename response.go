// Package client provides cient request itilities.
package client

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

// Response is a request response wrapping the original *http.Request.
type Response struct {
	*http.Response
}

func (r *Response) String() string {
	return fmt.Sprint(r.Response)
}

// Close closes the response body.
// It must be called after body is no longer used.
func (r *Response) Close() error {
	return r.Body.Close()
}

// Cookie returns the named cookie provided in the request or http.ErrNoCookie if not found.
// If multiple cookies match the given name, only one cookie will be returned.
func (r *Response) Cookie(name string) (*http.Cookie, error) {
	for _, c := range r.Cookies() {
		if c.Name == name {
			return c, nil
		}
	}
	return nil, http.ErrNoCookie
}

// BodyBytes returns the response body as bytes.
func (r *Response) BodyBytes() ([]byte, error) {
	return ioutil.ReadAll(r.Body)
}

// BodyString returns the response body as a string.
func (r *Response) BodyString() (string, error) {
	b, err := r.BodyBytes()
	return string(b), err
}

// JSON decodes a JSON body into v.
func (r *Response) JSON(v interface{}) error {
	err := json.NewDecoder(r.Body).Decode(v)
	if err != nil && err != io.EOF {
		return err
	}
	return nil
}

// Path returns the final response path.
func (r *Response) Path() string {
	return r.Request.URL.Path
}
