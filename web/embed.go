package web

import "embed"

// Build the frontend before compiling Go: npm --prefix web run build.
//
//go:embed dist
var Files embed.FS
