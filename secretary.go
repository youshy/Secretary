package main

type Provider interface {
	GetSecret(key string) (string, error)
	PutSecret(key, value string) error
}
