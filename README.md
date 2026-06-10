# dbd-ping-overlay

A small, always-on-top desktop overlay for Windows that shows the
approximate ping to your current Dead by Daylight game server.

The overlay runs as a completely separate desktop application. It looks at
the operating system's own process list and connection table from the
outside and never touches the game itself.

## What this tool does

- Looks for a running `DeadByDaylight-Win64-Shipping.exe` process (the name
  is configurable).
- Lists that process's active network connections using the operating
  system's connection table.
- Picks the most likely game server address from those connections (public
  IPv4, preferring established TCP / active UDP connections).
- Pings that address with the standard `ping` command and parses the
  average latency and packet loss.
- Displays the result in a small always-on-top window that updates every
  few seconds.

## What this tool does NOT do

This project is intentionally limited in scope for safety and
maintainability:

- **No DLL injection.** Nothing is ever loaded into the game process.
- **No reading of game process memory.** The game's memory is never opened,
  read, or written.
- **No modification of game files.** Nothing on disk belonging to the game
  is touched.
- **No hooks inside the game process.** The overlay does not attach to,
  hook, or otherwise interfere with the game's execution in any way.

## Why this approach is safer

Everything the overlay does is information that is already visible to any
other program on your machine through normal operating system APIs and
command-line tools (`tasklist`, `netstat`, `ping`). The application:

- Runs as its own separate window, not as an in-game overlay rendered
  inside the game's process or graphics context.
- Only reads process and connection metadata (PID, remote IP/port,
  connection state) — never process memory contents.
- Shells out to well-known, unprivileged system commands with bounded
  timeouts, so a slow or failing command cannot hang the UI.

Because no code is ever injected into or read from the game process, this
tool carries none of the risks (crashes, anti-cheat flags, save corruption)
associated with injection- or memory-reading-based tools.

## Requirements

- Windows 10/11
- [Go](https://go.dev/) (latest stable release)

## Running in development

```powershell
go run ./cmd/dbd-ping-overlay
```

This runs with default settings: it looks for
`DeadByDaylight-Win64-Shipping.exe`, refreshes every 2 seconds, and pings
with 4 packets at a 1 second timeout each.

## Building a release binary

```powershell
go build -ldflags="-H windowsgui" -o dbd-ping-overlay.exe ./cmd/dbd-ping-overlay
```

The `-H windowsgui` flag prevents a console window from appearing alongside
the overlay.

## Configuration

All settings have sensible defaults and can be overridden with command line
flags:

| Flag                 | Default                              | Description                                      |
| -------------------- | ------------------------------------- | ------------------------------------------------ |
| `--process-name`     | `DeadByDaylight-Win64-Shipping.exe`   | Executable name of the game process to monitor    |
| `--refresh-interval` | `2s`                                  | How often to refresh connection and ping data     |
| `--ping-count`       | `4`                                   | Number of ICMP echo requests sent per ping check  |
| `--ping-timeout`     | `1s`                                  | Per-reply timeout passed to the `ping` command    |

Example:

```powershell
dbd-ping-overlay.exe --process-name DeadByDaylight-Win64-Shipping.exe --refresh-interval 5s --ping-count 6 --ping-timeout 1500ms
```

## The overlay window

A small (300x140) always-on-top window shows:

- **DBD status** — whether the game process is currently running, and its
  PID
- **Server** — the detected game server IP address
- **Ping** — the average round-trip time, in milliseconds
- **Loss** — packet loss percentage from the last ping check
- **Updated** — the time of the last refresh
- A status/error message describing the current state

Network checks run on a background goroutine, so the window stays
responsive even if `tasklist`, `netstat`, or `ping` are slow to respond.
Closing the window shuts the application down cleanly.

## Troubleshooting

- **DBD process not found** — the overlay shows "Dead by Daylight is not
  running". Make sure the game is running, or that `--process-name` matches
  the actual executable name.
- **Server not detected** — the game process is running, but no public
  remote address could be identified yet. This can happen briefly during
  matchmaking or if all current connections are to private/local addresses.
- **ICMP blocked** — if the average latency cannot be determined, the
  overlay shows "Ping unavailable / ICMP may be blocked". Some servers or
  networks drop ICMP echo requests even though the game connection itself
  works fine.
- **UDP remote address shown as `*:*`** — `netstat` reports `*:*` for UDP
  sockets that have not (yet) connected to a specific remote peer. Such
  entries are ignored when picking a server candidate.

## Roadmap

- Passive RTT estimation using Npcap/gopacket instead of (or in addition to)
  ICMP ping.
- Smarter server candidate scoring (e.g. weighting by traffic volume over
  time).
- A system tray icon with show/hide and quit actions.
- A transparent, click-through overlay mode.
- Region detection based on the server IP's ASN/geolocation.
- Packaged Windows releases (installer or standalone zip).

## Project layout

```text
dbd-ping-overlay/
  cmd/dbd-ping-overlay/   application entry point
  internal/app/           wires everything together and runs the poll loop
  internal/process/       finds the game process via tasklist
  internal/netstat/       inspects connections via netstat and picks a server
  internal/ping/          runs and parses "ping"
  internal/overlay/       the always-on-top window (Windows only)
  internal/config/        configuration and flag parsing
  internal/model/         shared status types
  internal/logging/       logger setup
```

## License

[MIT](LICENSE)
