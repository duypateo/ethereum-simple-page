package main

import (
	"github.com/duypateo/ethereum-simple-page/etherscan"
	"github.com/duypateo/ethereum-simple-page/renderer"
	"log"
	"net/http"
	"os"
	_ "github.com/heroku/x/hmetrics/onload"
)

// Index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	pageData := map[string]interface{}{
		"Title":   "Welcome to Ethereum simple gate",
		"IsValid": true,
		"Address": "",
	}

	validKeys, ok := r.URL.Query()["invalid"]
	if ok && len(validKeys[0]) != 0 {
		pageData["IsValid"] = false
	}
	log.Printf("is valid address: %v", pageData["IsValid"])

	addrKeys, ok := r.URL.Query()["address"]
	if ok && len(addrKeys[0]) != 0 {
		pageData["Address"] = addrKeys[0]
	}

	// testing line
	renderer.RenderTemplate(w, "index", pageData)
}

// Index page
func checkHandler(w http.ResponseWriter, r *http.Request) {
	pageData := map[string]interface{}{
		"Title": "Your Ethereum balance snapshot",
	}

	address := r.FormValue("ether_address")
	if !etherscan.IsValidAddress(address) {
		http.Redirect(w, r, "index?invalid=true&address="+address, http.StatusFound)
		return
	}

	log.Println(address)
	pageData["Address"] = address

	balanceResp, err := etherscan.GetBalanceByAddress(address)
	log.Println(balanceResp)
	if err == nil {
		pageData["BalanceResp"] = balanceResp
	}

	blockResp, err := etherscan.GetBlocksByAddress(address)
	log.Println(blockResp)
	if err == nil {
		pageData["BlockResp"] = blockResp
	}

	renderer.RenderTemplate(w, "check", pageData)
}

// main
func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	port = ":" + port

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/check", checkHandler)

	serverErr := http.ListenAndServe(port, nil)
	log.Fatal(serverErr)
}
