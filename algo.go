package knapsack

import (
	"fmt"

	"github.com/cockroachdb/apd/v2"
)

type Transaction struct {
	// a UUID of transaction
	ID string
	// in USD, typically a value betwen 0.01 and 1000 USD.
	Amount *apd.Decimal
	// a 2-letter country code of where the bank is located
	BankCountryCode string
}

// Prioritize returns a list of transactions (sub-list of transactions param) that can be served within defined costLimit.
//
// transactions is a list of transactions to be prioritized.
// Each transaction in transactions list has a certain cost.
// costLimit defines the limit the resulting transactions list (summing up all the transactions om this list) must not exceed.
// costs defines a mapping of bank country code -> transaction cost that prioritize relies upon when choosing between transactions.
//
// This function is based on a solution for a knapsack problem - https://en.wikipedia.org/wiki/Knapsack_problem.
func Prioritize(transactions []Transaction, costLimit int, costs map[string]int) ([]Transaction, error) {
	n := len(transactions)
	totalCapacity := costLimit

	// solutions is a [n+1][totalCapacity+1] matrix of already discovered solutions.
	// The row 0 and the column 0 in this matrix do not have any real meaning behind them, and are used to bootstrap the algorithm.
	solutions := make([][]*apd.Decimal, n+1)
	for i := range solutions {
		solutions[i] = make([]*apd.Decimal, totalCapacity+1)
		for j := range solutions[i] {
			solutions[i][j] = apd.New(0, 0)
		}
	}

	for tx := 1; tx <= n; tx++ {
		for capacity := 1; capacity <= totalCapacity; capacity++ {
			maxValWithoutCurr := apd.New(0, 0).Set(solutions[tx-1][capacity])
			maxValWithCurr := apd.New(0, 0)

			weightOfCurr := costs[transactions[tx-1].BankCountryCode]
			if capacity >= weightOfCurr {
				maxValWithCurr := maxValWithCurr.Set(transactions[tx-1].Amount)

				remainingCapacity := capacity - weightOfCurr
				_, err := apd.BaseContext.Add(maxValWithCurr, maxValWithCurr, solutions[tx-1][remainingCapacity])
				if err != nil {
					return nil, fmt.Errorf("add to maxValWithCurr: %w", err)
				}
			}

			solutions[tx][capacity] = max(maxValWithoutCurr, maxValWithCurr)
		}
	}

	// Build a list of transactions to be included in the result by iterating over solutions matrix in reverse to the algo used above.
	var result []Transaction
	for tx, capacity := n, totalCapacity; tx > 0 && capacity > 0; tx-- {
		if solutions[tx][capacity].Cmp(solutions[tx-1][capacity]) != 0 {
			// txIdx is an index of this tx in the transactions list.
			txIdx := tx - 1

			cost := costs[transactions[txIdx].BankCountryCode]
			capacity -= cost
			result = append(result, transactions[txIdx])
		}
	}

	return result, nil
}

func max(left *apd.Decimal, right *apd.Decimal) *apd.Decimal {
	switch left.Cmp(right) {
	case 1:
		return left
	default:
		return right
	}
}
