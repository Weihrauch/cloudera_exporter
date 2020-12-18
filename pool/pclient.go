package pool

import (
    "sync"
	"net/http"
	"crypto/tls"
)
var once sync.Once

var instance *PClient

var C_api_request_type = "https"
var C_max_connexions = 0 
var C_req_per_seconds = 0

func GetPClient() *PClient {
    once.Do(func() {
		var httpClient *http.Client

		if(C_api_request_type == "https") {
			tr := &http.Transport{
				TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
			}
			httpClient = &http.Client{Transport: tr}
		}else{
			httpClient = http.DefaultClient
		}
        instance = NewPClient(httpClient, C_max_connexions, C_req_per_seconds)
        // instance = NewPClient()
    })
    return instance
}

func Init(api_request_type string, max_connexions int, req_per_seconds int) {
	C_api_request_type = api_request_type
	C_max_connexions = max_connexions
	C_req_per_seconds = req_per_seconds

	GetPClient()
}