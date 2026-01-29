# Safe-Install

A unified, cross-platform CLI wrapper that enforces security policies (like internal registries and lockfiles) across multiple package managers (pip, npm, etc.).

## Features

- **Security First**: Enforces virtual environment requirements and blocks interactive prompts
- **Registry Enforcement**: Forces use of approved internal registries
- **Cross-Platform**: Runs on Windows, Linux, and macOS
- **Extensible**: Built to support multiple package managers (pip now, npm coming soon)
- **Agent-Friendly**: Designed to work seamlessly with AI development agents

## Installation

Download the latest release for your platform from the [Releases](https://github.com/yourusername/safe-install/releases) page.

Place the binary in your PATH and make it executable (on Unix-like systems):

```bash
chmod +x safe-install-linux
sudo mv safe-install-linux /usr/local/bin/safe-install
```

## Configuration

Safe-Install reads its configuration from a JSON file located at:

- **Windows**: `C:\ProgramData\SafeInstall\config.json`
- **Linux**: `/etc/safe-install/config.json`
- **macOS**: `/Library/Application Support/SafeInstall/config.json`

### Configuration Schema

```json
{
  "common": {
    "block_interactive": true
  },
  "managers": {
    "pip": {
      "enabled": true,
      "require_venv": true,
      "registry_url": "https://artifactory.local/api/pypi/simple",
      "binary_path": "python"
    },
    "npm": {
      "enabled": false,
      "registry_url": "https://artifactory.local/api/npm/npm-repo"
    }
  }
}
```

### Configuration Options

#### Common Options
- `block_interactive`: If true, always inject `--no-input` / `--yes` flags to block interactive prompts

#### Pip Options
- `enabled`: Enable or disable pip support
- `require_venv`: If true, require a virtual environment to be active
- `registry_url`: The internal registry URL to enforce
- `binary_path`: Path to Python/pip binary (e.g., "python" or "pip3")

#### NPM Options
- `enabled`: Enable or disable npm support (coming soon)
- `registry_url`: The internal registry URL to enforce

## Usage

### Basic Usage

```bash
safe-install pip install requests
```

This command will:
1. Load and validate the configuration
2. Check if a virtual environment is active (if required)
3. Block any attempts to use `--index-url` or `--extra-index-url`
4. Inject the configured registry URL and safety flags
5. Execute the pip command with all security policies applied

### Help

```bash
safe-install --help
safe-install pip --help
safe-install npm --help
```

## Integration with AI Agents

To force AI agents to use this tool instead of direct package manager commands, update your agent's settings:

```json
{
  "permissions": {
    "allow": [
      "Bash(safe-install *)"
    ],
    "deny": [
      "Bash(pip *)",
      "Bash(python -m pip *)",
      "Bash(npm *)"
    ]
  }
}
```

This configuration ensures that the agent can only interact with package managers through the safe-install wrapper, closing the "direct access" loophole and enforcing your organization's security policies.

## Building from Source

### Prerequisites
- Go 1.21 or later

### Build

```bash
go mod download
go build -o safe-install
```

### Cross-Platform Build

```bash
GOOS=linux GOARCH=amd64 go build -o safe-install-linux
GOOS=darwin GOARCH=amd64 go build -o safe-install-darwin
GOOS=windows GOARCH=amd64 go build -o safe-install.exe
```

## Architecture

The tool uses a **Subcommand Pattern**:

- `safe-install pip [args...]` → Wraps Python `pip`
- `safe-install npm [args...]` → Wraps Node `npm` (Coming Soon)
- `safe-install config` → Shows current policy (Coming Soon)

### Directory Structure

```
safe-install/
├── cmd/
│   ├── root.go       # Entry point (Cobra/CLI logic)
│   ├── pip.go        # "pip" subcommand definition
│   └── npm.go        # "npm" subcommand (stub for now)
├── internal/
│   ├── config/       # Config loading & validation
│   ├── policy/       # Rules engine (Venv check, Repo URL)
│   └── executor/     # Safe subprocess runner (Injects flags)
├── main.go           # Calls cmd.Execute()
└── README.md
```

## Security Features

1. **Virtual Environment Enforcement**: Prevents accidental system-wide package installations
2. **Registry Locking**: Enforces use of approved internal registries
3. **Interactive Prompt Blocking**: Prevents hanging on user input in automated environments
4. **Argument Sanitization**: Blocks dangerous flags like `--index-url` override attempts

## Contributing

Contributions are welcome! Please ensure all code follows existing patterns and conventions.

## License

[Your License Here]