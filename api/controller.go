package api

import (
	"VRFChainlink/database"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"strconv"
	"time"
)

type PageFilter struct {
	Page int
	Size int
}

func (p *PageFilter) Check(c *gin.Context) error {
	page, err := strconv.Atoi(c.DefaultQuery("page", "1"))
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("%s", err))
	}
	p.Page = page
	if p.Page <= 0 {
		return fmt.Errorf("error: page must be greater than 0")
	}

	size, err := strconv.Atoi(c.DefaultQuery("size", "10"))
	if err != nil {
		return fmt.Errorf(fmt.Sprintf("%s", err))
	}
	p.Size = size
	if p.Size <= 0 {
		return fmt.Errorf("error: size must be greater than 0")
	}

	return nil
}

func SearchByTime(c *gin.Context, query *bun.SelectQuery) error {
	strTimeFrom, isTimeFrom := c.GetQuery("from_time")
	strTimeTo, isTimeTo := c.GetQuery("to_time")
	if !isTimeFrom && !isTimeTo {
		return nil
	}

	if !isTimeFrom {
		timeTo, err := time.Parse(time.RFC3339, strTimeTo)
		if err != nil {
			return fmt.Errorf("invalid time format. Use RFC3339 format")
		}

		query = query.Where("time <= ?", timeTo)
		return nil
	}

	if !isTimeTo {
		timeFrom, err := time.Parse(time.RFC3339, strTimeFrom)
		if err != nil {
			return fmt.Errorf("invalid time format. Use RFC3339 format")
		}

		query = query.Where("time >= ?", timeFrom)
		return nil
	}

	timeTo, err := time.Parse(time.RFC3339, strTimeTo)
	if err != nil {
		return fmt.Errorf("invalid time format. Use RFC3339 format")
	}

	timeFrom, err := time.Parse(time.RFC3339, strTimeFrom)
	if err != nil {
		return fmt.Errorf("invalid time format. Use RFC3339 format")
	}

	if timeFrom.After(timeTo) {
		return fmt.Errorf("invalid time value. to_time must be greater than from_time")
	}

	query = query.Where("time >= ?", timeFrom).Where("time <= ?", timeTo)
	return nil
}

func GetRequestTransactionByHash(c *gin.Context) ([]database.RequestRandom, error) {
	responseData := new([]database.RequestRandom)
	hash := c.Param("hash")
	query := db.NewSelect().Model(responseData).
		Where("transaction_hash = ?", hash)

	err := SearchByTime(c, query)
	if err != nil {
		return nil, err
	}

	err = query.Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return *responseData, err
}

func GetResponseTransactionByHash(c *gin.Context) ([]database.ResponseRandom, error) {
	responseData := new([]database.ResponseRandom)
	hash := c.Param("hash")
	query := db.NewSelect().Model(responseData).
		Where("transaction_hash = ?", hash)

	err := SearchByTime(c, query)
	if err != nil {
		return nil, err
	}

	err = query.Scan(context.Background())
	if err != nil {
		return nil, err
	}

	return *responseData, err
}

func SearchByWalletAddress(c *gin.Context, query *bun.SelectQuery) {
	walletAddress, ok := c.GetQuery("wallet_address")
	if !ok {
		return
	}

	query = query.Where("wallet_address = ?", walletAddress)
	return
}

func SearchByTxHash(c *gin.Context, query *bun.SelectQuery) {
	txHash, ok := c.GetQuery("transaction_hash")
	if !ok {
		return
	}

	query = query.Where("transaction_hash = ?", txHash)
	return
}

func SortByAmount(c *gin.Context, query *bun.SelectQuery) error {
	sort, ok := c.GetQuery("sort")
	if !ok {
		return nil
	}

	if sort == "desc" {
		query = query.Order("amount DESC")
		return nil
	}

	if sort == "asc" {
		query = query.Order("amount ASC")
		return nil
	}

	return fmt.Errorf("error: invalid value for sort_by_amount, only esc or desc")
}

