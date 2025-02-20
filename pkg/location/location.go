package location

import (
	"encoding/json"
	"math"
	"strconv"
	"strings"
)

const Pi = 3.14159265358979323846

// New creates a new Location instance.
func New(latitude, longitude float64) Location {
	return Location{
		state: locationState{
			Coordinates: &Coordinates{
				Latitude:  latitude,
				Longitude: longitude,
			},
		},
	}
}

// Location represents a geographical coordinate.
type Location struct {
	state locationState
}

type locationState struct {
	Coordinates *Coordinates `json:"coordinates"`
	CountryCode string       `json:"country_code"`
	CityCode    string       `json:"city_code"`
}

func (l Location) IsEmpty() bool {
	return l.state.Coordinates.IsEmpty()
}

func (l *Location) Coordinates() *Coordinates {
	return l.state.Coordinates
}

func (l *Location) CountryCode() string {
	return l.state.CountryCode
}

func (l *Location) CityCode() string {
	return l.state.CityCode
}

func (l *Location) WithCountryCode(countryCode string) *Location {
	l.state.CountryCode = countryCode
	return l
}

func (l *Location) WithCityCode(cityCode string) *Location {
	l.state.CityCode = cityCode
	return l
}

// DistanceTo calculates the distance between two locations in kilometers.
func (l *Location) DistanceTo(other Location) (float64, bool) {
	if l.state.Coordinates == nil || other.state.Coordinates == nil {
		return 0, false
	}

	lat1 := degToRad(l.state.Coordinates.Latitude)
	lon1 := degToRad(l.state.Coordinates.Longitude)
	lat2 := degToRad(other.state.Coordinates.Latitude)
	lon2 := degToRad(other.state.Coordinates.Longitude)

	return calculateHaversineDistance(lat1, lon1, lat2, lon2)
}

func (l *Location) Equals(other Location) bool {
	return l.state.Coordinates.Latitude == other.state.Coordinates.Latitude &&
		l.state.Coordinates.Longitude == other.state.Coordinates.Longitude
}

func (l *Location) UnmarshalJSON(data []byte) error {
	var state locationState
	if err := json.Unmarshal(data, &state); err != nil {
		return err
	}

	l.state = state
	return nil
}

func (l Location) MarshalJSON() ([]byte, error) {
	return json.Marshal(l.state)
}

func (l Location) String() string {
	var sb strings.Builder
	sb.WriteString("{")

	if l.Coordinates() != nil {
		sb.WriteString("coordinates: ")
		sb.WriteString(l.Coordinates().String())
	}

	if l.state.CountryCode != "" {
		sb.WriteString(", country_code: ")
	}

	if l.state.CityCode != "" {
		sb.WriteString(", city_code: ")
	}

	sb.WriteString("}")

	return sb.String()
}

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', -1, 64)
}

func calculateHaversineDistance(lat1, lon1, lat2, lon2 float64) (float64, bool) {
	const earthRadiusKm = 6371
	const multiplier = 2.0

	dLat := lat2 - lat1
	dLon := lon2 - lon1

	a := (squared(math.Sin(dLat/2)) +
		math.Cos(lat1)*math.Cos(lat2)*squared(math.Sin(dLon/2)))

	c := multiplier * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c, true
}

func squared(x float64) float64 {
	return x * x
}

func degToRad(deg float64) float64 {
	return deg * (Pi / 180)
}
