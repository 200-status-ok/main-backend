package main

import (
	"flag"
	"fmt"
	"github.com/200-status-ok/main-backend/src/MainService/Cmd/Migrate"
)

func main() {
	migrateFlag := flag.Bool("migrate", false, "Run database migration")
	dropFlag := flag.Bool("drop", false, "Drop database")
	esSetupFlag := flag.Bool("esSetup", false, "Setup elastic search")
	insertAllPostersInESFlag := flag.Bool("bulk", false, "Insert all posters in elastic search")

	flag.Parse()

	if *migrateFlag {
		Migrate.ModelsMigrate()
		fmt.Println("Migrate done")
	} else if *dropFlag {
		Migrate.DropModels()
		fmt.Println("Drop done")
	} else if *esSetupFlag {
		Migrate.ESSetup()
		fmt.Println("ES setup done")
	} else if *insertAllPostersInESFlag {
		Migrate.InsertAllPostersInES()
		fmt.Println("Insert all posters in elastic search done")
	} else {
		fmt.Println("No flags")
	}
}
