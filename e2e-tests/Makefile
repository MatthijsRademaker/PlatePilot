# PlatePilot iOS E2E Tests Makefile

.PHONY: help test test-flow test-debug clean setup install-deps

help:
	@echo "PlatePilot iOS E2E Tests"
	@echo ""
	@echo "Available targets:"
	@echo "  make test              - Run all E2E tests"
	@echo "  make test-flow name=X  - Run specific test flow (e.g., make test-flow name=home)"
	@echo "  make test-debug        - Run tests in debug mode"
	@echo "  make test-clean        - Clean build and run tests"
	@echo "  make test-quick        - Run tests (skip build if possible)"
	@echo "  make setup             - Set up backend services only"
	@echo "  make clean             - Clean up everything"
	@echo "  make install-deps      - Install all dependencies (Maestro, XcodeGen)"
	@echo ""

# Install dependencies
install-deps:
	@echo "Installing dependencies..."
	@command -v xcodegen >/dev/null 2>&1 || brew install xcodegen
	@command -v maestro >/dev/null 2>&1 || (curl -Ls "https://get.maestro.mobile.dev" | bash)
	@echo "✅ Dependencies installed"

# Run all tests
test:
	./scripts/run-ios-e2e.sh

# Run specific test flow
test-flow:
	@if [ -z "$(name)" ]; then \
		echo "❌ Error: Please specify flow name with name=<flow>"; \
		echo "Available flows:"; \
		ls -1 flows/*.yaml | xargs -n 1 basename | sed 's/.yaml//'; \
		exit 1; \
	fi
	./scripts/run-ios-e2e.sh --flow $(name)

# Run tests in debug mode
test-debug:
	./scripts/run-ios-e2e.sh --debug

# Clean build and test
test-clean:
	./scripts/run-ios-e2e.sh --clean

# Quick test (skip build)
test-quick:
	./scripts/run-ios-e2e.sh --no-build

# Setup backend only
setup:
	./scripts/setup-backend.sh

# Clean up everything
clean:
	./scripts/cleanup.sh --volumes
	@echo "Cleaning iOS build artifacts..."
	rm -rf ../src/ios/PlatePilot/build
	rm -rf ~/Library/Developer/Xcode/DerivedData/PlatePilot-*
	@echo "✅ Cleanup complete"

# CI test (for GitHub Actions)
ci-test:
	./scripts/ci-test.sh
