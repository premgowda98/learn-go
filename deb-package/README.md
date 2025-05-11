# Gowda - A Simple Message Printer

This is a simple Go program packaged as a Debian package that prints the message provided as a command-line argument.

## Usage

```
gowda <message>
```

Example:
```
gowda "Hello, world!"
```

## Building from Source

To build the Go program from source:

```bash
go build -o gowda
```

## Understanding Debian Packages

A Debian package (.deb) is the standard software packaging format used by Debian and Debian-based Linux distributions like Ubuntu. It contains:

1. The compiled program
2. Metadata about the package
3. Installation instructions
4. Scripts to run during installation/removal

### Key Files in a Debian Package

The `debian/` directory contains all the configuration files needed to build a Debian package:

- **control**: Contains metadata about the package, including:
  - Package name, version, and architecture
  - Dependencies
  - Description
  - Maintainer information

  **Detailed Explanation of the Control File's Dependencies:**
  
  In the `control` file, you'll typically see this important line:
  ```
  Depends: ${shlibs:Depends}, ${misc:Depends}
  ```
  
  These are **substitution variables** that get automatically resolved during the package build process:
  
  - `${shlibs:Depends}`: Automatically detects shared library dependencies that your executables need
    - The `dpkg-shlibdeps` tool analyzes your binary and determines which libraries it uses
    - For Go programs, this may be minimal if statically linked, or more extensive if using CGO
    - Example output: `libc6 (>= 2.29), libstdc++6 (>= 9.3.0)`
  
  - `${misc:Depends}`: Collects dependencies required by debhelper scripts used in your package build
    - Added by various debhelper commands when they need specific packages to function
    - Handles special cases like dependencies on particular package versions
    - May be empty in simple packages

  This automatic dependency resolution ensures your package will work on the target system without manually tracking libraries, making your package more maintainable and accurate.

- **rules**: A Makefile that specifies how to build and install the package.
  - Tells the build system how to compile your software
  - Handles the installation of files to the correct locations
  
  **Detailed Explanation of the Rules File:**
  
  The `rules` file is essentially a Makefile that controls how the package is built and installed. Here's what our current rules file does:
  
  ```bash
  #!/usr/bin/make -f
  
  %:
      dh $@
  
  override_dh_auto_build:
      go build -o gowda
  
  override_dh_auto_install:
      install -D -m 0755 gowda $(CURDIR)/debian/gowda/usr/bin/gowda
  ```
  
  - `#!/usr/bin/make -f`: This is a shebang line indicating this is a makefile to be processed by GNU make
  
  - `%: dh $@`: This is a catch-all rule that forwards all targets to the `dh` (debhelper) command. The `dh` command is a helper that will execute a sequence of debhelper commands for each stage of the package build process.
  
  - `override_dh_auto_build:`: This overrides the automatic build step with our custom commands
    - `go build -o gowda`: Builds the Go program and names the executable 'gowda'
  
  - `override_dh_auto_install:`: This overrides the automatic installation step
    - `install -D -m 0755 gowda $(CURDIR)/debian/gowda/usr/bin/gowda`: Installs the gowda executable to usr/bin in the package with file permissions 0755 (readable and executable by everyone, writable by owner)
    - `$(CURDIR)` is a variable that represents the current directory

- **changelog**: Records the version history of the package.
  - Each entry includes version, distribution, changes, and maintainer details
  - The most recent entry determines the package version

- **copyright**: Contains copyright and license information.

- **compat**: Specifies the debhelper compatibility level.

- **source/format**: Defines the source package format (commonly 3.0 (quilt)).

### Package Maintainer Scripts

Debian packages can include special scripts that run at different points during package installation, upgrade, and removal. These are known as maintainer scripts:

1. **preinst**: Runs BEFORE a package is installed
2. **postinst**: Runs AFTER a package is installed
3. **prerm**: Runs BEFORE a package is removed
4. **postrm**: Runs AFTER a package is removed

**Example preinst script** (debian/preinst):
```bash
#!/bin/sh
set -e

# Create a backup of any existing configuration if it exists
if [ -f /etc/gowda/config.conf ]; then
    echo "Backing up existing configuration..."
    cp /etc/gowda/config.conf /etc/gowda/config.conf.backup
fi

# Make sure the required directory exists
if [ ! -d /var/log/gowda ]; then
    echo "Creating log directory..."
    mkdir -p /var/log/gowda
    chmod 755 /var/log/gowda
fi

# Exit with success
exit 0
```

