package bureau

import (
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"
)

func TestDateSuffix(t *testing.T) {
	now := time.Date(2025, 10, 2, 12, 0, 0, 0, time.UTC)

	s, err := DateSuffix(0, now)
	if err != nil {
		t.Fatal(err)
	}
	if s != "2025-10-02" {
		t.Fatalf("got %q", s)
	}

	s, _ = DateSuffix(1, now)
	if s != "2025-10-02b" {
		t.Fatalf("got %q", s)
	}

	s, _ = DateSuffix(2, now)
	if s != "2025-10-02c" {
		t.Fatalf("got %q", s)
	}

	s, _ = DateSuffix(24, now)
	if s != "2025-10-02y" {
		t.Fatalf("got %q", s)
	}

	s, _ = DateSuffix(25, now)
	if s != "2025-10-02z026" {
		t.Fatalf("got %q", s)
	}

	s, _ = DateSuffix(26, now)
	if s != "2025-10-02z027" {
		t.Fatalf("got %q", s)
	}

	s, _ = DateSuffix(99, now)
	if s != "2025-10-02z100" {
		t.Fatalf("got %q", s)
	}

	s, _ = DateSuffix(999, now)
	if s != "2025-10-02z1000" {
		t.Fatalf("got %q", s)
	}
}

func TestParseTaskDirName(t *testing.T) {
	got, ok := ParseTaskDirName("2025-10-01-some-urgent-task")
	if !ok {
		t.Fatal("expected ok")
	}
	if got.DatePrefix != "2025-10-01" || got.Slug != "some-urgent-task" {
		t.Fatalf("got %#v", got)
	}

	got, ok = ParseTaskDirName("2025-10-01b-second-task")
	if !ok {
		t.Fatal("expected ok")
	}
	if got.DatePrefix != "2025-10-01b" || got.Slug != "second-task" {
		t.Fatalf("got %#v", got)
	}

	got, ok = ParseTaskDirName("2025-10-01z026-many-tasks")
	if !ok {
		t.Fatal("expected ok")
	}
	if got.DatePrefix != "2025-10-01z026" || got.Slug != "many-tasks" {
		t.Fatalf("got %#v", got)
	}

	if _, ok := ParseTaskDirName("invalid-name"); ok {
		t.Fatal("expected not ok")
	}
}

func TestGetAllTaskDirs(t *testing.T) {
	root := t.TempDir()
	tasksRoot := filepath.Join(root, "_tasks")

	if err := os.MkdirAll(filepath.Join(tasksRoot, "2025-10-01-one"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tasksRoot, "2025-10-01b-two"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tasksRoot, "not-a-task"), 0o755); err != nil {
		t.Fatal(err)
	}

	dirs, err := GetAllTaskDirs(tasksRoot)
	if err != nil {
		t.Fatal(err)
	}

	if len(dirs) != 2 || dirs[0] != "2025-10-01-one" || dirs[1] != "2025-10-01b-two" {
		t.Fatalf("got %#v", dirs)
	}
}

func TestGetCurrentTaskDir(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("symlinks require elevated permissions on some Windows setups")
	}

	root := t.TempDir()
	tasksRoot := filepath.Join(root, "_tasks")
	if err := os.MkdirAll(tasksRoot, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tasksRoot, "2025-10-01-test-task"), 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.Symlink("2025-10-01-test-task", CurrentLinkPath(tasksRoot)); err != nil {
		t.Fatal(err)
	}

	dir, ok, err := GetCurrentTaskDir(tasksRoot)
	if err != nil {
		t.Fatal(err)
	}
	if !ok {
		t.Fatal("expected ok")
	}
	if dir != "2025-10-01-test-task" {
		t.Fatalf("got %q", dir)
	}
}

func TestFindTaskDirBySlug(t *testing.T) {
	root := t.TempDir()
	tasksRoot := filepath.Join(root, "_tasks")
	if err := os.MkdirAll(filepath.Join(tasksRoot, "2025-10-01-my-task"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tasksRoot, "2025-10-02-another-task"), 0o755); err != nil {
		t.Fatal(err)
	}

	dir, ok, err := FindTaskDirBySlug(tasksRoot, "my-task")
	if err != nil {
		t.Fatal(err)
	}
	if !ok || dir != "2025-10-01-my-task" {
		t.Fatalf("got %q ok=%v", dir, ok)
	}

	_, ok, err = FindTaskDirBySlug(tasksRoot, "non-existent")
	if err != nil {
		t.Fatal(err)
	}
	if ok {
		t.Fatal("expected not ok")
	}
}

