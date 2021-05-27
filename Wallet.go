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

// func (myTime Time) String() string {

// 	return fmt.Sprintf("%d-%02d-%02d %02d:%02d\n",
// 		myTime.Year, myTime.Month, myTime.Day,
// 		myTime.Hour, myTime.Minute)
// }

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

func newWallet(c echo.Context) error {

	if c.Request().Method != "POST" {

		return c.JSON(http.StatusMethodNotAllowed, ErrorResponse{
			StatusCode: 405,
			Message:    "Method not allowed.",
		})
	}

	reqBody, _ := ioutil.ReadAll(c.Request().Body) // read the body of POST req in to slice bytes.
	//strReqBody := string(reqBody)
	//splitBody := strings.Split(strReqBody, ",")   // maybe needed for error handeling.
	newWallet := new(Wallet)
	current_time := time.Now()
	err := json.Unmarshal(reqBody, newWallet)

	// must check if the newPOSTWallet.name was in the public list.

	if findWalletByName(newWallet.Name) > -1 {

		return c.JSON(http.StatusNotAcceptable, ErrorResponse{
			StatusCode: 406,
			Message:    "Not acceptable.",
		})
	}

	(*newWallet).Last_updated = fmt.Sprintf("%d-%02d-%02d %02d:%02d",
		current_time.Year(), int(current_time.Month()), current_time.Day(),
		current_time.Hour(), current_time.Minute())

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
		Message:      "Food added successfully!",
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

	if index == -1 {

		return c.JSON(http.StatusNotFound, ErrorResponse{
			StatusCode: 404,
			Message:    "Not found.",
		})
	}

	reqBody, _ := ioutil.ReadAll(c.Request().Body) // read the body of POST req in to slice bytes.
	newWallet := new(Wallet)
	err := json.Unmarshal(reqBody, newWallet)

	if err != nil {

		return c.JSON(404, ErrorResponse{
			StatusCode: 404,
			Message:    "Invalid Types.",
		})
	}
	walletRecords[index].Name = newWallet.Name

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
	newCoin := new(Coin)
	err := json.Unmarshal(reqBody, newCoin)

	walletRecords[index].Coins = append(walletRecords[index].Coins, *newCoin)

	if err != nil {

		return c.JSON(404, ErrorResponse{
			StatusCode: 404,
			Message:    "Invalid Types.",
		})
	}

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
	updateCoin := new(Coin)
	err := json.Unmarshal(reqBody, updateCoin)

	if err != nil {

		return c.JSON(404, ErrorResponse{
			StatusCode: 404,
			Message:    "Invalid Types.",
		})
	}

	//------------- updating coin in wallet part
	walletRecords[index].Coins[indexCoinInWallet].Name = updateCoin.Name
	walletRecords[index].Coins[indexCoinInWallet].Symbol = updateCoin.Symbol
	walletRecords[index].Coins[indexCoinInWallet].Amount = updateCoin.Amount
	walletRecords[index].Coins[indexCoinInWallet].Rate = updateCoin.Rate

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

}

func searchForCoin(wallet Wallet, symbol string) int {

	for i := 0; i < len(wallet.Coins); i++ {
		if wallet.Coins[i].Symbol == symbol {
			return i
		}
	}
	return -1
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

	//e.POST("/:{wname}/:{symbol}", newCoinInWallet)
	e.POST("/:{wname}/coins", newCoinInWallet)
	e.GET("/:{wname}", getWalletInfo)
	e.PUT("/:{wname}/:{symbol}", updateCoinInWallet)
	e.DELETE("/:{wname}/:{symbol}", deleteCoinFromWallet)

	e.Logger.Fatal(e.Start(PORT))
}
