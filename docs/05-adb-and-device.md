# ADB, install, logcat & device loop

Source of truth for the desired developer loop:  
[Wayer `docs/ADB_GUIDE.md` on branch `refactor/native-modules`](https://github.com/ojilon/Wayer/blob/refactor/native-modules/docs/ADB_GUIDE.md).

Penit wraps the same operations with project-aware defaults (package IDs, main activity, logcat filters).

## Locating adb

1. Config `android_sdk` / per-project override.
2. `$ANDROID_HOME` / `$ANDROID_SDK_ROOT` → `platform-tools/adb[.exe]`.
3. `PATH`.
4. Fail via `penit doctor` with install hints.

## Commands

| CLI | Purpose |
|-----|---------|
| `penit devices` | `adb devices -l` (mark unauthorized / offline) |
| `penit install <key> [--debug\|--release] [--serial S]` | `adb install -r <apk>` |
| `penit uninstall <key> [--debug\|--release] [--serial S]` | `adb uninstall <package_id>` |
| `penit launch <key> [--debug] [--serial S]` | `adb shell am start -n <pkg>/<activity>` |
| `penit logcat <key> [--clear] [--serial S] [--filter …]` | filtered live logcat |
| `penit run-as <key> <shell-cmd…>` | `adb shell run-as <debug-pkg> …` (debuggable only) |

TUI bar mirrors the most common ones: `install`, `uninstall`, `logcat`, `devices`.

## Install details

- Prefer the packaged or just-built APK path reported by the last build/package.
- Debug package ID and release package ID come from project config (Wayer uses `.debug` suffix).
- On `INSTALL_FAILED_UPDATE_INCOMPATIBLE` / signature mismatch → print the uninstall command and offer to run it.
- Always echo the exact `adb` command for transparency.

## Logcat

Default filters per project (from config), e.g. for Wayer:

```
WayerNative|AndroidRuntime
```

Implementation:

```
adb logcat -c          # if --clear
adb logcat             # stream, client-side filter or adb -s … '*:S' Tag:V …
```

In TUI the stream goes into the viewport; `esc` stops. Full unfiltered capture can be written to  
`<data_root>/projects/<key>/logs/logcat-<timestamp>.log`.

## Typical daily loop (what Penit should make one command sequence)

```bash
penit build wayer debug
penit install wayer --debug
penit logcat wayer --clear
# interact on device; watch filtered logs
```

Or in TUI after `open wayer`:

```
> build debug
> install
> logcat --clear
```

## Signed release test

```bash
penit build wayer release
penit sign wayer
penit install wayer --release
penit logcat wayer
```

Unsigned release APKs are refused by the system; Penit should detect "unsigned" (or missing certs) and refuse to install, pointing the user at `sign`.
