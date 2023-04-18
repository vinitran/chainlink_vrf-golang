package main

import (
	"VRFChainlink/event"
	"github.com/joho/godotenv"
	"log"
	"math/big"
	"os"
	"strconv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	//db, err := database.ConnectDatabase()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//
	//gin := api.NewGin(db)
	//gin.Run()

	trackingTx, err := event.NewEventTracking(os.Getenv("RPC"), os.Getenv("CONTRACT_ADDRESS"))
	if err != nil {
		log.Fatal(err)
	}

	fromBlock, err := strconv.ParseInt(os.Getenv("FROM_BLOCK"), 10, 64)
	if err != nil {
		log.Fatal(err)
	}

	trackingTx.GetEventFromBlockNumber(db, big.NewInt(fromBlock))
}
