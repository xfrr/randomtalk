package location_test

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfrr/randomtalk/internal/shared/location"
)

func TestNewAndEmpty(t *testing.T) {
	empty := location.New(0, 0)
	assert.True(t, empty.IsEmpty())

	loc := location.New(12.34, 56.78)
	assert.False(t, loc.IsEmpty())
	assert.InDelta(t, 12.34, loc.Coordinates.Latitude, 0.0001)
	assert.InDelta(t, 56.78, loc.Coordinates.Longitude, 0.0001)
}

func TestWithCodesAndEquals(t *testing.T) {
	base := location.New(1, 2)
	a := base.WithCountryCode("US").WithCityCode("NYC")
	b := location.New(1, 2).WithCountryCode("US").WithCityCode("NYC")
	assert.True(t, a.Equals(b))

	c := base.WithCountryCode("CA")
	assert.False(t, a.Equals(c))
}

func TestCoordinatesValidation(t *testing.T) {
	valid := location.Coordinates{Latitude: 90, Longitude: -180}
	assert.True(t, valid.IsValid())

	invalidLat := location.Coordinates{Latitude: -91, Longitude: 0}
	assert.False(t, invalidLat.IsValid())

	invalidLon := location.Coordinates{Latitude: 0, Longitude: 181}
	assert.False(t, invalidLon.IsValid())
}

func TestDistanceTo(t *testing.T) {
	origin := location.New(0, 0)
	// self-distance
	d0, err := origin.DistanceTo(origin)
	require.NoError(t, err)
	assert.InDelta(t, 0.0, d0, 0.0001)

	// 1° latitude ≈ 111.319 km
	oneDeg := location.New(1, 0)
	d1, err := origin.DistanceTo(oneDeg)
	require.NoError(t, err)
	assert.InDelta(t, 111.319, d1, 0.15)

	// invalid coordinate
	bad := location.New(0, 200)
	_, err = origin.DistanceTo(bad)
	assert.ErrorIs(t, err, location.ErrInvalidCoordinates)
}

func TestJSONRoundTrip(t *testing.T) {
	loc := location.New(1.23, 4.56).WithCountryCode("US").WithCityCode("LA")
	data, err := json.Marshal(loc)
	require.NoError(t, err)

	expectedJSON := `{
	  "coordinates": {"latitude":1.23,"longitude":4.56},
	  "country_code":"US",
	  "city_code":"LA"
	}`
	assert.JSONEq(t, expectedJSON, string(data))

	var unmarshaled location.Location
	err = json.Unmarshal(data, &unmarshaled)
	require.NoError(t, err)
	assert.True(t, loc.Equals(unmarshaled))
}

func TestString(t *testing.T) {
	loc := location.New(1, 2).WithCountryCode("C").WithCityCode("D")
	s := loc.String()
	assert.Contains(t, s, "coordinates:")
	assert.Contains(t, s, "latitude: 1")
	assert.Contains(t, s, "longitude: 2")
	assert.Contains(t, s, `country_code: "C"`)
	assert.Contains(t, s, `city_code: "D"`)
}
