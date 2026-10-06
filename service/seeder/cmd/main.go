package main

import (
	"fmt"
	"os"
	"os/exec"
)

var domains = []string{
	"user",
	"role",
	"merchant",
	"merchant_detail",
	"merchant_award",
	"merchant_policy",
	"merchant_business",
	"category",
	"product",
	"order",
	"shipping_address",
	"review",
	"transaction",
	"slider",
	"banner",
}

func main() {
	for _, d := range domains {
		fmt.Printf("== seeding %s ==\n", d)
		cmd := exec.Command("go", "run", "./service/"+d+"/cmd/seeder")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, "seeder %s failed: %v\n", d, err)
			os.Exit(1)
		}
	}
	fmt.Println("All domains seeded successfully")
}
