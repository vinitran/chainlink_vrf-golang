package event

import (
	"VRFChainlink/database"
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/uptrace/bun"
	"math/big"
	"time"
)

type TrackingEvent struct {
	Client  *ethclient.Client
	Address common.Address
}

var (
	requestCreatedHash  = crypto.Keccak256Hash([]byte("RequestCreated(address,uint256,uint256)")).Hex()
	responseCreatedHash = crypto.Keccak256Hash([]byte("ResponseCreated(address,uint256,uint256[])")).Hex()
)

func NewEventTracking(rpc, address string) (*TrackingEvent, error) {
	client, err := ethclient.Dial(rpc)
	if err != nil {
		return nil, err
	}

	addr := common.HexToAddress(address)
	tracking := &TrackingEvent{
		Client:  client,
		Address: addr,
	}

	return tracking, nil
}

func (tracking *TrackingEvent) GetEventFromBlockNumber(db *bun.DB, number *big.Int) {
	for i := number.Int64(); ; i = i + 5000 {
		blockNumber := big.NewInt(i)

		req, res, err := tracking.GetEventByBlockNumber(blockNumber)
		if err != nil {
			fmt.Println(blockNumber, err)
			err = database.InsertBlockErrorToDb(db, int(blockNumber.Int64()))
			if err != nil {
				fmt.Println("insert block error to db:", err)
			}
			continue
		}

		err = database.InsertRequestRandomToDb(db, req)
		if err != nil {
			fmt.Println("insert request to db:", err)
			err = database.InsertBlockErrorToDb(db, int(blockNumber.Int64()))
			if err != nil {
				fmt.Println("insert block error to db:", err)
			}
			continue
		}

		err = database.InsertResponseRandomToDb(db, res)
		if err != nil {
			err = database.InsertBlockErrorToDb(db, int(blockNumber.Int64()))
			if err != nil {
				fmt.Println("insert block error to db:", err)
			}
			fmt.Println("insert response to db:", err)
		}

		fmt.Println(i)
	}
}

func (tracking *TrackingEvent) GetEventByBlockNumber(number *big.Int) ([]database.RequestRandom, []database.ResponseRandom, error) {
	lastestBlockNumber, err := tracking.GetLatestBlockNumber()
	if err != nil {
		return nil, nil, err
	}

	query := ethereum.FilterQuery{
		FromBlock: number,
		ToBlock:   new(big.Int).SetInt64(number.Int64() + 5000),
		Addresses: []common.Address{
			tracking.Address,
		},
	}

	if number.Int64()+5000 > lastestBlockNumber.Int64() {
		timeDelay := 60
		time.Sleep(time.Duration(timeDelay) * time.Second)
		query = ethereum.FilterQuery{
			FromBlock: number,
			ToBlock:   nil,
			Addresses: []common.Address{
				tracking.Address,
			},
		}
	}

	logs, err := tracking.Client.FilterLogs(context.Background(), query)
	if err != nil {
		return nil, nil, err
	}

	var request []database.RequestRandom
	var response []database.ResponseRandom
	for _, vLog := range logs {
		switch vLog.Topics[0].Hex() {
		case requestCreatedHash:
			timeStamp, err := tracking.GetTimeOfBlock(big.NewInt(int64(vLog.BlockNumber)))
			if err != nil {
				return nil, nil, err
			}

			request = append(request, database.RequestRandom{
				User:      common.HexToAddress(vLog.Topics[1].Hex()).String(),
				RequestId: new(big.Int).SetBytes(vLog.Topics[2].Bytes()).String(),
				Amount:    int(new(big.Int).SetBytes(vLog.Data).Int64()),
				TxHash:    vLog.TxHash.String(),
				Index:     int(vLog.Index),
				Time:      timeStamp,
			})
		case responseCreatedHash:
			var prizeIds []int
			for i := 64; i+32 < len(vLog.Data); i = i + 32 {
				prizeIds = append(prizeIds, int(new(big.Int).SetBytes(vLog.Data[i:i+32]).Int64()))
			}

			timeStamp, err := tracking.GetTimeOfBlock(big.NewInt(int64(vLog.BlockNumber)))
			if err != nil {
				return nil, nil, err
			}

			response = append(response, database.ResponseRandom{
				User:      common.HexToAddress(vLog.Topics[1].Hex()).String(),
				RequestId: new(big.Int).SetBytes(vLog.Topics[2].Bytes()).String(),
				PrizeIds:  prizeIds,
				TxHash:    vLog.TxHash.String(),
				Index:     int(vLog.Index),
				Time:      timeStamp,
			})
		}
	}
	return request, response, nil
}

func (tracking *TrackingEvent) GetLatestBlockNumber() (*big.Int, error) {
	header, err := tracking.Client.HeaderByNumber(context.Background(), nil)
	if err != nil {
		return nil, err
	}

	return header.Number, nil
}

func (tracking *TrackingEvent) GetTimeOfBlock(block *big.Int) (time.Time, error) {
	header, err := tracking.Client.HeaderByNumber(context.Background(), block)
	if err != nil {
		return time.Time{}, err
	}

	timeStamp := time.Unix(int64(header.Time), 0)
	return timeStamp, nil
}
