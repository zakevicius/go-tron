package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/coingate/tron-withdrawal-service/protos"
	"google.golang.org/grpc/metadata"
)

type Balance struct {
	Data struct {
		Balance int64 `json:"data[0].balance"`
	}
}

func CheckBalance(address string) bool {
	conn := NewGrpcClient("")

	account := new(protos.Account)
	account.Address = []byte(address)

	var err error

	ctx, cancel := conn.getContext()

	defer cancel()

	acc, err := conn.Client.GetAccount(ctx, account)
	if err != nil {
		fmt.Println("GET ACCOUNT FAIL")
		fmt.Println(err)
		return false
	}
	fmt.Println(account.Type, account.Address)
	if !bytes.Equal(acc.Address, account.Address) {
		fmt.Println("ACCOUNT NOT FOUND")
		return false
	}

	result := make(map[string]interface{})
	result["address"] = address
	result["type"] = acc.GetType()
	result["balance"] = float64(acc.GetBalance()) / 1000000
	asJSON, _ := json.Marshal(result)
	fmt.Println(JSONPrettyFormat(string(asJSON)))
	fmt.Println("CHECK BALANCE SUCCESS")

	return false
}

func Transfer(address string) (string, error) {

	return "", nil
}


func JSONPrettyFormat(in string) string {
	var out bytes.Buffer
	err := json.Indent(&out, []byte(in), "", "  ")
	if err != nil {
		return in
	}
	return out.String()
}

func (g *GrpcClient) getContext() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), conn.grpcTimeout)
	if len(conn.apiKey) > 0 {
		ctx = metadata.AppendToOutgoingContext(ctx, "TRON-PRO-API-KEY", conn.apiKey)
	}
	return ctx, cancel
}
