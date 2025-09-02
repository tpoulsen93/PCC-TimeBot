package main

import (
	"fmt"
	"log"
	"os"

	"github.com/tpoulsen/pcc-timebot/internal/admin"
)

func main() {
	if err := admin.AddTime(); err != nil {
		log.Fatalf("Error adding time: %v", err)
		os.Exit(1)
	}
	fmt.Println("Time added successfully!")
}
