package main

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"time"

	ct "github.com/blockchain-tps-test/samples/theta/accessors"
	"github.com/thetatoken/theta/common"
	"github.com/thetatoken/theta/crypto"
	"github.com/thetatoken/theta/ledger/types"
	"github.com/thetatoken/theta/rpc"
	"github.com/thetatoken/thetasubchain/eth/abi/bind"
	"github.com/ybbus/jsonrpc"
)

func init_token(c EthClient, privKeyString []string) {
	authKey, err := crypto.HexToECDSA(privKeyString[0])
	if err != nil {
		log.Fatal(err)
	}
	authPublicKey := authKey.Public()
	publicKeyECDSA, ok := authPublicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("error casting public key to ECDSA")
	}

	authAddress := pubkeyToAddress(*publicKeyECDSA)
	for _, priv := range privKeyString {
		privateKey, err := crypto.HexToECDSA(priv)
		if err != nil {
			log.Fatal(err)
		}

		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatal("error casting public key to ECDSA")
		}

		fromAddress := pubkeyToAddress(*publicKeyECDSA)
		nonce, err := c.client.PendingNonceAt(context.Background(), authAddress)
		transfer_token(authAddress, fromAddress, c, int64(nonce))

		subchainTNT20Address := common.HexToAddress("0x5C3159dDD2fe0F9862bC7b7D60C1875fa8F81337") // subchain 0x5C3159dDD2fe0F9862bC7b7D60C1875fa8F81337 mainchain 0x59AF421cB35fc23aB6C8ee42743e6176040031f4
		//erc20TokenBank, err := ct.NewTNT20TokenBank(common.HexToAddress(TokenBankAddress), c.client)
		subchainTNT20Instance, _ := ct.NewMockTNT20(subchainTNT20Address, c.client)
		// auth action
		gasPrice := c.getGasPriceSuggestion(context.Background())
		nonce, err = c.client.PendingNonceAt(context.Background(), authAddress)
		auth, err := bind.NewKeyedTransactorWithChainID(crypto.ECDSAToPrivKey(authKey), chainID)
		if err != nil {
			log.Fatal(err)
		}
		auth.Nonce = big.NewInt(int64(nonce))
		auth.Value = common.Big0
		auth.GasLimit = uint64(3000000) // in units
		auth.GasPrice = &gasPrice
		_, err1 := subchainTNT20Instance.Mint(auth, fromAddress, big.NewInt(9999999999999000))
		auth.Nonce = big.NewInt(int64(nonce + 1))
		if err1 != nil {
			fmt.Println(err1)
		}
		// account approve
		gasPrice = c.getGasPriceSuggestion(context.Background())
		nonce, err = c.client.PendingNonceAt(context.Background(), fromAddress)
		accountAuth, err := bind.NewKeyedTransactorWithChainID(crypto.ECDSAToPrivKey(privateKey), chainID)
		if err != nil {
			log.Fatal(err)
		}
		accountAuth.Nonce = big.NewInt(int64(nonce))
		accountAuth.Value = common.Big0
		accountAuth.GasLimit = uint64(3000000) // in units
		accountAuth.GasPrice = &gasPrice
		_, err1 = subchainTNT20Instance.Approve(accountAuth, common.HexToAddress(TokenBankAddress), big.NewInt(9999999999999000))
		if err1 != nil {
			fmt.Println(err1)
		}
		time.Sleep(500 * time.Millisecond)
		balance, _ := subchainTNT20Instance.BalanceOf(nil, fromAddress)
		allowance, _ := subchainTNT20Instance.Allowance(nil, fromAddress, common.HexToAddress(TokenBankAddress))
		fmt.Println(fromAddress, " token balance is", balance, " allowance is ", allowance)

	}

	fmt.Println("Init token done!")

}

func transfer_token(fromAddress, toAddress common.Address, c EthClient, nonce int64) {
	if fromAddress == toAddress {
		return
	}
	pri, err := hex.DecodeString(privs[0])
	thetaPrivateKey, err := crypto.PrivateKeyFromBytes(pri)
	if err != nil {
		log.Fatal(err)
	}
	theta := big.NewInt(0)
	tfuel := big.NewInt(1).Mul(big.NewInt(3e18), big.NewInt(100))
	fee := big.NewInt(3e17)
	inputs := []types.TxInput{{
		Address: fromAddress,
		Coins: types.Coins{
			TFuelWei: new(big.Int).Add(tfuel, fee),
			ThetaWei: theta,
		},
		Sequence: uint64(nonce + 1),
	}}
	outputs := []types.TxOutput{{
		Address: toAddress,
		Coins: types.Coins{
			TFuelWei: tfuel,
			ThetaWei: theta,
		},
	}}
	sendTx := &types.SendTx{
		Fee: types.Coins{
			ThetaWei: new(big.Int).SetUint64(0),
			TFuelWei: fee,
		},
		Inputs:  inputs,
		Outputs: outputs,
	}
	sig, err := thetaPrivateKey.Sign(sendTx.SignBytes("tsub360001")) //privatenet,testnet,tsub360777
	if err != nil {
		log.Fatalln("Failed to sign transaction: %v\n", err)
	}
	sendTx.SetSignature(fromAddress, sig)

	raw, err := types.TxToBytes(sendTx)
	if err != nil {
		log.Fatalln("Failed to encode transaction: %v\n", err)
	}
	signedTx := hex.EncodeToString(raw)
	var res *jsonrpc.RPCResponse
	res, err = c.rpcClient.Call("theta.BroadcastRawTransactionAsync", rpc.BroadcastRawTransactionArgs{TxBytes: signedTx})
	if err != nil {
		log.Fatalln("Failed to broadcast transaction: %v\n", err)
	}
	if res.Error != nil {
		log.Fatalln("Server returned error: %v\n", res.Error)
	}
	result := &rpc.BroadcastRawTransactionResult{}
	err = res.GetObject(result)
	if err != nil {
		log.Fatalln("Failed to parse server response: %v\n", err)
	}

}
