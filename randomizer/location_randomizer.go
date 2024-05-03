package randomizer

import "math/rand"

func GetRandomLatitude(locationSpread float64) float64 {
	minLat := 45.0 - (locationSpread / 2)
	maxLat := 45.0 + (locationSpread / 2)
	randomLat := minLat + rand.Float64()*(maxLat-minLat)
	return randomLat
}

func GetRandomLongitude(locationSpread float64) float64 {
	minLong := 90.0 - (locationSpread / 2)
	maxLong := 90.0 + (locationSpread / 2)
	randomLong := minLong + rand.Float64()*(maxLong-minLong)
	return randomLong
}
