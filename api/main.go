package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/uptrace/bun"
	"log"
	"os"
)

var db *bun.DB

type GinEngine struct {
	g *gin.Engine
}

func NewGin(database *bun.DB) *GinEngine {
	db = database
	return &GinEngine{g: gin.New()}
}

func (gin *GinEngine) Run() {
	gin.SetupRoutes()
	err := gin.g.Run(fmt.Sprintf(":%s", os.Getenv("PORT_SV")))
	if err != nil {
		log.Fatal(err)
	}
}
