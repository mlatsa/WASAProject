//go:build webui

package webui

import "embed"

// Embed the built frontend (dist/)
 //go:embed dist/* dist/**/*
var Dist embed.FS
