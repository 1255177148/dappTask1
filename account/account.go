package account

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
)

var WeiPerEth = new(big.Int).Exp(big.NewInt(10), big.NewInt(18), nil)

// TransferETH 交易ETH
//
// 参数:
//
//	privateKey  - 用户私钥
//	to          - 接收ETH地址
//	amount      - 交易的EHT数量，单位是EHT，例如0.1ETH
//	client      - 客户端
func TransferETH(privateKey string, to *common.Address, amount float64, client *ethclient.Client) {
	private, err := crypto.HexToECDSA(privateKey)
	if err != nil {
		log.Fatal(err)
	}
	publicKey := private.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}
	fromAddress := crypto.PubkeyToAddress(*publicKeyECDSA)
	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatal(err)
	}
	value := EthToWei(big.NewFloat(amount))
	gasLimit := uint64(21000)
	head, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	baseFee := head.BaseFee
	// 你愿意给矿工的小费（GasTipCap）
	tipCap := big.NewInt(2_000_000_000) // 2 gwei

	// GasFeeCap = BaseFee + TipCap * 2 作为一个保险上限
	feeCap := new(big.Int).Add(baseFee, new(big.Int).Mul(tipCap, big.NewInt(2)))
	chainID, _ := client.NetworkID(context.Background())
	txData := &types.DynamicFeeTx{
		ChainID:   chainID,
		Nonce:     nonce,
		GasTipCap: tipCap,
		GasFeeCap: feeCap,
		Gas:       gasLimit,
		To:        to,
		Value:     value,
		Data:      nil,
	}
	tx := types.NewTx(txData)
	signedTx, err := types.SignTx(tx, types.LatestSignerForChainID(chainID), private) // 对事务进行签名
	if err != nil {
		log.Fatal(err)
	}
	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("事务已发送：%s\n", signedTx.Hash().Hex())
	receipt, err := bind.WaitMined(context.Background(), client, signedTx)
	if err != nil {
		log.Fatal(err)
	}

	if receipt.Status == types.ReceiptStatusSuccessful {
		fmt.Printf("交易成功! tx hash: %s\n", signedTx.Hash().Hex())
	} else {
		fmt.Printf("交易失败! tx hash: %s\n", signedTx.Hash().Hex())
	}
}

// EthToWei ETH -> Wei
func EthToWei(eth *big.Float) *big.Int {
	wei := new(big.Float).Mul(eth, new(big.Float).SetInt(WeiPerEth))
	result := new(big.Int)
	wei.Int(result) // 截断小数部分
	return result
}
