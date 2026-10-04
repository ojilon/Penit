# Build pipeline & Gradle error capture

## Commands

| CLI | TUI bar | Behaviour |
|-----|---------|-----------|
| `penit build <key> debug` | `build debug` | `assembleDebug` |
| `penit build <key> release` | `build release` | `assembleRelease` |
| `penit test <key>` | `test` | `testDebugUnitTest` (default) |
| `penit test <key> --device` | `test --device` | instrumented (later) |

All run from the project root with `./gradlew` (or `gradlew.bat` on Windows). Optional `-g <gradle_user_home>` from config.

## Live streaming

- stdout + stderr are merged and streamed line-by-line to the TUI viewport / CLI.
- Simultaneously written to  
  `<data_root>/projects/<key>/logs/build-<timestamp>-<variant>.log`
- `esc` in TUI (or Ctrl+C in CLI) cancels the process group.

## Structured error extraction

On non-zero exit:

1. Scan the full log for high-signal patterns (order matters):
   - `* What went wrong:` … following paragraph
   - `FAILURE: Build failed with an exception.`
   - `Execution failed for task '…'`
   - `error:` / `Error:` lines that contain `.kt:`, `.java:`, `.xml:`
   - `e: file:///` Kotlin compiler errors
   - CMake / NDK lines when present (`CMake Error`, `ninja: error`)
2. Produce a **headline block** (max ~15 lines) shown immediately after the failure banner.
3. Always print:
   - Full log path
   - Expected APK paths for the requested variant
   - Last successful APK path from cache (if any)

Example output:

```
════════════════════════════════════════
BUILD FAILED  ·  wayer  ·  release
════════════════════════════════════════
* What went wrong:
Execution failed for task ':app:compileReleaseKotlin'.
> A failure occurred while executing …
  e: file:///…/NativeEngine.kt:42:5 Unresolved reference: foo

Full log : D:/PenitData/projects/wayer/logs/build-20261004-153012-release.log
Expected : …/app/build/outputs/apk/release/app-release-unsigned.apk
           …/app/build/outputs/apk/release/app-release.apk
```

## Success reporting (improvement over ReleaseForge today)

On success always print absolute paths of produced APKs and their sizes:

```
✓ assembleRelease finished
  app-release-unsigned.apk  12.4 MB
  → D:/Dev/Wayer/app/build/outputs/apk/release/app-release-unsigned.apk
```

## Package step

`penit package <key> [version]`

- Copies debug + release (unsigned) APKs into  
  `<data_root>/projects/<key>/releases/<version>/`
- Renames using `app_name` + version.
- Refuses to overwrite an existing version directory.

## Improvements vs legacy Python scripts

| Legacy | Penit |
|--------|-------|
| Hard-coded Conductino name | Config-driven `app_name` |
| No live stream / no full log path | Stream + persisted log + headline |
| Silent "Gradle failed" | Structured headlines + expected APK locations |
| Single-repo layout assumption | Multi-project registry |
| In-repo `release/` only | Data-root first; optional mirror |
