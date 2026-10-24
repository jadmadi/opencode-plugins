// check-catalog validates plugins.json against the plugin repos and the
// runtime baseline. It is read-only and needs no network.
//
// Usage, from the repo root:
//
//	go run script/check-catalog/main.go
//
// Exit codes: 0 all checks pass, 1 findings, 2 usage or parse error.
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	catalogPath  = "plugins.json"
	readmePath   = "README.md"
	baselinePath = "runtime-baseline.json"
	repoBase     = "https://github.com/jadmadi/"
)

type command struct {
	Name string `json:"name"`
}

type entry struct {
	Name          string    `json:"name"`
	Version       string    `json:"version"`
	LatestTag     string    `json:"latestTag"`
	BuiltUnder    string    `json:"builtUnder"`
	BuiltUnderTag string    `json:"builtUnderTag"`
	Repo          string    `json:"repo"`
	Directory     string    `json:"directory"`
	Visibility    string    `json:"visibility"`
	Published     bool      `json:"published"`
	Commands      []command `json:"commands"`
}

type catalog struct {
	Plugins []entry `json:"plugins"`
}

type baseline struct {
	Runtime string `json:"runtime"`
}

type packageJSON struct {
	Version    string `json:"version"`
	BuiltUnder string `json:"builtUnder"`
}

type checker struct {
	findings []string
	readme   string
	baseline string
}

func (c *checker) fail(format string, args ...any) {
	c.findings = append(c.findings, fmt.Sprintf(format, args...))
}

func (c *checker) checkEntry(e entry) {
	if e.Name == "" || e.Version == "" || e.Repo == "" || e.Directory == "" || e.Visibility == "" {
		c.fail("%s: missing a required field (name, version, repo, directory, visibility)", e.Name)
	}
	if e.LatestTag != "v"+e.Version {
		c.fail("%s: latestTag %q does not match version %q", e.Name, e.LatestTag, e.Version)
	}
	if e.BuiltUnderTag != "builtUnder-"+e.BuiltUnder {
		c.fail("%s: builtUnderTag %q does not match builtUnder %q", e.Name, e.BuiltUnderTag, e.BuiltUnder)
	}
	if e.BuiltUnder != c.baseline {
		c.fail("%s: builtUnder %q does not match the baseline runtime %q", e.Name, e.BuiltUnder, c.baseline)
	}
	if e.Repo != repoBase+e.Directory {
		c.fail("%s: repo %q does not match the directory %q", e.Name, e.Repo, e.Directory)
	}
	if e.Published && !strings.Contains(c.readme, e.Repo) {
		c.fail("%s: published but README.md does not link %s", e.Name, e.Repo)
	}

	dir := ""
	for _, candidate := range []string{filepath.Join("..", e.Directory), e.Directory} {
		if info, err := os.Stat(candidate); err == nil && info.IsDir() {
			dir = candidate
			break
		}
	}
	if dir == "" {
		c.fail("%s: directory %s not found (expected the sibling repo)", e.Name, e.Directory)
		return
	}
	raw, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		// A skills pack has no package.json; the version lives in the catalog.
		return
	}
	var pkg packageJSON
	if err := json.Unmarshal(raw, &pkg); err != nil {
		c.fail("%s: package.json does not parse: %v", e.Name, err)
		return
	}
	if pkg.Version != e.Version {
		c.fail("%s: package.json version %q does not match catalog version %q", e.Name, pkg.Version, e.Version)
	}
	if pkg.BuiltUnder != e.BuiltUnder {
		c.fail("%s: package.json builtUnder %q does not match catalog builtUnder %q", e.Name, pkg.BuiltUnder, e.BuiltUnder)
	}
}

func main() {
	dir := "."
	if len(os.Args) > 1 {
		dir = os.Args[1]
	}
	if err := os.Chdir(dir); err != nil {
		fmt.Printf("check-catalog: ERROR\n  %v\n", err)
		os.Exit(2)
	}

	c := &checker{}
	raw, err := os.ReadFile(catalogPath)
	if err != nil {
		fmt.Printf("check-catalog: ERROR\n  read %s: %v (run from the repo root)\n", catalogPath, err)
		os.Exit(2)
	}
	var cat catalog
	if err := json.Unmarshal(raw, &cat); err != nil {
		fmt.Printf("check-catalog: ERROR\n  parse %s: %v\n", catalogPath, err)
		os.Exit(2)
	}
	if len(cat.Plugins) == 0 {
		fmt.Printf("check-catalog: ERROR\n  %s has no plugins\n", catalogPath)
		os.Exit(2)
	}
	if err := c.loadBaseline(); err != nil {
		fmt.Printf("check-catalog: ERROR\n  %v\n", err)
		os.Exit(2)
	}
	readme, err := os.ReadFile(readmePath)
	if err != nil {
		fmt.Printf("check-catalog: ERROR\n  read %s: %v\n", readmePath, err)
		os.Exit(2)
	}
	c.readme = string(readme)

	seen := map[string]bool{}
	for _, e := range cat.Plugins {
		if seen[e.Name] {
			c.fail("%s: duplicate entry", e.Name)
		}
		seen[e.Name] = true
		c.checkEntry(e)
	}

	if len(c.findings) == 0 {
		fmt.Printf("check-catalog: OK\n  %d entries match the repos, the README, and the baseline %s\n", len(cat.Plugins), c.baseline)
		return
	}
	fmt.Println("check-catalog: FINDINGS")
	for _, f := range c.findings {
		fmt.Printf("  - %s\n", f)
	}
	os.Exit(1)
}

func (c *checker) loadBaseline() error {
	raw, err := os.ReadFile(baselinePath)
	if err != nil {
		return fmt.Errorf("read %s: %w", baselinePath, err)
	}
	var base baseline
	if err := json.Unmarshal(raw, &base); err != nil {
		return fmt.Errorf("parse %s: %w", baselinePath, err)
	}
	if base.Runtime == "" {
		return fmt.Errorf("%s has no runtime", baselinePath)
	}
	c.baseline = base.Runtime
	return nil
}
