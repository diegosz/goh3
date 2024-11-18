package h3

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const eps = 1e-4

// validH3Index resolution 5
const (
	validH3Index = H3Index(0x850dab63fffffff)
)

var validGeoCoord = GeoCoord{
	Latitude:  67.15092686397712,   // 67.1509268640
	Longitude: -168.39088858096966, // -168.3908885810
}

func almostEqual(t *testing.T, expected, actual interface{}, msgAndArgs ...interface{}) {
	assert.InEpsilon(t, expected, actual, eps, msgAndArgs...)
}

func assertGeoCoord(t *testing.T, expected, actual GeoCoord) {
	require.Equal(t, expected.Latitude, actual.Latitude)
	require.Equal(t, expected.Longitude, actual.Longitude)
	almostEqual(t, expected.Latitude, actual.Latitude, "latitude mismatch")
	almostEqual(t, expected.Longitude, actual.Longitude, "longitude mismatch")
}

func TestFromGeo(t *testing.T) {
	t.Parallel()
	h := FromGeo(GeoCoord{
		Latitude:  67.194013596,
		Longitude: 191.598258018,
	}, 5)
	assert.Equal(t, validH3Index, h, "expected %x but got %x", validH3Index, h)
}

func TestToGeo(t *testing.T) {
	t.Parallel()
	g := ToGeo(validH3Index)
	assertGeoCoord(t, validGeoCoord, g)
}

func TestToGeoBatch(t *testing.T) {
	t.Parallel()
	b := NewTBatch()
	g := b.ToGeo(validH3Index)
	assertGeoCoord(t, validGeoCoord, g)
}

func TestToGeoCaller(t *testing.T) {
	t.Parallel()
	c := NewCaller2()
	g, err := c.ToGeo(validH3Index)
	assert.NoError(t, err)
	assertGeoCoord(t, validGeoCoord, g)
}

func BenchmarkToGeo(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ToGeo(validH3Index)
	}
}
