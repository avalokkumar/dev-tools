// Package devx provides DevOps utilities (Dockerfile lint, .env parse +
// diff, lightweight Kubernetes YAML structural validation).
//
// External API:
//
//	LintDockerfile([]byte) (LintResult, error)
//	ParseEnv([]byte, ParseEnvOptions) (EnvResult, error)
//	DiffEnv([]byte, []byte) (EnvDiffResult, error)
//	ValidateK8s([]byte, K8sOptions) (K8sResult, error)
package devx

import (
	"bufio"
	"bytes"
	"fmt"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/devforge/devforge/pkg/engine"
)

// ---------- Dockerfile linter ----------

// LintResult is a list of findings keyed by rule code + line number.
type LintResult struct {
	Diagnostics []engine.Diagnostic `json:"diagnostics"`
}

// LintDockerfile applies a small set of best-practice rules.
func LintDockerfile(input []byte) (LintResult, error) {
	res := LintResult{}
	sc := bufio.NewScanner(bytes.NewReader(input))
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	lineNum := 0
	hasUser := false
	hasFrom := false
	entrypointCount := 0
	cmdCount := 0
	for sc.Scan() {
		lineNum++
		raw := sc.Text()
		line := strings.TrimSpace(raw)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		upper := strings.ToUpper(line)
		switch {
		case strings.HasPrefix(upper, "FROM "):
			hasFrom = true
			tag := strings.TrimSpace(line[5:])
			if !strings.Contains(tag, ":") {
				res.Diagnostics = append(res.Diagnostics, mkDiag("DOCKER.LINT.NO_TAG",
					fmt.Sprintf("FROM image %q has no explicit tag (defaults to :latest)", tag),
					engine.SevWarn, lineNum))
			} else if strings.HasSuffix(tag, ":latest") {
				res.Diagnostics = append(res.Diagnostics, mkDiag("DOCKER.LINT.LATEST_TAG",
					"FROM uses :latest; pin a specific version for reproducibility",
					engine.SevWarn, lineNum))
			}
		case strings.HasPrefix(upper, "USER "):
			hasUser = true
		case strings.HasPrefix(upper, "ENTRYPOINT"):
			entrypointCount++
		case strings.HasPrefix(upper, "CMD"):
			cmdCount++
		case strings.HasPrefix(upper, "ADD "):
			res.Diagnostics = append(res.Diagnostics, mkDiag("DOCKER.LINT.PREFER_COPY",
				"ADD has surprising URL/tar handling; prefer COPY for files",
				engine.SevWarn, lineNum))
		case strings.HasPrefix(upper, "RUN "):
			if strings.Contains(upper, "APT-GET INSTALL") && !strings.Contains(upper, "RM -RF /VAR/LIB/APT") {
				res.Diagnostics = append(res.Diagnostics, mkDiag("DOCKER.LINT.APT_LIST_LEFTOVER",
					"apt-get install without cleaning /var/lib/apt/lists bloats the image",
					engine.SevWarn, lineNum))
			}
			if strings.Contains(upper, "APT-GET INSTALL") && !strings.Contains(upper, "--NO-INSTALL-RECOMMENDS") {
				res.Diagnostics = append(res.Diagnostics, mkDiag("DOCKER.LINT.APT_RECOMMENDS",
					"apt-get install should use --no-install-recommends to keep images small",
					engine.SevInfo, lineNum))
			}
			if strings.Contains(upper, "CURL ") && !strings.Contains(upper, "-FSSL") && !strings.Contains(upper, "-SSF") {
				res.Diagnostics = append(res.Diagnostics, mkDiag("DOCKER.LINT.CURL_FLAGS",
					"prefer curl -fsSL for non-interactive downloads (fail on HTTP error)",
					engine.SevInfo, lineNum))
			}
		}
	}
	if !hasFrom {
		res.Diagnostics = append(res.Diagnostics, mkDiag("DOCKER.LINT.NO_FROM",
			"Dockerfile has no FROM instruction", engine.SevError, 0))
	}
	if !hasUser {
		res.Diagnostics = append(res.Diagnostics, mkDiag("DOCKER.LINT.RUNS_AS_ROOT",
			"no USER directive — container will run as root", engine.SevWarn, 0))
	}
	if entrypointCount > 1 {
		res.Diagnostics = append(res.Diagnostics, mkDiag("DOCKER.LINT.MULTI_ENTRYPOINT",
			fmt.Sprintf("found %d ENTRYPOINT directives; only the last takes effect", entrypointCount),
			engine.SevWarn, 0))
	}
	if cmdCount > 1 {
		res.Diagnostics = append(res.Diagnostics, mkDiag("DOCKER.LINT.MULTI_CMD",
			fmt.Sprintf("found %d CMD directives; only the last takes effect", cmdCount),
			engine.SevWarn, 0))
	}
	return res, nil
}

func mkDiag(code, msg string, sev engine.Severity, line int) engine.Diagnostic {
	d := engine.Diagnostic{Code: code, Message: msg, Severity: sev}
	if line > 0 {
		d.Span = &engine.Span{StartLine: line, EndLine: line}
	}
	return d
}

// ---------- .env parser + diff ----------

// ParseEnvOptions tunes ParseEnv.
type ParseEnvOptions struct {
	// AllowExport accepts lines like "export FOO=1".
	AllowExport bool `json:"allowExport,omitempty"`
}

