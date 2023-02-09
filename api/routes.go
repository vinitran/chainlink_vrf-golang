package api

func (gin *GinEngine) SetupRoutes() {
	client := gin.g.Group("/api")
	{
		client.GET("/randoms/request", GetRequestRandom)
		client.GET("/randoms/response", GetResponseRandom)
		client.GET("/randoms/request/:request_id", GetRequestRandomById)
		client.GET("/randoms/response/:request_id", GetResponseRandomById)
		client.GET("/spinning/total", GetTotalSpinning)
		client.GET("/spinning/total/:address", GetSpinningCountByAddress)
		client.GET("/spinning/prize", GetSpinningPrize)
		client.GET("/spinning/prize/:request_id", GetSpinningPrizeById)
		client.GET("/spinning/prize/total", GetSpinningTotalPrize)
		client.GET("/spinning/prize/total/:address", GetSpinningTotalPrizeByAddress)
	}
	//select wallet_address, array_agg(prize_ids) from response_random where wallet_address = '0xAdfD8DAa41c23c18064074416d3428a3086e1621' group by wallet_address;

}
