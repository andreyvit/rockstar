<div align="center">

# Rockstar

Collection of AI skills and extensions to rock your development.

_Compatible with Claude Code and Codex_

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

</div>

Rockstar is a collection of skills, commands and tools that enhance AI-assisted development workflows in Claude Code and Codex.

Contains:

* report-keeping
* star-team


## Installation

**Claude Code** — use my marketplace:

```shell
/plugin marketplace add andreyvit/claude-code-plugins
/plugin install rockstar
```

**Codex** — clone and run the installer:

```shell
git clone https://github.com/andreyvit/rockstar.git
cd rockstar
./install
```

See also: [Codex skills docs](https://developers.openai.com/codex/skills/create-skill).


## Using the skills

If you want to use a skill occasionally, just tell your agent to e.g. `use star-team skill` when you want to.

To enable by default, add the following to your `AGENTS.md`:

```
Use report-keeping and star-team skills for all future user requests and tasks.
```


## `report-keeping` skill

> Teaches your AI agents to store your input and write report files after each step. This solves agent amnesia and helps medium-term steering.

When an agent works in multiple iterative steps, it's helpful to keep a structured set of "reports":

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

Bureau provides:

- A skill that tells agents to read and write these reports
- A CLI tool that ensures correct naming and handling of files

You can install `bureau` tool by running `skills/report-keeping/scripts/bureau --install-symlink`.

When run without arguments, `bureau` prints the current task directory and a list of report files:

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


## `star-team` skill

WIP.


---


## Hacking on Rockstar

```shell
claude --plugin-dir /path/to/rockstar
```

Bureau self-test:

```bash
skills/report-keeping/scripts/bureau-selftest.sh
```

---

## License

MIT
