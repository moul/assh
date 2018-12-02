package ratelimit // import "moul.io/assh/pkg/ratelimit"

// Package ratelimit based on http://hustcat.github.io/rate-limit-example-in-go/

import (
	"fmt"
	"io"
	"time"

	"golang.org/x/time/rate"
)

type reader struct {
	r       io.Reader
	limiter *rate.Limiter
}

// NewReader returns a reader that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func NewReader(r io.Reader, l *rate.Limiter) io.Reader {
	return &reader{
		r:       r,
		limiter: l,
	}
}

func (r *reader) Read(buf []byte) (int, error) {
	n, err := r.r.Read(buf)
	if n <= 0 {
		return n, err
	}

	now := time.Now()
	rv := r.limiter.ReserveN(now, n)
	if !rv.OK() {
		return 0, fmt.Errorf("exceeds limiter's burst")
	}
	delay := rv.DelayFrom(now)
	//fmt.Printf("Read %d bytes, delay %d\n", n, delay)
	time.Sleep(delay)
	return n, err
}

type writer struct {
	w       io.Writer
	limiter *rate.Limiter
}

// NewWriter returns a writer that is rate limited by
// the given token bucket. Each token in the bucket
// represents one byte.
func NewWriter(w io.Writer, l *rate.Limiter) io.Writer {
	return &writer{
		w:       w,
		limiter: l,
	}
}

func (w *writer) Write(buf []byte) (int, error) {
	n, err := w.w.Write(buf)
	if n <= 0 {
		return n, err
	}

	now := time.Now()
	rv := w.limiter.ReserveN(now, n)
	if !rv.OK() {
		return 0, fmt.Errorf("exceeds limiter's burst")
	}
	delay := rv.DelayFrom(now)
	//fmt.Printf("Write %d bytes, delay %d\n", n, delay)
	time.Sleep(delay)
	return n, err
}
