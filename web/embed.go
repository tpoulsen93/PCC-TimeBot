// Package web embeds the built React single-page application so it can be
// served directly from the Go binary. The real assets are produced by
// `vite build` into app/dist; a placeholder index.html is committed so that
// `go build` always succeeds even before the frontend has been built.
package web

import "embed"

//go:embed all:app/dist
var distFS embed.FS

// DistFS returns the embedded filesystem containing the built SPA.
func DistFS() embed.FS {
	return distFS
}
