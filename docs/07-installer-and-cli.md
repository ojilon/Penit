# Installer, CLI surface & packaging

## Binary name

`penit` (Windows: `penit.exe`)

## Installer goals (mirror ReleaseForge spirit)

- One-command install from GitHub releases (or local build).
- Cross-platform: Linux amd64/arm64, macOS amd64/arm64, Windows amd64/arm64.
- Prefer user-local install (`~/.local/bin` or `%LOCALAPPDATA%\Penit`) with optional system-wide.

Suggested scripts (to be written under `installer/`):

```
installer/
  install.sh          # curl | bash style
  install.ps1         # Windows
  README.md
```

Example user flow:

```bash
# Linux / macOS
curl -fsSL https://raw.githubusercontent.com/ojilon/Penit/main/installer/install.sh | bash

# or from a release asset
curl -L …/penit-linux-amd64.tar.gz | tar xz
sudo mv penit /usr/local/bin/
```

```powershell
# Windows
irm https://raw.githubusercontent.com/ojilon/Penit/main/installer/install.ps1 | iex
```

## Build / release of Penit itself

Makefile or just Go + goreleaser later. Minimal Makefile targets:

```
build-local
build          # multi GOOS/GOARCH into dist/
install        # build-local + move to /usr/local/bin or equivalent
test
clean
```

Version injected via ldflags (`main.Version`, `BuildTime`, `GitCommit`).

## Cobra command tree (v1)

```
penit
  init                 # create data-root + default config
  tui                  # interactive (default when no subcommand?)
  projects
    list
    add
    remove
  build    <key> [debug|release]
  test     <key>
  package  <key> [version]
  sign     <key>
  install  <key> [--debug|--release] [--serial]
  uninstall <key> …
  launch   <key> …
  devices
  logcat   <key> …
  version  <key> [--set]
  doctor
  logs     <key> …
  help / --help
```

Flags that apply widely:

- `--data-root`
- `--config`
- `--serial` (adb)
- `--json` (machine-readable status for scripting, later)

## Doctor

`penit doctor` (and TUI `doctor`) checks:

- data-root writable
- each configured project path exists + has gradlew
- ANDROID_HOME / sdk
- adb present and at least one device (or warn)
- apksigner present
- keystore files exist (not passwords)

Exit non-zero if anything critical is missing.

## Migration note for README

Top of README must state:

> **Penit is a focused Android APK build/sign/install tool.**  
> It will be merged into [ReleaseForge](https://github.com/ojilon/releaseforge) over time.  
> Use it standalone until that integration lands.
