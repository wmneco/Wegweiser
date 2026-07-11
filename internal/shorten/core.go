// Package shorten implements the shorten-service hexagon: creating a short
// code for a target URL. Only the driving side (HTTP adapter, inbound port,
// core) is wired here; the outbound store port arrives with storage.
package shorten

import (
	"context"
	"errors"
)

// ErrNotImplemented is returned by Core.Shorten until the storage change
// fills in real shortening logic.
var ErrNotImplemented = errors.New("shorten: not implemented")

// Shortener is the inbound port for the shorten hexagon: create a short code
// for a target URL.
type Shortener interface {
	Shorten(ctx context.Context, url string) (string, error)
}

// Core is the shorten hexagon's domain core. It implements Shortener with a
// placeholder until the storage change wires in an outbound save port.
type Core struct{}

// NewCore builds a shorten Core.
func NewCore() *Core {
	return &Core{}
}

// Shorten is a placeholder; it always returns ErrNotImplemented.
func (c *Core) Shorten(ctx context.Context, url string) (string, error) {
	return "", ErrNotImplemented
}
