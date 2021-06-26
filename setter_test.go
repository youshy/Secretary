package main

import (
	"errors"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

var (
	key1   = "KEY_ONE"
	value1 = "VALUE_ONE"

	key2   = "KEY_TWO"
	value2 = "VALUE_TWO"

	key3   = "KEY_THREE"
	value3 = "VALUE_THREE"

	key4   = "KEY_FOUR"
	value4 = "VALUE_FOUR"

	keyNoValue = "KEY_NO_VALUE"

	faultyConnectionError = errors.New("Faulty connection")
)

type mockProvider struct {
	Output map[string]string
}

func (m *mockProvider) GetSecret(key string) (string, error) {
	return "", nil
}

func (m *mockProvider) PutSecret(key, value string) error {
	m.Output[key] = value
	return nil
}

type mockProviderWithLimit struct {
	Output   map[string]string
	Iterator int
}

func (m *mockProviderWithLimit) GetSecret(key string) (string, error) {
	return "", nil
}

func (m *mockProviderWithLimit) PutSecret(key, value string) error {
	if m.Iterator == 2 {
		return faultyConnectionError
	}
	m.Output[key] = value
	m.Iterator++
	return nil
}

type mockProviderWithFault struct {
}

func (m *mockProviderWithFault) GetSecret(key string) (string, error) {
	return "", nil
}

func (m *mockProviderWithFault) PutSecret(key, value string) error {
	return faultyConnectionError
}

func TestSetKeys(t *testing.T) {
	os.Setenv(key1, value1)
	defer os.Unsetenv(key1)

	os.Setenv(key2, value2)
	defer os.Unsetenv(key2)

	os.Setenv(key3, value3)
	defer os.Unsetenv(key3)

	os.Setenv(key4, value4)
	defer os.Unsetenv(key4)

	t.Run("SetKeys puts secrets with value and not throw error", func(t *testing.T) {
		outputMap := make(map[string]string, 0)
		outputMap[key1] = value1
		outputMap[key2] = value2
		outputMap[key3] = value3
		outputMap[key4] = value4

		mock := &mockProvider{}
		mock.Output = make(map[string]string, 0)
		err := setKeys(mock, key1, key2, key3, key4)

		assert.NoError(t, err)
		assert.Equal(t, outputMap, mock.Output)
	})

	t.Run("SetKeys puts secret for 3 keys with value, doesn't save one without value", func(t *testing.T) {
		outputMap := make(map[string]string, 0)
		outputMap[key1] = value1
		outputMap[key2] = value2
		outputMap[key3] = value3

		mock := &mockProvider{}
		mock.Output = make(map[string]string, 0)
		err := setKeys(mock, key1, key2, key3, keyNoValue)

		assert.NoError(t, err)
		assert.Equal(t, outputMap, mock.Output)
	})

	t.Run("SetKeys saves first two keys, throws an error with third one after provider error", func(t *testing.T) {
		outputMap := make(map[string]string, 0)
		outputMap[key1] = value1
		outputMap[key2] = value2

		mock := &mockProviderWithLimit{}
		mock.Output = make(map[string]string, 0)

		err := setKeys(mock, key1, key2, key3, key4)

		assert.ErrorIs(t, err, faultyConnectionError)
		assert.Equal(t, mock.Iterator, 2)
		assert.Equal(t, outputMap, mock.Output)
	})

	t.Run("SetKeys doesn't save any key after first one throws an error", func(t *testing.T) {
		mock := &mockProviderWithFault{}

		err := setKeys(mock, key1, key2)

		assert.ErrorIs(t, err, faultyConnectionError)
	})
}
