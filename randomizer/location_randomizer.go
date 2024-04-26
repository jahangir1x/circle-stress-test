package randomizer

import "math/rand"

func GetRandomLatitude() float64 {
	minLat := 45.0
	maxLat := 45.05
	randomLat := minLat + rand.Float64()*(maxLat-minLat)
	return randomLat
}

func GetRandomLongitude() float64 {
	minLong := 45.0
	maxLong := 45.05
	randomLong := minLong + rand.Float64()*(maxLong-minLong)
	return randomLong
}
