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
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
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
    lumberjackLogger := &lumberjack.Logger{
        Filename:   "logs/shoreline.log",
        MaxSize:    10, // megabytes
        MaxBackups: 5,
        MaxAge:     28,   // days
        Compress:   true, // gzip
    }

    encoderConfig := zap.NewProductionEncoderConfig()
    encoderConfig.TimeKey = "ts"
    encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

    core := zapcore.NewCore(
        zapcore.NewConsoleEncoder(encoderConfig),
        zapcore.AddSync(lumberjackLogger),
        zap.InfoLevel,
    )

    sl.logger = zap.New(core)
	return nil
}

type responseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       *bytes.Buffer
}

func (r *responseRecorder) WriteHeader(code int) {
	r.statusCode = code
	r.ResponseWriter.WriteHeader(code)
}

func (r *responseRecorder) Write(b []byte) (int, error) {
	r.body.Write(b)
	return r.ResponseWriter.Write(b)
}

// ServeHTTP implements caddyhttp.MiddlewareHandler.
func (sl *ShorelineLogger) ServeHTTP(w http.ResponseWriter, r *http.Request, next caddyhttp.Handler) error {
	rec := &responseRecorder{
		ResponseWriter: w,
		statusCode:     200,
		body:           &bytes.Buffer{},
	}

	// Build a curl command string representing the incoming HTTP request.
	var curlCmd strings.Builder

	serverName := "";
	serverIface := r.Context().Value(caddyhttp.ServerCtxKey)
	if serverIface != nil {
		server, ok := serverIface.(*caddyhttp.Server)
		if ok {
			serverName = server.Name()

		}
	}

	curlCmd.WriteString(serverName);
	curlCmd.WriteString(":")
	curlCmd.WriteString(r.Host);
	curlCmd.WriteString(" REQUEST\n")
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
	err := next.ServeHTTP(rec, r)

	// Log the response status code and body
	if err == nil {
		var responseContent strings.Builder

		responseContent.WriteString(serverName);
		responseContent.WriteString(":")
		responseContent.WriteString(r.Host);
		responseContent.WriteString(" RESPONSE\n")
		responseContent.WriteString(fmt.Sprintf("%d %s\n", rec.statusCode, http.StatusText(rec.statusCode)))
		for name, values := range r.Header {
			if name == "Content-Length" {
				continue
			}
			for _, val := range values {
				responseContent.WriteString(fmt.Sprintf("%s: %s\n", name, val))
			}
		}
		responseContent.WriteString("\nBODY: ");
		responseContent.WriteString(rec.body.String());
		responseContent.WriteString("\n");
		sl.logger.Info(responseContent.String())
	}

	return err;
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
