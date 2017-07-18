// Package client provides an HTTP client for clean requests.
package client

import (
	"bytes"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var httpNoRedirectClient = &http.Client{
	CheckRedirect: func(_ *http.Request, _ []*http.Request) error {
		return http.ErrUseLastResponse
	},
}

type body struct {
	buf *bytes.Buffer
	w   *multipart.Writer
}

type field struct{ key, value string }

// Request is a client request.
type Request interface {
	Value(key, value string) Request
	File(key, filename string, file io.Reader) Request
	OpenFile(key, filepath string) Request
	Header(key, value string) Request
	Cookie(*http.Cookie) Request
	UserAgent(value string) Request
	DisableRedirect() Request
	ForceMultipart() Request
	Do() (*Response, error)
	String() string
}

type request struct {
	method  string
	url     string
	header  http.Header
	fields  []*field // Keep ordered non-multipart values until we know the request content type.
	cookies []*http.Cookie
	body    *body
	err     error // Deferred error for easy chaining.
	noRedir bool  // Disable redirection following.
}

// New makes a new request for method and URL.
func New(method, url string) Request {
	return &request{method: method, url: url}
}

// Get makes a GET request.
func Get(url string) Request {
	return New(http.MethodGet, url)
}

// Post makes a POST request.
func Post(url string) Request {
	return New(http.MethodPost, url)
}

// Put makes a PUT request.
func Put(url string) Request {
	return New(http.MethodPut, url)
}

// Patch makes a PATCH request.
func Patch(url string) Request {
	return New(http.MethodPatch, url)
}

// Delete makes a DELETE request.
func Delete(url string) Request {
	return New(http.MethodDelete, url)
}

// Head makes a HEAD request.
func Head(url string) Request {
	return New(http.MethodHead, url)
}

// Value adds a form value.
func (r *request) Value(key, value string) Request {
	if r.err != nil {
		return r
	}
	if r.body == nil {
		r.fields = append(r.fields, &field{key, value})
		return r
	}
	err := r.body.w.WriteField(key, value)
	if err != nil {
		r.err = err
	}
	return r
}

// File adds a multipart form file.
func (r *request) File(key, filename string, file io.Reader) Request {
	if r.err != nil {
		return r
	}
	r.ForceMultipart()
	filew, err := r.body.w.CreateFormFile(key, filename)
	if err != nil {
		r.err = err
		return r
	}
	_, err = io.Copy(filew, file)
	if err != nil {
		r.err = err
	}
	return r
}

// OpenFile opens a file with filepath and adds it to the request form.
func (r *request) OpenFile(key, filepath string) Request {
	if r.err != nil {
		return r
	}
	file, err := os.Open(filepath)
	if err != nil {
		r.err = err
		return r
	}
	defer file.Close()
	return r.File(key, file.Name(), file)
}

// Header adds a header to the response.
func (r *request) Header(key, value string) Request {
	if r.err != nil {
		return r
	}
	if r.header == nil {
		r.header = make(http.Header)
	}
	r.header.Add(key, value)
	return r
}

// Cookie adds a cookie to the request.
func (r *request) Cookie(c *http.Cookie) Request {
	if r.err != nil {
		return r
	}
	r.cookies = append(r.cookies, c)
	return r
}

// UserAgent sets the User-Agent header.
func (r *request) UserAgent(value string) Request {
	r.Header("User-Agent", value)
	return r
}

// DisableRedirect avoids following response redirections.
func (r *request) DisableRedirect() Request {
	if r.err != nil {
		return r
	}
	r.noRedir = true
	return r
}

// ForceMultipart forces a "multipart/form-data" response, even with no files.
func (r *request) ForceMultipart() Request {
	if r.err != nil {
		return r
	}
	if r.body == nil {
		buf := new(bytes.Buffer)
		r.body = &body{buf, multipart.NewWriter(buf)}
	}
	for _, f := range r.fields { // We know request will be multipart: write fields to body.
		r.Value(f.key, f.value)
	}
	r.fields = nil
	return r
}

// Do sends the request and returns the response.
func (r *request) Do() (*Response, error) {
	if r.err != nil {
		return nil, r.err
	}

	var body io.Reader
	if r.body != nil {
		if err := r.body.w.Close(); err != nil {
			return nil, err
		}
		body = r.body.buf
	} else if len(r.fields) > 0 {
		uv := make(url.Values)
		for _, f := range r.fields {
			uv[f.key] = append(uv[f.key], f.value)
		}
		body = strings.NewReader(uv.Encode())
	}

	req, err := http.NewRequest(r.method, r.url, body)
	if err != nil {
		return nil, err
	}
	if len(r.header) > 0 {
		req.Header = r.header
	}
	if r.body != nil {
		req.Header.Set("Content-Type", r.body.w.FormDataContentType())
	}
	for _, c := range r.cookies {
		req.AddCookie(c)
	}

	var res *http.Response
	if r.noRedir {
		res, err = httpNoRedirectClient.Do(req)
	} else {
		res, err = http.DefaultClient.Do(req)
	}
	return &Response{Response: res}, err
}

func (r *request) String() string {
	s := r.method + " " + r.url + "\n"
	s += "\tHeader:\n"
	s += "\t\tContent-Type: "
	if r.body == nil {
		s += "application/x-www-form-urlencoded"
	} else {
		s += r.body.w.FormDataContentType()
	}
	s += "\n"
	for k, v := range r.header {
		s += "\t\t" + k + ": " + strings.Join(v, ", ") + "\n"
	}
	if len(r.cookies) > 0 {
		s += "\tCookies:\n"
		for _, v := range r.cookies {
			s += "\t\t" + fmt.Sprint(v)
		}
	}
	return s
}
