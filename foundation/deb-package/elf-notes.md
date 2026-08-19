# **Comprehensive Linux Notes: Dynamic Linking, RPATH, and Debian Packaging**

## **1. ELF Binaries and Dynamic Linking**

### **What is ELF?**
- **ELF (Executable and Linkable Format)** is the standard file format for executables, shared libraries, and object files on Linux
- Contains multiple sections including headers, code, data, and metadata
- **Dynamic linking** allows binaries to use shared libraries at runtime instead of embedding code statically

### **Types of Linking:**
```bash
# Static linking - everything embedded in binary
gcc -static program.c -o program

# Dynamic linking - uses shared libraries at runtime  
gcc program.c -o program
```

---

## **2. Shared Libraries (.so files)**

### **What are Shared Libraries?**
- Files with `.so` (shared object) extension
- Code shared between multiple programs
- Loaded into memory at runtime
- Reduces disk space and memory usage

### **Library Naming Convention:**
```bash
libname.so.major.minor.patch
# Example: libc.so.6.2.35
```

### **Common Library Locations:**
```bash
/lib/                          # Essential system libraries
/usr/lib/                      # Additional system libraries  
/lib/x86_64-linux-gnu/         # Architecture-specific (64-bit)
/usr/lib/x86_64-linux-gnu/     # Architecture-specific user libraries
/usr/local/lib/                # Locally installed libraries
```

---

## **3. Dynamic Linker/Loader (ld.so)**

### **What is ld.so?**
- **Runtime dynamic linker** that loads shared libraries when a program starts
- Usually `/lib64/ld-linux-x86-64.so.2` on 64-bit systems
- Resolves symbol dependencies and maps libraries into memory

### **Library Search Order:**
1. **RPATH** (embedded in binary)
2. **LD_LIBRARY_PATH** (environment variable)
3. **RUNPATH** (embedded in binary, newer than RPATH)
4. **System default paths** (`/lib`, `/usr/lib/`, etc.)
5. **ldconfig cache** (`/etc/ld.so.cache`)

---

## **4. RPATH Deep Dive**

### **What is RPATH?**
- **Runtime Path** - directories where the dynamic linker should search for libraries
- **Embedded in the ELF binary itself**
- Takes precedence over most other search methods
- **Survives binary moves** - travels with the binary

### **RPATH vs RUNPATH:**
```bash
# RPATH (older, higher precedence)
- Searched before LD_LIBRARY_PATH
- Cannot be overridden by LD_LIBRARY_PATH

# RUNPATH (newer, lower precedence)  
- Searched after LD_LIBRARY_PATH
- Can be overridden by LD_LIBRARY_PATH
```

### **Setting RPATH:**

#### **At Compile Time:**
```bash
# GCC/G++
gcc -Wl,-rpath,/path1:/path2 program.c -o program

# With multiple paths
gcc -Wl,-rpath,'/path1:/path2:/path3' program.c -o program
```

#### **Post-Compilation (using patchelf):**
```bash
# Set RPATH
patchelf --set-rpath "/path1:/path2:/path3" binary

# Show current RPATH
patchelf --print-rpath binary

# Remove RPATH
patchelf --remove-rpath binary

# Force RPATH (prevent conversion to RUNPATH)
patchelf --set-rpath "/path1:/path2" --force-rpath binary
```

#### **Post-Compilation (using chrpath):**
```bash
# Set RPATH (can only shrink existing RPATH)
chrpath -r "/path1:/path2" binary

# Show current RPATH
chrpath binary

# Remove RPATH
chrpath -d binary
```

### **Viewing RPATH:**
```bash
# Method 1: readelf
readelf -d binary | grep -E "(RPATH|RUNPATH)"

# Method 2: objdump
objdump -p binary | grep -E "(RPATH|RUNPATH)"

# Method 3: patchelf
patchelf --print-rpath binary

# Method 4: chrpath
chrpath binary
```

---

## **5. ldd Command**

