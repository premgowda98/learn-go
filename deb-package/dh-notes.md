# 📦 `debian/rules` — Overview

The `debian/rules` file is essentially a **Makefile** that defines how your Debian package is built.

A typical minimal version:

```make
#!/usr/bin/make -f
%:
	dh $@
```

This tells `dh` to run a standard sequence of build steps, defined by the `debhelper` toolset.

You can override any of those steps with `override_dh_<step>` functions.

---

# 🧭 Default Build Sequence (`dh $@`)

This is the **default execution order** (simplified and most commonly relevant steps):

| Step                      | Can Be Overridden As            | Purpose                                              |
| ------------------------- | ------------------------------- | ---------------------------------------------------- |
| 1. `dh_clean`             | `override_dh_clean`             | Cleans up build artifacts                            |
| 2. `dh_autoreconf`        | `override_dh_autoreconf`        | Re-runs `autotools`, if used                         |
| 3. `dh_auto_configure`    | `override_dh_auto_configure`    | Configures the package (runs `./configure`)          |
| 4. `dh_auto_build`        | `override_dh_auto_build`        | Compiles the code (runs `make`)                      |
| 5. `dh_auto_test`         | `override_dh_auto_test`         | Runs tests                                           |
| 6. `dh_auto_install`      | `override_dh_auto_install`      | Installs into the staging area (`debian/<pkgname>/`) |
| 7. `dh_install`           | `override_dh_install`           | Moves files into the final package layout            |
| 8. `dh_installdocs`       | `override_dh_installdocs`       | Installs documentation (README, etc.)                |
| 9. `dh_installchangelogs` | `override_dh_installchangelogs` | Installs changelogs                                  |
| 10. `dh_strip`            | `override_dh_strip`             | Strips symbols (⚠️ may affect RPATH)                 |
| 11. `dh_compress`         | `override_dh_compress`          | Compresses man pages, docs                           |
| 12. `dh_fixperms`         | `override_dh_fixperms`          | Fixes file permissions                               |
| 13. `dh_shlibdeps`        | `override_dh_shlibdeps`         | Calculates shared library dependencies               |
| 14. `dh_installdeb`       | `override_dh_installdeb`        | Prepares control files like postinst, preinst        |
| 15. `dh_gencontrol`       | `override_dh_gencontrol`        | Generates the final control file                     |
| 16. `dh_md5sums`          | `override_dh_md5sums`           | Creates MD5 checksums for files                      |
| 17. `dh_builddeb`         | `override_dh_builddeb`          | Builds the `.deb` package                            |

---

# ✍️ Commonly Overridden Sections With Notes

Here’s a detailed look at what you may override and why:

---

### ✅ `override_dh_auto_build`

Used to customize how your software is compiled.

```make
override_dh_auto_build:
	go build -o bin/mybinary ./cmd/mybinary
```

---

### ✅ `override_dh_auto_install`

Used to copy compiled files and dependencies into the package layout.

```make
override_dh_auto_install:
	install -d debian/yourpkg/usr/bin
	install -m 0755 bin/mybinary debian/yourpkg/usr/bin/
```

---

### ✅ `override_dh_strip`

Used to prevent or customize stripping of binaries.

```make
override_dh_strip:
	dh_strip --strip-debug  # Safer; keeps RPATH and symbol version info
```

---

### ✅ `override_dh_shlibdeps`

Used to pass extra paths to `dpkg-shlibdeps` for non-standard library locations.

```make
override_dh_shlibdeps:
	dh_shlibdeps -l/var/kitecyber/lib
```

---

### ✅ `override_dh_builddeb`

Used to customize where `.deb` files are output or post-process after building.

```make
override_dh_builddeb:
	dh_builddeb --destdir=../builds/output
```

---

### ✅ `override_dh_buildinfo`

Optional; generates metadata about the build (host, time, etc.)

```make
override_dh_buildinfo:
	dh_buildinfo
```

This step is typically used if the `buildinfo` tool is installed.

---

### ✅ `override_dh_clean`

Clean custom build artifacts that `dh_clean` misses.

```make
override_dh_clean:
	dh_clean
	rm -rf build/ dist/
```

---

# ⚙️ Less Common but Useful Overrides

| Step                 | Use-case                                     |
| -------------------- | -------------------------------------------- |
| `dh_installsystemd`  | If you install a systemd service file        |
| `dh_installexamples` | To install example configs/scripts           |
| `dh_usrlocal`        | To avoid installing into `/usr/local`        |
| `dh_prep`            | If you need to tweak before packaging starts |

---

# 🛡️ Best Practices

* Use `override_dh_strip` to avoid breaking RPATH or symbol metadata
* Use `dh_shlibdeps` with `-l<your_lib_path>` if you're packaging with local `.so` files
* Don’t use `postinst` for things like `patchelf` unless absolutely necessary
* Build binaries with `debug` and `strip` flags controlled, not hardcoded
* Prefer placing all files in `debian/<package_name>/` manually in `auto_install`

---

# 🧪 Debugging the Sequence

To see exactly what steps `dh` is running:

```bash
debuild -us -uc -b -d -v
```

Or if you're calling `dpkg-buildpackage` directly:

```bash
dpkg-buildpackage -us -uc -b -j1
```

You can also add `set -x` inside `debian/rules` to trace what your overrides are doing.

---

## ✅ TL;DR: Minimal but Powerful `debian/rules`

```make
#!/usr/bin/make -f

%:
	dh $@

override_dh_auto_build:
	go build -o bin/clientagent ./cmd/clientagent

override_dh_auto_install:
	install -d debian/clientagent/var/kitecyber
	install -m 0755 bin/clientagent debian/clientagent/var/kitecyber/clientagent
	patchelf --set-rpath /var/kitecyber/lib:/usr/lib:/lib debian/clientagent/var/kitecyber/clientagent

override_dh_strip:
	dh_strip --strip-debug

override_dh_shlibdeps:
	dh_shlibdeps -l/var/kitecyber/lib
```