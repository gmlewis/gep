// Copyright 2026 Google Inc. All rights reserved.
// Use of this source code is governed by the Apache 2.0
// license that can be found in the LICENSE file.

package scenarios

import (
	"bytes"
	_ "embed"

	designscenarios "github.com/gmlewis/gep/v2/design/scenarios"
)

//go:embed testdata/set_smoke.json
var fixtureSetJSON []byte

// LoadFixtureSet returns the embedded reusable voxel scenario set by decoding
// it through [designscenarios.LoadScenarioSet].
func LoadFixtureSet() (*designscenarios.ScenarioSet, error) {
	return designscenarios.LoadScenarioSet(bytes.NewReader(fixtureSetJSON))
}
