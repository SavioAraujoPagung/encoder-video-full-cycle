package utils_test

import (
	"encoder/framework/utils"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIsJson(t *testing.T) {
	json := `{
		"id": "123",
		"nome": "teste",
		"idade": "20"
	}`

	err := utils.IsJson(json)
	require.Nil(t, err)
}

func TestNotJson(t *testing.T) {
	json := `{
		"id": "123",
		"nome": "teste",
		"idade": "20",
	}`

	err := utils.IsJson(json)
	require.Error(t, err)
}