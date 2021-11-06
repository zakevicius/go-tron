package transaction

import (
	"github.com/coingate/tron-withdrawal-service/cmd/sub"
	api "github.com/coingate/tron-withdrawal-service/protos"
	"github.com/ethereum/go-ethereum/accounts/keystore"
)

type sender struct {
	ks      *keystore.KeyStore
	account *keystore.Account
}

type Controller struct {
	executionError error
	resultError    error
	client         *cmd.GrpcClient
	tx             *api.Transaction
	sender         sender
	Behavior       behavior
	Result         *api.Return
	Receipt        *api.TransactionInfo
}

type behavior struct {
	DryRun               bool
	SigningImpl          SignerImpl
	ConfirmationWaitTime uint32
}

type SignerImpl int

const (
	Software SignerImpl = iota
	Ledger
)

func NewController(
	client *cmd.GrpcClient,
	senderKs *keystore.KeyStore,
	senderAcct *keystore.Account,
	tx *api.Transaction,
	options ...func(*Controller),
) *Controller {

	ctrlr := &Controller{
		executionError: nil,
		resultError:    nil,
		client:         client,
		sender: sender{
			ks:      senderKs,
			account: senderAcct,
		},
		tx:       tx,
		Behavior: behavior{false, Software, 0},
	}
	for _, option := range options {
		option(ctrlr)
	}
	return ctrlr
}
