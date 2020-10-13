package knapsack

import (
	"testing"

	"github.com/cockroachdb/apd/v2"
	"github.com/stretchr/testify/assert"
)

func Test_prioritize(t *testing.T) {
	costs := map[string]int{
		"code1":1,
		"code2":2,
		"code3":3,
	}

	tx1 := Transaction{
		ID: "1",
		Amount: apd.New(1,0),
		BankCountryCode:"code1",
	}
	tx2 := Transaction{
		ID: "2",
		Amount: apd.New(2,0),
		BankCountryCode:"code2",
	}
	tx3 := Transaction{
		ID: "3",
		Amount: apd.New(3,0),
		BankCountryCode:"code3",
	}

	transactions := []Transaction{
		tx1, tx2, tx3,
	}

	t.Run("1", func(t *testing.T) {
		got, err := prioritize(transactions, 10, costs)

		assert.NoError(t, err)
		want := []Transaction{
			tx3, tx2, tx1,
		}
		assert.Equal(t, want, got)
	})
	t.Run("2", func(t *testing.T) {
		got, err := prioritize(transactions, 6, costs)

		assert.NoError(t, err)
		want := []Transaction{
			tx3, tx2, tx1,
		}
		assert.Equal(t, want, got)
	})
	t.Run("3", func(t *testing.T) {
		got, err := prioritize(transactions, 5, costs)

		assert.NoError(t, err)
		want := []Transaction{
			tx3, tx2,
		}
		assert.Equal(t, want, got)
	})
	t.Run("4", func(t *testing.T) {
		got, err := prioritize(transactions, 3, costs)

		assert.NoError(t, err)
		want := []Transaction{
			tx2, tx1,
		}
		assert.Equal(t, want, got)
	})
	t.Run("5", func(t *testing.T) {
		got, err := prioritize(transactions, 1, costs)

		assert.NoError(t, err)
		want := []Transaction{
			tx1,
		}
		assert.Equal(t, want, got)
	})
	t.Run("6", func(t *testing.T) {
		got, err := prioritize(transactions, 0, costs)

		assert.NoError(t, err)
		var want []Transaction
		assert.Equal(t, want, got)
	})
}
