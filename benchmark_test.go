package trimmedmean

import (
	"testing"
)

func BenchmarkMean100(b *testing.B) {
	data := make([]float64, 100)

	for i := range data {
		data[i] = float64(i + 1)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Mean(data, 0.05)

		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMean10000(b *testing.B) {
	data := make([]float64, 10000)

	for i := range data {
		data[i] = float64(i + 1)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Mean(data, 0.05)

		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkMeanAsymmetric(b *testing.B) {
	data := make([]float64, 1000)

	for i := range data {
		data[i] = float64(i + 1)
	}

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		_, err := Mean(data, 0.05, 0.10)

		if err != nil {
			b.Fatal(err)
		}
	}
}
