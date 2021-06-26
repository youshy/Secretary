package main

import (
	"log"
	"os"
)

// setKeys grabs every key defined by the call to the function.
// If the key exists, it gets the value of the key and saves it
// with the provider.
func setKeys(provider Provider, keys ...string) error {
	for _, key := range keys {
		value := os.Getenv(key)
		if value == "" {
			log.Printf("Value for key %s does not exist.\n", key)
			continue
		}

		err := provider.PutSecret(key, value)
		if err != nil {
			log.Printf("Error putting value: %v\n", err)
			return err
		}
	}

	return nil
}
