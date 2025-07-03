# Mess Manager (messmgr)

A CLI tool that helps developers manage messy projects with multiple repositories and applications.

## Overview

Mess Manager is designed to help developers deal with complex project structures that involve multiple repositories and applications, such as:

- Systems with separated frontend and backend applications
- Micro-services backends
- Modular systems

## Installation

### Prerequisites

- Go 1.18 or later
- Git (for repository cloning)

### Building from Source

```bash
git clone <this-repository>
cd messmgr
go build -o messmgr .
```

You can then move the `messmgr` binary to a directory in your `PATH` for global access.

## Configuration

Mess Manager uses a `mess.config` JSON file to define your project structure. This file should be placed in your project root directory.

### Configuration Structure

```json
{
  "name": "your-project-name",
  "repos": [
    {
      "name": "unique-repo-name",
      "url": "https://github.com/example/repo.git"
    }
  ],
  "applications": [
    {
      "name": "unique-app-name",
      "repos": ["repo-name-1", "repo-name-2"],
      "scripts": {
        "script-name": "command-to-execute",
        "parallel-script": ["command-1", "command-2"]
      }
    }
  ]
}
```

### Configuration Fields

- **name**: Project name (required)
- **repos**: Array of repository definitions
  - **name**: Unique repository name (required)
  - **url**: Git repository URL (required)
- **applications**: Array of application definitions
  - **name**: Unique application name (required)
  - **repos**: Array of repository names that this application depends on
  - **scripts**: Dictionary of script names and their commands
    - Script values can be either a string (single command) or array of strings (parallel commands)

## Commands

### Global Flags

- `-f, --file <path>`: Specify custom config file path (default: `mess.config`)

### Initialize Project

```bash
# Create a sample mess.config file
messmgr init
```

### Repository Management

```bash
# Add a new repository
messmgr repo add <repo-name> <repo-url>

# Remove a repository
messmgr repo remove <repo-name>
messmgr repo rm <repo-name>          # alias

# Clone a repository
messmgr repo get <repo-name>
messmgr repo clone <repo-name>       # alias
```

### Application Management

```bash
# Add a new application
messmgr app add <app-name>
messmgr application add <app-name>   # alias

# Link an application with a repository
messmgr app link <app-name> <repo-name>
messmgr application link <app-name> <repo-name>  # alias

# Setup an application (clone repos and create symlinks)
messmgr app setup <app-name>
messmgr application setup <app-name>  # alias

# Run a script for an application
messmgr app run <app-name> <script-name>
messmgr application run <app-name> <script-name>  # alias
```

## Directory Structure

When you use Mess Manager, it creates the following directory structure in your project:

```
your-project/
├── mess.config
├── repos/
│   ├── frontend/          # Cloned repositories
│   ├── backend/
│   └── api/
└── applications/
    ├── web-app/
    │   ├── frontend/      # Symbolic links to repos
    │   └── backend/
    └── mobile-app/
        └── api/
```

## Usage Examples

### 1. Initialize a New Project

```bash
# Create initial configuration
messmgr init

# This creates a mess.config with sample data
```

### 2. Set Up a Multi-Repository Project

```bash
# Add repositories
messmgr repo add frontend https://github.com/company/frontend.git
messmgr repo add backend https://github.com/company/backend.git
messmgr repo add api https://github.com/company/api.git

# Create application
messmgr app add web-app

# Link repositories to application
messmgr app link web-app frontend
messmgr app link web-app backend

# Setup the application (clone repos and create symlinks)
messmgr app setup web-app
```

### 3. Add Scripts and Run Them

Edit your `mess.config` to add scripts:

```json
{
  "name": "my-project",
  "repos": [...],
  "applications": [
    {
      "name": "web-app",
      "repos": ["frontend", "backend"],
      "scripts": {
        "start": "npm run dev",
        "build": ["npm run build:frontend", "npm run build:backend"],
        "test": "npm test"
      }
    }
  ]
}
```

Then run scripts:

```bash
# Run a single command
messmgr app run web-app start

# Run parallel commands
messmgr app run web-app build
```

### 4. Working with Different Config Files

```bash
# Use a different config file
messmgr -f ./configs/production.config init
messmgr -f ./configs/production.config repo add prod-db https://github.com/company/prod-db.git
```

## How It Works

1. **Repository Management**: Repositories are cloned to `repos/<repo-name>` directories
2. **Application Setup**: When you run `app setup`, it:
   - Clones any missing repositories linked to the application
   - Creates symbolic links in `applications/<app-name>/` pointing to the corresponding repositories
3. **Script Execution**: Scripts run in the application directory (`applications/<app-name>/`) where symlinks provide access to all linked repositories

## Error Handling

- Duplicate repository/application names are prevented
- Missing repositories/applications are detected and reported
- Repository removal checks for usage in applications and prompts for confirmation
- Clear error messages guide users to fix configuration issues

## Contributing

This tool was built using:
- Go programming language
- [Cobra](https://github.com/spf13/cobra) for CLI framework
- Standard library for JSON parsing and file operations

## License

[Add your license information here] 