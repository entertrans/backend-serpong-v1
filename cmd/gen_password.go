// gen_password.go
package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

func main() {
	// Password yang benar (NIP)
	passwords := []string{
		"198501012010011001",
		"198702152011012002",
	}

	for _, pwd := range passwords {
		hash, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		fmt.Printf("NIP: %s\n", pwd)
		fmt.Printf("HASH: '%s'\n\n", string(hash))
	}
}
