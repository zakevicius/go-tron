package cmd

import (
	"fmt"
	api "github.com/coingate/tron-withdrawal-service/protos"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"os"
	"time"
)

var (
	//node 		string
	//timeout		uint32
	withTLS     bool
	apiKey      string
	conn 		GrpcClient
)

type GrpcClient struct {
	Address     string
	Conn        *grpc.ClientConn
	Client      api.WalletClient
	grpcTimeout time.Duration
	opts        []grpc.DialOption
	apiKey      string
}


func NewGrpcClient(address string) *GrpcClient {
	conn := &GrpcClient{
		Address:     address,
		grpcTimeout: 5 * time.Second,
	}

	// load grpc options
	opts := make([]grpc.DialOption, 0)
	if withTLS {
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(nil)))
	} else {
		opts = append(opts, grpc.WithInsecure())
	}

	// check for env API Key
	if trongridKey := os.Getenv("TRONGRID_APIKEY"); len(trongridKey) > 0 {
	apiKey = trongridKey
	}
	// set API
	conn.apiKey = apiKey

	if err := conn.Start(opts...); err != nil {
		fmt.Println(err)
		return nil
	}

	//if len(signer) > 0 {
	//var err error
	//if signerAddress, err = findAddress(signer); err != nil {
	//return err
	//}
	//}

	//var err error
	//passphrase, err = getPassphrase()
	//if err != nil {
	//return err
	//}

	//if len(defaultKeystoreDir) > 0 {
	//// set default directory
	//store.SetDefaultLocation(defaultKeystoreDir)
	//}

	return conn
}

func (g *GrpcClient) Start(opts ...grpc.DialOption) error {
	var err error
	if len(g.Address) == 0 {
		g.Address = "grpc.shasta.trongrid.io:50051"
	}
	g.opts = opts
	g.Conn, err = grpc.Dial(g.Address, opts...)

	if err != nil {
		return fmt.Errorf("Connecting GRPC Client: %v", err)
	}
	g.Client = api.NewWalletClient(g.Conn)

	return nil
}

