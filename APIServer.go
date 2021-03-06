package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

type Coin struct {
	Name   string  `json:"name"`
	Symbol string  `json:"symbol"`
	Amount float64 `json:"amount"`
	Rate   float64 `json:"rate"`
}

type Wallet struct {
	Name         string  `json:"name"`
	Balance      float64 `json:"balance"`
	Coins        []Coin  `json:"coins"`
	Last_updated string  `json:"last_updated"`
}

var walletRecords []Wallet

type POSTWallet struct {
	Name         string  `json:"name"`
	Balance      float64 `json:"balance"`
	Coins        []Coin  `json:"coins"`
	Last_updated string  `json:"last_updated"`
	StatusCode   int     `json:"code"`
	Message      string  `json:"message"`
}

type POSTCoin struct {
	Name       string  `json:"name"`
	Symbol     string  `json:"symbol"`
	Amount     float64 `json:"amount"`
	Rate       float64 `json:"rate"`
	StatusCode int     `json:"code"`
	Message    string  `json:"message"`
}

type GETWalletsResponse struct {
	Size       int      `json:"size"`
	Wallets    []Wallet `json:"wallets"`
	StatusCode int      `json:"code"`
	Message    string   `json:"message"`
}

type ErrorResponse struct {
	StatusCode int    `json:"code"`
	Message    string `json:"message"`
}

// ------------------------ Wallet Handlers ------------------------------------------------------
func newWallet(c echo.Context) error {

	if c.Request().Method != "POST" {

		return c.JSON(http.StatusMethodNotAllowed, ErrorResponse{
			StatusCode: 405,
			Message:    "Method not allowed.",
		})
	}

	reqBody, _ := ioutil.ReadAll(c.Request().Body) // read the body of POST req in to slice bytes.
	strReqBody := string(reqBody)
	splitBody := strings.Split(strReqBody, ",") // maybe needed for error handeling.

	if len(splitBody) > 1 {

		return c.JSON(http.StatusNotAcceptable, ErrorResponse{
			StatusCode: 406,
			Message:    "Request not acceptable.",
		})
	}
	newWallet := new(Wallet)
	err := json.Unmarshal(reqBody, newWallet)

	if findWalletByName(newWallet.Name) > -1 {

		return c.JSON(http.StatusNotAcceptable, ErrorResponse{
			StatusCode: 406,
			Message:    "Not acceptable.",
		})
	}

	(*newWallet).Last_updated = updateDate()
	(*newWallet).Balance = 0.0
	(*newWallet).Coins = make([]Coin, 0)

	if err != nil {

		return c.JSON(404, ErrorResponse{
			StatusCode: 404,
			Message:    "Invalid Types.",
		})
	}

	newPOSTWallet := POSTWallet{

		Name:         newWallet.Name,
		Balance:      newWallet.Balance,
		Coins:        newWallet.Coins,
		Last_updated: newWallet.Last_updated,
		StatusCode:   200,
		Message:      "Wallet added successfully!",
	}
	walletRecords = append(walletRecords, *newWallet)

	return c.JSON(http.StatusOK, newPOSTWallet)
}

func getWallets(c echo.Context) error {

	if c.Request().Method != "GET" {

		return echo.NewHTTPError(404, ErrorResponse{
			StatusCode: http.StatusMethodNotAllowed,
			Message:    "Method not allowed.",
		})

	}
	newGETWalletsResponse := GETWalletsResponse{
		Size:       len(walletRecords),
		Wallets:    walletRecords,
		StatusCode: 200,
		Message:    "All wallets received successfully!",
	}
	return c.JSON(http.StatusOK, newGETWalletsResponse)
}

func updateWallet(c echo.Context) error {

	wname := c.Param("{wname}")

	initial := wname
	walletName := strings.TrimLeft(strings.TrimRight(initial, "}"), "{")
	// fmt.Println(walletName)
	index := findWalletByName(walletName)

	if index == -1 { // for check is that wallet is in the walletRecords or not.

		return c.JSON(http.StatusNotFound, ErrorResponse{
			StatusCode: 404,
			Message:    "Not found.",
		})
	}

	reqBody, _ := ioutil.ReadAll(c.Request().Body)

	strReqBody := string(reqBody)
	splitBody := strings.Split(strReqBody, ",") // maybe needed for error handeling.

	if len(splitBody) > 1 {

		return c.JSON(http.StatusMethodNotAllowed, ErrorResponse{
			StatusCode: 406,
			Message:    "Request not acceptable.",
		})
	}

	newWallet := new(Wallet)
	err := json.Unmarshal(reqBody, newWallet)

	if err != nil {

		return c.JSON(404, ErrorResponse{
			StatusCode: 404,
			Message:    "Invalid Types.",
		})
	}

	if findWalletByName(newWallet.Name) > -1 { // for duplicate in walletRecords from request body.

		return c.JSON(http.StatusMethodNotAllowed, ErrorResponse{
			StatusCode: 406,
			Message:    "New wallet name is duplicate!",
		})
	}

	//------------------ updating the wallet------------------------
	walletRecords[index].Name = newWallet.Name
	walletRecords[index].Last_updated = updateDate()

	return c.JSON(http.StatusOK, POSTWallet{
		Name:         walletRecords[index].Name,
		Balance:      walletRecords[index].Balance,
		Coins:        walletRecords[index].Coins,
		Last_updated: walletRecords[index].Last_updated, // maybe need update.
		StatusCode:   200,
		Message:      "Wallet name changed successfully!",
	})
}

