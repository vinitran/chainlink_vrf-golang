package database

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/uptrace/bun"
	"time"
)

type RequestRandom struct {
	bun.BaseModel `bun:"table:request_random,alias:req"`
	Id            int       `bun:"id,pk,autoincrement" json:"id"`
	User          string    `bun:"wallet_address,notnull" json:"user"`
	RequestId     string    `bun:"request_id,notnull" json:"requestId"`
	Amount        int       `bun:"amount,notnull" json:"amount"`
	TxHash        string    `bun:"transaction_hash,notnull" json:"txHash"`
	Index         int       `bun:"index,notnull" json:"index"`
	Time          time.Time `bun:"time,notnull" json:"time"`
}

type ResponseRandom struct {
	bun.BaseModel `bun:"table:response_random,alias:res"`
	Id            int       `bun:"id,pk,autoincrement" json:"id"`
	User          string    `bun:"wallet_address,notnull" json:"user"`
	RequestId     string    `bun:"request_id,notnull" json:"requestId"`
	PrizeIds      []int     `bun:"prize_ids,notnull" json:"prizeIds"`
	TxHash        string    `bun:"transaction_hash,notnull" json:"txHash"`
	Index         int       `bun:"index,notnull" json:"index"`
	Time          time.Time `bun:"time,notnull" json:"time"`
}

type BlockError struct {
	bun.BaseModel `bun:"table:error_block"`
	Id            int `bun:"id,pk,autoincrement" json:"id"`
	Block         int `bun:"block,notnull" json:"block""`
}

func (e *RequestRandom) String() (*string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	dataStr := string(data)
	return &dataStr, nil
}

func (e *ResponseRandom) String() (*string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return nil, err
	}

	dataStr := string(data)
	return &dataStr, nil
}

func CreateTable(db *bun.DB) error {
	err := createRequestRandomTable(db)
	if err != nil {
		return err
	}

	err = createResponseRandomTable(db)
	if err != nil {
		return err
	}

	err = createBlockErrorTable(db)
	if err != nil {
		return err
	}

	return nil
}

func InsertRequestRandomToDb(db *bun.DB, data []RequestRandom) error {
	if data == nil {
		return nil
	}
	//
	//var events []RequestRandom
	//for _, event := range data {
	//	events = append(events, RequestRandom{
	//		User:      event.User,
	//		RequestId: event.RequestId,
	//		Amount:    event.Amount,
	//		TxHash:    event.TxHash,
	//		Index:     event.Index,
	//	})
	//}
	_, err := db.NewInsert().
		Model(&data).
		Exec(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("database: inserted to db")
	return nil
}

func InsertResponseRandomToDb(db *bun.DB, data []ResponseRandom) error {
	if data == nil {
		return nil
	}

	//var events []ResponseRandom
	//for _, event := range data {
	//	events = append(events, ResponseRandom{
	//		User:      event.User,
	//		RequestId: event.RequestId,
	//		PrizeIds:  event.PrizeIds,
	//		TxHash:    event.TxHash,
	//		Index:     event.Index,
	//	})
	//}

	_, err := db.NewInsert().
		Model(&data).
		Exec(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("database: inserted to db")
	return nil
}

func InsertBlockErrorToDb(db *bun.DB, block int) error {
	blockErr := BlockError{
		Block: block,
	}

	_, err := db.NewInsert().
		Model(&blockErr).
		Exec(context.Background())
	if err != nil {
		return err
	}

	fmt.Println("database: inserted to db")
	return nil
}

func createRequestRandomTable(db *bun.DB) error {
	_, err := db.NewCreateTable().
		Model((*RequestRandom)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func createResponseRandomTable(db *bun.DB) error {
	_, err := db.NewCreateTable().
		Model((*ResponseRandom)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func createBlockErrorTable(db *bun.DB) error {
	_, err := db.NewCreateTable().
		Model((*BlockError)(nil)).
		IfNotExists().
		Exec(context.Background())
	if err != nil {
		return err
	}

	return nil
}
