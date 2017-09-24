package api

import (
	"dynamic-dns/config"
	"dynamic-dns/helpers"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

var dnsConf = &config.Config{}

type WebserviceResponse struct {
	Success bool
	Message string
}

func Init(configPath string) {
	dnsConf.LoadConfig(configPath)
}

func RequestHandlerForAdd(w http.ResponseWriter, r *http.Request) {
	RequestHandlerCommon(w, r, "add")
}

func RequestHandlerForUpdate(w http.ResponseWriter, r *http.Request) {
	RequestHandlerCommon(w, r, "update")
}

func RequestHandlerForDelete(w http.ResponseWriter, r *http.Request) {
	RequestHandlerCommon(w, r, "delete")
}

func RequestHandlerCommon(w http.ResponseWriter, r *http.Request, operation string) {

	response := WebserviceResponse{}

	var SecretKey string
	var domain string
	var address string

	vals := r.URL.Query()
	SecretKey = vals["secret"][0]
	domain = vals["domain"][0]
	address = vals["ip"][0]

	if SecretKey != dnsConf.SecretKey {
		log.Println(fmt.Sprintf("Invalid shared secret: %s Original[%s]", SecretKey, dnsConf.SecretKey))
		response.Success = false
		response.Message = "Invalid Credentials"
		json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	var addrType string
	if helpers.IsValidIP4(address) {
		addrType = "A"
	} else if helpers.IsValidIP6(address) {
		addrType = "AAAA"
	} else {
		response.Success = false
		response.Message = fmt.Sprintf("%s is neither a valid IPv4 nor IPv6 address", address)
	}

	if addrType != "" {
		if domain == "" {
			response.Success = false
			response.Message = fmt.Sprintf("Domain not set %s", address)
			log.Println(fmt.Sprintf("Domain not set"))
			return
		}

		result := processRecord(domain, address, addrType, operation)

		if result == "" {
			response.Success = true
			response.Message = fmt.Sprintf("Operation [%s] on %s record for %s to IP address %s executed successfully.", addrType, domain, address)
		} else {
			response.Success = false
			response.Message = result
		}
	}

	json.NewEncoder(w).Encode(response)
}

func processRecord(domain string, address string, addrType string, operation string) string {

	var response string

	switch operation {
	case "add":
		AddRecord(domain, address, addrType)
	case "update":
		UpdateRecord(domain, address, addrType)
	case "delete":
		DeleteRecord(domain, address, addrType)
	default:
		response = "Invalid Operation"
	}
	return response
}