**Example postinst script** (debian/postinst):
```bash
#!/bin/sh
set -e

# Configure the package after installation
case "$1" in
    configure)
        # Create user if it doesn't exist
        if ! getent passwd gowda > /dev/null; then
            echo "Creating gowda user..."
            useradd --system --no-create-home --home-dir /nonexistent gowda
        fi

        # Set permissions on log directory
        if [ -d /var/log/gowda ]; then
            chown gowda:gowda /var/log/gowda
        fi
        
        # Create default config if it doesn't exist
        if [ ! -f /etc/gowda/config.conf ]; then
            echo "Creating default configuration..."
            mkdir -p /etc/gowda
            echo "# Default configuration for gowda" > /etc/gowda/config.conf
            echo "LOG_LEVEL=info" >> /etc/gowda/config.conf
            chown root:root /etc/gowda/config.conf
            chmod 644 /etc/gowda/config.conf
        fi

        # Restart service if it was running (for upgrades)
        if [ -x "/etc/init.d/gowda" ]; then
            if [ -x "$(command -v invoke-rc.d)" ]; then
                invoke-rc.d gowda restart || true
            else
                /etc/init.d/gowda restart || true
            fi
        fi
    ;;

    abort-upgrade|abort-remove|abort-deconfigure)
        # Nothing to do here
    ;;

    *)
        echo "postinst called with unknown argument \`$1'" >&2
        exit 1
    ;;
esac

# Exit with success
exit 0
```

**Example prerm script** (debian/prerm):
```bash
#!/bin/sh
set -e

# Clean up before package removal
case "$1" in
    remove|upgrade|deconfigure)
        # Stop the service if it's running
        if [ -x "/etc/init.d/gowda" ]; then
            if [ -x "$(command -v invoke-rc.d)" ]; then
                invoke-rc.d gowda stop || true
            else
                /etc/init.d/gowda stop || true
            fi
        fi
    ;;

    failed-upgrade)
        # Nothing to do
    ;;

    *)
        echo "prerm called with unknown argument \`$1'" >&2
        exit 1
    ;;
esac

# Exit with success
exit 0
```

**Example postrm script** (debian/postrm):
```bash
#!/bin/sh
set -e

# Clean up after package removal
case "$1" in
    purge)
        # Remove all configuration and data
        if [ -d /etc/gowda ]; then
            echo "Removing configuration directory..."
            rm -rf /etc/gowda
        fi
        
        if [ -d /var/log/gowda]; then
            echo "Removing log directory..."
            rm -rf /var/log/gowda
        fi
        
        # Remove the user if it exists
        if getent passwd gowda > /dev/null; then
            echo "Removing gowda user..."
            userdel gowda || true
        fi
    ;;
    
    remove|upgrade|failed-upgrade|abort-install|abort-upgrade|disappear)
        # Nothing special to do
    ;;

    *)
        echo "postrm called with unknown argument \`$1'" >&2
        exit 1
    ;;
esac

# Exit with success
exit 0
```

These scripts handle various aspects of the installation, upgrade, and removal process:

- **preinst**: Prepares the system before installing the package (creates directories, backs up existing files)
- **postinst**: Configures the package after installation (creates users, sets permissions, creates default configs)
- **prerm**: Prepares for package removal (stops services)
- **postrm**: Cleans up after package removal (removes configuration, log files, users)

When these scripts are included in your package, the Debian package system will automatically run them at the appropriate times during installation, upgrade, and removal.

### Benefits of Using Debian Packages

- System-wide installation with proper file placement
- Dependency management
- Easy installation, upgrade, and removal
- Integration with system package management

## Creating a Debian Package

### Prerequisites

You need to have the following packages installed:

```bash
sudo apt update
sudo apt install build-essential debhelper devscripts dh-make
```

### Building the Debian Package

1. Navigate to the project directory:
   ```bash
   cd /path/to/gowda
   ```

2. Build the Debian package:
   ```bash
   dpkg-buildpackage -b -us -uc
   ```
   
   This command:
   - `-b`: Build a binary-only package (no source package)
   - `-us`: Do not sign the source package
   - `-uc`: Do not sign the .changes file
   
   The process:
   1. Reads the `debian/control` file for package information
   2. Executes the build rules in `debian/rules`
   3. Creates a `.deb` package in the parent directory

3. The .deb package will be created in the parent directory.

## Installing the Package

### From a Local .deb File

