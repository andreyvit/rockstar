package bureau

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"
)

var taskDirNameRe = regexp.MustCompile(`^(\d{4}-\d{2}-\d{2}(?:[b-y]|z\d{3,})?)-(.+)$`)
var reportFileRe = regexp.MustCompile(`^\d+-.*\.md$`)
var reportNumberRe = regexp.MustCompile(`^(\d+)-`)

type ParsedTaskDirName struct {
	DatePrefix string
	Slug       string
}

func DateSuffix(index int, now time.Time) (string, error) {
	if index < 0 {
		return "", errors.New("index must be >= 0")
	}

	dateStr := now.UTC().Format("2006-01-02")
	if index == 0 {
		return dateStr, nil
	}
	if index <= 24 {
		return dateStr + string(rune('a'+index)), nil
	}
	return dateStr + "z" + fmt.Sprintf("%03d", index+1), nil
}

func ParseTaskDirName(dirName string) (ParsedTaskDirName, bool) {
	m := taskDirNameRe.FindStringSubmatch(dirName)
	if m == nil {
		return ParsedTaskDirName{}, false
	}
	return ParsedTaskDirName{DatePrefix: m[1], Slug: m[2]}, true
}

func CurrentLinkPath(tasksRoot string) string {
	return filepath.Join(tasksRoot, "current")
}

func GetAllTaskDirs(tasksRoot string) ([]string, error) {
	if err := os.MkdirAll(tasksRoot, 0o755); err != nil {
		return nil, err
	}

	entries, err := os.ReadDir(tasksRoot)
	if err != nil {
		return []string{}, nil
	}

	var dirs []string
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if _, ok := ParseTaskDirName(name); !ok {
			continue
		}
		dirs = append(dirs, name)
	}
	sort.Strings(dirs)
	return dirs, nil
}

func GetRecentTaskDirs(tasksRoot string, now time.Time) ([]string, error) {
	allDirs, err := GetAllTaskDirs(tasksRoot)
	if err != nil {
		return nil, err
	}

	thirtyDaysAgo := now.AddDate(0, 0, -30)
	cutoffStr := thirtyDaysAgo.UTC().Format("2006-01-02")

	var recent []string
	for _, dirName := range allDirs {
		parsed, ok := ParseTaskDirName(dirName)
		if !ok {
			continue
		}

		dateOnly := parsed.DatePrefix
		if len(dateOnly) > 0 {
			last := dateOnly[len(dateOnly)-1]
			if last >= 'b' && last <= 'z' {
				dateOnly = dateOnly[:len(dateOnly)-1]
			}
		}

		if dateOnly >= cutoffStr {
			recent = append(recent, dirName)
		}
	}
	return recent, nil
}

func GetCurrentTaskDir(tasksRoot string) (string, bool, error) {
	target, err := os.Readlink(CurrentLinkPath(tasksRoot))
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return "", false, nil
		}
		return "", false, err
	}

	targetPath := target
	if !filepath.IsAbs(targetPath) {
		targetPath = filepath.Join(tasksRoot, targetPath)
	}
	return filepath.Base(targetPath), true, nil
}

func GetReportFiles(tasksRoot, taskDir string) ([]string, error) {
	taskPath := filepath.Join(tasksRoot, taskDir)
	entries, err := os.ReadDir(taskPath)
	if err != nil {
		return []string{}, nil
	}

	var reportFiles []string
	for _, entry := range entries {
		name := entry.Name()
		if reportFileRe.MatchString(name) {
			reportFiles = append(reportFiles, name)
		}
	}
	sort.Strings(reportFiles)

	if len(reportFiles) <= 50 {
		return reportFiles, nil
	}
	return append(append([]string{}, reportFiles[:20]...), reportFiles[len(reportFiles)-30:]...), nil
}

func FindNextTaskDirName(tasksRoot, slug string, now time.Time) (string, error) {
	allDirs, err := GetAllTaskDirs(tasksRoot)
	if err != nil {
		return "", err
	}
	existing := make(map[string]struct{}, len(allDirs))
	for _, dir := range allDirs {
		existing[dir] = struct{}{}
	}

	today := now.UTC().Format("2006-01-02")

	for i := 0; i < 1000; i++ {
		datePrefix, err := DateSuffix(i, now)
		if err != nil {
			return "", err
		}
		if !strings.HasPrefix(datePrefix, today) {
			break
		}

		candidate := datePrefix + "-" + slug
		if _, ok := existing[candidate]; !ok {
			return candidate, nil
		}
	}

	return "", errors.New("too many tasks for today (max 1000)")
}

func UpdateCurrentSymlink(tasksRoot, taskDir string) error {
	if err := os.MkdirAll(tasksRoot, 0o755); err != nil {
		return err
	}

	currentLink := CurrentLinkPath(tasksRoot)
	_ = os.Remove(currentLink)

	return os.Symlink(taskDir, currentLink)
}

func FindNextReportNumber(tasksRoot, taskDir string) (int, error) {
	reportFiles, err := GetReportFiles(tasksRoot, taskDir)
	if err != nil {
		return 0, err
	}
	if len(reportFiles) == 0 {
		return 1, nil
	}

	maxNumber := 0
	for _, name := range reportFiles {
		m := reportNumberRe.FindStringSubmatch(name)
		if m == nil {
			continue
		}
		n, err := strconv.Atoi(m[1])
		if err != nil {
			continue
		}
		if n > maxNumber {
			maxNumber = n
		}
	}
	return maxNumber + 1, nil
}

func FindTaskDirBySlug(tasksRoot, slug string) (string, bool, error) {
	allDirs, err := GetAllTaskDirs(tasksRoot)
	if err != nil {
		return "", false, err
	}

	for i := len(allDirs) - 1; i >= 0; i-- {
		parsed, ok := ParseTaskDirName(allDirs[i])
		if ok && parsed.Slug == slug {
			return allDirs[i], true, nil
		}
	}
	return "", false, nil
}

func NewReportFileName(nextNumber int, suffix string) (string, error) {
	if nextNumber <= 0 {
		return "", errors.New("nextNumber must be >= 1")
	}
	if strings.TrimSpace(suffix) == "" {
		return "", errors.New("suffix must be non-empty")
	}
	return fmt.Sprintf("%03d-%s.md", nextNumber, suffix), nil
}
