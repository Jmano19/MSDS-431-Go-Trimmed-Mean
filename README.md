# Go Trimmed Mean

A reusable Go package for calculating symmetric and asymmetric trimmed means from integer and floating-point datasets.

Developed for Northwestern University's MSDS 431 programming assignment.

**GitHub Repository:** https://github.com/Jmano19/MSDS-431-Go-Trimmed-Mean

## Overview

The `trimmedmean` package calculates the arithmetic mean after removing a specified proportion of observations from the lower and upper ends of a sorted dataset.

A trimmed mean is a measure of central tendency that reduces the influence of extreme observations and outliers.

The package supports:

- Integer and floating-point observations.
- Symmetric trimming, with equal proportions removed from both tails.
- Asymmetric trimming, with different proportions removed from each tail.
- Ordinary arithmetic means when no trimming is specified.
- Generic numeric types using Go generics.
- Input validation and descriptive errors.
- Preservation of the original input slice.

## Project Structure

```text
MSDS 431 Go Trimmed Mean/
├── go.mod
├── trimmedmean.go
├── trimmedmean_test.go
├── benchmark_test.go
└── README.md
```

## Requirements

- Go 1.18 or later (generics support required).
- Git for downloading the package.

## Installation

Install the package using Go modules:

```bash
go get github.com/Jmano19/MSDS-431-Go-Trimmed-Mean
```

Import the package into a Go program:

```go
import trimmedmean "github.com/Jmano19/MSDS-431-Go-Trimmed-Mean"
```

## Function Reference

### Mean

```go
func Mean[T Number](data []T, trim ...float64) (float64, error)
```

The `Mean` function accepts a numeric slice and zero, one, or two trimming proportions.

| Arguments | Behavior |
|---|---|
| `Mean(data)` | Ordinary arithmetic mean |
| `Mean(data, 0.05)` | Remove 5% from each tail |
| `Mean(data, 0.10)` | Remove 10% from each tail |
| `Mean(data, 0.05, 0.10)` | Remove 5% from the lower tail and 10% from the upper tail |

The first trimming argument represents the lower-tail proportion. The second represents the upper-tail proportion.

When only one argument is supplied, the same proportion is applied to both tails.

### Number Type

The package defines a generic `Number` constraint supporting:

- Signed integers: `int`, `int8`, `int16`, `int32`, `int64`
- Unsigned integers: `uint`, `uint8`, `uint16`, `uint32`, `uint64`
- Floating-point values: `float32`, `float64`

## Usage Examples

### Example 1: Integer observations

```go
package main

import (
    "fmt"
    trimmedmean "github.com/Jmano19/MSDS-431-Go-Trimmed-Mean"
)

func main() {
    data := []int{
        1, 2, 3, 4, 5,
        6, 7, 8, 9, 10,
    }

    result, err := trimmedmean.Mean(data, 0.20)
    if err != nil {
        panic(err)
    }

    fmt.Printf("Trimmed mean: %.2f\n", result)
}
```

Expected output:

```text
Trimmed mean: 5.50
```

### Example 2: Floating-point observations

```go
data := []float64{
    1.5, 2.5, 3.5, 4.5, 5.5,
}

result, err := trimmedmean.Mean(data, 0.20)
if err != nil {
    panic(err)
}

fmt.Printf("Trimmed mean: %.2f\n", result)
```

### Example 3: Asymmetric trimming

```go
data := []float64{
    1, 2, 3, 4, 5,
    6, 7, 8, 9, 10,
}

result, err := trimmedmean.Mean(data, 0.10, 0.20)
if err != nil {
    panic(err)
}

fmt.Printf("Trimmed mean: %.2f\n", result)
```

Expected output:

```text
Trimmed mean: 5.00
```

This removes one observation from the lower tail and two observations from the upper tail.

## Calculation Method

The package follows these steps:

