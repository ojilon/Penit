# Project registry & configuration

## Known projects (v1)

Penit is intentionally narrow. The registry is a map in the global config:

```json
{
  "data_root": "D:/PenitData",
  "android_sdk": null,
  "projects": {
    "wayer": {
      "path": "D:/Dev/Wayer",
      "app_name": "Wayer",
      "package_id_debug": "com.example.wayer.debug",
      "package_id_release": "com.example.wayer",
      "main_activity": "com.example.wayer.core.MainActivity",
      "keystore_path": "wayer-release.jks",
      "keystore_alias": "wayer",
      "gradle_user_home": null,
      "logcat_filters": ["WayerNative", "AndroidRuntime"]
    },
    "conductino_android": {
      "path": "D:/Dev/conductino_android",
      "app_name": "Conductino-Study",
      "package_id_debug": "...",
      "package_id_release": "...",
      "main_activity": "...",
      "keystore_path": "conductino-release.jks",
      "keystore_alias": "conductino",
      "logcat_filters": ["AndroidRuntime"]
    },
    "fdroid": { "...": "..." },
    "cooda":  { "...": "..." }
  }
}
```

Paths are absolute. Relative keystore paths are resolved against the project root.

## Detection helpers (per project open)

When a project is selected:

1. Confirm `path` exists and is a directory.
2. Look for `gradlew` / `gradlew.bat` and `settings.gradle` / `settings.gradle.kts`.
3. Read `gradle.properties` for `app.versionCode` / `app.versionName` (or project-specific keys if documented later).
4. Locate last-known APKs under:
   - `app/build/outputs/apk/debug/app-debug.apk`
   - `app/build/outputs/apk/release/app-release-unsigned.apk`
   - `app/build/outputs/apk/release/app-release.apk` (if already signed by Gradle)
5. Cache a small `projects/<key>/cache.json` under data-root (version, last build time, last APK paths).

No deep recursive scan of the whole tree (unlike ReleaseForge).

## Adding a new project

```bash
penit projects add <key> --path /abs/path --app-name "MyApp"
# then edit config for package IDs, activity, keystore, filters
```

Or hand-edit the JSON; `penit doctor` validates.

## Version source

Default: `gradle.properties` keys `app.versionCode` and `app.versionName` (same as legacy scripts).

`penit version <key>` → print  
`penit version <key> --set 1.2.0` → write name + increment code by 1  
`penit version <key> --code-only` → increment code only

## Artifact naming

Never hard-code "Conductino-Study". Use `app_name` from config:

```
{app_name}-{version}-debug.apk
{app_name}-{version}-release.apk
{app_name}-{version}-release-signed.apk
```

Packaged copies live under:

```
<data_root>/projects/<key>/releases/<version>/
```

Optional in-repo `release/<version>/` mirror only if `mirror_in_repo: true` is set for that project.
