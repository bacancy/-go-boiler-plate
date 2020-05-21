package main

import (
	"bacancy/go-boiler-plate/app/common"
	"bacancy/go-boiler-plate/app/router"
)

func main() {
	common.ConnectToDatabase()
	router.ConfigureRouter()
	router.CreateRouter()
	router.RunRouter()
}
