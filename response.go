// Package client provides cient request itilities.
package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"mime"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"
)

// Response is a request response wrapping the original *http.Request.
type Response struct {
	*http.Response
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

// Dump is for debug purpose.
// It prints the request info, writes the body in a file and opens it in the browser.
// It panics on error.
func (r *Response) Dump() {
	log.Println("-", r)

	var ext string
	exts, _ := mime.ExtensionsByType(r.Header.Get("Content-Type"))
	if len(exts) > 0 {
		ext = exts[0]
	}
	name := fmt.Sprintf("response-dump-%d%s", time.Now().Unix(), ext)
	f, err := os.Create(name)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if strings.HasPrefix(r.Header.Get("Content-Type"), "text/html") {
		buf := bytes.NewBufferString("<pre style=\"background:#000;color:#0f0;font:13px/1.2 monospace;padding:20px\">")
		log.New(buf, "", log.LstdFlags).Print(" - ", r)
		buf.WriteString("</pre>")
		if _, err = io.Copy(f, buf); err != nil {
			panic(err)
		}
	}

	if _, err = io.Copy(f, r.Body); err != nil {
		panic(err)
	}

	openFile(name)
}

func openFile(url string) error {
	var cmd string
	var args []string
	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default:
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

func (r *Response) String() string {
	s := r.Status + " - " + r.Proto + " " + r.Request.Method + " " + r.Request.URL.String() + "\n"
	if len(r.Header) > 0 {
		s += "\tHeader:\n"
		for k, v := range r.Header {
			s += "\t\t" + k + ": " + strings.Join(v, ", ") + "\n"
		}
	}
	if len(r.Cookies()) > 0 {
		s += "\tCookies:\n"
		for _, v := range r.Cookies() {
			s += "\t\t" + fmt.Sprint(v)
		}
	}
	return s
}
