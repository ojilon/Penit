# Architecture

## High-level layout

```
main.go
  └── cmd/                    # Cobra commands
        root, tui, init, build, test, package, sign,
        install, devices, logcat, uninstall, doctor, version, …

internal/
  android/                    # APK paths, apksigner, adb wrappers
  build/                      # Gradle runner + error extractor
  config/                     # global + per-project JSON
  project/                    # known projects registry + path resolution
  storage/                    # data-root layout (logs, artifacts, cache)
  log/                        # live stream + persist
  git/                        # optional: branch/status for header
  version/                    # gradle.properties versionCode/Name R/W
  tui/                        # styles (light-blue theme), widgets, help
  app/                        # Bubble Tea model (command bar + viewport)

installer/                    # install scripts / packaging helpers
docs/                         # this plan set + future guides
scripts/                      # (optional) legacy Python kept for reference only
```

## Design rules

1. **Known projects only** – config maps keys (`wayer`, `conductino_android`, …) to absolute paths. No filesystem walk of arbitrary folders.
2. **Scan is cheap** – on open / build, verify the path still exists, read `gradle.properties` for version, detect `gradlew`, locate last APKs. Cache under data-root.
3. **Cobra and TUI share the same domain packages** – no duplicated pipelines.
4. **Data root is mandatory after `init`** – logs, packaged APKs, and command history live outside the git trees.
5. **Password never on disk** – keystore path + alias in config; password only via interactive prompt (huh / terminal).
6. **Structured Gradle failures** – parse stdout/stderr for `* What went wrong:`, `FAILURE:`, file:line patterns; surface a short headline + path to full log file + expected APK locations.

## Data flow (build → sign → install)

```
User: penit build wayer release
  → config.Load → project.Resolve("wayer")
  → build.RunGradle(assembleRelease)  // streams to log + TUI
  → on success: report APK path(s)
  → on failure: extract headlines, write full log, exit non-zero

User: penit sign wayer
  → locate unsigned release APK (or packaged copy)
  → android.FindApksigner
  → prompt password
  → sign + verify
  → report signed path

User: penit install wayer --release
  → android.ADB.Devices → pick serial if multiple
  → adb install -r <signed-apk>
  → optional: am start

User: penit logcat wayer
  → adb logcat -c then filtered stream (configurable tags)
```

## Error handling principles

- Missing project path → clear message + hint to edit config.
- Gradle non-zero → never just "failed"; always:
  - 3–10 headline lines
  - full log path under data-root
  - expected APK output paths (so user knows where to look)
- adb "unauthorized" / no devices → actionable text from ADB_GUIDE patterns.
- Signature mismatch on install → suggest uninstall first.

## Relationship to ReleaseForge packages

| ReleaseForge | Penit equivalent | Notes |
|--------------|------------------|-------|
| `internal/project` (any-folder scan) | `internal/project` (registry only) | Simpler |
| `internal/android` | `internal/android` | Expanded: adb, logcat, install |
| `internal/build` | `internal/build` | Gradle-focused + error parser |
| `internal/tui` + `app` | same | Light-blue theme, leaner command set |
| `internal/storage` | same idea | Smaller layout |
| Generic release notes / gh | deferred | Not v1 |
