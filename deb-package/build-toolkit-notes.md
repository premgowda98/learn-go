# 📦 𝗗𝗲𝗯𝗶𝗮𝗻 𝗣𝗮𝗰𝗸𝗮𝗴𝗲 𝗧𝗼𝗼𝗹𝗸𝗶𝘁

---

## 🧱 1. **Building Debian Packages**

| Tool                | Purpose                                             | Example                        |
| ------------------- | --------------------------------------------------- | ------------------------------ |
| `debuild`           | Most common frontend to build `.deb` from `debian/` | `debuild -us -uc`              |
| `dpkg-buildpackage` | Lower-level build command (used by `debuild`)       | `dpkg-buildpackage -us -uc -b` |
| `fakeroot`          | Used with build commands to simulate root           | `fakeroot debian/rules binary` |
| `dh`                | Debhelper wrapper that runs all build steps         | `dh $@` inside `debian/rules`  |
| `make`              | Often used to invoke custom build logic             | `make deb`                     |

---

## 🧪 2. **Inspecting `.deb` Files**

### 📁 See what’s inside

| Tool          | Purpose                                                           | Example                                          |
| ------------- | ----------------------------------------------------------------- | ------------------------------------------------ |
| `dpkg-deb -c` | List contents of a `.deb` file                                    | `dpkg-deb -c clientagent_2.0.0_amd64.deb`        |
| `dpkg-deb -R` | Extracts `.deb` file                                              | `dpkg-deb -R clientagent_2.0.0_amd64.deb ./dest` |
| `dpkg-deb -b` | Build package again from extracted dir                            | `dpkg-deb -b ./dest clientagent_2.0.0_amd64.deb` |
| `dpkg -x`     | Extract the contents to a directory                               | `dpkg -x clientagent.deb /tmp/pkg`               |
| `ar`          | `.deb` is an ar archive – list internal files                     | `ar t clientagent.deb`                           |
| `ar x`        | Extract `.deb`'s `control.tar.gz`, `data.tar.xz`, `debian-binary` | `ar x clientagent.deb`                           |

---

## 📋 3. **Inspecting Control Metadata**

