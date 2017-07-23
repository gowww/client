package client_test

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gowww/client"
)

func Example() {
	file, err := os.Open("data/one.txt")
	if err != nil {
		panic(err)
	}
	defer file.Close()

	req := client.Post("http://example.com").
		DisableRedirect().
		ForceMultipart().
		Header("Accept-Language", "en").
		UserAgent("Googlebot/2.1 (+http://www.google.com/bot.html)").
		Cookie(&http.Cookie{Name: "session", Value: "123"}).
		Value("id", "123").
		Value("name", "Doe").
		File("file", "one.txt", file).
		OpenFile("picture", "one.png").
		OpenFile("picture", "two.png")

	res, err := req.Do()
	if err != nil {
		panic(err)
	}
	defer res.Close()

	fmt.Println(res)
}
