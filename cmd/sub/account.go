package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/coingate/tron-withdrawal-service/cmd/sub/transaction"
	"github.com/coingate/tron-withdrawal-service/protos/core"
	"github.com/golang/protobuf/proto"
	"math"
	"os"
)

func Withdraw(toAddress string, amount float64) error {

	if toAddress == "" {
		return fmt.Errorf("no receive address specified")
	}

	// get amount
	valueInt := int64(amount * math.Pow10(6))

	var err error
	fromAddress := os.Getenv("TRON_SENDER_ADDRESS")

	contract := &core.TransferContract{}

	if contract.OwnerAddress, err = DecodeCheck(fromAddress); err != nil {
		return err
	}

	if contract.ToAddress, err = DecodeCheck(toAddress); err != nil {
		return err
	}

	contract.Amount = valueInt

	ctx, cancel := conn.getContext()
	defer cancel()

	tx, err := conn.Client.CreateTransaction2(ctx, contract)

	if err != nil {
		return err
	}
	if proto.Size(tx) == 0 {
		return fmt.Errorf("bad transaction")
	}

	if tx.GetResult().GetCode() != 0 {
		return fmt.Errorf("%s", tx.GetResult().GetMessage())
	}



	var ctrlr *transaction.Controller

	ks, acct, err := store.UnlockedKeystore(signerAddress.String(), passphrase)
	if err != nil {
		return err
	}

	ctrlr = transaction.NewController(&conn, ks, acct, tx.Transaction, opts)
	if err = ctrlr.ExecuteTransaction(); err != nil {
		return err
	}

	result := make(map[string]interface{})
	result["from"] = fromAddress
	result["to"] = toAddress
	result["amount"] = amount
	result["txID"] = common.BytesToHexString(tx.GetTxid())
	result["blockNumber"] = ctrlr.Receipt.BlockNumber
	result["message"] = string(ctrlr.Result.Message)
	result["receipt"] = map[string]interface{}{
		"fee":      ctrlr.Receipt.Fee,
		"netFee":   ctrlr.Receipt.Receipt.NetFee,
		"netUsage": ctrlr.Receipt.Receipt.NetUsage,
	}

	asJSON, _ := json.Marshal(result)
	fmt.Println(JSONPrettyFormat(string(asJSON)))
	return nil
}

