package main

import (
	"encoding/csv"
	"fmt"
	"io"
	"os"

	"github.com/cockroachdb/apd/v2"

	"knapsack"
)

func main()  {
	transactions, err := parseTransactions("data/transactions.csv")
	if err != nil {
		panic(err)
	}

	costs := map[string]int{
		"ae":80,
		"ar":87,
		"au":250,
		"be":46,
		"bh":82,
		"br":37,
		"ca":12,
		"ch":55,
		"cl":83,
		"cn":115,
		"cy":77,
		"de":48,
		"es":56,
		"fi":50,
		"fj":360,
		"fr":53,
		"gi":61,
		"gr":66,
		"hk":130,
		"id":227,
		"ie":42,
		"il":79,
		"it":62,
		"jp":122,
		"ky":30,
		"ma":88,
		"mx":14,
		"ng":102,
		"nl":47,
		"no":46,
		"nz":350,
		"pl":49,
		"ro":51,
		"ru":55,
		"sa":78,
		"se":47,
		"sg":130,
		"th":133,
		"tr":99,
		"ua":52,
		"uk":45,
		"us":10,
		"vn":129,
		"za":105,
	}

	const totalTimeLimitMs = 1000

	result, err := knapsack.Prioritize(transactions, totalTimeLimitMs, costs)
	if err != nil {
		panic(err)
	}

	total, err := sumTransactions(result)
	if err != nil {
		panic(err)
	}

	fmt.Println(fmt.Sprintf("time limit ms: %d", totalTimeLimitMs))
	fmt.Println(total)
}

func sumTransactions(transactions []knapsack.Transaction) (*apd.Decimal, error) {
	result := apd.New(0, 0)

	for _, tx := range transactions {
		_, err := apd.BaseContext.Add(result, result, tx.Amount)
		if err != nil {
			return nil, fmt.Errorf("add to result: %w", err)
		}
	}

	return result, nil
}

func parseTransactions(file string) ([]knapsack.Transaction, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	csvr := csv.NewReader(f)

	var result []knapsack.Transaction
	for {
		row, err := csvr.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		amount, _, err := apd.NewFromString(row[1])
		if err != nil {
			return nil, fmt.Errorf("parse tx amount: %w", err)
		}
		tx := knapsack.Transaction{
			ID:              row[0],
			Amount:          amount,
			BankCountryCode: row[2],
		}

		result = append(result, tx)
	}

	return result, nil
}
