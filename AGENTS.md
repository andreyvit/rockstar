# Agent notes (bureau)

## Goals / constraints

- Portable Bash: must run on macOS `/bin/bash` 3.2 and on Linux Bash.
- Zero dependencies: do not add external tools/libs beyond typical POSIX userland (`date`, `sort`, `readlink`, `ln`, `mkdir`, `rm`, `printf`).
- Behavior stability matters: CLI output is treated as an interface; keep it byte-for-byte stable unless intentionally changing the UX.

## Code structure

- The whole tool is a single script: `./bureau`.
- `BUREAU_DIR` is the global “tasks root” and is the single source of truth. It is set once in `main()` from:
  - environment `BUREAU_DIR`
  - default `_tasks`
- Helper functions assume `BUREAU_DIR` is set.

## Bash compatibility rules

- Don’t use Bash 4+ features (e.g. associative arrays, `mapfile`, `${var,,}`, etc.).
- Prefer `printf` over `echo`; always quote variables.
- Keep strict mode (`set -euo pipefail`) enabled.
- Avoid date arithmetic; this tool only needs “today” (`date +%F`).
- Assume task/report names don’t contain spaces/backslashes; code uses simple word-splitting for readability.

## Tests

- Run: `./bureau-selftest.sh`
- The selftest asserts exact stdout strings; if you change formatting, update tests deliberately.

## Behavioral invariants to preserve

- Task “today” is computed in local time.
- Task directories are discovered leniently: any entry name starting with 4 digits is treated as a task.
- Report files are discovered leniently: any filename starting with a digit.
- Report file ordering is numeric by leading number.
- If >50 report files: list earliest 20 + latest 30 (numeric order).
- Next report number is computed from the listing’s max leading number.
- `-T/--list-tasks` lists up to 10 most recent task directories (lexicographic order).
