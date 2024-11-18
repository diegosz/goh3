package h3

import (
	"errors"
	"runtime"
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/akhenakh/goh3/ch3"
	"modernc.org/libc"
)

type Batch struct {
	*libc.TLS
}

func NewBatch() *Batch {
	return &Batch{TLS: libc.NewTLS()}
}

func (c *Batch) FromGeo(geo GeoCoord, res int) H3Index {
	return fromGeo(c.TLS, geo, res)
}

func (c *Batch) FromLatLng(lat, lng float64, res int) H3Index {
	return c.FromGeo(GeoCoord{Latitude: lat, Longitude: lng}, res)
}

func (c *Batch) ToGeo(h H3Index) GeoCoord {
	return toGeo(c.TLS, h)
}

func (c *Batch) Close() {
	c.TLS.Close()
}

// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------

type TBatch struct {
	calls chan toGeoReq
}

func NewTBatch() *TBatch {
	c := &TBatch{
		calls: make(chan toGeoReq),
	}
	go c.backgroundThread(c.calls)
	return c
}

func (c *TBatch) ToGeo(h H3Index) GeoCoord {
	ret := make(chan GeoCoord, 1)
	r := toGeoReq{Cell: h, Return: ret}
	c.calls <- r
	return <-ret
}

type toGeoReq struct {
	Cell   H3Index
	Return chan<- GeoCoord
}

func (c *TBatch) backgroundThread(fcalls <-chan toGeoReq) {
	runtime.LockOSThread()
	tls := libc.NewTLS()
	for f := range fcalls {
		f.Return <- toGeo(tls, f.Cell)
	}
	runtime.UnlockOSThread()
}

// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------

var ErrClosed = errors.New("caller closed")

type Caller struct {
	calls chan call
}

type callFunc int

const (
	FuncClose callFunc = iota
	FuncFromGeo
	FuncToGeo
)

type call struct {
	Func   callFunc
	Args   []any
	Return chan<- GeoCoord
}

func NewCaller() *Caller {
	c := &Caller{
		calls: make(chan call),
	}
	ready := make(chan struct{})
	go func() {
		runtime.LockOSThread()
		c.backgroundThread(ready)
		runtime.UnlockOSThread()
	}()
	<-ready
	return c
}

func (c *Caller) ToGeo(cell H3Index) (GeoCoord, error) {
	ret := make(chan GeoCoord, 1)
	// c.calls <- call{Func: FuncToGeo, Args: []any{cell}, Return: ret}
	c.calls <- call{Func: FuncToGeo, Args: []any{cell}, Return: ret}
	r := <-ret
	// if len(r) != 2 {
	// 	return GeoCoord{Latitude: 1}, errors.New("invalid return length")
	// }
	// if r[1] != nil {
	// 	e, ok := r[1].(error)
	// 	if !ok {
	// 		return GeoCoord{Latitude: 1}, errors.New("invalid return error type")
	// 	}
	// 	return GeoCoord{Latitude: 1}, e
	// }
	// g, ok := r[0].(GeoCoord)
	// if !ok {
	// 	return GeoCoord{Latitude: 1}, errors.New("invalid return type")
	// }
	return r, nil
}

func (c *Caller) backgroundThread(ready chan struct{}) {
	// runtime.LockOSThread()
	calls := make(chan call)
	c.calls = calls
	close(ready)
	// tls := libc.NewTLS()
	for f := range calls {
		// switch f.Func {
		// case FuncClose:
		// 	return
		// case FuncFromGeo:
		// 	geo := f.Args[0].(GeoCoord)
		// 	res := f.Args[1].(int)
		// 	f.Return <- []any{fromGeo(tls, geo, res), nil}
		// case FuncToGeo:
		// 	cell, ok := f.Args[0].(H3Index)
		// 	if !ok {
		// 		f.Return <- []any{GeoCoord{Latitude: 1}, errors.New("invalid argument type")}
		// 		continue
		// 	}
		// 	f.Return <- []any{toGeo(tls, cell), nil}
		// default:
		// 	panic("invalid function")
		// }
		tls := libc.NewTLS()

		cell := f.Args[0].(H3Index)
		cg := ch3.TLatLng{}
		ch3.XcellToLatLng(tls, ch3.TH3Index(cell), uintptr(unsafe.Pointer(&cg)))
		g := GeoCoord{}
		g.Latitude = rad2deg * float64(cg.Flat)
		g.Longitude = rad2deg * float64(cg.Flng)
		f.Return <- g
		// f.Return <- GeoCoord{Latitude: g.Latitude, Longitude: 2}
		// f.Return <- GeoCoord{Latitude: float64(uint64(cell)), Longitude: 2}
		tls.Close()
	}
	// runtime.UnlockOSThread()
}

// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------
// -----------------------------------------------------------------------------

type Caller2 struct {
	calls  chan call2
	mu     sync.Mutex
	closed atomic.Bool
}

type callFunc2 int

const (
	Func2Close callFunc2 = iota
	Func2FromGeo
	Func2ToGeo
)

type call2 struct {
	Func   callFunc2
	Args   []any
	Return chan<- []any
}

func NewCaller2() *Caller2 {
	calls := make(chan call2)
	c := &Caller2{
		calls: calls,
	}
	go c.backgroundThread(c.calls)
	return c
}

func (c *Caller2) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed.Load() {
		return
	}
	c.closed.Store(true)
	close(c.calls)
}

func (c *Caller2) ToGeo(cell H3Index) (GeoCoord, error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed.Load() {
		return GeoCoord{}, ErrClosed
	}
	ret := make(chan []any, 1)
	c.calls <- call2{Func: Func2ToGeo, Args: []any{cell}, Return: ret}
	r := <-ret
	if len(r) != 2 {
		return GeoCoord{Latitude: 1}, errors.New("invalid return length")
	}
	if r[1] != nil {
		e, ok := r[1].(error)
		if !ok {
			return GeoCoord{Latitude: 1}, errors.New("invalid return error type")
		}
		return GeoCoord{Latitude: 1}, e
	}
	g, ok := r[0].(GeoCoord)
	if !ok {
		return GeoCoord{Latitude: 1}, errors.New("invalid return type")
	}
	return g, nil
}

func (c *Caller2) backgroundThread(calls <-chan call2) {
	runtime.LockOSThread()
	// tls := libc.NewTLS()
	// defer runtime.UnlockOSThread()
	for f := range calls {
		tls := libc.NewTLS()
		switch f.Func {
		case Func2Close:
			return
		case Func2FromGeo:
			geo := f.Args[0].(GeoCoord)
			res := f.Args[1].(int)
			f.Return <- []any{fromGeo(tls, geo, res), nil}
		case Func2ToGeo:
			cell, ok := f.Args[0].(H3Index)
			if !ok {
				f.Return <- []any{GeoCoord{Latitude: 1}, errors.New("invalid argument type")}
				continue
			}
			g := toGeo(tls, cell)
			f.Return <- []any{g, nil}
		}
		tls.Close()
	}
	runtime.UnlockOSThread()
}
