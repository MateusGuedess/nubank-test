# Capital Gains Tax Calculator
This is a simple calculator to calculate the capital gains tax.

## How to use

You can pass a txt file to the input of the program.

```bash
  go run main.go < example.txt    
```
The program will read one group of operations per line:
```bash
  [{"operation":"buy", "unit-cost":10.00, "quantity": 10000}, {"operation":"buy", "unit-cost":25.00, "quantity": 5000}]
```

I let one file as a model called `example.txt` in the root of the project.

## Running Test

```bash
    go test --cover # cover flag to show the coverage
```
# nubank-test
