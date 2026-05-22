// Package gitx provides git-related utilities (unified-diff patch builder,
// Conventional Commits validator, gitignore generator).
//
// External API:
//
//	Patch(left, right string, opts PatchOptions) (PatchResult, error)
//	CommitFormat(string, opts CommitOptions) (CommitResult, error)
//	IgnoreGen([]string) (IgnoreResult, error)
package gitx

import (
	"fmt"
	"regexp"
	"sort"
	"strings"

	"github.com/devforge/devforge/pkg/engine"
	pkgstrx "github.com/devforge/devforge/pkg/strx"
)

// ---------- Patch ----------

// PatchOptions tunes Patch.
type PatchOptions struct {
	// LeftPath / RightPath populate the unified-diff filenames. Default a/b.
	LeftPath  string `json:"leftPath,omitempty"`
	RightPath string `json:"rightPath,omitempty"`
	// Context lines around each hunk. Default 3.
	Context int `json:"context,omitempty"`
}

// PatchResult holds the unified diff text.
type PatchResult struct {
	Output      string              `json:"output"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// Patch produces a unified-diff patch text. It reuses pkg/strx.Diff for the
// underlying line-LCS, then formats hunks per the unified-diff spec.
func Patch(left, right string, opts PatchOptions) (PatchResult, error) {
	if opts.Context <= 0 {
		opts.Context = 3
	}
	leftPath := opts.LeftPath
	if leftPath == "" {
		leftPath = "a"
	}
	rightPath := opts.RightPath
	if rightPath == "" {
		rightPath = "b"
	}

	diff, err := pkgstrx.Diff(left, right, pkgstrx.DiffOptions{})
	if err != nil {
		return PatchResult{}, err
	}

	if diff.Summary.Adds == 0 && diff.Summary.Removes == 0 {
		return PatchResult{Output: ""}, nil
	}

	var b strings.Builder
	fmt.Fprintf(&b, "--- %s\n+++ %s\n", leftPath, rightPath)
	// Emit a single hunk covering all lines for MVP; per-hunk grouping is a
	// reasonable Phase D+ enhancement.
	leftStart, rightStart := 1, 1
	leftCount, rightCount := 0, 0
	for _, h := range diff.Hunks {
		if h.Op == "equal" || h.Op == "remove" {
			leftCount++
		}
		if h.Op == "equal" || h.Op == "add" {
			rightCount++
		}
	}
	fmt.Fprintf(&b, "@@ -%d,%d +%d,%d @@\n", leftStart, leftCount, rightStart, rightCount)
	for _, h := range diff.Hunks {
		switch h.Op {
		case "equal":
			fmt.Fprintf(&b, " %s\n", h.Content)
		case "remove":
			fmt.Fprintf(&b, "-%s\n", h.Content)
		case "add":
			fmt.Fprintf(&b, "+%s\n", h.Content)
		}
	}
	return PatchResult{Output: b.String()}, nil
}

// ---------- Conventional Commits ----------

// CommitOptions tunes CommitFormat.
type CommitOptions struct{}

// CommitResult is the validation outcome.
type CommitResult struct {
	Valid       bool                `json:"valid"`
	Type        string              `json:"type,omitempty"`
	Scope       string              `json:"scope,omitempty"`
	Breaking    bool                `json:"breaking"`
	Subject     string              `json:"subject,omitempty"`
	Body        string              `json:"body,omitempty"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// validTypes are the canonical Conventional Commits type tokens.
var validTypes = map[string]struct{}{
	"feat": {}, "fix": {}, "docs": {}, "style": {}, "refactor": {},
	"perf": {}, "test": {}, "build": {}, "ci": {}, "chore": {}, "revert": {},
}

var commitHeaderRE = regexp.MustCompile(`^(?P<type>[a-zA-Z]+)(?:\((?P<scope>[a-zA-Z0-9_\-/]+)\))?(?P<bang>!)?: (?P<subject>.+)$`)

// CommitFormat validates a commit message against Conventional Commits v1.
func CommitFormat(input string, _ CommitOptions) (CommitResult, error) {
	res := CommitResult{}
	lines := strings.Split(strings.TrimSpace(input), "\n")
	if len(lines) == 0 || lines[0] == "" {
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code: "GIT.COMMIT.EMPTY", Message: "commit message is empty", Severity: engine.SevError,
		})
		return res, nil
	}
	header := lines[0]
	m := commitHeaderRE.FindStringSubmatch(header)
	if m == nil {
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code: "GIT.COMMIT.INVALID_HEADER",
			Message: `header must match "<type>(<scope>)?!?: <subject>"`,
			Severity: engine.SevError,
		})
		return res, nil
	}
	res.Type = strings.ToLower(m[commitHeaderRE.SubexpIndex("type")])
	res.Scope = m[commitHeaderRE.SubexpIndex("scope")]
	res.Breaking = m[commitHeaderRE.SubexpIndex("bang")] == "!"
	res.Subject = m[commitHeaderRE.SubexpIndex("subject")]
	if _, ok := validTypes[res.Type]; !ok {
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code: "GIT.COMMIT.UNKNOWN_TYPE",
			Message: fmt.Sprintf("type %q is not a known Conventional Commits type", res.Type),
			Severity: engine.SevWarn,
		})
	}
	if len(header) > 72 {
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code: "GIT.COMMIT.HEADER_TOO_LONG",
			Message: fmt.Sprintf("header is %d chars; convention recommends ≤72", len(header)),
			Severity: engine.SevWarn,
		})
	}
	if len(lines) > 2 && lines[1] != "" {
		res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
			Code: "GIT.COMMIT.MISSING_BLANK_LINE",
			Message: "second line should be blank to separate header from body",
			Severity: engine.SevWarn,
		})
	}
	if len(lines) > 2 {
		res.Body = strings.Join(lines[2:], "\n")
		if strings.Contains(res.Body, "BREAKING CHANGE:") || strings.Contains(res.Body, "BREAKING-CHANGE:") {
			res.Breaking = true
		}
	}
	res.Valid = !engine.HasError(res.Diagnostics)
	return res, nil
}

