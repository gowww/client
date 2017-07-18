package client

import (
	"fmt"
	"net/http"
	"testing"
)

func TestRequestString(t *testing.T) {
	req := Post("http://example.com").
		DisableRedirect().
		ForceMultipart().
		Header("Accept-Language", "en").
		Header("Accept-Language", "fr").
		UserAgent("Googlebot/2.1 (+http://www.google.com/bot.html)").
		Cookie(&http.Cookie{Name: "session", Value: "123"}).
		Value("id", "123").
		Value("name", "Doe").
		OpenFile("picture", "one.png").
		OpenFile("picture", "two.png")

	fmt.Println(req)
}
