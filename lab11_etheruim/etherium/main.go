package main

import (
	"context"
	"fmt"

	"log"
	"math/big"

	firebase "firebase.google.com/go"
	"github.com/ethereum/go-ethereum/ethclient"
	"google.golang.org/api/option"
)

func main() {
	client, err := ethclient.Dial("ваш ключ инфуры")
	if err != nil {
		log.Fatalln(err)
	}
	header, err := client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(header.Number.String()) // The lastes block in blockchain because nil pointer in header
	blockNumber := big.NewInt(header.Number.Int64())
	block, err := client.BlockByNumber(context.Background(), blockNumber) //get block with this number
	if err != nil {
		log.Fatal(err)
	}
	// all info about block
	fmt.Println(block.Number().Uint64())
	fmt.Println(block.Time())
	fmt.Println(block.Difficulty().Uint64())
	fmt.Println(block.Hash().Hex())
	fmt.Println(len(block.Transactions()))
	ctx := context.Background()

	// configure database URL
	conf := &firebase.Config{
		DatabaseURL: "ЮРЛ realtime database",
	}

	// fetch service account key
	opt := option.WithCredentialsFile("путь до ключа firebase)

	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		log.Fatalln("error in initializing firebase app: ", err)
	}

	clienteth, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("error in creating firebase DB client: ", err)
	}

	// create ref at path user_scores/:userId
	ref := clienteth.NewRef(fmt.Sprint(block.Number().Uint64()))

	if err := ref.Set(context.TODO(), map[string]interface{}{"blockTime": block.Time(), "blockDifficulty": block.Difficulty().Uint64(), "blockHashHex": block.Hash().Hex(), "blockTransactions": block.Transactions()}); err != nil {
		log.Fatal(err)
	}

	fmt.Println("score added/updated successfully!")
}
