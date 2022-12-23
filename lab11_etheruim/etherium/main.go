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
	fmt.Println("\033[38;5;159m BLOCK NUMBER: \u001B[0m" + header.Number.String())
	blockNumber := big.NewInt(header.Number.Int64())
	block, err := client.BlockByNumber(context.Background(), blockNumber)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Print("\033[38;5;159m Block time: \u001B[0m")
	fmt.Println(block.Time())
	fmt.Print("\033[38;5;159m Block difficulty: \u001B[0m")
	fmt.Println(block.Difficulty().Uint64())
	fmt.Print("\033[38;5;159m Block hash: \u001B[0m")
	fmt.Println(block.Hash().Hex())
	fmt.Print("\033[38;5;159m Block transactions length: \u001B[0m")
	fmt.Println(len(block.Transactions()))
	ctx := context.Background()

	conf := &firebase.Config{
		DatabaseURL: "ваш датабейз юрл",
	}

	opt := option.WithCredentialsFile("путь до вашего ключа")
	app, err := firebase.NewApp(ctx, conf, opt)
	if err != nil {
		fmt.Print("\033[38;5;159m Block time:\u001B[0m")
		log.Fatalln("error in initializing firebase app: ", err)
	}

	clienteth, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("error in creating firebase DB client: ", err)
	}

	ref := clienteth.NewRef(fmt.Sprint(block.Number().Uint64()))

	if err := ref.Set(context.TODO(), map[string]interface{}{"blockTime": block.Time(), "blockDifficulty": block.Difficulty().Uint64(), "blockHashHex": block.Hash().Hex(), "blockTransactions": block.Transactions()}); err != nil {
		log.Fatal(err)
	}

	fmt.Print("\033[38;5;28m Block added successfully!\u001B[0m")

}
