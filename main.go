package main

import (
	"log"
	"net/http"
	"os"

	"github.com/duypateo/ethereum-simple-page/database"
	"github.com/duypateo/ethereum-simple-page/etherscan"
	"github.com/duypateo/ethereum-simple-page/renderer"
	_ "github.com/heroku/x/hmetrics/onload"
)

// Index page
func indexHandler(w http.ResponseWriter, r *http.Request) {
	pageData := map[string]interface{}{
		"Title":           "Welcome to Ethereum simple gate",
		"IsValid":         "true",
		"IsChecksumIssue": "false",
		"ChecksumAddress": "",
		"Address":         "",
	}

	validKeys, ok := r.URL.Query()["isvalid"]
	if ok && len(validKeys[0]) != 0 {
		pageData["IsValid"] = validKeys[0]
	}

	addrKeys, ok := r.URL.Query()["address"]
	if ok && len(addrKeys[0]) != 0 {
		pageData["Address"] = addrKeys[0]
	}

	checksumKeys, ok := r.URL.Query()["checksum"]
	if ok && len(validKeys[0]) != 0 {
		pageData["IsChecksumIssue"] = checksumKeys[0]

		if checksumKeys[0] == "true" {
			address := pageData["Address"].(string)
			pageData["ChecksumAddress"] = etherscan.ChecksumAddress(address)
		}
	}

	pageData["recentAddrs"] = database.GetHistories()
	log.Println(pageData["recentAddrs"])

	// testing line
	renderer.RenderTemplate(w, "index", pageData)
}

// Index page
func checkHandler(w http.ResponseWriter, r *http.Request) {
	pageData := map[string]interface{}{
		"Title": "Your Ethereum balance snapshot",
	}

	addressFound, ok := r.URL.Query()["ether_address"]
	if !ok || len(addressFound[0]) == 0 {
		http.Redirect(w, r, "index", http.StatusFound)
		return
	}

	address := addressFound[0]
	if !etherscan.IsValidAddress(address) {
		http.Redirect(w, r, "index?isvalid=false&address="+address, http.StatusFound)
		return
	}

	checksumAddr := etherscan.ChecksumAddress(address)
	if checksumAddr != address {
		http.Redirect(w, r, "index?isvalid=false&address="+address+"&checksum=true", http.StatusFound)
		return
	}

	log.Println(address)
	pageData["Address"] = address

	// save or update to db
	database.InsertToHistory(address)

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

	err := database.InitConnection()
	if err != nil {
		log.Println(err)
		log.Fatal(err)
	}

	defer database.DB.Close()

	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/check", checkHandler)

	log.Println(etherscan.ChecksumAddress("0xd40800cc8b4f853eaea90b2b14b1ddda5511755b"))

	serverErr := http.ListenAndServe(port, nil)
	log.Fatal(serverErr)
}
