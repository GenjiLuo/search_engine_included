package main

import (
	"app/databases"
	"app/routers"
)

func main() {
	databases.AutoMigrate()

	r := routers.Load()
	r.Run()
}
