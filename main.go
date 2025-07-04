package main

import (
	"context"
	"encoding/csv"
	"fmt"
	binance_connector "github.com/binance/binance-connector-go"
	"github.com/ethereum/go-ethereum/common"
	"github.com/joho/godotenv"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

var (
	apiKey    string
	secretKey string
	coin      string
	network   string
	baseURL   string
)

func init() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
		return
	}

	apiKey = os.Getenv("API_KEY")
	secretKey = os.Getenv("SECRET_KEY")
	coin = os.Getenv("COIN")
	network = os.Getenv("NETWORK")
	baseURL = os.Getenv("BASEURL")

	// 启动自检
	if apiKey == "" || secretKey == "" || coin == "" || network == "" || baseURL == "" {
		panic("请在.env文件中填写API_KEY、SECRET_KEY、COIN、NETWORK、BASEURL")
	}
}

func withdraw(client *binance_connector.Client, asset, address string, amount float64, network string) error {
	// Define the withdraw request parameters
	withdrawRequest := client.NewWithdrawService().
		Coin(asset).
		Address(address).
		Amount(amount).
		Network(network)

	// Send the withdraw request
	result, err := withdrawRequest.Do(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("Withdraw result: %v\n", result)
	return nil
}

func main() {
	// Initialise the client
	client := binance_connector.NewClient(apiKey, secretKey, baseURL)

	httpCleint := &http.Client{
		Transport: &http.Transport{
			// 设置代理，从环境变量中获取
			Proxy: http.ProxyFromEnvironment,
		},
	}
	client.HTTPClient = httpCleint

	// Read address and amount data from a CSV file and load them into a slice for looping
	file, err := os.Open("addresses.csv")
	if err != nil {
		fmt.Println("Error opening CSV file", err)
		return
	}
	defer file.Close()
	reader := csv.NewReader(file)
	var addressAmountPairs []struct {
		Address string
		Amount  float64
	}

	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			fmt.Println("Error reading CSV record", err)
			return
		}

		amount, err := strconv.ParseFloat(record[1], 64)
		if err != nil {
			fmt.Println("Error parsing amount from CSV", err)
			continue
		}

		addressAmountPair := struct {
			Address string
			Amount  float64
		}{
			Address: common.HexToAddress(record[0]).Hex(),
			Amount:  amount,
		}

		addressAmountPairs = append(addressAmountPairs, addressAmountPair)
	}
	//ctx := context.Background()
	//accountInfo, err := client.NewGetAccountService().Do(ctx)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}
	//
	//for _, bal := range accountInfo.Balances {
	//	f, err := strconv.ParseFloat(bal.Free, 64)
	//	if err != nil {
	//		fmt.Println("转换失败:", err)
	//		return
	//	}
	//	l, err := strconv.ParseFloat(bal.Locked, 64)
	//	if err != nil {
	//		fmt.Println("转换失败:", err)
	//		return
	//	}
	//
	//	if f == 0 && l == 0 {
	//		continue
	//	}
	//	fmt.Println(binance_connector.PrettyPrint(bal))
	//}
	//fmt.Println(binance_connector.PrettyPrint(accountInfo.Balances))

	//allCoinsInfo, err := client.NewGetAllCoinsInfoService().Do(ctx)
	//if err != nil {
	//	fmt.Println(err)
	//	return
	//}

	//for _, bal := range allCoinsInfo {
	//	f, err := strconv.ParseFloat(bal.Free, 64)
	//	if err != nil {
	//		fmt.Println("转换失败:", err)
	//		return
	//	}
	//	l, err := strconv.ParseFloat(bal.Locked, 64)
	//	if err != nil {
	//		fmt.Println("转换失败:", err)
	//		return
	//	}
	//
	//	if f == 0 && l == 0 {
	//		continue
	//	}
	//
	//	if bal.Name == "BNB" || bal.Coin == "BNB" {
	//		fmt.Println(binance_connector.PrettyPrint(bal))
	//	}
	//}

	//fmt.Println(binance_connector.PrettyPrint(allCoinsInfo))

	// Now you can loop over the addressAmountPairs slice for further processing
	for _, pair := range addressAmountPairs {
		fmt.Printf("Address: %s, Amount: %.2f\n", pair.Address, pair.Amount)
		err := withdraw(client, coin, pair.Address, pair.Amount, network)
		if err != nil {
			fmt.Println("Error withdrawing ", coin, err)
			return
		}
		time.Sleep(2 * time.Second)
	}

}
