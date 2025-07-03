# Mess Manager (mess)

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
cd mess
go build -o mess .
```

You can then move the `mess` binary to a directory in your `PATH` for global access.

## Configuration

Mess Manager uses a `mess.json` JSON file to define your project structure. This file should be placed in your project root directory.

### Configuration Structure

```json
{
  "name": "your-project-name",
  "repos": [
    {
      "name": "unique-repo-name",
      "url": "https://github.com/example/repo.git",
      "clone_params": "--depth=1"
    }
  ],
  "applications": [
    {
      "name": "unique-app-name",
      "repos": ["repo-name-1", "repo-name-2"],
      "scripts": {
        "script-name": "command-to-execute",
        "parallel-script": ["command-1", "command-2"]
      },
      "env": {
        "NODE_ENV": "development",
        "API_URL": "http://localhost:3000"
      },
      "pre-setup": "npm install",
      "post-setup": "npm run build"
    }
  ]
}
```

### Configuration Fields

- **name**: Project name (required)
- **repos**: Array of repository definitions
  - **name**: Unique repository name (required)
  - **url**: Git repository URL (required)
  - **clone_params**: Optional additional parameters for git clone command
- **applications**: Array of application definitions
  - **name**: Unique application name (required)
  - **repos**: Array of repository names that this application depends on
  - **scripts**: Dictionary of script names and their commands
    - Script values can be either a string (single command) or array of strings (parallel commands)
  - **env**: Optional dictionary of environment variables (key-value pairs)
  - **pre-setup**: Optional script command to run before setup
  - **post-setup**: Optional script command to run after setup

## Commands

### Global Flags

- `-f, --file <path>`: Specify custom config file path (default: `mess.json`)

### Initialize Project

```bash
# Create a sample mess.json file
mess init
```

### Repository Management

```bash
# Add a new repository
mess repo <repo-name> add <repo-url>

# Remove a repository
mess repo <repo-name> remove
mess repo <repo-name> rm          # alias

# Clone a repository
mess repo <repo-name> get

# Execute git commands on a repository
mess repo <repo-name> <git-command>
# Example: mess repo frontend status
# Example: mess repo backend pull
```

### Application Management

```bash
# Add a new application
mess application <app-name> init
mess app <app-name> init           # alias

# Link an application with repositories
mess application <app-name> link <repo-name> [repo-name...]
mess app <app-name> link <repo-name> [repo-name...]  # alias

# Setup an application (clone repos, create symlinks, run pre/post scripts)
mess application <app-name> setup
mess app <app-name> setup          # alias

# Clone application repositories and create symlinks
mess application <app-name> clone
mess app <app-name> clone          # alias

# Run a script for an application
mess application <app-name> run <script-name>
mess app <app-name> run <script-name>  # alias
```

## Directory Structure

When you use Mess Manager, it creates the following directory structure in your project:

```
your-project/
├── mess.json
├── repos/
│   ├── frontend/          # Cloned repositories
│   ├── backend/
│   └── api/
└── applications/          # Or custom path via MESS_APPLICATION_ROOT
    ├── web-app/
    │   ├── frontend/      # Symbolic links to repos
    │   └── backend/
    └── mobile-app/
        └── api/
```

## Environment Variables

- **MESS_APPLICATION_ROOT**: Custom path for applications directory (defaults to `{mess.json location}/applications`)

## Usage Examples

### 1. Initialize a New Project

```bash
# Create initial configuration
mess init

# This creates a mess.json with sample data
```

### 2. Set Up a Multi-Repository Project

```bash
# Add repositories
mess repo frontend add https://github.com/company/frontend.git
mess repo backend add https://github.com/company/backend.git
mess repo api add https://github.com/company/api.git

# Create application
mess app web-app init

# Link repositories to application
mess app web-app link frontend backend

# Setup the application (clone repos, create symlinks, run setup scripts)
mess app web-app setup
```

### 3. Add Scripts and Environment Variables

Edit your `mess.json` to add scripts and environment variables:

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
      },
      "env": {
        "NODE_ENV": "development",
        "API_URL": "http://localhost:3000"
      },
      "pre-setup": "npm install",
      "post-setup": "npm run build"
    }
  ]
}
```

Then run scripts:

```bash
# Run a single command
mess app web-app run start

# Run parallel commands
mess app web-app run build
```

### 4. Working with Git Commands

```bash
# Check status of a repository
mess repo frontend status

# Pull latest changes
mess repo backend pull

# Create a new branch
mess repo api checkout -b feature/new-feature
```

### 5. Working with Different Config Files

```bash
# Use a different config file
mess -f ./configs/production.json init
mess -f ./configs/production.json repo prod-db add https://github.com/company/prod-db.git
```

## How It Works

1. **Repository Management**: Repositories are cloned to `repos/<repo-name>` directories
2. **Application Setup**: When you run `app setup`, it:
   - Clones any missing repositories linked to the application
   - Creates the application directory in `MESS_APPLICATION_ROOT` (defaults to `applications/<app-name>/`)
   - Executes the `pre-setup` script if defined
   - Creates symbolic links in the application directory pointing to the corresponding repositories
   - Executes the `post-setup` script if defined
3. **Script Execution**: Scripts run in the application directory where symlinks provide access to all linked repositories
   - Single string commands are executed directly
   - Array of strings are executed in parallel as separate sub-processes
4. **Git Command Delegation**: Git commands are delegated directly to the `git` CLI for each repository

## Error Handling

- Duplicate repository/application names are prevented
- Missing repositories/applications are detected and reported
- Repository removal checks for usage in applications and prompts for confirmation
- Clear error messages guide users to fix configuration issues
- Validation ensures referenced repositories exist before linking to applications

## Contributing

This tool was built using:
- Go programming language
- [Cobra](https://github.com/spf13/cobra) for CLI framework
- Standard library for JSON parsing and file operations

## License

[Add your license information here] 