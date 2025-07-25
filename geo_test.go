package main

import (
	"math"
	"testing"
)

const Epsilon = 1e-7 // Precision for lat,lon comparisons, down to metre accuracy

func areFloat64Equal(a, b float64) bool {
	return math.Abs(a-b) < Epsilon
}

type testCoordinates struct {
	p     coordinates
	q     coordinates
	outKm float64
}

func getTestCoordinates() []testCoordinates {
	return []testCoordinates{
		{
			newCoordinates(22.55, 43.12),  // Rio de Janeiro, Brazil
			newCoordinates(13.45, 100.28), // Bangkok, Thailand
			6094.544408786774,
		},
		{
			newCoordinates(20.10, 57.30), // Port Louis, Mauritius
			newCoordinates(0.57, 100.21), // Padang, Indonesia
			5145.525771394785,
		},
		{
			newCoordinates(51.45, 1.15),  // Oxford, United Kingdom
			newCoordinates(41.54, 12.27), // Vatican, City Vatican City
			1389.1793118293067,
		},
		{
			newCoordinates(22.34, 17.05), // Windhoek, Namibia
			newCoordinates(51.56, 4.29),  // Rotterdam, Netherlands
			3429.89310043882,
		},
		{
			newCoordinates(63.24, 56.59), // Esperanza, Argentina
			newCoordinates(8.50, 13.14),  // Luanda, Angola
			6996.18595539861,
		},
		{
			newCoordinates(90.00, 0.00), // North/South Poles
			newCoordinates(48.51, 2.21), // Paris,  France
			4613.477506482742,
		},
		{
			newCoordinates(45.04, 7.42),  // Turin, Italy
			newCoordinates(3.09, 101.42), // Kuala Lumpur, Malaysia
			10078.111954385415,
		},
	}
}

func TestHaversineDistance(t *testing.T) {
	for _, input := range getTestCoordinates() {
		kilometers := Distance(input.p, input.q).Kilometers()

		if !areFloat64Equal(input.outKm, kilometers) {
			t.Errorf("fail: want %v %v -> %v got %v",
				input.p,
				input.q,
				input.outKm,
				kilometers,
			)
		}
	}
}
