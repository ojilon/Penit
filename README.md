# Penit

**Minimal Android APK build · sign · install · debug tool** for a focused set of local Gradle projects (`wayer`, `conductino_android`, `fdroid`, `cooda`, …).

> **Migration notice**  
> Penit is intentionally standalone so the Android pipeline can be used and improved in isolation.  
> **It will be moved / merged into [ReleaseForge](https://github.com/ojilon/releaseforge) over time.**  
> Until then, treat this repo as the working home for the APK workflow.

Single binary. Go. Cobra CLI + Charm Bubble Tea TUI (light-blue theme).

## What it does

| Capability | Description |
|------------|-------------|
| **Build** | `assembleDebug` / `assembleRelease` + unit tests via `gradlew` |
| **Error capture** | Streams Gradle output, extracts clear failure headlines, shows full log path and expected APK locations |
| **Package** | Copies named APKs into a versioned folder under a local data-root |
| **Sign + verify** | Finds `apksigner`, prompts for keystore password only, signs and verifies |
| **ADB loop** | devices · install -r · uninstall · launch · filtered logcat · run-as (from Wayer ADB guide patterns) |
| **TUI** | Modern terminal UI with light-blue accents and a persistent command bar |

## Docs (plan first)

| Doc | Purpose |
|-----|---------|
| [docs/00-vision.md](docs/00-vision.md) | Goals, scope, relationship to ReleaseForge |
| [docs/01-architecture.md](docs/01-architecture.md) | Packages and data flow |
| [docs/02-project-registry.md](docs/02-project-registry.md) | Known projects + config shape |
| [docs/03-build-and-errors.md](docs/03-build-and-errors.md) | Gradle runner + structured errors |
| [docs/04-sign-and-package.md](docs/04-sign-and-package.md) | apksigner + packaging |
| [docs/05-adb-and-device.md](docs/05-adb-and-device.md) | Install / logcat / device loop |
| [docs/06-tui-design.md](docs/06-tui-design.md) | Light-blue TUI layout and commands |
| [docs/07-installer-and-cli.md](docs/07-installer-and-cli.md) | Installer, Cobra tree, doctor |
| [docs/08-implementation-plan.md](docs/08-implementation-plan.md) | Ordered build slices |

## Quick start (after Phase 0+)

```bash
go build -o penit .
./penit init
# edit config: point projects.wayer.path (etc.) at your local clones
./penit doctor
./penit build wayer debug
./penit install wayer --debug
./penit logcat wayer --clear
./penit tui
```

## Legacy scripts

The original Python release scripts (sign, package, version, …) live under [ReleaseForge `scripts/`](https://github.com/ojilon/releaseforge/tree/main/scripts) and were the reference for this rewrite. Penit does **not** invoke them at runtime.

## License

MIT
