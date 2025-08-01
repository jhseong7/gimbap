# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## About GIMBAP

GIMBAP is a Go web framework focused on dependency injection and modular architecture. It provides automatic dependency injection management, flexible server engine switching (Gin, Fiber, Echo), and microservice support. The framework is currently in early development stage.

## Common Commands

### Development
- `go mod tidy` - Update dependencies
- `go test ./...` - Run all tests (uses Ginkgo/Gomega testing framework in dependency package)
- `go build -o gimbap ./sample-cmd/app` - Build sample application
- `go run ./sample-cmd/app` - Run sample application with Fiber engine
- `go run ./sample-cmd/https` - Run HTTPS sample application

### Testing
Tests are written using Ginkgo and Gomega testing frameworks. Only the `dependency` package currently has tests. Run specific package tests with:
- `go test ./dependency` - Run dependency injection tests

## Architecture Overview

GIMBAP follows a modular dependency injection pattern with these core concepts:

### Core Components
- **App Container**: Main application container that orchestrates all components
- **Modules**: Logical groupings of providers and dependencies
- **Providers**: Services/components that can be injected into other components
- **Controllers**: HTTP request handlers that implement REST endpoints
- **Server Engines**: Pluggable HTTP server implementations (Gin, Fiber, Echo)
- **Microservices**: Background services managed by the application container
- **Dependency Manager**: Handles automatic dependency injection using reflection

### Key Files
- `gimbap.go`: Main package exports and type aliases
- `app/app.go`: Core application container implementation
- `dependency/`: Dependency injection management with multiple backends (native, Fx, Gimbap)
- `engine/`: HTTP server engine abstractions and implementations
- `sample/`: Example implementations showing framework usage patterns
- `sample-cmd/`: Runnable example applications

### Dependency Injection Pattern
The framework uses constructor function injection where dependencies are automatically resolved based on function signatures. Controllers and services define their dependencies as constructor parameters, and the framework automatically provides them at runtime.

### Engine Switching
Applications can switch between different HTTP servers (Gin, Fiber, Echo) by changing the ServerEngine in AppOption. Each engine has its own controller implementations in the sample directories.

### Module Structure
The framework organizes code into modules that declare their providers and dependencies. Modules can depend on other modules, creating a dependency graph that the injection system resolves.