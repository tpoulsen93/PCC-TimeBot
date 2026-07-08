package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/tpoulsen/pcc-timebot/internal/admin"
)

func main() {
	_ = godotenv.Overload()
	if err := admin.AddEmployee(); err != nil {
		log.Fatalf("Error adding employee: %v", err)
		os.Exit(1)
	}
	fmt.Println("Employee added successfully!")
}
