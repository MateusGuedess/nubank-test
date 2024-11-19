package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

type Operation struct {
	OperationType string  `json:"operation"`
	UnitCost      float32 `json:"unit-cost"`
	Quantity      int     `json:"quantity"`
}

type Tax struct {
	Tax float32 `json:"tax"`
}

type TaxJSON struct {
	Tax json.Number `json:"tax"`
}

func main() {
	var reader io.Reader
	if len(os.Args) < 2 {
		reader = os.Stdin
	}

	operations, err := readOperations(reader)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	taxes := calculateTaxes(operations)
	printTaxes(taxes)
}

func readOperations(reader io.Reader) ([][]Operation, error) {
	scanner := bufio.NewScanner(reader)
	var operations [][]Operation

	for scanner.Scan() {
		input := scanner.Text()
		if len(strings.TrimSpace(input)) == 0 {
			continue
		}

		var ops []Operation

		if err := json.Unmarshal([]byte(input), &ops); err != nil {
			fmt.Printf("Error parsing JSON on line: %s\nError: %v\n", input, err)
			continue
		}

		operations = append(operations, ops)
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("error reading operations file: %v", err)
	}

	return operations, nil
}

func calculateWeightedAverage(currentTotalStockMarket int, currentWeightedAverage, newUnitCost float32, newQuantity int) float32 {
	totalQuantity := currentTotalStockMarket + newQuantity
	if totalQuantity == 0 {
		return 0
	}
	return ((float32(currentTotalStockMarket) * currentWeightedAverage) + (float32(newQuantity) * newUnitCost)) / float32(totalQuantity)
}

func calculateTaxes(operations [][]Operation) [][]Tax {
	var taxes [][]Tax
	var totalLoss float32

	for _, operation := range operations {
		var tax []Tax
		var weightedAverage float32
		var totalStockMarket int

		for _, op := range operation {
			if op.OperationType == "buy" {
				weightedAverage = calculateWeightedAverage(totalStockMarket, weightedAverage, op.UnitCost, op.Quantity)
				totalStockMarket += op.Quantity
				tax = append(tax, Tax{0.00})
			} else if op.OperationType == "sell" {
				profitOrLoss := (op.UnitCost - weightedAverage) * float32(op.Quantity)
				totalStockMarket -= op.Quantity

				if profitOrLoss < 0 {
					totalLoss += -profitOrLoss
					tax = append(tax, Tax{0.00})
				} else {
					if totalLoss > 0 {
						if totalLoss >= profitOrLoss {
							totalLoss -= profitOrLoss
							profitOrLoss = 0
						} else {
							profitOrLoss -= totalLoss
							totalLoss = 0
						}
					}

					if op.UnitCost*float32(op.Quantity) > 20000.00 && profitOrLoss > 0 {
						taxAmount := profitOrLoss * 0.20
						tax = append(tax, Tax{taxAmount})
					} else {
						tax = append(tax, Tax{0.00})
					}
				}
			}
		}
		taxes = append(taxes, tax)
	}

	return taxes
}

func (t Tax) MarshalJSON() ([]byte, error) {
	return json.Marshal(TaxJSON{
		Tax: json.Number(fmt.Sprintf("%.2f", t.Tax)),
	})
}

func printTaxes(taxes [][]Tax) {
	for _, tax := range taxes {
		taxJSON, err := json.Marshal(tax)
		if err != nil {
			fmt.Println("Error marshalling tax:", err)
			return
		}
		fmt.Println(string(taxJSON))
	}
}
