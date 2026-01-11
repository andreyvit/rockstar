<div align="center">

# Rockstar

**Claude Code & Codex plugins/skills to rock your development**

*by Andrey Tarantsov*

[![License: MIT](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE)

</div>

Rockstar is a collection of plugins and skills that enhance AI-assisted development workflows in Claude Code and Codex.

---

## Installation

<details open>
<summary><strong>Claude Code</strong></summary>

Add this marketplace and install plugins:

```shell
/plugin marketplace add andreyvit/rockstar
/plugin install bureau@rockstar
```

</details>

<details>
<summary><strong>Codex</strong></summary>

Clone and copy the skill + CLI:

```shell
git clone https://github.com/andreyvit/rockstar.git
./rockstar/bureau/bureau --install-symlink

# Copy the skill (user-scoped)
cp -r rockstar/bureau/skills/report-keeping ~/.codex/skills/

# Or for repo-scoped (travels with your project)
cp -r rockstar/bureau/skills/report-keeping .codex/skills/
```

See [Codex Skills documentation](https://developers.openai.com/codex/skills/create-skill) for more details.

</details>

---

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

## Development

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
