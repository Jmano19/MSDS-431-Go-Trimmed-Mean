package trimmedmean

import (
	"errors"
	"fmt"
	"math"
	"sort"
)

// Number represents the numeric types supported by the trimmed mean package.
type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 |
		~float32 | ~float64
}

// Mean calculates the arithmetic mean after trimming observations.
//
// With no trim argument, Mean calculates the ordinary arithmetic mean.
// With one trim argument, equal proportions are removed from both tails.
// With two trim arguments, the first specifies lower-tail trimming and
// the second specifies upper-tail trimming.
//
// The number of observations removed from each tail is calculated
// using floor(n * proportion), consistent with R's trimmed mean.
//
// Mean does not modify the original input slice.
func Mean[T Number](data []T, trim ...float64) (float64, error) {
	if len(data) == 0 {
		return 0, errors.New("input data cannot be empty")
	}

	if len(trim) > 2 {
		return 0, errors.New("at most two trimming proportions are allowed")
	}

	lowerTrim := 0.0
	upperTrim := 0.0

	if len(trim) >= 1 {
		lowerTrim = trim[0]
		upperTrim = trim[0]
	}

	if len(trim) == 2 {
		upperTrim = trim[1]
	}

	if err := validateTrim(lowerTrim); err != nil {
		return 0, fmt.Errorf("invalid lower trimming proportion: %w", err)
	}

	if err := validateTrim(upperTrim); err != nil {
		return 0, fmt.Errorf("invalid upper trimming proportion: %w", err)
	}

	// Convert observations to float64 and validate the values.
	sortedData := make([]float64, len(data))

	for i, value := range data {
		converted := float64(value)

		if math.IsNaN(converted) || math.IsInf(converted, 0) {
			return 0, fmt.Errorf("observation at index %d is not finite", i)
		}

		sortedData[i] = converted
	}

	// Sort a copy so the caller's original slice remains unchanged.
	sort.Float64s(sortedData)

	lowerCount := int(math.Floor(float64(len(data)) * lowerTrim))
	upperCount := int(math.Floor(float64(len(data)) * upperTrim))

	first := lowerCount
	last := len(sortedData) - upperCount

	if first >= last {
		return 0, errors.New("trimming must leave at least one observation")
	}

	// Scale values before summation to reduce overflow risk.
	scale := 0.0

	for _, value := range sortedData[first:last] {
		scale = math.Max(scale, math.Abs(value))
	}

	if scale == 0 {
		return 0, nil
	}

	// Use compensated summation to improve numerical accuracy.
	sum := 0.0
	compensation := 0.0

	for _, value := range sortedData[first:last] {
		scaledValue := value / scale

		adjusted := scaledValue - compensation
		nextSum := sum + adjusted

		compensation = (nextSum - sum) - adjusted
		sum = nextSum
	}

	count := float64(last - first)
	result := (sum / count) * scale

	if math.IsInf(result, 0) || math.IsNaN(result) {
		return 0, errors.New("calculated mean is not finite")
	}

	return result, nil
}

// validateTrim checks that a trimming proportion is finite
// and falls within the interval [0, 1).
func validateTrim(proportion float64) error {
	if math.IsNaN(proportion) || math.IsInf(proportion, 0) {
		return errors.New("proportion must be finite")
	}

	if proportion < 0 || proportion >= 1 {
		return errors.New("proportion must be at least 0 and less than 1")
	}

	return nil
}
