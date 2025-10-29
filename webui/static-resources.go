//go:build !webui

package webui

import "embed"

// Dist is empty when not embedding the built UI.
// The web server can still run API-only in this mode.
var Dist embed.FS
