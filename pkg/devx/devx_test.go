package devx

import (
	"testing"

	"github.com/devforge/devforge/pkg/engine"
)

func hasCode(d []engine.Diagnostic, code string) bool {
	for _, x := range d {
		if x.Code == code {
			return true
		}
	}
	return false
}

func TestLintDockerfile_LatestTag(t *testing.T) {
	t.Parallel()
	r, _ := LintDockerfile([]byte("FROM nginx:latest\n"))
	if !hasCode(r.Diagnostics, "DOCKER.LINT.LATEST_TAG") {
		t.Fatalf("missing LATEST_TAG: %+v", r.Diagnostics)
	}
}

func TestLintDockerfile_NoTag(t *testing.T) {
	t.Parallel()
	r, _ := LintDockerfile([]byte("FROM nginx\n"))
	if !hasCode(r.Diagnostics, "DOCKER.LINT.NO_TAG") {
		t.Fatalf("missing NO_TAG: %+v", r.Diagnostics)
	}
}

func TestLintDockerfile_NoUser(t *testing.T) {
	t.Parallel()
	r, _ := LintDockerfile([]byte("FROM alpine:3.20\nRUN echo hi\n"))
	if !hasCode(r.Diagnostics, "DOCKER.LINT.RUNS_AS_ROOT") {
		t.Fatalf("missing RUNS_AS_ROOT: %+v", r.Diagnostics)
	}
}

func TestLintDockerfile_MultiCmd(t *testing.T) {
	t.Parallel()
	r, _ := LintDockerfile([]byte("FROM alpine:3\nUSER nobody\nCMD [\"a\"]\nCMD [\"b\"]\n"))
	if !hasCode(r.Diagnostics, "DOCKER.LINT.MULTI_CMD") {
		t.Fatalf("missing MULTI_CMD: %+v", r.Diagnostics)
	}
}

func TestLintDockerfile_NoFrom(t *testing.T) {
	t.Parallel()
	r, _ := LintDockerfile([]byte("RUN echo hi\n"))
	if !hasCode(r.Diagnostics, "DOCKER.LINT.NO_FROM") {
		t.Fatalf("missing NO_FROM: %+v", r.Diagnostics)
	}
}

func TestParseEnv_HappyPath(t *testing.T) {
	t.Parallel()
	r, _ := ParseEnv([]byte("# comment\nFOO=1\nBAR=\"hello world\"\n"), ParseEnvOptions{})
	if r.Values["FOO"] != "1" || r.Values["BAR"] != "hello world" {
		t.Fatalf("got %+v", r.Values)
	}
}

func TestParseEnv_DuplicateKey(t *testing.T) {
	t.Parallel()
	r, _ := ParseEnv([]byte("FOO=1\nFOO=2\n"), ParseEnvOptions{})
	if !hasCode(r.Diagnostics, "ENV.PARSE.DUPLICATE_KEY") {
		t.Fatalf("missing DUPLICATE_KEY")
	}
	if r.Values["FOO"] != "2" {
		t.Fatalf("last-write-wins not honored: %s", r.Values["FOO"])
	}
}

func TestParseEnv_InvalidKey(t *testing.T) {
	t.Parallel()
	r, _ := ParseEnv([]byte("12foo=1\n"), ParseEnvOptions{})
	if !hasCode(r.Diagnostics, "ENV.PARSE.INVALID_KEY") {
		t.Fatalf("missing INVALID_KEY")
	}
}

func TestDiffEnv(t *testing.T) {
	t.Parallel()
	left := []byte("A=1\nB=2\n")
	right := []byte("A=1\nC=3\nB=99\n")
	r, _ := DiffEnv(left, right)
	if len(r.Added) != 1 || r.Added[0] != "C" {
		t.Fatalf("Added = %v", r.Added)
	}
	if len(r.Changed) != 1 || r.Changed[0].Key != "B" {
		t.Fatalf("Changed = %v", r.Changed)
	}
}

func TestValidateK8s_MissingFields(t *testing.T) {
	t.Parallel()
	r, _ := ValidateK8s([]byte("metadata:\n  name: x\n"), K8sOptions{})
	if r.Valid {
		t.Fatalf("expected invalid")
	}
	if !hasCode(r.Diagnostics, "K8S.MISSING.API_VERSION") {
		t.Fatalf("missing api_version diag")
	}
}

func TestValidateK8s_Valid(t *testing.T) {
	t.Parallel()
	yamlDoc := []byte(`apiVersion: v1
kind: Pod
metadata:
  name: my-pod
`)
	r, _ := ValidateK8s(yamlDoc, K8sOptions{})
	if !r.Valid {
		t.Fatalf("expected valid: %+v", r.Diagnostics)
	}
	if r.Kind != "Pod" {
		t.Fatalf("Kind = %q", r.Kind)
	}
}
