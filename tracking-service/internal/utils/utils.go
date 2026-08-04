package utils

import (
	"errors"
	"math"
	"time"
)

func ValidateCoordinates(latitude float64, longitude float64) error {
	if latitude < -90 || latitude > 90 {
		return errors.New("latitude must be between -90 and 90")
	}
	if longitude < -180 || longitude > 180 {
		return errors.New("longitude must be between -180 and 180")
	}
	return nil
}

// using the haversine formula to calculate the distance between two points on the earth's surface
func CalculateDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadius = 6371 // Earth's radius in kilometers

	lat1Rad := lat1 * math.Pi / 180
	lon1Rad := lon1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	lon2Rad := lon2 * math.Pi / 180

	latDiff := lat2Rad - lat1Rad
	lonDiff := lon2Rad - lon1Rad

	a := math.Sin(latDiff/2)*math.Sin(latDiff/2) + math.Cos(lat1Rad)*math.Cos(lat2Rad)*math.Sin(lonDiff/2)*math.Sin(lonDiff/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

func CalculateETA(distanceKm float64, avgSpeedKmH float64) time.Duration {
	hours := distanceKm / avgSpeedKmH
	return time.Duration(hours * float64(time.Hour))
}
