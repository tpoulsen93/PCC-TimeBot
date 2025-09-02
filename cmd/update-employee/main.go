package main

import (
	"fmt"

	"github.com/tpoulsen/pcc-timebot/internal/admin"
)

func main() {
	admin.UpdateEmployee()
	fmt.Println("Employee updated successfully!")
}
