package block

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
)

func GetBlockInfo(client *ethclient.Client) {
	block, err := client.BlockByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	if block == nil {
		log.Fatal("nil block")
	}
	fmt.Printf("Block number: %d\n", block.Number())
	fmt.Printf("Block hash: %s\n", block.Hash().Hex())
	fmt.Printf("Transaction number: %d\n", block.Transactions().Len())
	fmt.Printf("Timestamp: %d\n", block.Time())
}
