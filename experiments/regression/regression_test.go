// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package regression_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

// pilots lists every applied-design pilot together with the test function that
// exercises its full evolve→decode→validate→promote→export→checkpoint flow.
var pilots = []struct {
	domain string
	pkg    string
	test   string
}{
	{
		domain: "circuit/half_adder",
		pkg:    "./experiments/circuit/half_adder",
		test:   "TestRunPilotPromotionCheckpointAndManifest",
	},
	{
		domain: "voxel/bracket",
		pkg:    "./experiments/voxel/bracket",
		test:   "TestRunPilotPromotionCheckpointAndManifest",
	},
	{
		domain: "control/mass_spring_damper",
		pkg:    "./experiments/control/mass_spring_damper",
		test:   "TestRunPilotPromotionCheckpointAndManifest",
	},
}

// TestCrossDomainPipelines is the cross-domain regression gate.
// It runs the full promote-export-checkpoint pipeline for each applied-design
// pilot in parallel and fails loudly if any pilot regresses.
//
// Each sub-test name corresponds to a domain pilot. Run a single domain with:
//
//	go test ./experiments/regression/... -run TestCrossDomainPipelines/circuit
func TestCrossDomainPipelines(t *testing.T) {
	moduleRoot := findModuleRoot(t)
	for _, p := range pilots {
		p := p
		t.Run(p.domain, func(t *testing.T) {
			t.Parallel()
			cmd := exec.Command("go", "test",
				"-count=1",
				"-timeout=300s",
				"-run="+p.test,
				p.pkg,
			)
			cmd.Dir = moduleRoot
			out, err := cmd.CombinedOutput()
			if err != nil {
				t.Fatalf("cross-domain regression for %q failed:\n%s", p.domain, out)
			}
		})
	}
}

// findModuleRoot walks up from the package directory until it finds a go.mod
// file, which marks the module root where `go test ./pkg/...` commands must be
// invoked.
func findModuleRoot(t *testing.T) string {
	t.Helper()
	dir, err := os.Getwd()
	if err != nil {
		t.Fatalf("os.Getwd: %v", err)
	}
	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			t.Fatal("could not find module root from working directory")
		}
		dir = parent
	}
}
