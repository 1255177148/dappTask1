package main

import (
	"fmt"
	"github.com/1255177148/dappTask1/account"
	"github.com/1255177148/dappTask1/block"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/27NL_0zlbK15k86qzL4emhJMRO-kDoPX")
	if err != nil {
		fmt.Println("连接ETH失败：", err)
	}
	block.GetBlockInfo(client)
	to := common.HexToAddress("0xa43d2b78416B4B1efce69136f41aeF1691378C9A")
	account.TransferETH("cabb9d1405205e92b2984ac19fbf28b17432d1f0af889d867a5df7e0e851cf4b",
		&to,
		0.05,
		client)
}
