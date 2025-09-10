package main

import (
	"github.com/pft/cmd"
	"github.com/pft/internal/database"
)


func main() {
	database.InitDB()
	defer database.DB.Close()

	cmd.Execute()
}