| Tool            | Purpose                                  | Example                       |                    |
| --------------- | ---------------------------------------- | ----------------------------- | ------------------ |
| `dpkg-deb -I`   | Show package metadata (control info)     | `dpkg-deb -I clientagent.deb` |                    |
| `dpkg -s <pkg>` | Show info for **installed** package      | `dpkg -s clientagent`         |                    |
| `dpkg -p <pkg>` | Query available package info (from repo) | `dpkg -p curl`                |                    |
| `dpkg -l`       | List installed packages                  | \`dpkg -l                     | grep clientagent\` |

---

## 🔍 4. **Inspecting ELF Binaries Inside `.deb`**

### 🧠 Use these to inspect Go/C/C++ binaries in your package:

| Tool         | Purpose                                         | Example                         |               |
| ------------ | ----------------------------------------------- | ------------------------------- | ------------- |
| `file`       | Tells you file type, strip status, architecture | `file /path/to/binary`          |               |
| `ldd`        | List dynamic dependencies (shared libraries)    | `ldd /path/to/binary`           |               |
| `readelf -d` | View dynamic section: RPATH, RUNPATH, NEEDED    | `readelf -d /path/to/binary`    |               |
| `readelf -h` | ELF header: arch, format, ABI                   | `readelf -h /path/to/binary`    |               |
| `objdump -x` | Extended headers and symbol info                | `objdump -x /path/to/binary`    |               |
| `nm`         | Show symbol table                               | \`nm /path/to/binary            | grep main\`   |
| `patchelf`   | Change/interrogate RPATH, SONAME, interpreter   | `patchelf --print-rpath binary` |               |
| `strings`    | Show printable strings in binary (quick debug)  | \`strings binary                | grep config\` |

---

## 🧰 5. **Inspecting Dependencies & Shared Libraries**

| Tool                        | Purpose                                           | Example                           |                     |
| --------------------------- | ------------------------------------------------- | --------------------------------- | ------------------- |
| `dpkg-shlibdeps`            | Computes shlib dependencies for binaries          | `dpkg-shlibdeps ./binary`         |                     |
| `ldconfig -p`               | List all available shared libraries               | \`ldconfig -p                     | grep libtesseract\` |
| `dpkg -S /path/to/file`     | Find which package owns a file                    | `dpkg -S /usr/lib/libc.so.6`      |                     |
| `apt-file search <libname>` | Find which package provides a file (if installed) | `apt-file search libtesseract.so` |                     |

---

## 🔧 6. **Manipulating Installed Packages**

| Tool                    | Purpose                                  | Example                              |
| ----------------------- | ---------------------------------------- | ------------------------------------ |
| `dpkg -i`               | Install `.deb` manually                  | `sudo dpkg -i clientagent.deb`       |
| `apt install ./pkg.deb` | Safer install with dependency resolution | `sudo apt install ./clientagent.deb` |
| `dpkg -r <pkg>`         | Remove a package                         | `sudo dpkg -r clientagent`           |
| `dpkg -P <pkg>`         | Purge (remove + config)                  | `sudo dpkg -P clientagent`           |

---

## 🧪 7. **Debugging `.deb` Build Process**

| Tool                       | Purpose                                       | Example                   |
| -------------------------- | --------------------------------------------- | ------------------------- |
| `debuild -us -uc -d -b -v` | Verbose build output                          | see full log              |
| `set -x` in `debian/rules` | Show each command that runs                   |                           |
| `fakeroot`                 | Run build without actual root permissions     |                           |
| `lintian`                  | Lint your package for Debian policy issues    | `lintian clientagent.deb` |
| `pbuilder` / `sbuild`      | Build in a clean chroot (like Docker)         | advanced usage            |
| `diffoscope`               | Compare two .deb builds (for reproducibility) |                           |

---

## 🔐 8. **RPATH / RUNPATH Debugging Tools**

| Tool                     | Purpose                             | Example                                     |             |
| ------------------------ | ----------------------------------- | ------------------------------------------- | ----------- |
| `patchelf --print-rpath` | View current RPATH                  | `patchelf --print-rpath ./binary`           |             |
| `patchelf --set-rpath`   | Set custom RPATH                    | `patchelf --set-rpath /var/mylibs ./binary` |             |
| `readelf -d`             | Also shows RPATH and RUNPATH        | \`readelf -d ./binary                       | grep PATH\` |
| `chrpath`                | Alternative to `patchelf` for RPATH | `chrpath -l ./binary`                       |             |

---

## 💡 Workflow Example

```bash
# 1. Build the .deb
make build     # build Go binary
debuild -us -uc

# 2. Inspect the contents
dpkg-deb -c clientagent_2.0.0_amd64.deb
dpkg-deb -I clientagent_2.0.0_amd64.deb

# 3. Extract and inspect ELF
dpkg -x clientagent_2.0.0_amd64.deb /tmp/clientagent
file /tmp/clientagent/var/kitecyber/clientagent
readelf -d /tmp/clientagent/var/kitecyber/clientagent

# 4. Lint and check dependencies
lintian clientagent_2.0.0_amd64.deb
ldd /tmp/clientagent/var/kitecyber/clientagent

# 5. Install and test
sudo apt install ./clientagent_2.0.0_amd64.deb
```

---

## ✅ TL;DR: Most Useful Commands Cheat Sheet

```bash
# Build
debuild -us -uc

# Install
sudo apt install ./pkg.deb

# Extract contents
dpkg -x pkg.deb ./outdir

# Metadata and files
dpkg-deb -I pkg.deb
dpkg-deb -c pkg.deb

# Inspect binary
file binary
ldd binary
readelf -d binary
patchelf --print-rpath binary

# Debug symbols
dpkg -s pkgname
nm binary | grep symbol
```