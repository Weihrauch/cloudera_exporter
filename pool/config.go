package pool


import (
    "sync"
)
var once sync.Once

type ce_config struct {
	Api_request_type string
	Max_connexions int
	Req_per_seconds int
}
var instance *ce_config
var Api_request_type = "https"
var Max_connexions = 0 
var Req_per_seconds = 0

func GetConfig() *ce_config {
    once.Do(func() {
		instance = &ce_config{}
		instance.Api_request_type = Api_request_type
		instance.Max_connexions = Max_connexions
		instance.Req_per_seconds = Req_per_seconds
	})
    return instance
}

func InitConfig(api_request_type string, max_connexions int, req_per_seconds int) {
	Api_request_type = api_request_type
	Max_connexions = max_connexions
	Req_per_seconds = req_per_seconds

	GetConfig()
}