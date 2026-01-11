---
name: star-team
description: Run complex tasks successfully by utilitizing a development team with famour personalities. Use when
---

Star Team is a specific AI workflow that helps coding agents produce better code and solve complex tasks. We build upon the common wisdom of clear separation of research, planning and execution steps.

Star Team workflow is based on subtasks and steps, and 3 principles.

Subtask: a part of the overall task with a specific deliverable. Between subtasks, the code must built, and all tests must pass.

Step: narrow activity towards completing a subtask. This can be one of the following:

- research
- planning
- review of plans
- execution - refactoring
- execution - writing tests
- execution - writing code
- execution - writing docs
- problem solving
- review of changes
- preserving and organizing knowledge


## Three Principles of Star Team

1. Split the entire task into subtasks and execute SEQUENTIALLY, making the first subtask deliverable before moving on to the second one. This is NOT a waterfall; we are open to learning from the first subtask and adjusting future subtasks.

2. Within a subtask, execute a series of steps, each focused on just one kind of activity, and uses a star skill for that activity. Do not mix activities within the step. If you're planning, you're not writing code. If you're writing code, you're not changing tests. If you're writing tests, you're not updating code. Etc.

3. Steps follow a specific workflow driven by planning and reviews.


## Subtasks

A user's request often involves doing multiple changes. Usually, AIs try to go 'depth first': they do a LOT of changes, produce loads of semi-broken code, and then tweak and adjust until that code works. That's counterproductive, and makes AI stupid.

Instead, we focus on delivering subtasks sequentially. This (1) helps us keep complexity manageable, (2) helps us learn faster; (3) allows user to review results easier and give feedback.

Make a WIP commit after each subtask. Specific format depends on project conventions, but something like `WIP on <task>: <subtask>` is a reasonable default format.


## Steps and Star Skills

Each activity, be it research, planning or coding, benefits from a particular set of instructions. We package those instructions as skills.

Use report-keeping skill to maintain long-term context throughout the task, across subtasks and steps.

When starting a step, load the corresponding skill, and re-read the appropriate reports.

When finishing a step, write a new report.


## Star Team Workflow

PREP phase:

1. Save user input if any, like `report-keeping` skill demands.

2. Research step: use `star-research` skill in a subagent to collect and trace the code relevant to the task.

PLANNING phase:

1. Tech Lead step: use `star-techlead` skill to make high-level decisions about the task. Is it finished? Is research complete? Is the spec clear? What exactly needs to be done? How should we approach that? If task is not done, Tech Lead determines the next step; in most cases, this is Tech Spec step, but the lead might request more research, or code review, etc.

2. Tech Spec step: use `star-techspec` skill to create a detailed specification for the task. This should mention the EXACT code-level changes that need to be done, with as much detail as we can, but without writing actual code. NO code in the report file, please, this is ONE STEP ABOVE the specific code level.

3. Plan Review step: use `star-planreview` skill to review the plan created in by Tech Lead and Tech Step steps. If review results in 100% approval, proceed to execution. If not, Tech Lead determines the next step.

EXECUTION phase:

1. Tests step: use `star-tests` skill to update tests for the current tasks. If stuck, call `star-problemsolver`.

2. Code step: use `star-impl` skill to write the code for the current task. If stuck, call `star-problemsolver`.

3. Documentation step: use `star-docs` skill to update the documentation for the current task.

4. Code Review step: use `star-codereview` skill to review the code written in Code step.

5. After execution, we ALWAYS go back to Tech Lead to determine the next step.

MANUAL TESTING phase (tech lead calls for this if user explicitly asked for it) using `star-browsertesting` skill.

FINALLY, if tech leads determines we're done on a subtask, it tells which subtask to switch to next, or if the entire task is done.

NEVER STOP until tech lead says the entire task is completely, 100% finished with no remaining outstanding issues. Otherwise, CONTINUE; if tech lead finished with no clear indication of the next steps, call them again asking for specific instructions on what to do next.
