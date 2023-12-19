package helper

import "math/rand"

func GenerateRandomID() string {
	charSet := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	// Initialize an empty string to store the result
	result := make([]byte, 10)

	// Generate random characters
	for i := 0; i < 10; i++ {
		result[i] = charSet[rand.Intn(len(charSet))]
	}

	return string(result)
}
