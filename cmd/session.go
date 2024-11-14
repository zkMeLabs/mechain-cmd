package main

import (
	"context"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/evmos/evmos/v12/x/evm/precompiles/storage"
)

const (
	EmptyEvmAddress     = "0x0000000000000000000000000000000000000000"
	BankAddress         = "0x0000000000000000000000000000000000001000"
	AuthAddress         = "0x0000000000000000000000000000000000001001"
	GovAddress          = "0x0000000000000000000000000000000000001002"
	StakingAddress      = "0x0000000000000000000000000000000000001003"
	DistributionAddress = "0x0000000000000000000000000000000000001004"
	SlashingAddress     = "0x0000000000000000000000000000000000001005"
	EvidenceAddress     = "0x0000000000000000000000000000000000001006"
	EpochsAddress       = "0x0000000000000000000000000000000000001007"
	AuthzAddress        = "0x0000000000000000000000000000000000001008"
	FeemarketAddress    = "0x0000000000000000000000000000000000001009"
	VirtualGroupAddress = "0x0000000000000000000000000000000000002000"
	StorageAddress      = "0x0000000000000000000000000000000000002001"
	SpAddress           = "0x0000000000000000000000000000000000002002"
)

const (
	DefaultGasLimit = 180000
)

func CreateTxOpts(ctx context.Context, client *ethclient.Client, hexPrivateKey string, chain *big.Int, gasLimit uint64, nonce uint64) (*bind.TransactOpts, error) {
	// create private key
	privateKey, err := crypto.HexToECDSA(hexPrivateKey)
	if err != nil {
		return nil, err
	}

	// Build transact tx opts with private key
	txOpts, err := bind.NewKeyedTransactorWithChainID(privateKey, chain)
	if err != nil {
		return nil, err
	}

	// set gas limit and gas price
	txOpts.GasLimit = gasLimit
	gasPrice, err := client.SuggestGasPrice(ctx)
	if err != nil {
		return nil, err
	}
	txOpts.GasPrice = gasPrice

	txOpts.Nonce = big.NewInt(int64(nonce))

	return txOpts, nil
}

func CreateStorageSession(client *ethclient.Client, txOpts bind.TransactOpts, contractAddress string) (*storage.IStorageSession, error) {
	storageContract, err := storage.NewIStorage(common.HexToAddress(contractAddress), client)
	if err != nil {
		return nil, err
	}
	session := &storage.IStorageSession{
		Contract: storageContract,
		CallOpts: bind.CallOpts{
			Pending: false,
		},
		TransactOpts: txOpts,
	}
	return session, nil
}
