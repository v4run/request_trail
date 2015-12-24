package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/fatih/color"
)

func main() {
	var maxRedirects int
	var url, method string
	var help bool

	flag.IntVar(&maxRedirects, "redirects", 10, "Maximum number of allowed redirects. -1 for no limit")
	flag.IntVar(&maxRedirects, "r", 10, "Maximum number of allowed redirects. -1 for no limit")
	flag.StringVar(&url, "url", "https://www.google.com", "Url to check")
	flag.StringVar(&url, "u", "https://www.google.com", "Url to check")
	flag.StringVar(&method, "method", "GET", "Request method to use")
	flag.StringVar(&method, "m", "GET", "Request method to use")
	flag.BoolVar(&help, "help", false, "Show help")
	flag.BoolVar(&help, "h", false, "Show help")
	flag.Parse()
	if help {
		flag.Usage()
		os.Exit(0)
	}

	urlColor := color.New(color.FgWhite).SprintFunc()
	successStatusColor := color.New(color.FgHiGreen).SprintFunc()
	redirectStatusColor := color.New(color.FgYellow).SprintFunc()
	errorSatusColor := color.New(color.FgHiRed).SprintFunc()
	infoStatusColor := color.New(color.FgHiBlue).SprintFunc()

	var statusCodeColor = func(code int) string {
		if code > 399 {
			return errorSatusColor(code)
		}
		if code > 299 {
			return redirectStatusColor(code)
		}
		if code > 199 {
			return successStatusColor(code)
		}
		return infoStatusColor(code)
	}

	t := &transportWrapper{}
	client := &http.Client{
		Transport:     t,
		CheckRedirect: checkRedirect(maxRedirects),
	}
	_, err := client.Get(url)
	if err != nil {
		log.Println(err)
	}
	for _, resp := range t.redirectTrail {
		fmt.Printf("%s %s\n", statusCodeColor(resp.code), urlColor(resp.url))
	}
}
