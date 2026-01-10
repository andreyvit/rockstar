#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
bureau="$repo_root/bureau"

mktemp_dir() {
  local d
  d="$(mktemp -d 2>/dev/null)" || d="$(mktemp -d -t bureau)"
  printf '%s' "$d"
}

date_utc_days_ago() {
  local days="$1"
  if date -u -d "$days days ago" +%F >/dev/null 2>&1; then
    date -u -d "$days days ago" +%F
    return 0
  fi
  date -u -v-"$days"d +%F
}

fail() {
  printf 'FAIL: %s\n' "$*" >&2
  exit 1
}

assert_eq() {
  local got="$1"
  local want="$2"
  if [[ "$got" != "$want" ]]; then
    printf -- '--- got ---\n%s\n--- want ---\n%s\n' "$got" "$want" >&2
    fail "assert_eq failed"
  fi
}

test_status_no_current() {
  local tmp
  tmp="$(mktemp_dir)"
  (cd "$tmp" && mkdir -p _tasks)

  local got
  got="$(cd "$tmp" && "$bureau")"

  local want
  want=$'No current task selected.\n\n[to start new report file] bureau -n <report-suffix>\n[to switch to a new task]  bureau -N <task-suffix>\n[to switch to prior task]  bureau -S <YYYY-MM-DD-task-suffix>\n[to see recent tasks]      bureau -T'

  assert_eq "$got" "$want"
}

test_status_lists_reports() {
  local tmp
  tmp="$(mktemp_dir)"

  local task_dir="2025-10-02-refactor-something"
  mkdir -p "$tmp/_tasks/$task_dir"
  ln -s "$task_dir" "$tmp/_tasks/current"
  touch \
    "$tmp/_tasks/$task_dir/001-user-request.md" \
    "$tmp/_tasks/$task_dir/002-plan.md" \
    "$tmp/_tasks/$task_dir/003-implementation.md" \
    "$tmp/_tasks/$task_dir/004-tests.md"

  local got
  got="$(cd "$tmp" && "$bureau")"

  local want
  want=$'Current task reports dir: _tasks/2025-10-02-refactor-something\n\n[4 reports found]\n001-user-request.md\n002-plan.md\n003-implementation.md\n004-tests.md\n\n[to start new report file] bureau -n <report-suffix>\n[to switch to a new task]  bureau -N <task-suffix>\n[to switch to prior task]  bureau -S <YYYY-MM-DD-task-suffix>\n[to see recent tasks]      bureau -T'

  assert_eq "$got" "$want"
}

test_new_report_file_prints_next() {
  local tmp
  tmp="$(mktemp_dir)"

  local task_dir="2025-10-02-refactor-something"
  mkdir -p "$tmp/_tasks/$task_dir"
  ln -s "$task_dir" "$tmp/_tasks/current"
  touch \
    "$tmp/_tasks/$task_dir/001-user-request.md" \
    "$tmp/_tasks/$task_dir/002-plan.md" \
    "$tmp/_tasks/$task_dir/003-implementation.md" \
    "$tmp/_tasks/$task_dir/004-tests.md"

  local got
  got="$(cd "$tmp" && "$bureau" -n your-suffix)"

  local want
  want=$'Write your report to:\n_tasks/2025-10-02-refactor-something/005-your-suffix.md'

  assert_eq "$got" "$want"
}

test_list_recent_tasks() {
  local tmp
  tmp="$(mktemp_dir)"

  local d1 today
  d1="$(date_utc_days_ago 10)"
  today="$(date -u +%F)"

  mkdir -p \
    "$tmp/_tasks/${d1}-implement-feature" \
    "$tmp/_tasks/${d1}b-fix-bug" \
    "$tmp/_tasks/${today}-refactor-something"

  local got
  got="$(cd "$tmp" && "$bureau" -T)"

  local want
  want="[3 tasks in past month]
${d1}-implement-feature
${d1}b-fix-bug
${today}-refactor-something"

  assert_eq "$got" "$want"
}

test_new_task_uses_today_and_b_suffix() {
  local tmp
  tmp="$(mktemp_dir)"

  local today
  today="$(date -u +%F)"

  local out1 got1
  out1="$(cd "$tmp" && "$bureau" -N first-task)"
  got1="${out1%%$'\n'*}"
  assert_eq "$got1" "Switched to new task $today-first-task."

  local out2 got2
  out2="$(cd "$tmp" && "$bureau" -N second-task)"
  got2="${out2%%$'\n'*}"
  assert_eq "$got2" "Switched to new task ${today}b-second-task."
}

test_switch_task() {
  local tmp
  tmp="$(mktemp_dir)"

  mkdir -p "$tmp/_tasks/2025-10-01-implement-feature"
  mkdir -p "$tmp/_tasks/2025-10-01b-fix-bug"
  ln -s "2025-10-01-implement-feature" "$tmp/_tasks/current"

  local out got
  out="$(cd "$tmp" && "$bureau" -S 2025-10-01b-fix-bug)"
  got="${out%%$'\n'*}"

  assert_eq "$got" "Switched to preexisting task 2025-10-01b-fix-bug."
}

test_help_output_contains_status() {
  local tmp
  tmp="$(mktemp_dir)"
  (cd "$tmp" && mkdir -p _tasks)

  local got
  got="$(cd "$tmp" && "$bureau" -h)"

  local want
  want=$'Bureau - cli tool for managing AI agent report files.\n\nNo current task selected.\n\n[to start new report file] bureau -n <report-suffix>\n[to switch to a new task]  bureau -N <task-suffix>\n[to switch to prior task]  bureau -S <YYYY-MM-DD-task-suffix>\n[to see recent tasks]      bureau -T\n\nAdditional options:\n\nOverride _task subdir name via -d <subdir> or by setting BUREAU_DIR=<subdir>.\nThis can be a name, a relative path or an absolute path.'

  assert_eq "$got" "$want"
}

main() {
  [[ -x "$bureau" ]] || fail "not executable: $bureau"

  test_status_no_current
  test_status_lists_reports
  test_new_report_file_prints_next
  test_list_recent_tasks
  test_new_task_uses_today_and_b_suffix
  test_switch_task
  test_help_output_contains_status

  printf 'OK\n'
}

main "$@"
