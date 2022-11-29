package sheets

import (
	"context"
	"github.com/itsahedge/blacklist/cmd/blacklists"
	"github.com/thirdweb-dev/go-sdk/v2/thirdweb"
	"testing"
)

func TestNewSheetsClient(t *testing.T) {
	c, err := NewSheetsClient()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(c)
}

func TestClient_ReadSheet(t *testing.T) {
	c, err := NewSheetsClient()
	if err != nil {
		t.Fatal(err)
	}

	values, err := c.ReadSheet("")
	for _, row := range values {
		t.Logf("%s, %s, %s \n", row[0], row[1], row[2])
	}
}

func TestClient_WriteSheet(t *testing.T) {
	c, err := NewSheetsClient()
	if err != nil {
		t.Fatal(err)
	}

	resp, err := c.WriteSheetRand()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(resp)
}

func Test_ClientWriteWeb3(t *testing.T) {
	// INITIALIZE SHEETS AND CLIENT
	c, err := NewSheetsClient()
	if err != nil {
		t.Fatal(err)
	}

	sdk, err := thirdweb.NewThirdwebSDK("mainnet", nil)
	if err != nil {
		t.Fatal(err)
	}
	contract, err := sdk.GetContractFromAbi(blacklists.USDC_TOKEN, blacklists.USDC_TOKEN_ABI)

	// HANDLE EVENTS
	toBlock := uint64(16072377)
	queryOptions := thirdweb.EventQueryOptions{
		FromBlock: 7581399,
		ToBlock:   &toBlock,
	}
	eventsBlacklist, _ := contract.Events.GetEvents(context.Background(), "Blacklisted", queryOptions)
	if err != nil {
		t.Fatal(err)
	}

	eventsUnblacklist, _ := contract.Events.GetEvents(context.Background(), "UnBlacklisted", queryOptions)
	if err != nil {
		t.Fatal(err)
	}

	for i, e := range eventsBlacklist {
		t.Logf("%v) Event: %v. Block: %v. Address: %v | tx: %v", i, e.EventName, e.Transaction.BlockNumber, e.Data["_account"], e.Transaction.TxHash)
	}
	for i, e := range eventsUnblacklist {
		t.Logf("%v) Event: %v. Block: %v. Address: %v | tx: %v", i, e.EventName, e.Transaction.BlockNumber, e.Data["_account"], e.Transaction.TxHash)
	}

	// WRITE
	_, err = c.WriteSheetWeb3(eventsBlacklist)
	if err != nil {
		t.Fatal(err)
	}
	//resp.UpdatedData./

	// this overrides the first couple rows since it starts at A2..
	//_, err = c.WriteSheetWeb3(eventsUnblacklist)
	//if err != nil {
	//	t.Fatal(err)
	//}
	t.Logf("wrote to sheets")
}