// EnvResult is the parsed key/value map plus warnings.
type EnvResult struct {
	Values      map[string]string   `json:"values"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// ParseEnv parses a dotenv-style file. Reports duplicate keys + malformed lines.
func ParseEnv(input []byte, opts ParseEnvOptions) (EnvResult, error) {
	res := EnvResult{Values: map[string]string{}}
	sc := bufio.NewScanner(bytes.NewReader(input))
	sc.Buffer(make([]byte, 0, 64*1024), 1<<20)
	lineNum := 0
	for sc.Scan() {
		lineNum++
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if opts.AllowExport && strings.HasPrefix(line, "export ") {
			line = strings.TrimSpace(strings.TrimPrefix(line, "export "))
		}
		eq := strings.Index(line, "=")
		if eq < 0 {
			res.Diagnostics = append(res.Diagnostics, mkDiag("ENV.PARSE.NO_EQUALS",
				fmt.Sprintf("line %d has no '='", lineNum), engine.SevError, lineNum))
			continue
		}
		key := strings.TrimSpace(line[:eq])
		val := strings.TrimSpace(line[eq+1:])
		// Strip matching surrounding quotes.
		if len(val) >= 2 && (val[0] == '"' && val[len(val)-1] == '"' || val[0] == '\'' && val[len(val)-1] == '\'') {
			val = val[1 : len(val)-1]
		}
		if !validKey(key) {
			res.Diagnostics = append(res.Diagnostics, mkDiag("ENV.PARSE.INVALID_KEY",
				fmt.Sprintf("key %q is not a valid identifier", key), engine.SevError, lineNum))
			continue
		}
		if _, exists := res.Values[key]; exists {
			res.Diagnostics = append(res.Diagnostics, mkDiag("ENV.PARSE.DUPLICATE_KEY",
				fmt.Sprintf("duplicate key %q (later value wins)", key),
				engine.SevWarn, lineNum))
		}
		res.Values[key] = val
	}
	return res, nil
}

func validKey(k string) bool {
	if k == "" {
		return false
	}
	for i, r := range k {
		ok := r == '_' ||
			(r >= 'A' && r <= 'Z') ||
			(r >= 'a' && r <= 'z') ||
			(i > 0 && r >= '0' && r <= '9')
		if !ok {
			return false
		}
	}
	return true
}

// EnvDiffResult lists per-key differences between two env files.
type EnvDiffResult struct {
	Added       []string            `json:"added"`
	Removed     []string            `json:"removed"`
	Changed     []EnvChange         `json:"changed"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// EnvChange is a single key whose value changed.
type EnvChange struct {
	Key   string `json:"key"`
	Left  string `json:"left"`
	Right string `json:"right"`
}

// DiffEnv computes per-key add/remove/change between two .env blobs.
func DiffEnv(left, right []byte) (EnvDiffResult, error) {
	l, _ := ParseEnv(left, ParseEnvOptions{AllowExport: true})
	r, _ := ParseEnv(right, ParseEnvOptions{AllowExport: true})
	res := EnvDiffResult{}
	keys := map[string]struct{}{}
	for k := range l.Values {
		keys[k] = struct{}{}
	}
	for k := range r.Values {
		keys[k] = struct{}{}
	}
	keyList := make([]string, 0, len(keys))
	for k := range keys {
		keyList = append(keyList, k)
	}
	sort.Strings(keyList)
	for _, k := range keyList {
		lv, lok := l.Values[k]
		rv, rok := r.Values[k]
		switch {
		case lok && !rok:
			res.Removed = append(res.Removed, k)
		case !lok && rok:
			res.Added = append(res.Added, k)
		case lv != rv:
			res.Changed = append(res.Changed, EnvChange{Key: k, Left: lv, Right: rv})
		}
	}
	return res, nil
}

// ---------- K8s YAML structural validate ----------

// K8sOptions tunes ValidateK8s.
type K8sOptions struct{}

// K8sResult is the validation outcome.
type K8sResult struct {
	Valid       bool                `json:"valid"`
	Kind        string              `json:"kind,omitempty"`
	APIVersion  string              `json:"apiVersion,omitempty"`
	Diagnostics []engine.Diagnostic `json:"diagnostics,omitempty"`
}

// ValidateK8s performs a structural check: well-formed YAML, presence of
// apiVersion / kind / metadata.name. Full schema validation is deferred to
// a Phase D+ ADR (would require a sizeable bundled OpenAPI subset).
func ValidateK8s(input []byte, _ K8sOptions) (K8sResult, error) {
	res := K8sResult{Valid: true}
	docs := bytes.Split(input, []byte("\n---\n"))
	for i, d := range docs {
		if len(bytes.TrimSpace(d)) == 0 {
			continue
		}
		var raw map[string]any
		if err := yaml.Unmarshal(d, &raw); err != nil {
			res.Valid = false
			res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
				Code: "K8S.YAML.PARSE",
				Message: fmt.Sprintf("doc %d: %v", i+1, err),
				Severity: engine.SevError,
			})
			continue
		}
		api, _ := raw["apiVersion"].(string)
		kind, _ := raw["kind"].(string)
		if api == "" {
			res.Valid = false
			res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
				Code: "K8S.MISSING.API_VERSION",
				Message: fmt.Sprintf("doc %d: apiVersion is required", i+1),
				Severity: engine.SevError,
			})
		}
		if kind == "" {
			res.Valid = false
			res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
				Code: "K8S.MISSING.KIND",
				Message: fmt.Sprintf("doc %d: kind is required", i+1),
				Severity: engine.SevError,
			})
		}
		md, _ := raw["metadata"].(map[string]any)
		name, _ := md["name"].(string)
		if name == "" {
			res.Valid = false
			res.Diagnostics = append(res.Diagnostics, engine.Diagnostic{
				Code: "K8S.MISSING.NAME",
				Message: fmt.Sprintf("doc %d: metadata.name is required", i+1),
				Severity: engine.SevError,
			})
		}
		if i == 0 {
			res.APIVersion = api
			res.Kind = kind
		}
	}
	return res, nil
}
