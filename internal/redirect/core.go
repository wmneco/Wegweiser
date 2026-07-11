// Package redirect implements the redirect-service hexagon: resolving a
// short code to its target URL. Only the driving side (HTTP adapter, inbound
// port, core) is wired here; the outbound store port arrives with storage.
package redirect

import (
	"context"
	"errors"
)

// ErrNotImplemented is returned by Core.Resolve until the storage change
// fills in real resolution logic.
var ErrNotImplemented = errors.New("redirect: not implemented")

// Resolver is the inbound port for the redirect hexagon: resolve a short
// code to its target URL.
type Resolver interface {
	Resolve(ctx context.Context, code string) (string, error)
}

// Core is the redirect hexagon's domain core. It implements Resolver with a
// placeholder until the storage change wires in an outbound lookup port.
type Core struct{}

// NewCore builds a redirect Core.
func NewCore() *Core {
	return &Core{}
}

// Resolve is a placeholder; it always returns ErrNotImplemented.
func (c *Core) Resolve(ctx context.Context, code string) (string, error) {
	return "", ErrNotImplemented
}
