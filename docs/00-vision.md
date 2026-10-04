# Penit Vision

## Purpose

**Penit** is a minimal, single-binary Go tool for building, signing, packaging, installing, and debugging Android APKs from a focused set of local Gradle projects.

It is the specialised Android pipeline that will later be absorbed into [ReleaseForge](https://github.com/ojilon/releaseforge). Until then it stands alone so the APK workflow can be iterated and used in isolation.

## Target projects

| Key | Typical local path / notes |
|-----|----------------------------|
| `wayer` | Wayer Android (incl. `refactor/native-modules` branch) |
| `conductino_android` | Original home of the Python `scripts/` pipeline |
| `fdroid` | F-Droid related Android app repo |
| `cooda` | Cooda Android app repo |

Penit does **not** try to be "any project on disk". It knows these four (and any additional entries you add to its config) and applies a consistent Gradle + apksigner + adb flow to them.

## Core capabilities

1. **Build** – `assembleDebug` / `assembleRelease` (and unit tests) via `gradlew` / `gradlew.bat`.
2. **Capture Gradle errors** – stream output, extract headline failure lines, store full log, show clear location of APKs (or why they are missing).
3. **Package** – copy named APKs into a versioned release folder (no hard-coded "Conductino-Study" strings; use project config).
4. **Sign + verify** – locate `apksigner`, prompt for keystore password only (never store it), sign release APK, verify.
5. **ADB device loop** (from Wayer `docs/ADB_GUIDE.md`):
   - devices / serial selection
   - install `-r` (debug or signed release)
   - launch activity
   - filtered logcat (`WayerNative|AndroidRuntime` style filters, configurable)
   - uninstall / clean install
   - `run-as` inspect for debuggable builds
6. **Modern terminal UI** – light-blue themed Bubble Tea TUI + Cobra CLI, modelled on ReleaseForge but leaner and Android-focused.
7. **Installer** – same spirit as ReleaseForge’s installer (one-command install of the binary).

## Relationship to ReleaseForge

- Penit re-implements and improves the legacy Python scripts that live under ReleaseForge’s `scripts/` (and originally under Conductino / Penit).
- README and docs state clearly: **this tool will be moved / merged into ReleaseForge over time**.
- Shared ideas: data-root for logs/artifacts, structured errors, no password storage, live log streaming, Cobra + Bubble Tea.
- Differences: Penit is *only* the Android pipeline for the listed repos; no generic project scanner, no Wails/CMake, no full release-notes/GitHub orchestration in v1.

## Non-goals (v1)

- Generic "any folder" detection (that stays in ReleaseForge).
- Storing keystore passwords.
- Web UI / Electron.
- Full GitHub release create (optional later; focus is local build → sign → install → debug).
- Multi-user / server mode.

## Success criteria

- From a clean machine (after `penit init` + config of project paths) you can:
  1. `penit build wayer debug` → clear APK path or structured Gradle error.
  2. `penit sign wayer` → signed APK + verify.
  3. `penit install wayer --debug` → adb install -r.
  4. `penit logcat wayer` → filtered live logs.
  5. TUI with light-blue theme feels modern and usable for the daily loop.
