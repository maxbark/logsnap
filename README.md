# logsnap

A CLI tool for capturing and diffing structured log output across deployments.

---

## Installation

```bash
go install github.com/yourusername/logsnap@latest
```

Or build from source:

```bash
git clone https://github.com/yourusername/logsnap.git && cd logsnap && go build -o logsnap .
```

---

## Usage

Capture a snapshot of your current log output:

```bash
logsnap capture --input ./logs/app.log --out snapshot-v1.json
```

Diff two snapshots across deployments:

```bash
logsnap diff snapshot-v1.json snapshot-v2.json
```

Example output:

```
[+] new error: "database connection timeout" (3 occurrences)
[-] resolved: "cache miss on startup" (was 12 occurrences)
[~] changed: "request latency" avg 120ms → 340ms
```

### Flags

| Flag | Description |
|------|-------------|
| `--input` | Path to structured log file (JSON, logfmt) |
| `--out` | Output path for the snapshot file |
| `--format` | Log format: `json` or `logfmt` (default: `json`) |
| `--threshold` | Minimum occurrences to include in snapshot (default: `1`) |

---

## License

MIT © 2024 yourusername