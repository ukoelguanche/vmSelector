package httpClient

import (
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"apodeiktikos.com/fbtest/util"
)

type RequestOptions struct {
	Body string
}

func DoRequest(method string, path string, options *RequestOptions) *http.Response {
	tokenID := util.ContextStorage.PveTokenId
	tokenSecret := util.ContextStorage.PveSecret
	host := util.ContextStorage.PveHosst

	apiToken := fmt.Sprintf("PVEAPIToken=%s=%s", tokenID, tokenSecret)

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: 10 * time.Second,
	}

	var bodyReader io.Reader
	if options != nil && options.Body != "" {
		bodyReader = strings.NewReader(options.Body)
	}

	URL := fmt.Sprintf("%s/api2/json%s", host, path)
	request, _ := http.NewRequest(method, URL, bodyReader)
	request.Header.Add("Authorization", apiToken)

	response, err := client.Do(request)
	if err != nil {
		log.Fatal("Error requesting [%s] [%s] %v", method, path, err)
	}

	return response
}
