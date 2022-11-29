package sheets

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"github.com/thirdweb-dev/go-sdk/v2/thirdweb"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/sheets/v4"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
)

type Client struct {
	SpreadsheetId string
	Srv           *sheets.Service
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

// should be able to set it between multiple docs
func NewSheetsClient() (*Client, error) {
	data, err := ioutil.ReadFile("web3-serviceaccount.json")
	if err != nil {
		return nil, err
	}
	conf, err := google.JWTConfigFromJSON(data, sheets.SpreadsheetsScope)
	if err != nil {
		return nil, err
	}
	client := conf.Client(context.Background())
	srv, err := sheets.New(client)
	if err != nil {
		return nil, err
	}
	spreadsheetID := os.Getenv("SPREADSHEET_ID")
	c := &Client{
		SpreadsheetId: spreadsheetID,
		Srv:           srv,
	}
	return c, nil
}

func (c *Client) ReadSheet(page string) ([][]interface{}, error) {
	//readRange := "A2:D2"
	var readRange string
	if page != "" {
		// "USDC!A2:D30"
		readRange = fmt.Sprintf("%s!A2:D2", page)
	}
	readRange = "A2:D2"
	//readRange
	resp, err := c.Srv.Spreadsheets.Values.Get(c.SpreadsheetId, readRange).Do()
	if err != nil {
		return nil, err
	}

	if len(resp.Values) == 0 {
		fmt.Println("No data found.")
		return nil, nil
	}
	return resp.Values, nil
}

func (c *Client) WriteSheetWeb3(events []thirdweb.ContractEvent) (*sheets.UpdateValuesResponse, error) {
	writeRange := "web3!A2"
	var vr sheets.ValueRange

	// create row data: event_name, block, address, tx_hash
	for i, e := range events {
		fmt.Printf("%v) Event: %v. Block: %v. Address: %v | tx: %v \n", i, e.EventName, e.Transaction.BlockNumber, e.Data["_account"], e.Transaction.TxHash)
		txLink := fmt.Sprintf("https://etherscan.io/tx/%s", e.Transaction.TxHash)
		// create the row
		rowData := []interface{}{e.EventName, e.Transaction.BlockNumber, e.Data["_account"], txLink}
		vr.Values = append(vr.Values, rowData)
	}

	// save the rows
	resp, err := c.Srv.Spreadsheets.Values.Update(c.SpreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (c *Client) WriteSheetRand() (*sheets.UpdateValuesResponse, error) {
	// Starting at A2..create 10 sequential entries
	writeRange := "Sheet1!A2"
	var vr sheets.ValueRange
	for i := 1; i <= 10; i++ {
		rowData := []interface{}{rand.Int(), rand.Int(), rand.Int()}
		vr.Values = append(vr.Values, rowData)
	}
	// save the rows
	resp, err := c.Srv.Spreadsheets.Values.Update(c.SpreadsheetId, writeRange, &vr).ValueInputOption("RAW").Do()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
