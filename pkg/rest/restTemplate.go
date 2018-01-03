package rest

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

// RequestConfig rest request config
type RequestConfig struct {
	Headers map[string]string
	Body    interface{}
}

// Post Post request to Restful endpoint
func Post(url string, config RequestConfig, target interface{}) (resp *http.Response, err error) {
	resp, err = do(http.MethodPost, url, config, target)
	return
}

// Get Get request to Restful endpoint
func Get(url string, config RequestConfig, target interface{}) (resp *http.Response, err error) {
	resp, err = do(http.MethodGet, url, config, target)
	return
}

func do(method string, url string, config RequestConfig, target interface{}) (resp *http.Response, err error) {
	client := &http.Client{}

	body, err := json.Marshal(config.Body)

	if err != nil {
		log.Println(err)
		return
	}

	req, err := http.NewRequest(method, url, bytes.NewBuffer(body))

	if err != nil {
		log.Println(err)
		return
	}

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/json")

	if config.Headers != nil {
		for key, value := range config.Headers {
			req.Header.Add(key, value)
		}
	}

	resp, err = client.Do(req)

	responseBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return
	}

	defer resp.Body.Close()

	if target != nil {
		err = json.Unmarshal(responseBody, target)
		if err != nil {
			log.Println(err)
		}
	}

	return

}
