package main

import (
	"context"
	"log"
	"math/big"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

func main() {
	appCtx, cancel := context.WithCancel(context.Background())

	rpcUrl := "wss://wss-testnet.bitkubchain.io"
	rpcClient, err := rpc.DialOptions(appCtx, rpcUrl)
	if err != nil {
		panic(err)
	}

	evmClient := ethclient.NewClient(rpcClient)

	blockNumber, err := evmClient.BlockNumber(appCtx)
	if err != nil {
		panic(err)
	}
	log.Println("Block number:", blockNumber)

	blockDetail, err := evmClient.BlockByNumber(appCtx, big.NewInt(int64(blockNumber)))
	if err != nil {
		panic(err)
	}
	log.Println("Gas used:", blockDetail.GasUsed())

	headerChannel := make(chan *types.Header)
	newBlockSubscriber, err := evmClient.SubscribeNewHead(appCtx, headerChannel)
	if err != nil {
		panic(err)
	}

	for {
		select {
		case header := <-headerChannel:
			log.Println("New block number:", header.Number)
			blockDetail, err := evmClient.BlockByNumber(appCtx, header.Number)
			if err != nil {
				log.Println("Error:", err)
				continue
			}

			log.Println("Gas used:", blockDetail.GasUsed())
		case <-appCtx.Done():
			log.Println("Application is terminated")
			cancel()
			evmClient.Close()
			newBlockSubscriber.Unsubscribe()
			return
		}
	}
}
