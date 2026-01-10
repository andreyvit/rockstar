#!/usr/bin/env bash

set -euo pipefail

repo_root="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
bureau="$repo_root/bureau"

mktemp_dir() {
  local d
  d="$(mktemp -d 2>/dev/null)" || d="$(mktemp -d -t bureau)"
  printf '%s' "$d"
}

fail() { printf 'FAIL: %s\n' "$*" >&2; exit 1; }

assert_success() {
  local exit_code="$1"
  if [[ "$exit_code" -ne 0 ]]; then
    fail "expected success, got exit code $exit_code"
  fi
}

assert_contains() {
  local haystack="$1"
  local needle="$2"
  if [[ "$haystack" != *"$needle"* ]]; then
    fail "expected output to contain: $needle"
  fi
}

assert_not_contains() {
  local haystack="$1"
  local needle="$2"
  if [[ "$haystack" == *"$needle"* ]]; then
    fail "expected output to not contain: $needle"
  fi
}

assert_symlink_target() {
  local link_path="$1"
  local want="$2"
  [[ -L "$link_path" ]] || fail "expected symlink: $link_path"
  local got
  # readlink prints the symlink target without resolving it.
  got="$(readlink "$link_path")"
  [[ "$got" == "$want" ]] || fail "symlink $link_path points to $got, want $want"
}

assert_dir_exists() { [[ -d "$1" ]] || fail "expected directory: $1"; }

test_status_no_current_does_not_error() {
  local tmp
  tmp="$(mktemp_dir)"
  (cd "$tmp" && mkdir -p _tasks)

  local out
  out="$(cd "$tmp" && "$bureau")"
  assert_contains "$out" "No current task"
}

test_status_with_current_does_not_error() {
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

  local out
  out="$(cd "$tmp" && "$bureau")"
  assert_contains "$out" "$task_dir"
}

test_new_report_file_returns_path_and_does_not_create_file() {
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

  local out path
  out="$(cd "$tmp" && "$bureau" -n your-suffix)"
  path="${out##*$'\n'}"

  [[ "$path" == */005-your-suffix.md ]] || fail "unexpected report path: $path"
  if [[ "$path" == /* ]]; then
    if [[ -e "$path" ]]; then
      fail "report file should not be created: $path"
    fi
  else
    if [[ -e "$tmp/$path" ]]; then
      fail "report file should not be created: $path"
    fi
  fi
}

test_list_tasks_returns_last_10() {
  local tmp
  tmp="$(mktemp_dir)"

  mkdir -p "$tmp/_tasks"
  local i
  for ((i = 1; i <= 12; i++)); do
    mkdir -p "$tmp/_tasks/2025-10-$(printf '%02d' "$i")-task-$i"
  done

  local out
  out="$(cd "$tmp" && "$bureau" -T)"

  assert_not_contains "$out" "2025-10-01-task-1"
  assert_not_contains "$out" "2025-10-02-task-2"
  assert_contains "$out" "2025-10-03-task-3"
  assert_contains "$out" "2025-10-12-task-12"
}

test_new_task_creates_dir_and_updates_current() {
  local tmp
  tmp="$(mktemp_dir)"

  local today
  today="$(date +%F)"

  local out
  out="$(cd "$tmp" && "$bureau" -N first-task)"
  assert_contains "$out" "Switched to new task"

  assert_symlink_target "$tmp/_tasks/current" "$today-first-task"
  assert_dir_exists "$tmp/_tasks/$today-first-task"

  (cd "$tmp" && "$bureau" -N second-task >/dev/null)
  assert_symlink_target "$tmp/_tasks/current" "${today}b-second-task"
  assert_dir_exists "$tmp/_tasks/${today}b-second-task"
}

test_switch_task_updates_current() {
  local tmp
  tmp="$(mktemp_dir)"

  mkdir -p "$tmp/_tasks/2025-10-01-implement-feature"
  mkdir -p "$tmp/_tasks/2025-10-01b-fix-bug"
  ln -s "2025-10-01-implement-feature" "$tmp/_tasks/current"
  assert_symlink_target "$tmp/_tasks/current" "2025-10-01-implement-feature"

  (cd "$tmp" && "$bureau" -S 2025-10-01b-fix-bug >/dev/null)
  assert_symlink_target "$tmp/_tasks/current" "2025-10-01b-fix-bug"
}

test_help_mentions_bureau_dir() {
  local tmp
  tmp="$(mktemp_dir)"
  (cd "$tmp" && mkdir -p _tasks)

  local out
  out="$(cd "$tmp" && "$bureau" -h)"
  assert_contains "$out" "Bureau - cli tool"
  assert_contains "$out" "BUREAU_DIR"
}

main() {
  [[ -x "$bureau" ]] || fail "not executable: $bureau"

  test_status_no_current_does_not_error
  test_status_with_current_does_not_error
  test_new_report_file_returns_path_and_does_not_create_file
  test_list_tasks_returns_last_10
  test_new_task_creates_dir_and_updates_current
  test_switch_task_updates_current
  test_help_mentions_bureau_dir

  printf 'OK\n'
}

main "$@"
