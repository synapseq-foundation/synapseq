# Compilation

This guide covers how to compile SynapSeq from source on macOS, Linux, Windows, and for WebAssembly.

You need Go v1.26 or later and `make` installed before building.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Install Go and Make](#install-go-and-make)
- [Compile SynapSeq](#compile-synapseq)
- [Installing the Binary](#installing-the-binary)
- [Additional Make Targets](#additional-make-targets)

## Prerequisites

Before compiling, make sure the following tools are available:

- Go v1.26 or later
- `make`
- `git`

## Install Go and Make

### macOS

Install Go with Homebrew or MacPorts:

```bash
# Using Homebrew
brew install go

# Using MacPorts
sudo port install go
```

### Linux (Ubuntu/Debian)

Install Go and `make` with `apt`, or install a newer Go release with `snap`:

```bash
# Update package list
sudo apt update

# Install Go
sudo apt install golang-go make

# Or install a newer version using snap
sudo snap install go --classic
```

### Linux (CentOS/RHEL/Fedora)

Install the required packages with your system package manager:

```bash
# For Fedora
sudo dnf install golang make

# For CentOS/RHEL
sudo yum install golang make
```

### Windows

Use Git Bash instead of PowerShell or CMD, since the Makefile relies on Unix-like shell commands.

1. Install [Git for Windows](https://git-scm.com/download/win) (includes Git Bash).
   After installation, you will have both **Git Bash** and **PowerShell** available.

2. Install [Scoop](https://scoop.sh/).
   Open PowerShell and run:

   ```powershell
   Set-ExecutionPolicy -ExecutionPolicy RemoteSigned -Scope CurrentUser
   Invoke-RestMethod -Uri https://get.scoop.sh | Invoke-Expression
   ```

3. Install Go and `make` using Scoop.
   In PowerShell, run:

   ```powershell
   scoop update
   scoop install go make
   ```

4. Open Git Bash and verify that everything is available:

```bash
go version
make --version
```

## Compile SynapSeq

### Clone the repository

If you have not cloned the project yet:

```bash
git clone https://github.com/ruanklein/synapseq.git
cd synapseq
```

### Build for the current platform

On macOS and Linux, the default build target creates a binary for the current operating system and architecture:

```bash
make
```

This creates the output binary in the `bin/` directory.

### Build for Windows

Use the platform-specific Windows targets to preserve the `.exe` extension, application icon, and Windows-specific command-line behavior:

```bash
make build-windows-amd64    # Windows 64-bit (Intel/AMD) - Recommended
make build-windows-arm64    # Windows 64-bit (ARM)
```

Do not use `make build` on Windows, as it creates a binary without the `.exe` extension and without the Windows-specific resource metadata.

The Windows build automatically generates resource metadata such as icon and version info using `goversioninfo` during the build.

### Cross-compile for other platforms

You can build for different operating systems and architectures:

```bash
# Linux
make build-linux-amd64      # Linux 64-bit (Intel/AMD)
make build-linux-arm64      # Linux 64-bit (ARM)

# macOS
make build-macos            # macOS ARM64 (Apple Silicon)
```

### Build for WebAssembly

To compile SynapSeq for use in web browsers:

```bash
make build-wasm
```

This generates the following files in the `wasm/` directory:

- `synapseq.wasm` - The WebAssembly binary
- `wasm_exec.js` - The Go WASM runtime copied from the local Go installation

## Installing the Binary

After compilation, you can install the binary system-wide.

### macOS and Linux

```bash
sudo make install
```

This installs SynapSeq to `/usr/local/bin/synapseq`.

### Windows

Using Git Bash as Administrator:

```bash
mkdir -p "/c/Program Files/SynapSeq"
cp bin/synapseq-windows-amd64.exe "/c/Program Files/SynapSeq/synapseq.exe"
```

After copying the executable, add `C:\Program Files\SynapSeq` to your `PATH` environment variable.

1. Open **Start Menu** and search for "Environment Variables"
2. Click **Edit the system environment variables**
3. Click **Environment Variables...**
4. Under **User variables** or **System variables**, select **Path**
5. Click **Edit...**
6. Click **New**
7. Add `C:\Program Files\SynapSeq`
8. Confirm all dialogs with **OK**

Restart Git Bash or PowerShell and verify:

```bash
synapseq -h
```

## Additional Make Targets

Useful maintenance targets from the Makefile:

```bash
make test     # Run the Go test suite
make clean    # Remove build artifacts
```
