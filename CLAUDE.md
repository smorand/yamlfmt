# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

yamlfmt is an extensible Go CLI tool and library for formatting YAML files. It uses a pluggable formatter architecture with a custom fork of yaml.v3 as the default formatter.

## Essential Commands

```bash
# Build
make build                    # Build to dist/yamlfmt with version info
make install                  # Install locally

# Test
make test                     # Run all unit tests
make test_v                   # Verbose tests
go test -run TestName ./...   # Run single test

# Integration tests (requires building first)
make integrationtest          # Run all integration tests
make integrationtest TESTNAME="pattern"  # Filter by test name
make integrationtest_update   # Update golden files
make command_test_case TESTNAME="name"   # Create new test case

# Lint/Vet
make vet                      # Run go vet
make addlicense_check         # Check Apache 2.0 license headers
make addlicense               # Add license headers

# Run the tool
./dist/yamlfmt -dry .         # Dry run (show diffs)
./dist/yamlfmt -lint .        # Check formatting (exit 1 if diffs)
cat file.yaml | ./dist/yamlfmt -in  # Format stdin
./dist/yamlfmt -debug paths,config .  # Debug output
```

## Architecture

### Core Components

- **`/cmd/yamlfmt/`** - CLI entry point, flag parsing, config loading
- **`/command/`** - Command orchestration, config merging, operation modes
- **`/engine/`** - File processing engine (consecutive), output formatting
- **`/formatters/basic/`** - Default formatter using custom yaml.v3 fork
- **`/formatters/kyaml/`** - Alternative KYAML-based formatter
- **`/pkg/yaml/`** - Custom fork of yaml.v3 (excluded from standard vet/test)

### Key Patterns

**Factory/Registry Pattern** (`formatter.go`): Formatters register via `Factory` interface, looked up from `Registry` at runtime.

**Feature System** (`feature.go`): Composable before/after hooks for formatting operations via `FeatureFunc` and `FeatureList`.

**Path Collection** (`path_collector.go`): Three modes - standard (recursive walk), doublestar (glob patterns), gitignore (respects .gitignore).

**Operation Modes**: Format (default), Lint (`-lint`), Dry Run (`-dry`), Stdin (`-in`), Print Config (`-print_conf`).

### Config Priority

1. `-conf` flag (explicit path)
2. Working directory (`.yamlfmt`, `yamlfmt.yml`, `yamlfmt.yaml`, `.yamlfmt.yaml`, `.yamlfmt.yml`)
3. Parent directories (walking up)
4. System config (`~/.config/yamlfmt/.yamlfmt`)
5. Defaults

## Testing

**Unit tests**: Standard `*_test.go` files throughout. Note: `pkg/yaml` is excluded from normal test runs.

**Integration tests** (`integrationtest/command/`):
- Build-tag gated with `//go:build integration_test`
- Requires `YAMLFMT_BIN` env var pointing to built binary
- Golden file testing with `testdata/{test_name}/before/`, `after/`, `stdout/` directories
- Use `make integrationtest_update` to regenerate golden files

## Important Notes

- The `pkg/yaml` package is a custom fork excluded from `make vet` and `make test` to avoid conflicts
- Integration tests require a pre-built binary via `make build` first
- Doublestar patterns need shell quoting to prevent expansion
- Debug codes: `paths`, `config`, `diffs` (diffs is performance-intensive)
- Line endings default to platform-specific (crlf on Windows, lf elsewhere)

## Fork Information (smorand/yamlfmt)

This fork adds the following features on top of upstream (google/yamlfmt):

### `force_block_scalar_style` option

Forces multiline strings to use block scalar style (literal `|` or folded `>`).

```yaml
# .yamlfmt config
formatter:
  force_block_scalar_style: literal  # or "folded"
```

**Files added/modified:**
- `formatters/basic/features/force_block_scalar.go` - Feature implementation
- `formatters/basic/config.go` - Config option
- `formatters/basic/features.go` - Feature wiring

### Git Workflow

- `origin` → `google/yamlfmt` (upstream)
- `smo` → `smorand/yamlfmt` (fork)
- `main` branch tracks upstream, `smo` branch contains our modifications

```bash
# Sync with upstream
git checkout main && git pull origin main
git checkout smo && git rebase main
git push smo smo --force-with-lease
```
