package location

// Coordinates represents a geographical coordinate.
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

func (c Coordinates) String() string {
	return "{latitude: " + formatFloat(c.Latitude) + ", longitude: " + formatFloat(c.Longitude) + "}"
}

func (c Coordinates) IsEmpty() bool {
	return c.Latitude == 0 && c.Longitude == 0
}

func (c Coordinates) IsValid() bool {
	return c.Latitude >= -90 && c.Latitude <= 90 &&
		c.Longitude >= -180 && c.Longitude <= 180
}
