package main

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestCalculateWeightedAverage(t *testing.T) {
	tests := []struct {
		name                    string
		currentTotalStockMarket int
		currentWeightedAverage  float32
		newUnitCost             float32
		newQuantity             int
		expected                float32
	}{
		{
			name:                    "Initial buy",
			currentTotalStockMarket: 0,
			currentWeightedAverage:  0,
			newUnitCost:             10,
			newQuantity:             100,
			expected:                10,
		},
		{
			name:                    "Additional buy with new price",
			currentTotalStockMarket: 100,
			currentWeightedAverage:  10,
			newUnitCost:             20,
			newQuantity:             100,
			expected:                15,
		},
		{
			name:                    "Additional buy with 3x for each parameter",
			currentTotalStockMarket: 300,
			currentWeightedAverage:  30,
			newUnitCost:             60,
			newQuantity:             300,
			expected:                45,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := calculateWeightedAverage(test.currentTotalStockMarket, test.currentWeightedAverage, test.newUnitCost, test.newQuantity)
			if result != test.expected {
				t.Errorf("Expected %.2f, got %.2f", test.expected, result)
			}
		})
	}
}

func TestCalculateTaxes(t *testing.T) {
	tests := []struct {
		name     string
		input    [][]Operation
		expected [][]Tax
	}{
		{
			name: "Buy one and sell two",
			input: [][]Operation{
				{
					{OperationType: "buy", UnitCost: 10.00, Quantity: 100},
					{OperationType: "sell", UnitCost: 15.00, Quantity: 50},
					{OperationType: "sell", UnitCost: 15.00, Quantity: 50},
				},
			},
			expected: [][]Tax{
				{
					{Tax: 0.00},
					{Tax: 0.00},
					{Tax: 0.00},
				},
			},
		},
		{
			name: "Buy one and sell two",
			input: [][]Operation{
				{
					{OperationType: "buy", UnitCost: 10.00, Quantity: 10000},
					{OperationType: "sell", UnitCost: 20.00, Quantity: 5000},
					{OperationType: "sell", UnitCost: 5.00, Quantity: 5000},
				},
			},
			expected: [][]Tax{
				{
					{Tax: 0.00},
					{Tax: 10000.00},
					{Tax: 0.00},
				},
			},
		},
		{
			name: "The two previous cases together",
			input: [][]Operation{
				{
					{OperationType: "buy", UnitCost: 10.00, Quantity: 100},
					{OperationType: "sell", UnitCost: 15.00, Quantity: 50},
					{OperationType: "sell", UnitCost: 15.00, Quantity: 50},
				},
				{
					{OperationType: "buy", UnitCost: 10.00, Quantity: 10000},
					{OperationType: "sell", UnitCost: 20.00, Quantity: 5000},
					{OperationType: "sell", UnitCost: 5.00, Quantity: 5000},
				},
			},
			expected: [][]Tax{
				{
					{Tax: 0.00},
					{Tax: 0.00},
					{Tax: 0.00},
				},
				{
					{Tax: 0.00},
					{Tax: 10000.00},
					{Tax: 0.00},
				},
			},
		},
		{
			name: "buy one sell two",
			input: [][]Operation{
				{
					{OperationType: "buy", UnitCost: 10.00, Quantity: 10000},
					{OperationType: "sell", UnitCost: 5.00, Quantity: 5000},
					{OperationType: "sell", UnitCost: 20.00, Quantity: 3000},
				},
			},
			expected: [][]Tax{
				{
					{Tax: 0.00},
					{Tax: 0.00},
					{Tax: 1000.00},
				},
			},
		},
		{
			name: "buy two times before the sells",
			input: [][]Operation{
				{
					{OperationType: "buy", UnitCost: 10.00, Quantity: 10000},
					{OperationType: "buy", UnitCost: 25.00, Quantity: 5000},
					{OperationType: "sell", UnitCost: 15.00, Quantity: 10000},
					{OperationType: "sell", UnitCost: 25.00, Quantity: 5000},
				},
			},
			expected: [][]Tax{
				{
					{Tax: 0.00},
					{Tax: 0.00},
					{Tax: 0.00},
					{Tax: 10000.00},
				},
			},
		},
		{
			name: "Buy, sell, buy, sell",
			input: [][]Operation{
				{
					{OperationType: "buy", UnitCost: 10.00, Quantity: 10000},
					{OperationType: "sell", UnitCost: 2.00, Quantity: 5000},
					{OperationType: "sell", UnitCost: 20.00, Quantity: 2000},
					{OperationType: "sell", UnitCost: 20.00, Quantity: 2000},
					{OperationType: "sell", UnitCost: 25.00, Quantity: 1000},
					{OperationType: "buy", UnitCost: 20.00, Quantity: 10000},
					{OperationType: "sell", UnitCost: 15.00, Quantity: 5000},
					{OperationType: "sell", UnitCost: 30.00, Quantity: 4350},
					{OperationType: "sell", UnitCost: 30.00, Quantity: 650},
				},
			},
			expected: [][]Tax{
				{
					{Tax: 0.00},
					{Tax: 0.00},
					{Tax: 0.00},
					{Tax: 0.00},
					{Tax: 3000.00},
					{Tax: 0.00},
					{Tax: 0.00},
					{Tax: 3700.00},
					{Tax: 0.00},
				},
			},
		},
		{
			name: "Loss and then sell with profit",
			input: [][]Operation{
				{
					{OperationType: "buy", UnitCost: 10, Quantity: 10000},
					{OperationType: "sell", UnitCost: 5, Quantity: 5000},
					{OperationType: "sell", UnitCost: 20, Quantity: 5000},
				},
			},
			expected: [][]Tax{
				{
					{Tax: 0.00},
					{Tax: 0.00},
					{Tax: 5000.00},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := calculateTaxes(test.input)
			if !reflect.DeepEqual(result, test.expected) {
				t.Errorf("Expected %v, got %v", test.expected, result)
			}
		})
	}
}

func TestReadOperations(t *testing.T) {
	mockInput := `

	[{"operation":"buy", "unit-cost":10.00, "quantity":10000}, {"operation":"sell", "unit-cost":20.00, "quantity":5000}]

	[{"operation":"buy", "unit-cost":10.00, "quantity": 10000}, {"operation":"sell", "unit-cost":50.00, "quantity": 10000}, {"operation":"buy", "unit-cost":20.00, "quantity": 10000}, {"operation":"sell", "unit-cost":50.00, "quantity": 10000}]
	
	[{"operation":"buy", "unit-cost":10.00, "quantity": 10000},{"operation":"sell", "unit-cost":2.00, "quantity": 5000},{"operation":"sell", "unit-cost":20.00, "quantity": 2000},{"operation":"sell", "unit-cost":20.00, "quantity": 2000},{"operation":"sell", "unit-cost":25.00, "quantity": 1000},{"operation":"buy", "unit-cost":20.00, "quantity": 10000},{"operation":"sell", "unit-cost":15.00, "quantity": 5000},{"operation":"sell", "unit-cost":30.00, "quantity": 4350},{"operation":"sell", "unit-cost":30.00, "quantity": 650}]
	
	[{"operation":"buy", "unit-cost":10.00, "quantity": 10000},{"operation":"sell", "unit-cost":2.00, "quantity": 5000},{"operation":"sell", "unit-cost":20.00, "quantity": 2000},{"operation":"sell", "unit-cost":20.00, "quantity": 2000},{"operation":"sell", "unit-cost":25.00, "quantity": 1000}]

	[]
`

	reader := strings.NewReader(mockInput)
	operations, err := readOperations(reader)
	if err != nil {
		t.Fatalf("readOperations failed: %v", err)
	}

	expected := [][]Operation{
		{
			{OperationType: "buy", UnitCost: 10.00, Quantity: 10000},
			{OperationType: "sell", UnitCost: 20.00, Quantity: 5000},
		},
		{
			{OperationType: "buy", UnitCost: 10.00, Quantity: 10000},
			{OperationType: "sell", UnitCost: 50.00, Quantity: 10000},
			{OperationType: "buy", UnitCost: 20.00, Quantity: 10000},
			{OperationType: "sell", UnitCost: 50.00, Quantity: 10000},
		},
		{
			{OperationType: "buy", UnitCost: 10.00, Quantity: 10000},
			{OperationType: "sell", UnitCost: 2.00, Quantity: 5000},
			{OperationType: "sell", UnitCost: 20.00, Quantity: 2000},
			{OperationType: "sell", UnitCost: 20.00, Quantity: 2000},
			{OperationType: "sell", UnitCost: 25.00, Quantity: 1000},
			{OperationType: "buy", UnitCost: 20.00, Quantity: 10000},
			{OperationType: "sell", UnitCost: 15.00, Quantity: 5000},
			{OperationType: "sell", UnitCost: 30.00, Quantity: 4350},
			{OperationType: "sell", UnitCost: 30.00, Quantity: 650},
		},
		{
			{OperationType: "buy", UnitCost: 10.00, Quantity: 10000},
			{OperationType: "sell", UnitCost: 2.00, Quantity: 5000},
			{OperationType: "sell", UnitCost: 20.00, Quantity: 2000},
			{OperationType: "sell", UnitCost: 20.00, Quantity: 2000},
			{OperationType: "sell", UnitCost: 25.00, Quantity: 1000},
		},
		{},
	}

	if !reflect.DeepEqual(operations, expected) {
		t.Errorf("Expected %v, got %v", expected, operations)
	}
}

func TestMarshalJSON(t *testing.T) {
	tax := Tax{Tax: 123.456}
	expected := `{"tax":123.46}`

	result, err := json.Marshal(tax)
	if err != nil {
		t.Fatalf("Failed to marshal tax: %v", err)
	}

	if string(result) != expected {
		t.Errorf("Expected %s, got %v", expected, string(result))
	}
}
