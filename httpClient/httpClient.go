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

var client = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
		MaxIdleConns:      100,
		IdleConnTimeout:   90 * time.Second,
		DisableKeepAlives: false, // Queremos reutilizar conexiones
	},
	Timeout: 10 * time.Second,
}

func DoRequestOO(method string, path string, options *RequestOptions) *http.Response {
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

func DoRequest(method string, path string, options *RequestOptions) *http.Response {
	tokenID := util.ContextStorage.PveTokenId
	tokenSecret := util.ContextStorage.PveSecret
	host := util.ContextStorage.PveHosst

	apiToken := fmt.Sprintf("PVEAPIToken=%s=%s", tokenID, tokenSecret)

	URL := fmt.Sprintf("%s/api2/json%s", host, path)

	var bodyReader io.Reader
	if options != nil && options.Body != "" {
		bodyReader = strings.NewReader(options.Body)
	}

	request, err := http.NewRequest(method, URL, bodyReader)
	if err != nil {
		log.Fatalln("error creating request[%s] [%s]: %w", method, path, err)
	}

	request.Header.Add("Authorization", apiToken)

	// Usamos el cliente global
	response, err := client.Do(request)
	if err != nil {
		// Evita log.Fatal aquí, es mejor devolver el error para manejarlo
		log.Fatalln("error requesting [%s] [%s]: %w", method, path, err)
	}

	return response
}
