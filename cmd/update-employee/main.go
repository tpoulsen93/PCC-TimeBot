package main

import (
	"fmt"

	"github.com/joho/godotenv"
	"github.com/tpoulsen/pcc-timebot/internal/admin"
)

func main() {
	_ = godotenv.Overload()

	admin.UpdateEmployee()
	fmt.Println("Employee updated successfully!")
}