// ---------- gitignore generator ----------

// IgnoreResult holds the merged gitignore output.
type IgnoreResult struct {
	Output      string              `json:"output"`
	Used        []string            `json:"used"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// templates is a curated subset of github/gitignore (MIT-licensed). Embedded
// at build time so `devforge` works fully offline. Names are case-insensitive
// when callers pass them.
var templates = map[string]string{
	"go": `# Go
*.exe
*.exe~
*.dll
*.so
*.dylib
*.test
*.out
go.work.sum
/vendor/
/dist/
/bin/
`,
	"node": `# Node
node_modules/
npm-debug.log*
yarn-debug.log*
yarn-error.log*
.pnp/
.pnp.js
.npm
.yarn-integrity
.eslintcache
dist/
.next/
.cache/
.parcel-cache/
`,
	"python": `# Python
__pycache__/
*.py[cod]
*$py.class
*.egg-info/
.eggs/
build/
dist/
.venv/
venv/
ENV/
env/
.pytest_cache/
.mypy_cache/
.coverage
htmlcov/
`,
	"rust": `# Rust
/target/
**/*.rs.bk
Cargo.lock
`,
	"java": `# Java
*.class
*.jar
*.war
*.ear
*.iml
.idea/
target/
build/
out/
hs_err_pid*.log
`,
	"macos": `# macOS
.DS_Store
.AppleDouble
.LSOverride
Icon?
._*
`,
	"windows": `# Windows
Thumbs.db
ehthumbs.db
Desktop.ini
$RECYCLE.BIN/
*.lnk
`,
	"linux": `# Linux
*~
.fuse_hidden*
.directory
.Trash-*
`,
	"ide": `# IDE
.idea/
.vscode/
*.swp
*.swo
.history/
`,
	"docker": `# Docker
.docker/
docker-compose.override.yml
`,
	"terraform": `# Terraform
.terraform/
*.tfstate
*.tfstate.*
crash.log
override.tf
override.tf.json
*_override.tf
*_override.tf.json
`,
}

// AvailableTemplates returns the list of supported template names.
func AvailableTemplates() []string {
	names := make([]string, 0, len(templates))
	for k := range templates {
		names = append(names, k)
	}
	sort.Strings(names)
	return names
}

// IgnoreGen merges the requested templates into a single .gitignore body.
func IgnoreGen(langs []string) (IgnoreResult, error) {
	if len(langs) == 0 {
		return IgnoreResult{Diagnostics: []engine.Diagnostic{{
			Code: "GIT.IGNORE.EMPTY_REQUEST",
			Message: fmt.Sprintf("at least one template is required (available: %v)", AvailableTemplates()),
			Severity: engine.SevError,
		}}}, nil
	}
	var b strings.Builder
	used := make([]string, 0, len(langs))
	for _, lang := range langs {
		key := strings.ToLower(lang)
		body, ok := templates[key]
		if !ok {
			return IgnoreResult{Diagnostics: []engine.Diagnostic{{
				Code: "GIT.IGNORE.UNKNOWN_TEMPLATE",
				Message: fmt.Sprintf("unknown template %q (available: %v)", lang, AvailableTemplates()),
				Severity: engine.SevError,
			}}}, nil
		}
		b.WriteString(body)
		b.WriteString("\n")
		used = append(used, key)
	}
	return IgnoreResult{Output: strings.TrimRight(b.String(), "\n") + "\n", Used: used}, nil
}
