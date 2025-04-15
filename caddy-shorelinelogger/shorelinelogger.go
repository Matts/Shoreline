package shorelinelogger

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"github.com/caddyserver/caddy/v2/caddyconfig/httpcaddyfile"
	"github.com/caddyserver/caddy/v2/modules/caddyhttp"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(ShorelineLogger{})
	// Register the Caddyfile directive "shorelinelogger" to set up this middleware
	httpcaddyfile.RegisterHandlerDirective("shorelinelogger", parseCaddyfile)
}

// ShorelineLogger is a Caddy module that logs HTTP requests in a curl command format.
type ShorelineLogger struct {
	logger *zap.Logger
}

// CaddyModule returns the Caddy module information.
func (ShorelineLogger) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "http.handlers.shorelinelogger",
		New: func() caddy.Module { return new(ShorelineLogger) },
	}
}

// Provision implements caddy.Provisioner.
func (sl *ShorelineLogger) Provision(ctx caddy.Context) error {
	// Obtain a logger for this module
	sl.logger = ctx.Logger(sl)
	return nil
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (sl *ShorelineLogger) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	// Build a curl command string representing the incoming HTTP request.
	var curlCmd strings.Builder
	curlCmd.WriteString("curl -X ")
	curlCmd.WriteString(r.Method)

	// Include request body if present (read and then restore it for the next handler)
	var bodyBytes []byte
	if r.Body != nil {
		bodyBytes, _ = io.ReadAll(r.Body)
		r.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))  // restore body
	}
	if len(bodyBytes) > 0 {
		// Escape quotes and newlines in body for safe shell command inclusion
		bodyStr := string(bodyBytes)
		bodyStr = strings.ReplaceAll(bodyStr, "\"", "\\\"")
		bodyStr = strings.ReplaceAll(bodyStr, "\n", "\\n")
		curlCmd.WriteString(" -d \"")
		curlCmd.WriteString(bodyStr)
		curlCmd.WriteString("\"")
	}

	// Add headers as curl -H arguments (skip Host and Content-Length as they are implicit/auto)
	for name, values := range r.Header {
		if name == "Host" || name == "Content-Length" {
			continue
		}
		for _, val := range values {
			escapedVal := strings.ReplaceAll(val, "\"", "\\\"")
			curlCmd.WriteString(" -H \"")
			curlCmd.WriteString(fmt.Sprintf("%s: %s", name, escapedVal))
			curlCmd.WriteString("\"")
		}
	}

	// Add the full URL (scheme://host/path?query)
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	curlCmd.WriteString(" ")
	curlCmd.WriteString(fmt.Sprintf("%s://%s%s", scheme, r.Host, r.RequestURI))

	// If client accepted compressed content, include --compressed flag for curl
	if _, ok := r.Header["Accept-Encoding"]; ok {
		curlCmd.WriteString(" --compressed")
	}

	// Log the constructed curl command
	sl.logger.Info(curlCmd.String())

	// Proceed with the next handler
	return next.ServeHTTP(w, r)
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
// Allows configuring this module via Caddyfile (no parameters for this directive).
func (sl *ShorelineLogger) UnmarshalCaddyfile(d *caddyfile.Dispenser) error {
	for d.Next() {
		if d.NextArg() {
			// No arguments are expected for shorelinelogger
			return d.ArgErr()
		}
	}
	return nil
}

// parseCaddyfile sets up the module from Caddyfile tokens.
func parseCaddyfile(h httpcaddyfile.Helper) (caddyhttp.MiddlewareHandler, error) {
	var sl ShorelineLogger
	if err := sl.UnmarshalCaddyfile(h.Dispenser); err != nil {
		return nil, err
	}
	return &sl, nil
}

// Interface guards to ensure struct satisfies required interfaces.
var (
	_ caddy.Provisioner           = (*ShorelineLogger)(nil)
	_ caddyhttp.MiddlewareHandler = (*ShorelineLogger)(nil)
	_ caddyfile.Unmarshaler       = (*ShorelineLogger)(nil)
)
