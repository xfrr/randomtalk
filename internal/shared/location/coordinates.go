// coordinates.go
package location

import (
	"fmt"
	"strconv"
)

// Coordinates holds a latitude/longitude pair.
type Coordinates struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// String implements fmt.Stringer.
func (c Coordinates) String() string {
	return fmt.Sprintf("{latitude: %s, longitude: %s}",
		strconv.FormatFloat(c.Latitude, 'f', -1, 64),
		strconv.FormatFloat(c.Longitude, 'f', -1, 64),
	)
}

// IsEmpty reports whether both lat and lon are exactly zero.
func (c Coordinates) IsEmpty() bool {
	return c.Latitude == 0 && c.Longitude == 0
}

// IsValid reports whether lat∈[-90,90] and lon∈[-180,180].
func (c Coordinates) IsValid() bool {
	return c.Latitude >= -90 && c.Latitude <= 90 &&
		c.Longitude >= -180 && c.Longitude <= 180
}
