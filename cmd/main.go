package main

import (
	"encoding/json"
	"fmt"
	"github.com/coingate/tron-withdrawal-service/cmd/sub"
	"github.com/coingate/tron-withdrawal-service/pkg/csv"
	"os"
	"sync"
)

type PrivateKeys []string
type Address string

var wg sync.WaitGroup

func main() {
	if len(os.Args) < 3 {
		fmt.Println("--json '[array of keys in json format]'")
		fmt.Println("--csv path/to/csv/file.csv")
		fmt.Println("Please provide additional arguments.")
		os.Exit(1)
	}

	var privateKeys PrivateKeys

	switch os.Args[1] {
	case "--json":
		err := json.Unmarshal([]byte(os.Args[2]), &privateKeys)
		if err != nil {
			fmt.Printf("encountered an error while unmarshaling json: %v\n", err)
			os.Exit(1)
		}
	case "--csv":
		csv.ExecuteWithCSV(os.Args[2])
	}

	processKeys(&privateKeys)

	wg.Wait()

	fmt.Println("##################################################")
	fmt.Println("DONE")
	fmt.Println("##################################################")
}


func processKeys(keys *PrivateKeys) {
	keysLength := len(*keys)
	for i := 0; i < keysLength; i++ {
		wg.Add(1)
		go processKey((*keys)[i])
	}

}


//c := make(chan string)
//go func() {
//	fmt.Println("Goroutines:\t", runtime.NumGoroutine())
//	for i := 0; i < keysLength; i++ {
//		c <- keys[i]
//		time.Sleep(100 * time.Millisecond)
//	}
//	close(c)
//}()


func processKey(privateKey string){
	address := cmd.DecodeAddress(privateKey)

	go processWithdrawal(address)
}


func processWithdrawal(address string) {
	wg.Done()
}
