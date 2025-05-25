// location.go
package location

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// ErrInvalidCoordinates is returned when a location has out-of-range latitude or longitude.
var ErrInvalidCoordinates = errors.New("invalid coordinates")

// Location represents a geographic point, with optional ISO country and city codes.
type Location struct {
	Coordinates Coordinates `json:"coordinates"`
	CountryCode string      `json:"country_code,omitempty"`
	CityCode    string      `json:"city_code,omitempty"`
}

// New returns a Location at the given latitude/longitude.
func New(lat, lon float64) Location {
	return Location{Coordinates: Coordinates{Latitude: lat, Longitude: lon}}
}

// IsEmpty reports whether both latitude and longitude are zero.
func (l Location) IsEmpty() bool {
	return l.Coordinates.IsEmpty()
}

// WithCountryCode returns a copy of l with CountryCode set.
func (l Location) WithCountryCode(code string) Location {
	l.CountryCode = code
	return l
}

// WithCityCode returns a copy of l with CityCode set.
func (l Location) WithCityCode(code string) Location {
	l.CityCode = code
	return l
}

// Equals reports whether two Locations have the same coordinates, country, and city.
func (l Location) Equals(o Location) bool {
	return l.Coordinates == o.Coordinates &&
		l.CountryCode == o.CountryCode &&
		l.CityCode == o.CityCode
}

// DistanceTo returns the great-circle distance between l and o (in km).
// If either location has invalid coordinates, ErrInvalidCoordinates is returned.
func (l Location) DistanceTo(o Location) (float64, error) {
	if !l.Coordinates.IsValid() || !o.Coordinates.IsValid() {
		return 0, ErrInvalidCoordinates
	}
	// convert degrees → radians
	lat1, lon1 := degToRad(l.Coordinates.Latitude), degToRad(l.Coordinates.Longitude)
	lat2, lon2 := degToRad(o.Coordinates.Latitude), degToRad(o.Coordinates.Longitude)
	// haversine formula
	dLat, dLon := lat2-lat1, lon2-lon1
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(dLon/2)*math.Sin(dLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	const earthRadiusKm = 6371
	return earthRadiusKm * c, nil
}

// String implements fmt.Stringer for easy debugging.
func (l Location) String() string {
	parts := []string{fmt.Sprintf("coordinates: %s", l.Coordinates.String())}
	if l.CountryCode != "" {
		parts = append(parts, fmt.Sprintf("country_code: %q", l.CountryCode))
	}
	if l.CityCode != "" {
		parts = append(parts, fmt.Sprintf("city_code: %q", l.CityCode))
	}
	return fmt.Sprintf("{%s}", strings.Join(parts, ", "))
}

func degToRad(deg float64) float64 {
	return deg * (math.Pi / 180)
}
