# ------------------------------ config ------------------------------
TAILWIND_VERSION := v4.3.3
TAILWIND_BIN := ./bin/tailwindcss
STYLES_ROOT := ./cmd/web/assets/styles

.PHONY: dev setup install-tailwind clean

dev: setup
		go tool templ generate --watch --proxy="http://localhost:8080" --cmd="go run ./cmd/web" & \
		$(TAILWIND_BIN) -i $(STYLES_ROOT)/input.css -o $(STYLES_ROOT)/index.css --watch & \
		wait

# Ensure the necessary dependencies (tools) exist
# install/download only upon absence
setup:
# 		@command -v templ >/dev/null 2>&1 || go install github.com/a-h/templ/cmd/templ@latest
		go get -tool github.com/a-h/templ/cmd/templ@latest
		@test -x $(TAILWIND_BIN) || $(MAKE) install-tailwind

install-tailwind:
		@printf '\033[32mDownloading Tailwind standalone CLI $(TAILWIND_VERSION)...\033[0m\n'
		@mkdir -p bin
		@OS=$$(uname -s | tr '[:upper:]' '[:lower:]'); \
			ARCH=$$(uname -m); \
			case "$$OS"		in darwin) OS=macos ;; esac; \
			case "$$ARCH" in x86_64) ARCH=x64 ;; aarch64) ARCH=arm64 ;; esac; \
			URL="https://github.com/tailwindlabs/tailwindcss/releases/download/$(TAILWIND_VERSION)/tailwindcss-$$OS-$$ARCH"; \
			echo "  -> $$URL"; \
			curl -sSL "$$URL" -o $(TAILWIND_BIN); \
			chmod +x $(TAILWIND_BIN)

# Remove the downloaded Tailwind binary (project-local, safe to delete)
clean:
			rm -rf bin
