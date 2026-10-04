package build

import (
	"strings"
	"testing"
)

func TestExtractHeadlines_WhatWentWrong(t *testing.T) {
	log := `
> Task :app:compileReleaseKotlin
e: file:///tmp/NativeEngine.kt:42:5 Unresolved reference: foo

FAILURE: Build failed with an exception.

* What went wrong:
Execution failed for task ':app:compileReleaseKotlin'.
> A failure occurred while executing org.jetbrains.kotlin.compilerRunner.GradleCompilerRunnerWithWorkers$GradleKotlinCompilerWorkAction

* Try:
> Run with --stacktrace option to get the stack trace.

BUILD FAILED in 12s
`
	h := ExtractHeadlines(log)
	if len(h) == 0 {
		t.Fatal("expected headlines, got none")
	}
	joined := strings.Join(h, "\n")
	if !strings.Contains(joined, "What went wrong") {
		t.Errorf("missing What went wrong:\n%s", joined)
	}
	if !strings.Contains(joined, "compileReleaseKotlin") {
		t.Errorf("missing task name:\n%s", joined)
	}
	if !strings.Contains(joined, "Unresolved reference") && !strings.Contains(joined, "e: file:///") {
		t.Logf("headlines:\n%s", joined)
	}
}

func TestExtractHeadlines_Empty(t *testing.T) {
	h := ExtractHeadlines("")
	if len(h) != 0 {
		t.Errorf("expected empty, got %v", h)
	}
}

func TestExpectedAPKPaths(t *testing.T) {
	debug := ExpectedAPKPaths("/proj", VariantDebug)
	if len(debug) != 1 || !strings.HasSuffix(debug[0], "app-debug.apk") {
		t.Errorf("debug paths: %v", debug)
	}
	rel := ExpectedAPKPaths("/proj", VariantRelease)
	if len(rel) < 2 {
		t.Errorf("release should list unsigned + signed candidates: %v", rel)
	}
}
