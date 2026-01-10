# bureau

A small portable Bash CLI that helps (AI) agents manage task-based report files in a predictable `_tasks/` directory: dated task folders, a `current` pointer, and sequentially-numbered Markdown reports.

This repo is a port of the *behavior* of `bureaumcp` (an MCP server) into a standalone CLI called `bureau`. MCP transport is intentionally out of scope here.

Zero third-party dependencies. Runs on macOS (old `/bin/bash`) and Linux.

## Why

When you (or an agent) work in multiple iterative steps, it’s helpful to keep a structured set of “reports”:

```
_tasks/
├── current -> 2025-10-02-refactor-something
├── 2025-10-01-implement-feature/
│   ├── 001-user-request.md
│   ├── 002-plan.md
│   ├── 003-implementation.md
│   └── 004-tests.md
└── 2025-10-02-refactor-something/
    ├── 001-analysis.md
    └── 002-plan.md
```

`bureau` automates:
- choosing the next task directory name for “today”
- tracking the active task via `_tasks/current`
- choosing the next report file number
- listing tasks and reports deterministically

## Install

```bash
./bureau
```

To install globally, copy `bureau` somewhere on your `PATH` (e.g. `~/bin/bureau`) and keep it executable.

## Usage

### Status (no args)

When run without arguments, `bureau` prints the current task directory and a (possibly truncated) list of report files:

```
Current task reports dir: _tasks/2025-10-02-refactor-something

[4 reports found]
001-user-request.md
002-plan.md
003-implementation.md
004-tests.md

[to start new report file] bureau -n <report-suffix>
[to switch to a new task]  bureau -N <task-suffix>
[to switch to prior task]  bureau -S <YYYY-MM-DD-task-suffix>
[to see recent tasks]      bureau -T
```

### Start a new report file (`-n`)

Prints the next sequential report filename for the current task (does not create the file):

```
Write your report to:
_tasks/2025-10-02-refactor-something/005-your-suffix.md
```

### List recent tasks (`-T`)

Prints up to 10 most recent task directories (oldest → newest among those shown):

```
[3 recent tasks:]
2025-10-01-implement-feature
2025-10-01b-fix-bug
2025-10-02-refactor-something
```

### Start a new task (`-N`)

Creates a new “today” task directory (with the first available suffix) and updates the `current` symlink:

```
Switched to new task YYYY-MM-DD-suffix.

<then the same output as running bureau with no args>
```

### Switch to an existing task (`-S`)

Switches the `current` symlink to an existing task directory:

```
Switched to preexisting task YYYY-MM-DD-suffix.

<then the same output as running bureau with no args>
```

### Tasks root (`BUREAU_DIR`)

By default, tasks live under `_tasks` in the current directory.

Override with:
- `BUREAU_DIR` environment variable

### Help (`-h` / `--help`)

`bureau -h` and `bureau --help` print:

- `Bureau - cli tool for managing AI agent report files.`
- the exact same output as `bureau` with no options
- additional notes about `BUREAU_DIR`

## Behavior (ported from `bureaumcp`)

This section is a precise spec of the underlying behavior, derived from `bureaumcp/index.js` and `bureaumcp/tools.test.js`.

### Task directory recognition

Directories are considered “task directories” only if they match:

Anything that begins with 4 digits (e.g. `2025...`) is treated as a task directory name.

### Generating a new task directory name

To create a new task for “today”, bureau picks the first unused directory name in this sequence:

1. `YYYY-MM-DD-<slug>` (first task of the day)
2. `YYYY-MM-DDb-<slug>` (second task of the day)
3. `YYYY-MM-DDc-<slug>` …
4. … through `YYYY-MM-DDy-<slug>` (25th task of the day)
5. `YYYY-MM-DDz026-<slug>` (26th task of the day)
6. `YYYY-MM-DDz027-<slug>` …

“Unused” means: there is no existing entry in `BUREAU_DIR` whose name starts with that date prefix (e.g. `2025-10-01b`).

Up to 1000 tasks per day are supported; beyond that, task creation fails.

Task dates are computed in local time.

### Current task tracking

The current task is stored as a symlink:

`<tasks-root>/current -> <task-dir-name>`

When switching tasks, the existing `current` symlink (if any) is removed and recreated. The symlink target is written as a relative name (the task directory name), but when reading the current task the basename is used, so an absolute target also works.

### Report files

Within a task directory, “report files” are any filenames matching:

any filename that starts with a digit

Examples:
- `001-user-request.md` ✅
- `42-not-padded.md` ✅
- `1` ✅
- `notes.md` ❌

The report listing is:
- sorted numerically by the leading number (like `sort -n`)
- returned in full if there are ≤ 50 report files
- truncated if there are > 50 report files: earliest 20 + latest 30

### Next report file number

To choose the next report number, bureau:
- scans the (possibly truncated) report listing
- extracts leading digits from each filename
- uses `max(prefixes) + 1` (gaps are ignored)

Because the listing is numerically sorted, the maximum prefix is always included in the “latest 30” slice, so truncation does not affect the next-number computation.

The filename format for new report files is:

`%03d-<suffix>.md`

### Recent tasks (`-T`)

`-T` lists up to 10 most recent task directories by lexicographic order (which matches chronological order for the `YYYY-MM-DD...` prefix).

## Development

Run regression tests:

```bash
./bureau-selftest.sh
```
