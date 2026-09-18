// check-runtime-api compares the installed OpenCode runtime and the published
// @opencode/plugin types with runtime-baseline.json.
//
// Usage, from the repo root:
//
//	go run script/check-runtime-api.go
//
// Exit codes: 0 baseline holds, 1 migrate, 2 usage error, 3 could not verify.
package main

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

const (
	baselinePath = "runtime-baseline.json"
	registryURL  = "https://registry.npmjs.org/@opencode%2Fplugin/latest"
	typesPath    = "package/dist/promise/plugin.d.ts"
)

type baseline struct {
	Runtime string `json:"runtime"`
	API     string `json:"api"`
	Checked string `json:"checked"`
	Note    string `json:"note"`
}

type registry struct {
	Version string `json:"version"`
	Dist    struct {
		Tarball string `json:"tarball"`
	} `json:"dist"`
}

var (
	semver       = regexp.MustCompile(`v?(\d+\.\d+\.\d+)`)
	catalogProp  = regexp.MustCompile(`(?m)^\s*readonly catalog\s*:`)
	modelProp    = regexp.MustCompile(`(?m)^\s*readonly model\s*:`)
	providerProp = regexp.MustCompile(`(?m)^\s*readonly provider\s*:`)
)

func readBaseline() (*baseline, error) {
	raw, err := os.ReadFile(baselinePath)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w (run from the repo root)", baselinePath, err)
	}
	var base baseline
	if err := json.Unmarshal(raw, &base); err != nil {
		return nil, fmt.Errorf("parse %s: %w", baselinePath, err)
	}
	return &base, nil
}

func installedRuntime() (string, error) {
	out, err := exec.Command("opencode", "--version").Output()
	if err != nil {
		return "", err
	}
	match := semver.FindStringSubmatch(string(out))
	if match == nil {
		return "", fmt.Errorf("unrecognized output %q", strings.TrimSpace(string(out)))
	}
	return match[1], nil
}

func fetchRegistry(client *http.Client) (*registry, error) {
	resp, err := client.Get(registryURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("registry returned %s", resp.Status)
	}
	var info registry
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return nil, err
	}
	if info.Version == "" || info.Dist.Tarball == "" {
		return nil, fmt.Errorf("registry response is missing the version or tarball")
	}
	return &info, nil
}

// detectAPI reads the plugin context type from the published package and
// reports which domain shape it declares.
func detectAPI(client *http.Client, tarball string) (string, error) {
	resp, err := client.Get(tarball)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("tarball returned %s", resp.Status)
	}
	gz, err := gzip.NewReader(resp.Body)
	if err != nil {
		return "", err
	}
	defer gz.Close()
	reader := tar.NewReader(gz)
	for {
		header, err := reader.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return "", err
		}
		if header.Name != typesPath {
			continue
		}
		content, err := io.ReadAll(reader)
		if err != nil {
			return "", err
		}
		text := string(content)
		switch {
		case catalogProp.MatchString(text):
			return "catalog", nil
		case modelProp.MatchString(text) && providerProp.MatchString(text):
			return "model-provider", nil
		default:
			return "unknown", nil
		}
	}
	return "", fmt.Errorf("%s not found in the tarball", typesPath)
}

func main() {
	base, err := readBaseline()
	if err != nil {
		fmt.Printf("runtime-baseline: ERROR\n  %v\n", err)
		os.Exit(2)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	installed, err := installedRuntime()
	if err != nil {
		fmt.Printf("runtime-baseline: UNKNOWN\n  could not read the installed OpenCode version: %v\n", err)
		os.Exit(3)
	}
	info, err := fetchRegistry(client)
	if err != nil {
		fmt.Printf("runtime-baseline: UNKNOWN\n  could not reach the npm registry: %v\n", err)
		os.Exit(3)
	}
	api, err := detectAPI(client, info.Dist.Tarball)
	if err != nil {
		fmt.Printf("runtime-baseline: UNKNOWN\n  could not read the published plugin types: %v\n", err)
		os.Exit(3)
	}

	var reasons []string
	if installed != base.Runtime {
		reasons = append(reasons, fmt.Sprintf("installed OpenCode is %s, baseline is %s", installed, base.Runtime))
	}
	if info.Version != base.Runtime {
		reasons = append(reasons, fmt.Sprintf("published package is %s, baseline is %s", info.Version, base.Runtime))
	}
	if api != base.API {
		reasons = append(reasons, fmt.Sprintf("published types declare the %s API, baseline is %s", api, base.API))
	}

	if len(reasons) > 0 {
		fmt.Println("runtime-baseline: MIGRATE")
		fmt.Printf("  installed opencode: %s\n", installed)
		fmt.Printf("  npm @opencode/plugin: %s (%s API)\n", info.Version, api)
		fmt.Printf("  baseline: %s %s, checked %s\n", base.Runtime, base.API, base.Checked)
		for _, reason := range reasons {
			fmt.Printf("  reason: %s\n", reason)
		}
		fmt.Println("  Follow the migration steps in AGENTS.md, then update runtime-baseline.json.")
		os.Exit(1)
	}

	fmt.Println("runtime-baseline: HOLD")
	fmt.Printf("  installed opencode: %s\n", installed)
	fmt.Printf("  npm @opencode/plugin: %s (%s API)\n", info.Version, api)
	fmt.Printf("  baseline: %s %s, checked %s\n", base.Runtime, base.API, base.Checked)
	fmt.Println("  Keep the plain { id, setup } object and the domains the baseline names.")
}