func TestGetRecentTaskDirs(t *testing.T) {
	now := time.Date(2025, 10, 2, 12, 0, 0, 0, time.UTC)
	recent := now.AddDate(0, 0, -10).Format("2006-01-02")
	old := now.AddDate(0, 0, -40).Format("2006-01-02")

	root := t.TempDir()
	tasksRoot := filepath.Join(root, "_tasks")
	if err := os.MkdirAll(filepath.Join(tasksRoot, recent+"-recent-task"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(tasksRoot, old+"-old-task"), 0o755); err != nil {
		t.Fatal(err)
	}

	dirs, err := GetRecentTaskDirs(tasksRoot, now)
	if err != nil {
		t.Fatal(err)
	}
	if len(dirs) != 1 || dirs[0] != recent+"-recent-task" {
		t.Fatalf("got %#v", dirs)
	}
}

func TestFindNextReportNumber(t *testing.T) {
	root := t.TempDir()
	tasksRoot := filepath.Join(root, "_tasks")
	taskDir := "2025-10-01-my-task"
	taskPath := filepath.Join(tasksRoot, taskDir)
	if err := os.MkdirAll(taskPath, 0o755); err != nil {
		t.Fatal(err)
	}

	next, err := FindNextReportNumber(tasksRoot, taskDir)
	if err != nil {
		t.Fatal(err)
	}
	if next != 1 {
		t.Fatalf("got %d", next)
	}

	if err := os.WriteFile(filepath.Join(taskPath, "001-start.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(taskPath, "002-plan.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(taskPath, "003-review.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	next, err = FindNextReportNumber(tasksRoot, taskDir)
	if err != nil {
		t.Fatal(err)
	}
	if next != 4 {
		t.Fatalf("got %d", next)
	}
}

func TestFindNextReportNumber_Gaps(t *testing.T) {
	root := t.TempDir()
	tasksRoot := filepath.Join(root, "_tasks")
	taskDir := "2025-10-01-my-task"
	taskPath := filepath.Join(tasksRoot, taskDir)
	if err := os.MkdirAll(taskPath, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(taskPath, "001-start.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(taskPath, "005-review.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	next, err := FindNextReportNumber(tasksRoot, taskDir)
	if err != nil {
		t.Fatal(err)
	}
	if next != 6 {
		t.Fatalf("got %d", next)
	}
}

func TestFindNextReportNumber_NonStandardPrefixes(t *testing.T) {
	root := t.TempDir()
	tasksRoot := filepath.Join(root, "_tasks")
	taskDir := "2025-10-01-my-task"
	taskPath := filepath.Join(tasksRoot, taskDir)
	if err := os.MkdirAll(taskPath, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(taskPath, "001-foo.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(taskPath, "002-bar.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(taskPath, "42-boz.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	next, err := FindNextReportNumber(tasksRoot, taskDir)
	if err != nil {
		t.Fatal(err)
	}
	if next != 43 {
		t.Fatalf("got %d", next)
	}
}

func TestGetReportFiles_VariableLengthNumericPrefixes(t *testing.T) {
	root := t.TempDir()
	tasksRoot := filepath.Join(root, "_tasks")
	taskDir := "2025-10-01-my-task"
	taskPath := filepath.Join(tasksRoot, taskDir)
	if err := os.MkdirAll(taskPath, 0o755); err != nil {
		t.Fatal(err)
	}

	if err := os.WriteFile(filepath.Join(taskPath, "1-first.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(taskPath, "11-eleventh.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(taskPath, "100-hundredth.md"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}

	files, err := GetReportFiles(tasksRoot, taskDir)
	if err != nil {
		t.Fatal(err)
	}
	if got, want := len(files), 3; got != want {
		t.Fatalf("got %d files: %#v", got, files)
	}
	if files[0] != "1-first.md" || files[1] != "100-hundredth.md" || files[2] != "11-eleventh.md" {
		t.Fatalf("got %#v", files)
	}

	next, err := FindNextReportNumber(tasksRoot, taskDir)
	if err != nil {
		t.Fatal(err)
	}
	if next != 101 {
		t.Fatalf("got %d", next)
	}
}

func TestGetReportFiles_Truncation(t *testing.T) {
	root := t.TempDir()
	tasksRoot := filepath.Join(root, "_tasks")
	taskDir := "2025-10-01-my-task"
	taskPath := filepath.Join(tasksRoot, taskDir)
	if err := os.MkdirAll(taskPath, 0o755); err != nil {
		t.Fatal(err)
	}

	for i := 1; i <= 100; i++ {
		name := filepath.Join(taskPath, fmt3(i)+"-file.md")
		if err := os.WriteFile(name, []byte("x"), 0o644); err != nil {
			t.Fatal(err)
		}
	}

	files, err := GetReportFiles(tasksRoot, taskDir)
	if err != nil {
		t.Fatal(err)
	}
	if len(files) != 50 {
		t.Fatalf("got %d", len(files))
	}
	if files[0] != "001-file.md" || files[19] != "020-file.md" || files[20] != "071-file.md" || files[49] != "100-file.md" {
		t.Fatalf("got %#v", files)
	}
}

func fmt3(n int) string {
	s := "000" + strconv.Itoa(n)
	return s[len(s)-3:]
}
