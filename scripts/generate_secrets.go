package main

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
)

// Simple script to generate secure random secrets for JWT
func main() {
	fmt.Println("Generating secure JWT secrets...\n")

	accessSecret, err := generateSecret(32)
	if err != nil {
		log.Fatalf("Failed to generate access secret: %v", err)
	}

	refreshSecret, err := generateSecret(32)
	if err != nil {
		log.Fatalf("Failed to generate refresh secret: %v", err)
	}

	fmt.Println("Add these to your .env file:")
	fmt.Println("================================")
	fmt.Printf("JWT_ACCESS_SECRET=%s\n", accessSecret)
	fmt.Printf("JWT_REFRESH_SECRET=%s\n", refreshSecret)
	fmt.Println("================================")
}

func generateSecret(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}
