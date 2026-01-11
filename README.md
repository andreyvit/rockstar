<div align="center">

# Rockstar

Collection of AI skills and extensions to rock your development

_for Claude Code & Codex_

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

</div>

Rockstar is a collection of plugins and skills that enhance AI-assisted development workflows in Claude Code and Codex.


## Installation

Add the following to your `AGENTS.md`:

```
Use report-keeping skill for all future user requests and tasks.
```

and then install into your coding agent.

### Claude Code

```shell
/plugin marketplace add andreyvit/rockstar
/plugin install bureau@rockstar
```

### Codex

Clone and run the installer:

```shell
git clone https://github.com/andreyvit/rockstar.git
./rockstar/install-codex
```

See [Codex skills documentation](https://developers.openai.com/codex/skills/create-skill) for more details.


## Bureau Plugin

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
- A **skill** that tells agents to read and write these reports
- A **CLI tool** that ensures correct naming and handling of files


---


## Hacking on Rockstar

### Setup

Clone the repo, then add it as a local marketplace:

```shell
/plugin marketplace add /path/to/rockstar
/plugin install bureau@rockstar
```

Or for temporary per-session loading:

```shell
claude --plugin-dir /path/to/rockstar/bureau
```

### Testing

```bash
./bureau/bureau-selftest.sh
```

---

## License

MIT
