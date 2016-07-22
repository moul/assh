package main

import (
	"github.com/haklop/gnotifier"
	"net/http"
)

var urls = []string{
	"http://google.com/",
	"http://golang.org/",
	"http://twitter.com/",
}

type httpReponse struct {
	url      string
	response *http.Response
	err      error
}

func asyncHTTPGets(urls []string) <-chan *httpReponse {
	ch := make(chan *httpReponse, len(urls)) // buffered
	for _, url := range urls {
		go func(url string) {
			resp, err := http.Get(url)

			if err == nil {
				resp.Body.Close()
			}

			ch <- &httpReponse{url, resp, err}
		}(url)
	}
	return ch
}

func main() {
	results := asyncHTTPGets(urls)
	for _ = range urls {
		result := <-results

		var message string
		if result.err != nil {
			message = result.err.Error()
		} else {
			message = result.response.Status
		}

		notification := gnotifier.Notification(result.url, message)
		notification.GetConfig().Expiration = 2000
		notification.GetConfig().ApplicationName = "monitor-app"
		notification.Push()
	}
}
