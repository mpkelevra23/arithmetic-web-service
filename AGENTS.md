# Agent Guidelines

This repository contains a distributed service for evaluating arithmetic expressions written in Go. Below are recommendations and rules that agents should follow when working on this project.

## General approach

- All source files reside in the `internal` package and are organized by functionality (`agent`, `orchestrator`, `handler`, etc.). When adding new packages, keep to the existing structure.
- Format the code with `go fmt ./...` before committing.
- Run static analysis with `go vet ./...` and execute the tests with `go test ./...`. Commits must not break the tests.
- Write commit messages in Russian that concisely describe the changes.
- New functions and exported types must include comments in Russian.
- The project targets Go 1.20+. Ensure your changes are compatible with the version specified in `go.mod`.

## Logging and configuration

- The project uses the `zap` logging library. If new logging is required, use this library and avoid `fmt.Printf`.
- Application configuration is stored in the `.env` file and loaded via the `config` package. When adding new parameters, update the README example and keep backward compatibility.

## Verification and testing

```bash
# Code formatting
go fmt ./...

# Static analysis
go vet ./...

# Run all tests
go test ./...
```

Please ensure all of these commands run without errors before submitting changes.
