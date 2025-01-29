# k8spreview

A terminal-based Kubernetes YAML resource viewer with an interactive TUI interface.

## Features

- Interactive list view of all Kubernetes resources in a YAML file
- Detailed view of individual resources with YAML representation
- Resource relationship visualization
- Interactive graph view showing connections between resources
- Filtering capabilities to quickly find resources
- Color-coded resource types for better visibility
- Keyboard-based navigation

## Installation

### Prerequisites

- Go 1.21 or later (for building from source)
- Git (for building from source)

### Using Go Install

The easiest way to install k8spreview is using Go:

```bash
go install github.com/johnoct/k8spreview@latest  # Latest version
# or
go install github.com/johnoct/k8spreview@v0.1.1  # Specific version
```

### Binary Releases

Download pre-built binaries from the [releases page](https://github.com/johnoct/k8spreview/releases).

### Building from Source

1. Clone the repository:
```bash
git clone https://github.com/johnoct/k8spreview.git
cd k8spreview
```

2. Build the binary:
```bash
go build -o k8spreview cmd/main.go
```

3. (Optional) Install globally:
```bash
go install
```

### Development Setup

1. Install dependencies:
```bash
go mod download
```

2. Run tests:
```bash
go test ./...
```

3. Build and run with hot reload (requires [air](https://github.com/cosmtrek/air)):
```bash
air -c .air.toml
```

## Usage

```bash
# View a YAML file
k8spreview <yaml-file>

# Show version information
k8spreview -version
```

### Navigation

- Use arrow keys to navigate the list
- Press `Enter` to view resource details
- Press `g` to view the resource graph
- Press `/` to filter resources
- Press `q` to go back or quit

## Example

```bash
# View a complex example with multiple resources and relationships
k8spreview examples/multi-resource.yaml
```

See the [examples](./examples) directory for more sample YAML files.

## Project Structure

```
.
├── cmd/
│   └── main.go           # Main application entry point
├── pkg/
│   ├── k8s/             # Kubernetes resource handling
│   │   ├── k8s.go       # Core resource types and functions
│   │   └── doc.go       # Package documentation
│   ├── ui/              # TUI components and styling
│   │   ├── app.go       # Application entry point
│   │   ├── model.go     # UI state and update logic
│   │   ├── styles.go    # UI styling definitions
│   │   └── doc.go       # Package documentation
│   └── version/         # Version information
│       └── version.go   # Version, commit, and build date
├── docs/                # Documentation
├── examples/            # Example YAML files
└── tests/              # Test files
```

## Development

### Running Tests

Run all tests:
```bash
go test ./...
```

Run tests with coverage:
```bash
go test -cover ./...
```

Generate coverage report:
```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Code Style

The project follows standard Go code style. Before committing, ensure your code:

1. Is formatted with `gofmt`:
```bash
gofmt -s -w .
```

2. Passes static analysis:
```bash
go vet ./...
```

3. Has no linting errors (requires [golangci-lint](https://golangci-lint.run/)):
```bash
golangci-lint run
```

### Release Process

1. Update version in `pkg/version/version.go`
2. Create and push a new tag:
```bash
git tag -a v0.1.1 -m "Release v0.1.1"
git push origin v0.1.1
```
3. GitHub Actions will automatically build and publish the release

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - Terminal UI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) - Style definitions
- [YAML v3](https://github.com/go-yaml/yaml) - YAML parsing
- [GoReleaser](https://goreleaser.com/) - Release automation

