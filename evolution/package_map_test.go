// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package evolution

import (
	"encoding/json"
	"errors"
	"io"
	"os/exec"
	"strings"
	"testing"
)

const modulePath = "github.com/gmlewis/gep/v2"

type packageMetadata struct {
	ImportPath string
	Imports    []string
}

func TestPhase4Milestone1_PackageBoundaries(t *testing.T) {
	metas := mustListPackages(t,
		modulePath+"/core",
		modulePath+"/evolution",
		modulePath+"/evolution/...",
	)

	foundCore := false
	foundEvolution := false
	for _, meta := range metas {
		allowedPrefixes := allowedImportPrefixes(meta.ImportPath)
		if len(allowedPrefixes) == 0 {
			continue
		}

		if meta.ImportPath == modulePath+"/core" {
			foundCore = true
		}
		if meta.ImportPath == modulePath+"/evolution" {
			foundEvolution = true
		}

		for _, imported := range meta.Imports {
			if !strings.HasPrefix(imported, modulePath+"/") {
				continue // standard library or external dependency
			}
			if hasAllowedPrefix(imported, allowedPrefixes) {
				continue
			}
			t.Errorf("%s imports %s; allowed internal prefixes are %v", meta.ImportPath, imported, allowedPrefixes)
		}
	}

	if !foundCore {
		t.Fatalf("expected to inspect %s/core", modulePath)
	}
	if !foundEvolution {
		t.Fatalf("expected to inspect %s/evolution", modulePath)
	}
}

func TestPhase4Milestone2_EnvironmentPackageBoundaries(t *testing.T) {
	metas := mustListPackages(t,
		modulePath+"/gymnasium",
		modulePath+"/gymnasium/...",
	)

	foundGymnasium := false
	for _, meta := range metas {
		allowedPrefixes := allowedImportPrefixes(meta.ImportPath)
		if len(allowedPrefixes) == 0 {
			continue
		}
		if meta.ImportPath == modulePath+"/gymnasium" {
			foundGymnasium = true
		}

		for _, imported := range meta.Imports {
			if !strings.HasPrefix(imported, modulePath+"/") {
				continue // standard library or external dependency
			}
			if hasAllowedPrefix(imported, allowedPrefixes) {
				continue
			}
			t.Errorf("%s imports %s; allowed internal prefixes are %v", meta.ImportPath, imported, allowedPrefixes)
		}
	}

	if !foundGymnasium {
		t.Fatalf("expected to inspect %s/gymnasium", modulePath)
	}
}

func TestPhase4Milestone3_LegacyModelPackageBoundaries(t *testing.T) {
	metas := mustListPackages(t,
		modulePath+"/model",
		modulePath+"/model/...",
	)

	foundModel := false
	for _, meta := range metas {
		allowedPrefixes := allowedImportPrefixes(meta.ImportPath)
		if len(allowedPrefixes) == 0 {
			continue
		}
		if meta.ImportPath == modulePath+"/model" {
			foundModel = true
		}

		for _, imported := range meta.Imports {
			if !strings.HasPrefix(imported, modulePath+"/") {
				continue // standard library or external dependency
			}
			if hasAllowedPrefix(imported, allowedPrefixes) {
				continue
			}
			t.Errorf("%s imports %s; allowed internal prefixes are %v", meta.ImportPath, imported, allowedPrefixes)
		}
	}

	if !foundModel {
		t.Fatalf("expected to inspect %s/model", modulePath)
	}
}

func allowedImportPrefixes(importPath string) []string {
	switch {
	case importPath == modulePath+"/core":
		return []string{modulePath + "/core"}
	case strings.HasPrefix(importPath, modulePath+"/evolution"):
		return []string{modulePath + "/core", modulePath + "/evolution"}
	case strings.HasPrefix(importPath, modulePath+"/gymnasium"):
		return []string{modulePath + "/common", modulePath + "/gymnasium"}
	case strings.HasPrefix(importPath, modulePath+"/model"):
		return []string{
			modulePath + "/common",
			modulePath + "/functions",
			modulePath + "/gene",
			modulePath + "/genome",
			modulePath + "/model",
		}
	default:
		return nil
	}
}

func hasAllowedPrefix(imported string, prefixes []string) bool {
	for _, prefix := range prefixes {
		if strings.HasPrefix(imported, prefix) {
			return true
		}
	}
	return false
}

func mustListPackages(t *testing.T, patterns ...string) []packageMetadata {
	t.Helper()

	args := append([]string{"list", "-json"}, patterns...)
	cmd := exec.Command("go", args...)
	out, err := cmd.Output()
	if err != nil {
		t.Fatalf("go list failed: %v", err)
	}

	dec := json.NewDecoder(strings.NewReader(string(out)))
	var result []packageMetadata
	for {
		var meta packageMetadata
		if err := dec.Decode(&meta); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			t.Fatalf("decode go list output: %v", err)
		}
		if meta.ImportPath == "" {
			t.Fatalf("decoded package with empty import path: %#v", meta)
		}
		result = append(result, meta)
	}

	if len(result) == 0 {
		t.Fatalf("go list returned no packages for %v", patterns)
	}
	return result
}
