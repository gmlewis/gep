// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package scenarios

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
)

// ScenarioID is the stable string identifier for one scenario.
type ScenarioID string

// ScenarioSplit is the name of an evaluation partition.
type ScenarioSplit string

const (
	// Train is the partition used during fitness optimization.
	Train ScenarioSplit = "train"
	// Validation is the held-out partition used to detect overfitting.
	Validation ScenarioSplit = "validation"
	// Test is the final held-out partition evaluated once, before promotion.
	Test ScenarioSplit = "test"
)

// Scenario is one individual scenario record within a [ScenarioSet].
type Scenario struct {
	// ID is the stable identifier for this scenario.
	ID ScenarioID `json:"id"`
	// Split is the evaluation partition this scenario belongs to.
	Split ScenarioSplit `json:"split"`
	// Tags are optional human-readable labels for grouping or filtering.
	Tags []string `json:"tags,omitempty"`
	// Params stores arbitrary scenario-specific JSON-serializable parameters.
	Params map[string]any `json:"params,omitempty"`
}

// ScenarioSet is a named, ordered collection of [Scenario] records from one
// source.
type ScenarioSet struct {
	// Name is the stable identifier for this set.
	Name string `json:"name"`
	// Source is an optional description of where this set originated
	// (e.g. a file path or dataset version string).
	Source string `json:"source,omitempty"`
	// Scenarios is the ordered list of scenario records.
	Scenarios []Scenario `json:"scenarios"`
}

// ScenarioRegistry aggregates one or more [ScenarioSet] values and validates
// that each [ScenarioID] belongs to at most one split.
type ScenarioRegistry struct {
	// Sets is the ordered list of scenario sets in this registry.
	Sets []ScenarioSet `json:"sets"`
}

// Validate returns an error if any [ScenarioID] appears more than once across
// all sets in the registry — whether in different splits or the same split.
// The first conflict found produces the error; the check order follows set
// declaration order and within-set scenario order.
func (r *ScenarioRegistry) Validate() error {
	seen := make(map[ScenarioID]ScenarioSplit)
	for _, set := range r.Sets {
		for _, sc := range set.Scenarios {
			prior, ok := seen[sc.ID]
			if !ok {
				seen[sc.ID] = sc.Split
				continue
			}
			if prior == sc.Split {
				return fmt.Errorf("scenario %q is listed more than once in split %q", sc.ID, sc.Split)
			}
			return fmt.Errorf("scenario %q belongs to both %q and %q splits", sc.ID, prior, sc.Split)
		}
	}
	return nil
}

// ByID returns the first [Scenario] with the given ID, searching sets in
// declaration order. The second return value is false when no scenario with
// that ID exists.
func (r *ScenarioRegistry) ByID(id ScenarioID) (Scenario, bool) {
	for _, set := range r.Sets {
		for _, sc := range set.Scenarios {
			if sc.ID == id {
				return sc, true
			}
		}
	}
	return Scenario{}, false
}

// BySplit returns all scenarios across all sets that belong to the given split,
// in set declaration order then insertion order within each set. The returned
// slice is a new allocation; mutating it does not affect the registry.
func (r *ScenarioRegistry) BySplit(split ScenarioSplit) []Scenario {
	var result []Scenario
	for _, set := range r.Sets {
		for _, sc := range set.Scenarios {
			if sc.Split == split {
				result = append(result, sc)
			}
		}
	}
	return result
}

var errNilSetReader = errors.New("nil scenario set reader")

// LoadScenarioSet decodes one [ScenarioSet] from JSON. It returns an error if
// r is nil, if the JSON is malformed, or if the stream contains more than one
// JSON object.
func LoadScenarioSet(r io.Reader) (*ScenarioSet, error) {
	if r == nil {
		return nil, errNilSetReader
	}
	dec := json.NewDecoder(r)
	dec.DisallowUnknownFields()
	var set ScenarioSet
	if err := dec.Decode(&set); err != nil {
		return nil, fmt.Errorf("decode scenario set JSON: %w", err)
	}
	if err := dec.Decode(&struct{}{}); !errors.Is(err, io.EOF) {
		return nil, errors.New("scenario set JSON must contain exactly one object")
	}
	return &set, nil
}

// LoadScenarioSetFile decodes one [ScenarioSet] from the named JSON file.
func LoadScenarioSetFile(filename string) (*ScenarioSet, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open scenario set file %q: %w", filename, err)
	}
	defer f.Close()
	return LoadScenarioSet(f)
}
