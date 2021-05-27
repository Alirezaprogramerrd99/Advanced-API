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
	Name   string
	Symbol string
	Amount float64
	rate   float64
}

type Time struct {
	Year     int
	Month    int
	Day      int
	Hour     int
	Minute   int
	TimeInfo string
}

func (myTime Time) String() string {

	return fmt.Sprintf("%d-%02d-%02d %02d:%02d\n",
		myTime.Year, myTime.Month, myTime.Day,
		myTime.Hour, myTime.Minute)
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

	// newTime := Time{
	// 	Year:   current_time.Year(),
	// 	Month:  int(current_time.Month()),
	// 	Day:    current_time.Day(),
	// 	Hour:   current_time.Hour(),
	// 	Minute: current_time.Minute(),
	// }

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

	fmt.Println(newWallet)

	newPOSTWallet := POSTWallet{

		Name:         newWallet.Name,
		Balance:      newWallet.Balance,
		Coins:        newWallet.Coins,
		Last_updated: newWallet.Last_updated,
		StatusCode:   200,
		Message:      "Food added successfully!",
	}

	walletRecords = append(walletRecords, *newWallet)

	// 	responseMsg := fmt.Sprintf(   // use c.String
	// 		`"{
	// 		name: "%s"
	// 		balance: %.1f,
	// 		coins: [],
	// 		Last_updated: "%s",
	// }"`, newPOSTWallet.Name, newPOSTWallet.Balance, newPOSTWallet.Last_updated)
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
	index := findWalletByName(walletName)

	if index == -1 {

		return c.JSON(http.StatusNotFound, ErrorResponse{
			StatusCode: 404,
			Message:    "Not found.",
		})
	}
	walletRecords[index].Name = walletName

	return c.JSON(http.StatusOK, POSTWallet{
		Name:         walletRecords[index].Name,
		Balance:      walletRecords[index].Balance,
		Coins:        walletRecords[index].Coins,
		Last_updated: walletRecords[index].Last_updated, // maybe need update.
		StatusCode:   200,
		Message:      "Wallet name changed successfully!",
	})
}

func findWalletByName(name string) int {

	for i := 0; i < len(walletRecords); i++ {
		if walletRecords[i].Name == name {
			return i
		}
	}
	return -1
}

func main() {

	walletRecords = []Wallet{}
	const PORT string = ":8080"
	// current_time := time.Now()
	// str := current_time.Year()
	// fmt.Printf("%T", str)
	// myTime := fmt.Sprintf("%d-%02d-%02d %02d:%02d\n",
	// 	current_time.Year(), current_time.Month(), current_time.Day(),
	// 	current_time.Hour(), current_time.Minute())
	// // individual elements of time can
	// // also be called to print accordingly

	// fmt.Println(myTime)

	e := echo.New()
	e.POST("/wallets", newWallet)
	e.GET("/wallets", getWallets)
	e.PUT("/wallets/:{wname}", updateWallet)

	// e.GET("/show", show)
	// e.GET("/users/:id", getUser)
	//e.POST("/save", save)

	e.Logger.Fatal(e.Start(PORT))
}
