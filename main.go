package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	client, err := ethclient.Dial("http://api.nightly-cluster-2.iotex.one:15014")
	if err != nil {
		log.Fatalf("connect rpc server error: %v", err)
	}

	keyBytes, err := hex.DecodeString("replace by your private key")
	if err != nil {
		log.Fatalf("decode private key error: %v", err)
	}
	key, err := crypto.ToECDSA(keyBytes)
	if err != nil {
		log.Fatalf("create esdsa private key from key bytes error: %v", err)
	}
	from := crypto.PubkeyToAddress(key.PublicKey)
	nonce, err := client.NonceAt(context.Background(), from, nil)
	if err != nil {
		log.Fatalf("get account nonce error: %v", err)
	}

	stakingABI, err := abi.JSON(bytes.NewReader([]byte(`[
		{
		  "inputs": [
			{
			  "internalType": "string",
			  "name": "candName",
			  "type": "string"
			},
			{
			  "internalType": "uint256",
			  "name": "amount",
			  "type": "uint256"
			},
			{
			  "internalType": "uint32",
			  "name": "duration",
			  "type": "uint32"
			},
			{
			  "internalType": "bool",
			  "name": "autoStake",
			  "type": "bool"
			},
			{
			  "internalType": "uint8[]",
			  "name": "data",
			  "type": "uint8[]"
			}
		  ],
		  "name": "createStake",
		  "outputs": [],
		  "stateMutability": "nonpayable",
		  "type": "function"
		}
	  ]
	`)))
	if err != nil {
		log.Fatalf("parse abi json error: %v", err)
	}

	amount, _ := new(big.Int).SetString("200000000000000000000", 10)
	data, err := stakingABI.Pack("createStake", "hashquark", amount, uint32(10), false, []uint8{})
	if err != nil {
		log.Fatalf("pack data error: %v", err)
	}

	rawTx := types.NewTransaction(
		nonce,
		common.HexToAddress("0x000000000000007374616b696e67437265617465"),
		big.NewInt(0),
		100000,
		big.NewInt(1000000000000),
		data,
	)

	sig, err := crypto.Sign(rawTx.Hash().Bytes(), key)
	if err != nil {
		log.Fatalf("sign tx error: %v", err)
	}
	signer := types.NewEIP155Signer(big.NewInt(4689))
	tx, err := rawTx.WithSignature(signer, sig)
	if err != nil {
		log.Fatalf("compose tx error: %v", err)
	}

	err = client.SendTransaction(context.Background(), tx)
	if err != nil {
		log.Fatalf("send tx error: %v", err)
	}
}