1. Validate the input dataset and trimming proportions.
2. Convert observations to `float64`.
3. Sort a copy of the observations in ascending order.
4. Calculate the number of observations removed from each tail using the floor of the sample size multiplied by the respective trimming proportion.
5. Remove the specified observations from both ends.
6. Calculate the arithmetic mean of the remaining observations.

The original input slice is not modified.

For a sample of size \(n\), lower trimming proportion \(p_L\), and upper trimming proportion \(p_U\):

\[
k_L = \lfloor n p_L \rfloor
\]

\[
k_U = \lfloor n p_U \rfloor
\]

The trimmed mean is:

\[
\bar{x}_{trim} =
\frac{\sum_{i=k_L+1}^{n-k_U}x_{(i)}}
{n-k_L-k_U}
\]

where \(x_{(i)}\) represents the sorted observations.

The implementation uses compensated summation to improve numerical accuracy and scaling to reduce intermediate overflow risk.

## Input Validation and Error Handling

The package returns errors for:

- Empty datasets.
- Negative trimming proportions.
- Trimming proportions greater than or equal to one.
- NaN or infinite trimming proportions.
- More than two trimming arguments.
- Trimming configurations that remove all observations.
- NaN or infinite observations.
- Non-finite calculated results.

A valid trimming proportion must be greater than or equal to zero and less than one.

## Testing

Run the unit tests:

```bash
go test -v ./...
```

Measure test coverage:

```bash
go test -cover ./...
```

Static analysis:

```bash
go vet ./...
```

Build the package:

```bash
go build ./...
```

### Test Coverage Results

The package achieved **98.9% statement coverage**.

Tests include:

- Ordinary arithmetic mean.
- Symmetric trimming.
- Asymmetric trimming.
- Integer and floating-point inputs.
- Negative and zero observations.
- Unsorted datasets.
- Input slice preservation.
- Invalid trimming proportions.
- Empty datasets.
- Non-finite observations.

## Performance Benchmarks

Benchmarks were conducted using Go's built-in testing framework.

Run the benchmarks:

```bash
go test -run "^$" -bench "^BenchmarkMean" -benchmem -benchtime=1s -count=1
```

### Benchmark Environment

| Component | Specification |
|---|---|
| Operating system | Windows |
| Architecture | amd64 |
| Go version | 1.27.1 |
| Processor | AMD Ryzen 7 5800X |
| Benchmark duration | 1 second |

### Benchmark Results

| Benchmark | ns/op | B/op | allocs/op |
|---|---:|---:|---:|
| Mean – 100 observations | 985.6 | 896 | 1 |
| Mean – 10,000 observations | 91,398 | 81,920 | 1 |
| Asymmetric – 1,000 observations | 8,979 | 8,192 | 1 |

The benchmarks demonstrate that the package can process datasets of different sizes while maintaining one allocation per operation in the tested environment.

Results are hardware- and environment-dependent.

## Design Considerations

The implementation uses Go generics to support multiple numeric types through a single function.

Additional design considerations include:

- Sorting a separate slice to preserve the caller's data.
- Returning descriptive errors for invalid inputs.
- Using `math.Floor` to determine the number of observations removed from each tail.
- Using compensated summation for improved numerical accuracy.
- Avoiding unnecessary external dependencies.

## GenAI Tools

ChatGPT (OpenAI) was used as a generative AI coding assistant during development.

Its assistance included:

- Discussing the package architecture and Go module organization.
- Developing an initial implementation using Go generics.
- Suggesting symmetric and asymmetric trimming approaches.
- Developing unit tests and benchmark functions.
- Reviewing error handling and edge cases.
- Assisting with README documentation and project organization.

The generated code was reviewed, executed, and tested locally. Test coverage, build results, static analysis, and benchmark measurements were obtained from the local Go environment.

The developer remains responsible for the final implementation, testing, and submission.

A development conversation log or link to the relevant ChatGPT conversation should accompany the assignment documentation.

## License

This project is intended as an open-source Go package. See the repository's LICENSE file for the applicable license.
