all:
	@echo "No default target. Please specify a target."

caddy:
	@echo "Building Caddy..."
	@xcaddy build v2.9.1 --with github.com/Matts/shoreline/caddy-shoreline=./caddy-shoreline
	@echo "Caddy build complete."

build-ui:
	@echo "Building UI..."
	@cd ui && npm install && npm run build
	@echo "UI build complete."