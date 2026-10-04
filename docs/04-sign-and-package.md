# Sign, verify & package

## Locating apksigner

Same strategy as legacy `sign.py`, improved for reliability:

1. Config override `apksigner_path` (global or per-project).
2. `PATH` (`apksigner` / `apksigner.bat`).
3. `$ANDROID_HOME` or `$ANDROID_SDK_ROOT` → `build-tools/<highest-version>/apksigner[.bat]`.
4. Fail with a clear doctor-style message if still missing.

## Sign flow

```
penit sign <key> [--apk <path>]
```

1. Resolve unsigned APK:
   - explicit `--apk`, or
   - packaged release APK under data-root, or
   - `app/build/outputs/apk/release/app-release-unsigned.apk`
2. Resolve keystore (absolute or relative to project root) + alias from project config.
3. Prompt for keystore password (never echoed, never stored).
4. Run:
   ```
   apksigner sign --ks <ks> --ks-key-alias <alias> --ks-pass pass:<pw> <apk>
   ```
5. Immediately verify:
   ```
   apksigner verify --verbose <apk>
   ```
6. Report signed path + certificate summary (from `--print-certs` if useful).

If the APK was already signed by Gradle’s `signingConfig`, Penit still allows a re-sign or skips with a warning.

## Package flow (recap)

```
penit package <key> [version]
```

- Creates `<data_root>/projects/<key>/releases/<version>/`
- Copies:
  - debug APK → `{app_name}-{version}-debug.apk`
  - unsigned release → `{app_name}-{version}-release.apk`
- After a successful `sign`, the signed file can be copied/renamed to  
  `{app_name}-{version}-release-signed.apk`

## Zip (optional, later)

Legacy had `zip_release.py`. v1 can defer; a simple `penit zip <key> <version>` that zips the release folder is easy to add once packaging is solid.

## Security notes

- Password only via interactive prompt (huh in TUI, terminal in CLI).
- No password in config files, env vars that are written to disk, or logs.
- Keystore path itself is fine to store; treat the file as secret on the user’s machine.
