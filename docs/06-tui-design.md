# TUI design — light-blue theme

## Inspiration

ReleaseForge’s command-bar + single viewport model (Claude Code / Gemini CLI style), but:

- Stronger visual identity (light blue / cyan accents).
- Leaner command surface focused on Android daily loop.
- Always-visible status of selected project + device.

## Colour palette (lipgloss)

| Role | Approx. | Usage |
|------|---------|-------|
| Primary accent | `#5B9BD5` / soft cyan-blue | borders, active prompt, success ticks |
| Secondary | `#7EC8E3` | secondary labels, sparkline |
| Background | terminal default (or very dark blue-grey if forced) | |
| Error | `#E06C75` | failure banners |
| Warning | `#E5C07B` | unauthorized device, unsigned APK |
| Muted | `#6B7B8A` | help text, paths |
| Text | default fg | body |

Keep contrast high for accessibility; avoid pure neon.

## Layout

```
┌─ Penit 0.1.0 · wayer · v0.3.1 (42) · device: R58M…  ● idle ─┐
│ branch main · clean · last build ✓ 2m ago                     │
├───────────────────────────────────────────────────────────────┤
│                                                               │
│  (scrollable viewport — command echo, build output,          │
│   logcat stream, help, error headlines)                       │
│                                                               │
├───────────────────────────────────────────────────────────────┤
│ > build debug                                                 │
│ build · test · install · logcat · devices · help · quit       │
└───────────────────────────────────────────────────────────────┘
```

- Top strip: tool version · project key · versionName (versionCode) · selected adb serial · state pill (`idle` / `building` / `installing` / `logcat`).
- Context line: git branch + dirty/clean + last-build summary.
- Viewport: ring buffer (~2000 lines), mouse wheel, auto-scroll while streaming.
- Input: always focused; ↑/↓ history; `esc` cancels running task; `ctrl+l` clears.
- Footer: context-sensitive hints.

Narrow terminals (<80 cols) collapse header to one line.

## Command bar verbs (v1)

| Input | Action |
|-------|--------|
| `open <key>` / `use <key>` | Select project from registry |
| `projects` | List configured keys + paths |
| `status` | Project + version + last APKs + device |
| `build debug\|release` | Async Gradle, streamed |
| `test` | Unit tests |
| `package [version]` | Copy APKs to data-root release folder |
| `sign` | Sign + verify (password prompt) |
| `install [--debug\|--release]` | adb install -r |
| `uninstall` | adb uninstall |
| `launch` | am start |
| `devices` | adb devices |
| `logcat [--clear]` | filtered stream |
| `logs [show\|tail]` | persisted build/logcat files |
| `version [--set X]` | show / set |
| `doctor` | SDK / adb / apksigner / gradlew checks |
| `help` / `?` / `-h` | help in viewport |
| `clear` / `quit` | |

Unknown → short error + "type help".

## Help view

Rendered inside the viewport (not a toast). Include the table above plus notes that full CLI flags live under `penit --help`.

## Libraries (same family as ReleaseForge)

- `charmbracelet/bubbletea`
- `charmbracelet/bubbles` (textinput, viewport)
- `charmbracelet/lipgloss`
- `charmbracelet/huh` (password, device picker, init forms)
- `spf13/cobra` for the CLI surface that the TUI also drives

## Differences from ReleaseForge TUI (current)

- Light-blue theme instead of neutral.
- Project selection is by **key** from a small registry, not free-path open.
- First-class `install` / `logcat` / `devices` in the bar (ReleaseForge keeps some of these CLI-only).
- No generic "scan any folder" flow.
