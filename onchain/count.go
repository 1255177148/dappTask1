package onchain

import (
	"context"
	"fmt"
	"github.com/1255177148/dappTask1/contract"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"log"
	"math/big"
	"time"
)

type Count struct {
	Instance        *contract.Contract
	ContractAddress *common.Address
}

var client *ethclient.Client

func NewCount(address string, _client *ethclient.Client) (*Count, error) {
	contractAddr := common.HexToAddress(address)
	token, err := contract.NewContract(contractAddr, _client)
	if err != nil {
		return nil, err
	}
	client = _client
	return &Count{Instance: token, ContractAddress: &contractAddr}, nil
}

func (count *Count) Accumulate(private string) (string, error) {
	privateKey, err := crypto.HexToECDSA(private) // 解析私钥
	if err != nil {
		panic(err)
	}
	ctx := context.Background()
	chainID, err := client.ChainID(ctx)
	if err != nil {
		return "", fmt.Errorf("获取链ID失败: %w", err)
	}
	// 创建托管账户的 TransactOpts 对象，绑定私钥和 chainID
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	if err != nil {
		return "", fmt.Errorf("创建Transactor失败: %w", err)
	}
	nonce, err := client.PendingNonceAt(ctx, auth.From)
	if err != nil {
		return "", fmt.Errorf("获取 pending nonce 失败: %w", err)
	}
	auth.Nonce = big.NewInt(int64(nonce))
	auth.GasLimit = 0
	tipCap, err := client.SuggestGasTipCap(ctx)
	if err != nil {
		return "", fmt.Errorf("获取建议小费单价失败: %w", err)
	}
	header, err := client.HeaderByNumber(ctx, nil)
	if err != nil {
		return "", fmt.Errorf("获取最新区块头失败: %w", err)
	}
	baseFee := header.BaseFee
	if baseFee == nil {
		baseFee = big.NewInt(1e9)
	}
	feeCap := new(big.Int).Add(new(big.Int).Mul(baseFee, big.NewInt(2)), tipCap)
	auth.GasTipCap = tipCap
	auth.GasFeeCap = feeCap
	tx, err := count.Instance.Accumulate(auth)
	if err != nil {
		return "", err
	}
	// ---------- 等待上链（打包） ----------
	waitCtx, cancel := context.WithTimeout(ctx, 2*time.Minute)
	receipt, err := bind.WaitMined(waitCtx, client, tx)
	cancel()
	if err != nil {
		log.Fatalf("等待上链失败: %v", err)
	}
	// ---------- 检查执行结果 ----------
	if receipt.Status != types.ReceiptStatusSuccessful {
		log.Fatalf("交易上链但执行失败，status=%d，区块=%d", receipt.Status, receipt.BlockNumber.Uint64())
	}
	fmt.Printf("交易成功！区块：%d，GasUsed：%d\n", receipt.BlockNumber.Uint64(), receipt.GasUsed)
	fmt.Printf("txHash: %s\n", tx.Hash().Hex())
	return tx.Hash().Hex(), nil
}

func (count *Count) GetCount() (string, error) {
	callOpts := &bind.CallOpts{
		Pending: false, // 表示查询最新区块
		Context: context.Background(),
		From:    common.Address{}, // 是solidity中的msg.sender，如果不需要可以不传
	}
	result, err := count.Instance.GetCount(callOpts)
	if err != nil {
		return "", err
	}
	return result.String(), nil
}
