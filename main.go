package main

import (
	"invoiceinaja/config"
	"invoiceinaja/database"
)

func main() {
	conf := config.InitConfiguration()
	database.InitDatabase(conf)
	//db := database.DB

}