### **What is ldd?**
- **List Dynamic Dependencies** - shows shared libraries required by a binary
- **NOT a binary analysis tool** - actually runs the program in a special mode
- Shows library locations and loading errors

### **Basic Usage:**
```bash
# Show dependencies
ldd /path/to/binary

# Verbose output
ldd -v /path/to/binary

# Show unused dependencies
ldd -u /path/to/binary
```

### **Sample Output:**
```bash
$ ldd /bin/ls
    linux-vdso.so.1 (0x00007fff123ab000)
    libselinux.so.1 => /lib/x86_64-linux-gnu/libselinux.so.1 (0x00007f8b4c123000)
    libc.so.6 => /lib/x86_64-linux-gnu/libc.so.6 (0x00007f8b4bf32000)
    /lib64/ld-linux-x86-64.so.2 (0x00007f8b4c456000)
```

### **Understanding Output:**
- **Left side**: Library name as referenced by binary
- **=>**: Points to actual file location
- **Right side**: Virtual memory address where loaded
- **Missing libraries**: Show as "not found"

### **Common Error Messages:**
```bash
# Library not found
libfoo.so.1 => not found

# Version issues  
./binary: ./binary: no version information available (required by ./binary)

# Missing interpreter
./binary: No such file or directory (wrong architecture)
```

---

## **6. Library Path Management**

### **Environment Variables:**
```bash
# LD_LIBRARY_PATH - runtime library search path
export LD_LIBRARY_PATH="/path1:/path2:$LD_LIBRARY_PATH"

# LD_DEBUG - debug library loading
LD_DEBUG=libs ./binary

# LD_PRELOAD - force load specific libraries
LD_PRELOAD="/path/to/lib.so" ./binary
```

### **System Configuration:**
```bash
# /etc/ld.so.conf - system library directories
echo "/usr/local/lib" >> /etc/ld.so.conf

# Update cache after changes
ldconfig

# View current cache
ldconfig -p

# Verbose cache update
ldconfig -v
```

---

## **7. Debugging Dynamic Linking**

### **Library Loading Debug:**
```bash
# Show library search process
LD_DEBUG=libs ./binary 2>&1 | less

# Show symbol binding
LD_DEBUG=bindings ./binary 2>&1 | less

# Show all debug info
LD_DEBUG=all ./binary 2>&1 | less
```

### **Finding Libraries:**
```bash
# Find library by name
find /lib /usr/lib -name "libname.so*" 2>/dev/null

# Find library using ldconfig
ldconfig -p | grep libname

# Check if library is available
pkg-config --libs libname
```

### **Architecture Compatibility:**
```bash
# Check binary architecture
file /path/to/binary

# Check library architecture  
file /path/to/library.so

# Check supported architectures
uname -m
```

---

## **8. Debian Packaging System**

### **Package Structure:**
```
debian/
├── control          # Package metadata
├── rules            # Build instructions  
├── postinst         # Post-installation script
├── preinst          # Pre-installation script
├── prerm            # Pre-removal script
├── postrm           # Post-removal script
├── changelog        # Version history
├── compat           # Debhelper compatibility level
└── install          # File installation mappings
```

### **Key Files:**

#### **debian/control:**
```bash
Source: package-name
Section: utils
Priority: optional
Maintainer: Your Name <email@example.com>
Build-Depends: debhelper (>= 10)
Standards-Version: 4.1.2

Package: package-name
Architecture: amd64
Depends: ${shlibs:Depends}, ${misc:Depends}
Description: Package description
```

#### **debian/rules:**
```makefile
#!/usr/bin/make -f

%:
	dh $@

override_dh_auto_build:
	# Custom build commands
	
override_dh_auto_install:
	# Custom installation commands
	
override_dh_shlibdeps:
	# Custom shared library dependency handling
	dh_shlibdeps --dpkg-shlibdeps-params=--ignore-missing-info
```

