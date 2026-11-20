package trimmedmean

import (
	"math"
	"reflect"
	"testing"
)

func TestMean(t *testing.T) {
	tests := []struct {
		name     string
		data     []float64
		trim     []float64
		expected float64
	}{
		{
			name:     "ordinary mean",
			data:     []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			expected: 5.5,
		},
		{
			name:     "symmetric trimming",
			data:     []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			trim:     []float64{0.2},
			expected: 5.5,
		},
		{
			name:     "asymmetric trimming",
			data:     []float64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10},
			trim:     []float64{0.1, 0.2},
			expected: 5,
		},
		{
			name:     "zero trimming",
			data:     []float64{2, 4, 6, 8},
			trim:     []float64{0},
			expected: 5,
		},
		{
			name:     "unsorted observations",
			data:     []float64{10, 2, 8, 4, 6},
			trim:     []float64{0.2},
			expected: 6,
		},
		{
			name:     "negative observations",
			data:     []float64{-10, -5, 0, 5, 10},
			expected: 0,
		},
		{
			name:     "all zero observations",
			data:     []float64{0, 0, 0, 0},
			trim:     []float64{0.25},
			expected: 0,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			actual, err := Mean(test.data, test.trim...)

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if math.Abs(actual-test.expected) > 1e-10 {
				t.Errorf("expected %.10f, got %.10f",
					test.expected, actual)
			}
		})
	}
}

func TestMeanInt(t *testing.T) {
	data := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	actual, err := Mean(data, 0.2)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := 5.5

	if math.Abs(actual-expected) > 1e-10 {
		t.Errorf("expected %.2f, got %.2f", expected, actual)
	}
}

func TestMeanFloat32(t *testing.T) {
	data := []float32{1.5, 2.5, 3.5, 4.5, 5.5}

	actual, err := Mean(data)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	expected := 3.5

	if math.Abs(actual-expected) > 1e-10 {
		t.Errorf("expected %.2f, got %.2f", expected, actual)
	}
}

func TestMeanDoesNotModifyInput(t *testing.T) {
	data := []float64{9, 2, 7, 1, 5}
	original := append([]float64(nil), data...)

	_, err := Mean(data, 0.2)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !reflect.DeepEqual(data, original) {
		t.Errorf("input data was modified: got %v, want %v",
			data, original)
	}
}

func TestMeanInvalidInputs(t *testing.T) {
	tests := []struct {
		name string
		data []float64
		trim []float64
	}{
		{
			name: "empty input",
			data: []float64{},
		},
		{
			name: "negative trimming",
			data: []float64{1, 2, 3, 4},
			trim: []float64{-0.1},
		},
		{
			name: "trimming greater than one",
			data: []float64{1, 2, 3, 4},
			trim: []float64{1.0},
		},
		{
			name: "NaN trimming",
			data: []float64{1, 2, 3, 4},
			trim: []float64{math.NaN()},
		},
		{
			name: "infinite trimming",
			data: []float64{1, 2, 3, 4},
			trim: []float64{math.Inf(1)},
		},
		{
			name: "too many trimming arguments",
			data: []float64{1, 2, 3, 4},
			trim: []float64{0.1, 0.1, 0.1},
		},
		{
			name: "all observations trimmed",
			data: []float64{1, 2, 3, 4},
			trim: []float64{0.5},
		},
		{
			name: "invalid upper trimming",
			data: []float64{1, 2, 3, 4},
			trim: []float64{0.1, 1.0},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			_, err := Mean(test.data, test.trim...)

			if err == nil {
				t.Error("expected an error, got nil")
			}
		})
	}
}

func TestMeanNonFiniteObservations(t *testing.T) {
	tests := [][]float64{
		{1, 2, math.NaN()},
		{1, math.Inf(1), 3},
		{1, math.Inf(-1), 3},
	}

	for _, data := range tests {
		_, err := Mean(data)

		if err == nil {
			t.Errorf("expected error for non-finite input: %v", data)
		}
	}
}
