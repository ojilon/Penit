package build

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

// Result is the outcome of a Gradle invocation.
type Result struct {
	Task       string
	Variant    Variant // empty for pure test tasks
	ExitCode   int
	LogPath    string
	LogBody    string
	Headlines  []string
	APKs       []APKInfo
	Expected   []string
	Cancelled  bool
	ProjectKey string
}

// Options configures a Gradle run.
type Options struct {
	ProjectRoot    string
	ProjectKey     string
	Task           string // e.g. assembleDebug, testDebugUnitTest
	Variant        Variant
	GradleUserHome string
	LogPath        string    // where to write the full log
	Stdout         io.Writer // live stream (defaults to os.Stdout)
	Stderr         io.Writer // unused; stdout+stderr are merged
}

// GradlewPath returns the platform-appropriate wrapper path, or "gradle" fallback.
func GradlewPath(projectRoot string) string {
	name := "gradlew"
	if runtime.GOOS == "windows" {
		name = "gradlew.bat"
	}
	p := filepath.Join(projectRoot, name)
	if st, err := os.Stat(p); err == nil && !st.IsDir() {
		return p
	}
	return "gradle"
}

// Run executes a Gradle task with live streaming and log capture.
func Run(opts Options) (*Result, error) {
	if opts.ProjectRoot == "" {
		return nil, fmt.Errorf("project root is empty")
	}
	if opts.Task == "" {
		return nil, fmt.Errorf("gradle task is empty")
	}
	out := opts.Stdout
	if out == nil {
		out = os.Stdout
	}

	gradlew := GradlewPath(opts.ProjectRoot)
	args := []string{opts.Task}
	if opts.GradleUserHome != "" {
		args = append(args, "-g", opts.GradleUserHome)
	}

	var cmd *exec.Cmd
	if runtime.GOOS == "windows" && strings.HasSuffix(strings.ToLower(gradlew), ".bat") {
		all := append([]string{"/C", gradlew}, args...)
		cmd = exec.Command("cmd", all...)
	} else {
		cmd = exec.Command(gradlew, args...)
	}
	cmd.Dir = opts.ProjectRoot

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, fmt.Errorf("stdout pipe: %w", err)
	}
	cmd.Stderr = cmd.Stdout // merge

	var logFile *os.File
	if opts.LogPath != "" {
		if err := os.MkdirAll(filepath.Dir(opts.LogPath), 0o755); err != nil {
			return nil, fmt.Errorf("mkdir logs: %w", err)
		}
		logFile, err = os.Create(opts.LogPath)
		if err != nil {
			return nil, fmt.Errorf("create log: %w", err)
		}
		defer logFile.Close()
	}

	var body strings.Builder
	var mu sync.Mutex

	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("start gradle: %w", err)
	}

	sc := bufio.NewScanner(stdout)
	buf := make([]byte, 0, 64*1024)
	sc.Buffer(buf, 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		mu.Lock()
		body.WriteString(line)
		body.WriteByte('\n')
		mu.Unlock()
		fmt.Fprintln(out, line)
		if logFile != nil {
			fmt.Fprintln(logFile, line)
		}
	}

	waitErr := cmd.Wait()
	exitCode := 0
	if waitErr != nil {
		if ee, ok := waitErr.(*exec.ExitError); ok {
			exitCode = ee.ExitCode()
		} else {
			exitCode = 1
		}
	}

	logBody := body.String()
	res := &Result{
		Task:       opts.Task,
		Variant:    opts.Variant,
		ExitCode:   exitCode,
		LogPath:    opts.LogPath,
		LogBody:    logBody,
		ProjectKey: opts.ProjectKey,
	}

	if opts.Variant != "" {
		res.Expected = ExpectedAPKPaths(opts.ProjectRoot, opts.Variant)
		for i, p := range res.Expected {
			if abs, err := filepath.Abs(p); err == nil {
				res.Expected[i] = abs
			}
		}
		if exitCode == 0 {
			res.APKs = FindAPKs(opts.ProjectRoot, opts.Variant)
		}
	}

	if exitCode != 0 {
		res.Headlines = ExtractHeadlines(logBody)
	}

	return res, nil
}

// TaskForVariant maps a variant to the standard assemble task.
func TaskForVariant(v Variant) string {
	switch v {
	case VariantDebug:
		return "assembleDebug"
	case VariantRelease:
		return "assembleRelease"
	default:
		return "assemble" + strings.ToUpper(string(v[:1])) + string(v[1:])
	}
}

// DefaultTestTask is the unit-test task used by `penit test`.
const DefaultTestTask = "testDebugUnitTest"
