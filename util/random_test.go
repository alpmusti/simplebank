package util

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestRandomInt(t *testing.T) {
	var min int64 = 0
	var max int64 = 100
	rand1 := RandomInt(min, max)
	require.GreaterOrEqual(t, rand1, min)
	require.LessOrEqual(t, rand1, max)
}

func TestRandomString(t *testing.T) {
	length := 12
	rand1 := RandomString(length)
	require.Len(t, rand1, length)
}

func TestRandomOwner(t *testing.T) {
	rand1 := RandomOwner()
	require.Len(t, rand1, 6)
}

func TestRandomMoney(t *testing.T) {
	var min int64 = 0
	var max int64 = 1000
	rand1 := RandomMoney()
	require.GreaterOrEqual(t, rand1, min)
	require.LessOrEqual(t, rand1, max)
}

func TestRandomCurrency(t *testing.T) {
	currencies := getAvailableCurrencies()
	rand1 := RandomCurrency()
	require.Contains(t, currencies, rand1)
}
