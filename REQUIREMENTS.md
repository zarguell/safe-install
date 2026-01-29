Renaming the tool to `safe-install` is a smart strategic move. It transforms the project from a "Python patch" into a **generic "Safety Harness" for any package manager**.

This modular architecture allows you to easily add `safe-install npm`, `safe-install gem`, or `safe-install cargo` later using the same config/policy engine.

Here is the **Revised Full Specification** designed for an AI Developer Agent.

***

# Project Specification: Safe-Install

**Tool Name:** `safe-install`
**Purpose:** A unified, cross-platform CLI wrapper that enforces security policies (like internal registries and lockfiles) across multiple package managers (pip, npm, etc.).
**Initial Scope:** Support `pip` (Python) immediately. Architecture must support `npm` (Node) as a future module.

## 1. High-Level Architecture
The tool uses a **Subcommand Pattern**.
*   `safe-install pip [args...]` -> Wraps Python `pip`
*   `safe-install npm [args...]` -> Wraps Node `npm` (Future)
*   `safe-install config` -> Shows current policy

### Directory Structure (Go)
```text
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

## 2. Detailed Functional Requirements

### A. Configuration (The "Policy Engine")
The tool must load a single JSON file that defines rules for *all* supported managers.

**Path:**
*   **Windows:** `C:\ProgramData\SafeInstall\config.json`
*   **Linux:** `/etc/safe-install/config.json`
*   **macOS:** `/Library/Application Support/SafeInstall/config.json`

**Schema (`config.json`):**
```json
{
  "common": {
    "block_interactive": true  // Always inject --no-input / --yes
  },
  "managers": {
    "pip": {
      "enabled": true,
      "require_venv": true,
      "registry_url": "https://artifactory.local/api/pypi/simple",
      "binary_path": "python" // or "pip3"
    },
    "npm": {
      "enabled": false,
      "registry_url": "https://artifactory.local/api/npm/npm-repo"
    }
  }
}
```

### B. The "Pip" Module (`safe-install pip`)
When the user runs `safe-install pip install requests`, the tool must:
1.  **Load Config:** Read the JSON. If `managers.pip.enabled` is false, abort.
2.  **Check Venv:** If `require_venv` is true, check `os.Getenv("VIRTUAL_ENV")`.
3.  **Sanitize:** Scan arguments. Fail if user provides `--index-url` or `--extra-index-url`.
4.  **Execute:** Run the underlying command with injected safety flags.
    *   *Input:* `safe-install pip install requests`
    *   *Executed:* `python -m pip install --no-input --index-url <URL> --no-index requests`

## 3. Implementation Plan (Agent Task List)

Copy this block to your AI Developer Agent:

```markdown
# Task: Build "Safe-Install" CLI

## Phase 1: Foundation & CLI Framework
- [ ] Initialize module: `go mod init safe-install`.
- [ ] Add **Cobra** library: `go get -u github.com/spf13/cobra` (Standard for Go CLIs).
- [ ] Setup `cmd/root.go`: Create the base command `safe-install`.
- [ ] Setup `internal/config`:
  - Create `Config` struct supporting `pip` and `npm` sections.
  - Implement OS-specific config path detection.
  - Fail fast if config is missing/invalid.

## Phase 2: The "Pip" Subcommand
- [ ] Create `cmd/pip.go`: Define the `pip` subcommand.
- [ ] Implement `internal/policy/pip_policy.go`:
  - Function `ValidatePipEnv()`: Check `VIRTUAL_ENV`.
  - Function `BuildPipArgs()`: Inject `--index-url`, `--no-index`, `--no-input`.
- [ ] Wire it up: `safe-install pip <args>` should run the transformed command.
- [ ] **Stream Output:** Ensure `stdout/stderr` form the subprocess are piped directly to the user (so AI agents see progress).

## Phase 3: Build & Release
- [ ] Create `.github/workflows/release.yml`:
  - Build for Windows (`safe-install.exe`), Linux (`safe-install-linux`), Mac (`safe-install-darwin`).
  - Trigger on `v*` tags.
- [ ] Create `README.md`:
  - Documentation for Admins (How to create `config.json`).
  - Documentation for Users (Usage examples).
  - Documentation for AI Agents (System Prompt / Settings snippet).

## Phase 4: Future Proofing (Stub)
- [ ] Create `cmd/npm.go`: A stub command that prints "NPM support coming soon" (proving modularity).
```

## 4. Documentation for "Claude Code"
This section goes in your `README.md` to help future you (or other devs) set this up for their agents.

### Integration with Claude Code
To force Claude to use this tool, update your `settings.json`:

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
*Note: This configuration forces the Agent to use your wrapper for both pip and (eventually) npm, effectively closing the "direct access" loophole.*