### **Build Process:**
```bash
# Build package
dpkg-buildpackage -us -uc

# Build binary package only
dpkg-buildpackage -us -uc -b

# Install built package
sudo dpkg -i package.deb

# Remove package
sudo dpkg -r package-name
```

---

## **9. Common Integration Issues**

### **RPATH and Packaging:**
- **Problem**: Debian packaging can modify or strip RPATH
- **Solution**: Use `override_dh_strip` in debian/rules to preserve RPATH
- **Best Practice**: Set RPATH consistently in build and post-install scripts

### **Architecture-Specific Paths:**
```bash
# Different architectures have different library paths
/lib/x86_64-linux-gnu/        # 64-bit Intel/AMD
/lib/aarch64-linux-gnu/       # 64-bit ARM  
/lib/arm-linux-gnueabihf/     # 32-bit ARM
```

### **Version Compatibility:**
- **glibc versions** affect dynamic linking behavior
- **Ubuntu 22.04**: glibc 2.35 (stricter)
- **Ubuntu 24.04**: glibc 2.39 (more flexible)

---

## **10. Best Practices**

### **RPATH Management:**
1. **Set complete RPATH** including all necessary system paths
2. **Be consistent** between build-time and runtime RPATH
3. **Include fallback paths** for different distributions
4. **Test on target systems** before deployment

### **Debian Packaging:**
1. **Handle dependencies properly** using `${shlibs:Depends}`
2. **Preserve RPATH** during packaging if needed
3. **Test installation/removal** thoroughly
4. **Use postinst scripts** for runtime configuration only

### **Debugging Workflow:**
1. **Check binary architecture** with `file`
2. **Verify RPATH** with `readelf` or `patchelf`
3. **Test dependencies** with `ldd`
4. **Debug loading** with `LD_DEBUG`
5. **Check system paths** with `ldconfig -p`

---

## **11. Quick Reference Commands**

```bash
# RPATH Management
patchelf --print-rpath binary
patchelf --set-rpath "/path1:/path2" binary
readelf -d binary | grep RPATH

# Library Dependencies  
ldd binary
LD_DEBUG=libs binary
ldconfig -p | grep libname

# Debian Packaging
dpkg-buildpackage -us -uc
dpkg -i package.deb
dpkg -l | grep package-name

# System Information
uname -m                    # Architecture
lsb_release -a             # Distribution info
/lib*/ld-linux*.so.* --version    # Linker version
```

---

## **12. What is patchelf?**

**`patchelf`** is a utility that **modifies ELF binaries** (executables and shared libraries) without needing to recompile them. It can:

- **Set/modify RPATH** and RUNPATH
- **Change interpreter** (dynamic linker)
- **Modify library dependencies**
- **Set/change SONAME** of shared libraries

### **Common patchelf Usage:**
```bash
# Set RPATH on a binary
patchelf --set-rpath "/path1:/path2:/path3" binary

# Show current RPATH  
patchelf --print-rpath binary

# Remove RPATH
patchelf --remove-rpath binary
```

---

## **13. Key Insights About RPATH**

### **RPATH Behavior:**
- **RPATH is embedded inside the binary itself** - it travels with the binary wherever you move it
- **Moving a binary does NOT affect its RPATH**
- **RPATH is stored in the ELF binary's headers**

### **Why Direct Binary Works vs Packaged Binary Issues:**
1. **Build time**: Binary gets RPATH set during compilation
2. **Direct testing**: Works because RPATH is intact and complete
3. **Packaging**: postinst script may re-patch with incomplete RPATH
4. **Runtime failure**: If postinst sets incomplete RPATH, libraries can't be found

### **Ubuntu Version Differences:**
- **Ubuntu 22.04**: Stricter dynamic linker (glibc 2.35) - requires explicit library paths in RPATH
- **Ubuntu 24.04**: More flexible dynamic linker (glibc 2.39) - can fall back to system paths even with incomplete RPATH

This comprehensive guide covers the fundamental concepts and their interactions in the Linux dynamic linking ecosystem, essential for understanding library management and packaging workflows.
