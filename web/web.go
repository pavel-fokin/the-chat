// Package web embeds the built frontend (dist) so the-chat-server can serve
// it without depending on files on disk at runtime.
package web

import "embed"

//go:embed all:dist
var Dist embed.FS
