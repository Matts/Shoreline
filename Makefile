all:
	@echo "No default target. Please specify a target."

caddy: caddy-shoreline/go.mod caddy-shoreline/shorelinelogger.go
	@echo "Building Caddy..."
	@xcaddy build v2.9.1 --with github.com/Matts/shoreline/caddy-shoreline=./caddy-shoreline
	@echo "Caddy build complete."

ui/.next/build-manifest.json:
	@echo "Building UI..."
	@cd ui && npm install && npm run build
	@echo "UI build complete."

ui: ui/.next/build-manifest.json ui/package.json ui/package-lock.json ui/next.config.js
	@echo "Built..."


clean:
	@echo "Cleaning up..."
	@rm -rf ui/.next
	@rm -rf caddy
	@echo "Clean complete."
