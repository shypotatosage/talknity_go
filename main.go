package main

import (
	"talknity/db"
	"talknity/routes"
)

func main() {
	db.Init()
	e := routes.Init()
	e.Logger.Fatal(e.Start(":9090"))
}