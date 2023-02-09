package api

import (
	"VRFChainlink/database"
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/uptrace/bun"
	"net/http"
	"time"
)

func GetRequestRandomById(c *gin.Context) {
	responseData := new(database.RequestRandom)
	requestId := c.Param("request_id")
	err := db.NewSelect().Model(responseData).
		Where("request_id = ?", requestId).
		Scan(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, render.JSON{Data: responseData})
	return
}

func GetResponseRandomById(c *gin.Context) {
	responseData := new(database.ResponseRandom)
	requestId := c.Param("request_id")
	err := db.NewSelect().Model(responseData).
		Where("request_id = ?", requestId).
		Scan(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, render.JSON{Data: responseData})
	return
}

func GetRequestRandom(c *gin.Context) {
	responseData := new([]database.RequestRandom)

	pageFilter := new(PageFilter)
	err := pageFilter.Check(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	offset := (pageFilter.Page - 1) * pageFilter.Size
	query := db.NewSelect().Model(responseData).
		Limit(pageFilter.Size).
		Offset(offset)

	SearchByWalletAddress(c, query)
	SearchByTxHash(c, query)

	err = SortByAmount(c, query)
	if err != nil {
		c.JSON(http.StatusUnsupportedMediaType, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	err = SearchByAmount(c, query)
	if err != nil {
		c.JSON(http.StatusUnsupportedMediaType, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	err = SearchByTime(c, query)
	if err != nil {
		c.JSON(http.StatusUnsupportedMediaType, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	err = query.Scan(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, render.JSON{Data: responseData})
	return
}

func GetResponseRandom(c *gin.Context) {
	responseData := new([]database.ResponseRandom)

	pageFilter := new(PageFilter)
	err := pageFilter.Check(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	offset := (pageFilter.Page - 1) * pageFilter.Size
	query := db.NewSelect().Model(responseData).
		Limit(pageFilter.Size).
		Offset(offset)

	SearchByWalletAddress(c, query)
	SearchByTxHash(c, query)

	err = SearchByTime(c, query)
	if err != nil {
		c.JSON(http.StatusUnsupportedMediaType, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	err = query.Scan(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, render.JSON{Data: responseData})
	return
}

func GetTotalSpinning(c *gin.Context) {
	pageFilter := new(PageFilter)
	err := pageFilter.Check(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	type Spinning struct {
		WalletAddress string `json:"wallet_address"`
		TotalAmount   int    `json:"total_amount"`
	}

	var spin []Spinning

	offset := (pageFilter.Page - 1) * pageFilter.Size
	query := db.NewSelect().Model(new([]database.RequestRandom)).
		ColumnExpr("sum(?) as total_amount", bun.Ident("req.amount")).
		ColumnExpr("wallet_address").
		GroupExpr("wallet_address").
		Limit(pageFilter.Size).
		Offset(offset)

	err = SortByTotalAmount(c, query)
	if err != nil {
		c.JSON(http.StatusUnsupportedMediaType, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	err = SearchByTotalAmount(c, query)
	if err != nil {
		c.JSON(http.StatusUnsupportedMediaType, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	err = SearchByTime(c, query)
	if err != nil {
		c.JSON(http.StatusUnsupportedMediaType, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	err = query.Scan(context.Background(), &spin)
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, render.JSON{Data: spin})
	return
}

func GetSpinningCountByAddress(c *gin.Context) {
	responseData := new(database.RequestRandom)
	amount := new(int)
	address := c.Param("address")

	query := db.NewSelect().Model(responseData).
		ColumnExpr("sum(amount)").
		Where("wallet_address = ?", address)

	err := SearchByTime(c, query)
	if err != nil {
		c.JSON(http.StatusUnsupportedMediaType, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	err = query.Scan(context.Background(), amount)
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	c.JSON(http.StatusOK, render.JSON{Data: amount})
	return
}

func GetTransactionByHash(c *gin.Context) {
	requestRandom, err := GetRequestTransactionByHash(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	responseRandom, err := GetResponseTransactionByHash(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	if requestRandom == nil && responseRandom == nil {
		c.JSON(http.StatusOK, render.JSON{Data: nil})
		return
	}

	if requestRandom == nil {
		c.JSON(http.StatusOK, render.JSON{Data: responseRandom})
		return
	}

	if responseRandom == nil {
		c.JSON(http.StatusOK, render.JSON{Data: requestRandom})
		return
	}

	type Random struct {
		RequestRandom  []database.RequestRandom  `json:"request_random"`
		ResponseRandom []database.ResponseRandom `json:"response_random"`
	}

	c.JSON(http.StatusOK, render.JSON{Data: Random{
		RequestRandom:  requestRandom,
		ResponseRandom: responseRandom,
	}})
	return
}

func GetSpinningTotalPrizeByAddress(c *gin.Context) {
	address := c.Param("address")
	query := db.NewSelect().Model(new(database.ResponseRandom)).
		ColumnExpr("json_agg(prize_ids)").
		Where("wallet_address = ?", address).
		GroupExpr("wallet_address")

	err := SearchByTime(c, query)
	if err != nil {
		c.JSON(http.StatusUnsupportedMediaType, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	var prize [][][]int
	err = query.Scan(context.Background(), &prize)
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	ticket, token := PrizeIdArrayToPrize(prize[0])

	type Spinning struct {
		Address string  `json:"address"`
		Ticket  int     `json:"ticket"`
		Token   float64 `json:"token"`
	}

	c.JSON(http.StatusOK, render.JSON{Data: Spinning{
		Address: address,
		Ticket:  ticket,
		Token:   token,
	}})
	return
}

func GetSpinningTotalPrize(c *gin.Context) {
	type Prize struct {
		WalletAddress string `json:"wallet_address"`
		Prize         [][]int
	}

	pageFilter := new(PageFilter)
	err := pageFilter.Check(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	offset := (pageFilter.Page - 1) * pageFilter.Size
	query := db.NewSelect().Model(new(database.ResponseRandom)).
		ColumnExpr("json_agg(prize_ids) as prize").
		ColumnExpr("wallet_address").
		GroupExpr("wallet_address").
		Limit(pageFilter.Size).
		Offset(offset)

	SearchByTxHash(c, query)

	err = SearchByTime(c, query)
	if err != nil {
		c.JSON(http.StatusUnsupportedMediaType, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	var prizes []Prize
	err = query.Scan(context.Background(), &prizes)
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	type TotalPrize struct {
		Address string  `json:"address"`
		Ticket  int     `json:"ticket"`
		Token   float64 `json:"token"`
	}
	var totalPrize []TotalPrize
	for _, prize := range prizes {
		ticket, token := PrizeIdArrayToPrize(prize.Prize)
		totalPrize = append(totalPrize, TotalPrize{
			Address: prize.WalletAddress,
			Ticket:  ticket,
			Token:   token,
		})
	}

	c.JSON(http.StatusOK, render.JSON{Data: totalPrize})
	return
}

func GetSpinningPrize(c *gin.Context) {
	var data []database.ResponseRandom

	pageFilter := new(PageFilter)
	err := pageFilter.Check(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	offset := (pageFilter.Page - 1) * pageFilter.Size
	query := db.NewSelect().Model(&data).
		Limit(pageFilter.Size).
		Offset(offset)

	SearchByWalletAddress(c, query)
	SearchByTxHash(c, query)

	err = SearchByTime(c, query)
	if err != nil {
		c.JSON(http.StatusUnsupportedMediaType, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	err = query.Scan(context.Background())
	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	type ResponseDetail struct {
		WalletAddress   string    `json:"wallet_address"`
		RequestId       string    `json:"request_id"`
		TransactionHash string    `json:"transaction_hash"`
		Index           int       `json:"index"`
		Time            time.Time `json:"time"`
		Ticket          int       `json:"ticket"`
		Token           float64   `json:"token"`
	}
	var responseData []ResponseDetail
	for _, dt := range data {
		var (
			ticket int
			token  float64
		)
		PrizeIdToPrize(dt.PrizeIds, &ticket, &token)
		responseData = append(responseData, ResponseDetail{
			WalletAddress:   dt.User,
			RequestId:       dt.RequestId,
			TransactionHash: dt.TxHash,
			Index:           dt.Index,
			Time:            dt.Time,
			Ticket:          ticket,
			Token:           token,
		})
	}

	c.JSON(http.StatusOK, render.JSON{Data: responseData})
	return
}

func GetSpinningPrizeById(c *gin.Context) {
	data := new(database.ResponseRandom)
	requestId := c.Param("request_id")

	err := db.NewSelect().Model(data).
		Where("request_id = ?", requestId).
		Scan(context.Background())

	if err != nil {
		c.JSON(http.StatusBadRequest, render.JSON{Data: fmt.Sprintf("%s", err)})
		fmt.Println(err)
		return
	}

	type ResponseDetail struct {
		WalletAddress   string    `json:"wallet_address"`
		RequestId       string    `json:"request_id"`
		TransactionHash string    `json:"transaction_hash"`
		Index           int       `json:"index"`
		Time            time.Time `json:"time"`
		Ticket          int       `json:"ticket"`
		Token           float64   `json:"token"`
	}

	var (
		ticket int
		token  float64
	)
	PrizeIdToPrize(data.PrizeIds, &ticket, &token)

	c.JSON(http.StatusOK, render.JSON{Data: ResponseDetail{
		WalletAddress:   data.User,
		RequestId:       data.RequestId,
		TransactionHash: data.TxHash,
		Index:           data.Index,
		Time:            data.Time,
		Ticket:          ticket,
		Token:           token,
	}})

	return
}
