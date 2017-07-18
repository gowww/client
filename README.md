# [![gowww](https://avatars.githubusercontent.com/u/18078923?s=20)](https://github.com/gowww) client [![GoDoc](https://godoc.org/github.com/gowww/client?status.svg)](https://godoc.org/github.com/gowww/client) [![Build](https://travis-ci.org/gowww/client.svg?branch=master)](https://travis-ci.org/gowww/client) [![Coverage](https://coveralls.io/repos/github/gowww/client/badge.svg?branch=master)](https://coveralls.io/github/gowww/client?branch=master) [![Go Report](https://goreportcard.com/badge/github.com/gowww/client)](https://goreportcard.com/report/github.com/gowww/client) ![Status Testing](https://img.shields.io/badge/status-testing-orange.svg)

Package [client](https://godoc.org/github.com/gowww/client) provides an HTTP client for clean requests.

## Installing

1. Get package:

	```Shell
	go get -u github.com/gowww/client
	```

2. Import it in your code:

	```Go
	import "github.com/gowww/client"
	```

## Usage

### Request

Use [Get](https://godoc.org/github.com/gowww/client#Get), [Post](https://godoc.org/github.com/gowww/client#Post), [Put](https://godoc.org/github.com/gowww/client#Put), [Patch](https://godoc.org/github.com/gowww/client#Patch), [Delete](https://godoc.org/github.com/gowww/client#Delete) or [Head](https://godoc.org/github.com/gowww/client#Head) with the destination URL to initiate a request.  
Options are chainable:

```Go
file, _ := os.Open("data/one.txt")
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
```

Finally, use [Do](https://godoc.org/github.com/gowww/client#Do) to send the requet and get the response or, eventually, the deferred first error of procedure:

```Go
res, err := req.Do()
if err != nil {
	panic(err)
}
defer res.Close()
```

Don't forget to close the response body when done.

### Response

A [Response](https://godoc.org/github.com/gowww/client#Response) wraps the standard [http.Response](https://golang.org/pkg/net/http/#Response) and provides some utility functions.

Use [Response.Cookie](https://godoc.org/github.com/gowww/client#Response.Cookie) to retrieve a single cookie:

```Go
c, err := res.Cookie(tokenCookieName)
if err != nil {
	// Generally, error is http.ErrNoCookie.
}
```

Use [Response.JSON](https://godoc.org/github.com/gowww/client#Response.JSON) to decode a JSON body into a variable:

```Go
jsres := new(struct{
	ID string `json:"id"`
})
res.JSON(jsres)
```