func SearchByTotalAmount(c *gin.Context, query *bun.SelectQuery) error {
	strAmountFrom, isAmountFrom := c.GetQuery("from_amount")
	strAmountTo, isAmountTo := c.GetQuery("to_amount")

	if !isAmountFrom && !isAmountTo {
		return nil
	}

	if !isAmountFrom {
		amountTo, err := strconv.ParseFloat(strAmountTo, 64)
		if err != nil {
			return fmt.Errorf("error: invalid type value for amount_to, only float type")
		}

		query = query.Having("sum(amount) <= ?", amountTo)
		return nil
	}

	if !isAmountTo {
		amountFrom, err := strconv.ParseFloat(strAmountFrom, 64)
		if err != nil {
			return fmt.Errorf("error: invalid type value for amount_from, only float type")
		}

		query = query.Having("sum(amount) >= ?", amountFrom)
		return nil
	}

	amountTo, err := strconv.ParseFloat(strAmountTo, 64)
	if err != nil {
		return fmt.Errorf("error: invalid type value for amount_to, only float type")
	}

	amountFrom, err := strconv.ParseFloat(strAmountFrom, 64)
	if err != nil {
		return fmt.Errorf("error: invalid type value for amount_from, only float type")
	}

	if amountFrom >= amountTo {
		return fmt.Errorf("error: invalid value for amount_from and amount_to, amount_to must be greater than amount_from")
	}

	query = query.Having("sum(amount) >= ?", amountFrom).Having("sum(amount) <= ?", amountTo)
	return nil
}

func SortByTotalAmount(c *gin.Context, query *bun.SelectQuery) error {
	sort, ok := c.GetQuery("sort")
	if !ok {
		return nil
	}

	if sort == "desc" {
		query = query.Order("total_amount DESC")
		return nil
	}

	if sort == "asc" {
		query = query.Order("total_amount ASC")
		return nil
	}

	return fmt.Errorf("error: invalid value for sort_by_amount, only asc or desc")
}

func SearchByAmount(c *gin.Context, query *bun.SelectQuery) error {
	strAmountFrom, isAmountFrom := c.GetQuery("from_amount")
	strAmountTo, isAmountTo := c.GetQuery("to_amount")

	if !isAmountFrom && !isAmountTo {
		return nil
	}

	if !isAmountFrom {
		amountTo, err := strconv.ParseFloat(strAmountTo, 64)
		if err != nil {
			return fmt.Errorf("error: invalid type value for amount_to, only float type")
		}

		query = query.Where("amount <= ?", amountTo)
		return nil
	}

	if !isAmountTo {
		amountFrom, err := strconv.ParseFloat(strAmountFrom, 64)
		if err != nil {
			return fmt.Errorf("error: invalid type value for amount_from, only float type")
		}

		query = query.Where("amount >= ?", amountFrom)
		return nil
	}

	amountTo, err := strconv.ParseFloat(strAmountTo, 64)
	if err != nil {
		return fmt.Errorf("error: invalid type value for amount_to, only float type")
	}

	amountFrom, err := strconv.ParseFloat(strAmountFrom, 64)
	if err != nil {
		return fmt.Errorf("error: invalid type value for amount_from, only float type")
	}

	if amountFrom >= amountTo {
		return fmt.Errorf("error: invalid value for amount_from and amount_to, amount_to must be greater than amount_from")
	}

	query = query.Where("amount >= ?", amountFrom).Where("amount <= ?", amountTo)
	return nil
}

func PrizeIdArrayToPrize(prizeIdsArray [][]int) (int, float64) {
	var ticket int
	var token float64
	for _, prizeIds := range prizeIdsArray {
		PrizeIdToPrize(prizeIds, &ticket, &token)
	}
	return ticket, token
}

var prizes = []float64{0, 0.1, 1, 0.25, 2, 0.5, 0.15, 2.5}

func PrizeIdToPrize(prizeIds []int, ticket *int, token *float64) {
	for _, prizeId := range prizeIds {
		if prizeId == 0 {
			continue
		}
		if prizeId == 2 || prizeId == 4 {
			*ticket = *ticket + int(prizes[prizeId])
			continue
		}
		*token = *token + prizes[prizeId]
	}

	// math.Round()
	*token, _ = strconv.ParseFloat(fmt.Sprintf("%.2f", *token), 2)
}
