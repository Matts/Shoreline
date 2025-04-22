all:
	@echo "No target specified. Available targets:"
	@echo "  caddy: Build Caddy with the Shoreline plugin."
	@echo "  ui: Build the UI."
	@echo "  clean: Clean up build artifacts."

caddy: caddy-shoreline/go.mod caddy-shoreline/shorelinelogger.go
	@echo "Building Caddy..."
	@xcaddy build v2.9.1 --with github.com/Matts/shoreline/caddy-shoreline=./caddy-shoreline
	@echo "Caddy build complete."

ui/.next/build-manifest.json:
	@echo "Building UI..."
	@cd ui && npm install && npm run build
	@echo "UI build complete."

ui: ui/.next/build-manifest.json
	@echo "Built..."

clean:
	@echo "Cleaning up..."
	@rm -rf ui/.next
	@rm -rf caddy
	@echo "Clean complete."

.PHONY: all ui clean help