You can install the package using either of these methods:

#### Method 1: Using apt (Recommended)
```bash
sudo apt install ./gowda_1.0.0_amd64.deb
```
This method:
- Automatically resolves and installs any dependencies
- Requires an internet connection to download dependencies
- The `./` prefix is important to tell apt this is a local file

#### Method 2: Using dpkg
```bash
sudo dpkg -i ../gowda_1.0.0_amd64.deb
```
This method:
- Directly installs the package without automatic dependency resolution
- If there are missing dependencies, installation may fail
- Works offline, but only if all dependencies are already installed

If you encounter dependency issues with the dpkg method, fix them with:
```bash
sudo apt-get install -f
```

## Go Install vs. Debian Package

### Go Install Method

Using `go install` installs Go programs in your Go path (usually `$HOME/go/bin`):

```bash
# From the project directory
go install .

# Or directly from a remote repository
go install github.com/username/gowda@latest
```

Benefits of `go install`:
- Simple and quick for Go developers
- Install directly from source repositories
- Installed in your user's Go bin directory, no system-wide changes

Limitations of `go install`:
- Only works for Go programs
- Doesn't manage non-Go dependencies
- Not integrated with the system package manager
- Files are only installed to Go-specific locations
- No standard way to remove installed packages

### Debian Package Method

Benefits over `go install`:
- System-wide installation with proper file placement in standard directories
- Integration with system package management (apt)
- Proper dependency handling
- Easy to distribute to non-Go users
- Standardized installation/removal process
- Can include additional files, documentation, and configuration

## Uninstalling the Package

```bash
sudo apt remove gowda
# To remove configuration files as well
sudo apt purge gowda
```

## Additional Resources

- [Debian New Maintainers' Guide](https://www.debian.org/doc/manuals/maint-guide/)
- [Debian Packaging Tutorial](https://wiki.debian.org/Packaging/Intro)
- [Go Packaging Best Practices](https://go.dev/doc/modules/developing)

## Privileged vs. Normal Debian Builds

When building Debian packages, you have two main approaches: privileged builds and normal (non-privileged) builds. Each has its advantages and use cases.

### Privileged Builds

Privileged builds are executed with root privileges and are typically done in a controlled environment.

**Characteristics:**
- Run as the root user
- Can interact with system-wide resources
- Often performed in a dedicated build environment (like pbuilder, sbuild, or a CI/CD pipeline)
- Used for official package distribution

**Example (using pbuilder):**
```bash
# Create a build environment
sudo pbuilder create --distribution bullseye

# Build the package
sudo pbuilder build ../gowda_1.0.0.dsc
```

**Advantages:**
- More accurate simulation of the real package installation environment
- Can test root-required operations in maintainer scripts
- Better isolation from the host system
- Controlled dependencies
- Reproducible builds

**Disadvantages:**
- Requires root access
- More complex setup
- Slower than normal builds

### Normal (Non-Privileged) Builds

Normal builds are executed as a regular user without root privileges.

**Example:**
```bash
# Build from the project directory
dpkg-buildpackage -b -us -uc -rfakeroot
```

The `-rfakeroot` option is crucial here - it uses the `fakeroot` utility to simulate root privileges for file operations without actually requiring them.

**Advantages:**
- Safer (no root privileges required)
- Simpler and faster setup
- Quicker iteration during development
- Can be run by any user without sudo access

**Disadvantages:**
- Cannot fully test all aspects of package installation
- May not catch all potential issues with maintainer scripts
- Less isolated from the host system

### Best Practices

1. **Development Workflow:**
   - Use normal builds with fakeroot during development for quick iteration
   - Use `lintian` to check your package for errors and policy violations:
     ```bash
     lintian ../gowda_1.0.0_amd64.changes
     ```

2. **Testing and Release:**
   - Test in a clean, privileged build environment before release
   - Consider using tools like pbuilder or sbuild for final testing
   - For maximum compatibility, test on the oldest supported distribution version

3. **Continuous Integration:**
   - Set up CI/CD pipelines to build packages in clean environments
   - Test installation, upgrade, and removal in isolated containers

### When to Use Each Method

- **Use normal builds with fakeroot when:**
  - Developing and debugging your package
  - Making frequent changes and needing quick feedback
  - Working on your personal developer machine

- **Use privileged builds when:**
  - Preparing for official distribution
  - Final testing before release
  - Building in CI/CD environments
  - Testing complex maintainer scripts that interact with system services