func deleteWallet(c echo.Context) error {

	wname := c.Param("{wname}")

	initial := wname
	walletName := trimBrackets(initial)
	index := findWalletByName(walletName)

	if index == -1 {

		return c.JSON(http.StatusNotFound, ErrorResponse{
			StatusCode: 404,
			Message:    "Not found.",
		})
	}

	PostWalletResponse := POSTWallet{
		Name:         walletRecords[index].Name,
		Balance:      walletRecords[index].Balance,
		Coins:        walletRecords[index].Coins,
		Last_updated: walletRecords[index].Last_updated, // maybe need update.
		StatusCode:   200,
		Message:      "Wallet deleted (logged out) successfully!",
	}

	walletRecords = RemoveWalletRecord(index)

	return c.JSON(http.StatusOK, PostWalletResponse)
}

func findWalletByName(name string) int {

	for i := 0; i < len(walletRecords); i++ {
		if walletRecords[i].Name == name {
			return i
		}
	}
	return -1
}

func trimBrackets(str string) string {
	return strings.TrimLeft(strings.TrimRight(str, "}"), "{")
}

func RemoveWalletRecord(index int) []Wallet {
	return append(walletRecords[:index], walletRecords[index+1:]...)
}

// ------------------------ Coin Handlers ------------------------------------------------------
func newCoinInWallet(c echo.Context) error { // here also we must update the balance of that wallet.
	wname := c.Param("{wname}")
	// symbolParam := c.Param("{symbol}")

	walletName := trimBrackets(wname)
	// symbol := trimBrackets(symbolParam)
	index := findWalletByName(walletName)

	if index == -1 {

		return c.JSON(404, ErrorResponse{
			StatusCode: 404,
			Message:    "Not found.",
		})
	}

	reqBody, _ := ioutil.ReadAll(c.Request().Body)
	strReqBody := string(reqBody)
	splitBody := strings.Split(strReqBody, ",") // maybe needed for error handeling.

	if len(splitBody) > 4 {

		return c.JSON(http.StatusMethodNotAllowed, ErrorResponse{
			StatusCode: 406,
			Message:    "Request not acceptable.",
		})
	}

	newCoin := new(Coin)
	err := json.Unmarshal(reqBody, newCoin)

	if err != nil {

		return c.JSON(404, ErrorResponse{
			StatusCode: 404,
			Message:    "Invalid Types.",
		})
	}

	if isSameInWallet(-1, walletRecords[index], newCoin.Name, newCoin.Symbol) {

		return c.JSON(406, ErrorResponse{
			StatusCode: 406,
			Message:    "duplicate coin in wallet!",
		})
	}
	walletRecords[index].Coins = append(walletRecords[index].Coins, *newCoin)
	updateBalance(&(walletRecords[index]))

	walletRecords[index].Last_updated = updateDate()

	return c.JSON(http.StatusOK, POSTCoin{
		Name:       newCoin.Name,
		Symbol:     newCoin.Symbol,
		Amount:     newCoin.Amount,
		Rate:       newCoin.Rate,
		StatusCode: 200,
		Message:    "Coin added successfully!",
	})
}

func getWalletInfo(c echo.Context) error {

	wname := c.Param("{wname}")
	walletName := trimBrackets(wname)
	index := findWalletByName(walletName)

	if index == -1 {

		return c.JSON(404, ErrorResponse{
			StatusCode: 404,
			Message:    "Not found.",
		})
	}

	return c.JSON(http.StatusOK, POSTWallet{

		Name:         walletRecords[index].Name,
		Balance:      walletRecords[index].Balance,
		Coins:        walletRecords[index].Coins,
		Last_updated: walletRecords[index].Last_updated,
		StatusCode:   200,
		Message:      "All coins received successfully!",
	})

}

