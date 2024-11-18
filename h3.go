package h3

import (
	"math"
	"runtime"
	"unsafe"

	"github.com/akhenakh/goh3/ch3"
	"modernc.org/libc"
)

var (
	deg2rad = math.Pi / 180.0
	rad2deg = 180.0 / math.Pi
)

type H3Index ch3.TH3Index

func (h H3Index) String() string {
	return uint64ToHex(uint64(h))
}

type GeoCoord struct {
	Latitude, Longitude float64
}

type GeoPolygon struct {
	Geofence []GeoCoord
	Holes    [][]GeoCoord
}

func FromGeo(geo GeoCoord, res int) H3Index {
	tls := libc.NewTLS()
	defer tls.Close()
	return fromGeo(tls, geo, res)
}

func FromLatLng(lat, lng float64, res int) H3Index {
	return FromGeo(GeoCoord{Latitude: lat, Longitude: lng}, res)
}

func ToGeo(h H3Index) GeoCoord {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	tls := libc.NewTLS()
	defer tls.Close()
	return toGeo(tls, h)
}

func fromGeo(tls *libc.TLS, geo GeoCoord, res int) H3Index {
	cgeo := ch3.TLatLng{
		Flat: deg2rad * geo.Latitude,
		Flng: deg2rad * geo.Longitude,
	}
	var i ch3.TH3Index
	ch3.XlatLngToCell(tls, uintptr(unsafe.Pointer(&cgeo)), int32(res), uintptr(unsafe.Pointer(&i)))
	return H3Index(i)
}

func toGeo(tls *libc.TLS, h H3Index) GeoCoord {
	cg := ch3.TLatLng{}
	ch3.XcellToLatLng(tls, ch3.TH3Index(h), uintptr(unsafe.Pointer(&cg)))
	g := GeoCoord{}
	g.Latitude = rad2deg * float64(cg.Flat)
	g.Longitude = rad2deg * float64(cg.Flng)
	return g
}
