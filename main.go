package main

import (
	"luketodd/dorsal/api"
	"luketodd/dorsal/database"
)

func main() {
	// migrations.Migrate()
	// migrations.MigrateTransactions()
	database.InitDatabase()
	api.StartApi()
}