func updateCoinInWallet(c echo.Context) error { // update balance

	wname := c.Param("{wname}")
	symbolParam := c.Param("{symbol}")

	walletName := trimBrackets(wname)
	symbol := trimBrackets(symbolParam)
	index := findWalletByName(walletName)

	if index == -1 {

		return c.JSON(404, ErrorResponse{
			StatusCode: 404,
			Message:    "Not found.",
		})
	}
	// search for a coin in wallet.
	indexCoinInWallet := searchForCoin(walletRecords[index], symbol)

	if indexCoinInWallet == -1 {

		return c.JSON(404, ErrorResponse{
			StatusCode: 404,
			Message:    "Coin is not in this wallet!",
		})
	}

	reqBody, _ := ioutil.ReadAll(c.Request().Body)

	strReqBody := string(reqBody)
	splitBody := strings.Split(strReqBody, ",") // maybe needed for error handeling.

	if len(splitBody) > 4 {

		return c.JSON(http.StatusMethodNotAllowed, ErrorResponse{
			StatusCode: 406,
			Message:    "Request not acceptable.",
		})
	}

	updateCoin := new(Coin)
	err := json.Unmarshal(reqBody, updateCoin)

	if err != nil {

		return c.JSON(404, ErrorResponse{
			StatusCode: 404,
			Message:    "Invalid Types.",
		})
	}

	if isSameInWallet(indexCoinInWallet, walletRecords[index], updateCoin.Name, updateCoin.Symbol) {

		return c.JSON(406, ErrorResponse{
			StatusCode: 406,
			Message:    "duplicate coin in wallet!",
		})
	}
	//------------- updating coin in wallet part
	walletRecords[index].Coins[indexCoinInWallet].Name = updateCoin.Name
	walletRecords[index].Coins[indexCoinInWallet].Symbol = updateCoin.Symbol
	walletRecords[index].Coins[indexCoinInWallet].Amount = updateCoin.Amount
	walletRecords[index].Coins[indexCoinInWallet].Rate = updateCoin.Rate

	updateBalance(&(walletRecords[index]))

	walletRecords[index].Last_updated = updateDate()

	return c.JSON(http.StatusOK, POSTCoin{

		Name:       updateCoin.Name,
		Symbol:     updateCoin.Symbol,
		Amount:     updateCoin.Amount,
		Rate:       updateCoin.Rate,
		StatusCode: 200,
		Message:    "Coin updated successfully!",
	})
}

func deleteCoinFromWallet(c echo.Context) error {

	wname := c.Param("{wname}")
	symbolParam := c.Param("{symbol}")

	walletName := trimBrackets(wname)
	symbol := trimBrackets(symbolParam)
	index := findWalletByName(walletName)

	if index == -1 {

		return c.JSON(404, ErrorResponse{
			StatusCode: 404,
			Message:    "Not found.",
		})
	}
	// search for a coin in wallet.
	indexCoinInWallet := searchForCoin(walletRecords[index], symbol)

	if indexCoinInWallet == -1 {

		return c.JSON(404, ErrorResponse{
			StatusCode: 404,
			Message:    "Coin is not in this wallet!",
		})
	}

	PostCoinResponse := POSTCoin{
		Name:       walletRecords[index].Coins[indexCoinInWallet].Name,
		Symbol:     walletRecords[index].Coins[indexCoinInWallet].Symbol,
		Amount:     walletRecords[index].Coins[indexCoinInWallet].Amount,
		Rate:       walletRecords[index].Coins[indexCoinInWallet].Rate,
		StatusCode: 200,
		Message:    "Coin deleted successfully!",
	}

	walletRecords[index].Coins = removeCoinFromWallet(walletRecords[index], indexCoinInWallet)
	updateBalance(&(walletRecords[index]))
	walletRecords[index].Last_updated = updateDate()

	return c.JSON(http.StatusOK, PostCoinResponse)
}

func searchForCoin(wallet Wallet, symbol string) int {

	for i := 0; i < len(wallet.Coins); i++ {
		if wallet.Coins[i].Symbol == symbol {
			return i
		}
	}
	return -1
}

func removeCoinFromWallet(wallet Wallet, index int) []Coin {
	return append(wallet.Coins[:index], wallet.Coins[index+1:]...)
}

func updateBalance(wallet *Wallet) {

	var product float64
	var sumAll float64 = 0.0

	for i := 0; i < len(wallet.Coins); i++ {

		product = wallet.Coins[i].Rate * wallet.Coins[i].Amount
		sumAll += product
	}
	wallet.Balance = sumAll
}

func isSameInWallet(invalidIndex int, wallet Wallet, name string, symbol string) bool {

	for i := 0; i < len(wallet.Coins); i++ {

		if i == invalidIndex {
			continue
		}

		if wallet.Coins[i].Name == name || wallet.Coins[i].Symbol == symbol {
			return true
		}
	}
	return false
}

func updateDate() string {
	current_time := time.Now()

	return fmt.Sprintf("%d-%02d-%02d %02d:%02d",
		current_time.Year(), int(current_time.Month()), current_time.Day(),
		current_time.Hour(), current_time.Minute())
}

func main() {

	walletRecords = []Wallet{}
	const PORT string = ":8080"

	e := echo.New()
	// --------- wallet methods -----------------
	e.POST("/wallets", newWallet)
	e.GET("/wallets", getWallets)
	e.PUT("/wallets/:{wname}", updateWallet)
	e.DELETE("/wallets/:{wname}", deleteWallet)

	// ------- coin methods ---------------------
	e.POST("/:{wname}/coins", newCoinInWallet)
	e.GET("/:{wname}", getWalletInfo)
	e.PUT("/:{wname}/:{symbol}", updateCoinInWallet)
	e.DELETE("/:{wname}/:{symbol}", deleteCoinFromWallet)

	e.Logger.Fatal(e.Start(PORT))
}
