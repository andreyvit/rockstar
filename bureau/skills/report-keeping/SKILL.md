---
name: report-keeping
description: Track task steps and write reports to preserve context across agent sessions and avoid forgetting user intentions.
allowed-tools: Bash, Read, Write
---

Bureau helps you succeed by providing longer-term memory for (1) user requests and (2) plans, research and big-picture path taken.

LLMs suffer from drift (when mundanity overshadows the big things) and compaction (when details are simply lost). This is one of the leading causes of LLM coding failures requiring human interventions and killing productivity. We don't want that, do we?

Bureau assists in keeping a per-task directory with numbered report files. You create and read these files, Bureau only takes care of paths.

Run as `${CLAUDE_PLUGIN_ROOT}/bureau` (Claude) or just `bureau` (Codex, others).


## THREE IMPORTANT CHANGES

1. Split your work into very clearly separated steps. Research, planning, reviewing plans, writing tests, writing implementation, reviewing code, updating docs, updating knowledge -- these should all be separate steps.

2. Before each step, re-read context from Bureau reports.

3. After each step, write a Bureau report file.


## Before each step

1. Run bureau without args to check current task and existing reports.

2. If current task doesn't match your work, run with `-N <slug>` to start a new task (or `-S <name>` to switch to an existing one IF user asked for that).

3. If user provided input, save it -- run with `-n user-something` (user-request, user-revision, user-feedback, etc.) and write their input VERBATIM, then append your interpretation.

4. Re-read recent reports to restore context.

  - Planning and review steps: read AS MANY reports as possible.

  - Execution steps: read all user requests + recent plans + all execution reports after the most recent plan.


## After each step

1. Run with `-n <suffix>` to get next report filename (e.g., plan, impl, fix, research).

2. Write a report including:
   - What you did and why
   - Decisions made and alternatives considered
   - Problems encountered and solutions
   - Open questions or blockers
   - Facts future agents should know
