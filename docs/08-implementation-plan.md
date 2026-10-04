# Implementation plan (ordered)

Work **one slice at a time**. Each step should leave a compilable binary.

## Phase 0 — Skeleton (this pass)

- [x] `docs/00`–`08` written
- [x] `go.mod`, `main.go`, empty Cobra root
- [x] `internal/config` – load/save global JSON
- [x] `internal/storage` – data-root layout helpers
- [x] `cmd/init`, `cmd/version` (tool version), `cmd/doctor`
- [x] README stating future move into ReleaseForge
- [x] Basic Makefile / build

## Phase 1 — Build + error capture

- [ ] `internal/build` – Gradle runner with live stream + cancel
- [ ] Gradle error headline extractor
- [ ] `cmd/build`, `cmd/test`
- [ ] Persist logs under data-root
- [ ] Report absolute APK paths on success / expected paths on failure

## Phase 2 — Version + package

- [ ] `internal/version` – read/write `gradle.properties`
- [ ] `cmd/version` (project version)
- [ ] `cmd/package` – copy & rename APKs using `app_name`

## Phase 3 — Sign

- [ ] `internal/android` – find apksigner, sign, verify
- [ ] Password prompt (CLI + huh for TUI later)
- [ ] `cmd/sign`

## Phase 4 — ADB loop

- [ ] `internal/android` – adb locate, devices, install, uninstall, launch, logcat, run-as
- [ ] `cmd/devices`, `install`, `uninstall`, `launch`, `logcat`
- [ ] Device serial selection when multiple devices present

## Phase 5 — TUI (light-blue)

- [ ] `internal/tui` – styles, theme constants
- [ ] `internal/app` – Bubble Tea model (header, viewport, command bar)
- [ ] Wire commands to the same domain packages
- [ ] History, clear, cancel, help-in-viewport
- [ ] `cmd/tui`

## Phase 6 — Installer & polish

- [ ] `installer/install.sh` + `install.ps1`
- [ ] Multi-arch build script / Makefile
- [ ] `penit doctor` fully implemented
- [ ] End-to-end smoke on Wayer (debug build → install → logcat)
- [ ] Optional: zip helper, simple notes template

## Phase 7 — Handoff readiness

- [ ] Document mapping of every package to the future ReleaseForge location
- [ ] Keep legacy Python under `scripts/` for reference only (or point at ReleaseForge’s copy)
- [ ] Tag a usable `v0.1.0`

## Coding conventions

- Go 1.22+ (match ReleaseForge’s 1.23 when practical).
- No third-party beyond Charm + Cobra for v1.
- Prefer absolute paths in user-facing messages.
- All user errors should be actionable (what to run / which file to edit).
- Tests for pure helpers (error parser, version R/W, path resolution); integration tests optional early on.

## Immediate next actions after docs

1. Write `README.md` (with the migration notice).
2. Scaffold `go.mod`, `main.go`, `cmd/root.go`.
3. Implement `config` + `storage` + `init`.
4. Then Phase 1 build runner.
