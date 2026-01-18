package main

import (
	"flag"
	"fmt"
	"log"
	"lunchorder/utils"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	action := flag.String("action", "", "encrypt, decrypt, or hash")
	input := flag.String("input", "", "text to process")
	flag.Parse()

	if *action == "" || *input == "" {
		fmt.Println("Usage: go run tools/crypto_tool.go -action=[encrypt|decrypt|hash] -input=\"text\"")
		os.Exit(1)
	}

	// Load .env from project root
	_ = godotenv.Load()

	key, err := utils.GetEncryptionKey()
	if err != nil {
		log.Fatal(err)
	}

	switch *action {
	case "encrypt":
		res, err := utils.Encrypt(*input, key)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Encrypted: %s\n", res)
	case "decrypt":
		res, err := utils.Decrypt(*input, key)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("Decrypted: %s\n", res)
	case "hash":
		res := utils.Hash(*input, key)
		fmt.Printf("Hash: %s\n", res)
	default:
		log.Fatal("Unknown action")
	}
}
