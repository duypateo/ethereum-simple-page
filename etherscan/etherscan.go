package etherscan

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"regexp"
)

// apiKeyToken for etherscan api
const apiKeyToken = "QWM1C5MJKDRPAWC4RTVMIFTRJ2FP8QFHAJ"

// apiOriginPath for etherscan api
const apiOriginPath = "https://api.etherscan.io/api?"

// Response represent for response of etherscan api
type Response struct {
	Status  string      `json:"status"`
	Message string      `json:"message"`
	Result  interface{} `json:"result"`
}

// fetchResponse - fetch response data
func fetchResponse(reqURL string) (Response, error) {
	resp := Response{}

	originResp, err := http.Get(reqURL)
	if err != nil {
		return resp, err
	}

	defer originResp.Body.Close()

	originRespBody, err := ioutil.ReadAll(originResp.Body)
	if err != nil {
		return resp, err
	}

	json.Unmarshal(originRespBody, &resp)
	return resp, nil
}

// GetBalanceByAddress - get balance by address
func GetBalanceByAddress(addr string) (Response, error) {
	reqURL := apiOriginPath
	reqURL += "module=account"
	reqURL += "&action=balance"
	reqURL += "&address=" + addr
	reqURL += "&tag=latest"
	reqURL += "&apikey=" + apiKeyToken

	return fetchResponse(reqURL)
}

// GetBlocksByAddress - get blocks by address
func GetBlocksByAddress(addr string) (Response, error) {
	reqURL := apiOriginPath
	reqURL += "module=account"
	reqURL += "&action=getminedblocks"
	reqURL += "&address=" + addr
	reqURL += "&blocktype=blocks"
	reqURL += "&page=1"
	reqURL += "&offset=10"
	reqURL += "&apikey=" + apiKeyToken

	return fetchResponse(reqURL)
}

// IsValidAddress - check if address is valid
func IsValidAddress(v string) bool {
	re := regexp.MustCompile("^0x[0-9a-fA-F]{40}$")
	return re.MatchString(v)
}
