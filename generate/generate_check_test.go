// Copyright 2022 Mineiros GmbH
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package generate_test

import (
	"testing"

	"github.com/madlambda/spells/assert"
	"github.com/mineiros-io/terramate/generate"
	tmstack "github.com/mineiros-io/terramate/stack"
	"github.com/mineiros-io/terramate/test/sandbox"
)

func TestCheckReturnsOutdatedStackFilenamesForGeneratedHCL(t *testing.T) {
	s := sandbox.New(t)

	stackEntry := s.CreateStack("stacks/stack")
	stack := stackEntry.Load()

	assertOutdatedFiles := func(want []string) {
		t.Helper()

		got, err := generate.CheckStack(s.RootDir(), stack)
		assert.NoError(t, err)
		assertEqualStringList(t, got, want)
	}

	// Checking detection when there is no config generated yet
	assertOutdatedFiles([]string{})
	stackEntry.CreateConfig(
		stackConfig(
			generateHCL(
				labels("test.tf"),
				terraform(
					str("required_version", "1.10"),
				),
			),
		).String())
	assertOutdatedFiles([]string{"test.tf"})

	s.Generate()

	assertOutdatedFiles([]string{})

	// Now checking when we have code + it gets outdated.
	stackEntry.CreateConfig(
		stackConfig(
			generateHCL(
				labels("test.tf"),
				terraform(
					str("required_version", "1.11"),
				),
			),
		).String())

	assertOutdatedFiles([]string{"test.tf"})

	s.Generate()

	// Changing generated filenames will trigger detection,
	// with new + old filenames.
	stackEntry.CreateConfig(
		stackConfig(
			generateHCL(
				labels("testnew.tf"),
				terraform(
					str("required_version", "1.11"),
				),
			),
		).String())

	assertOutdatedFiles([]string{"test.tf", "testnew.tf"})

	// Adding new filename to generation trigger detection
	stackEntry.CreateConfig(
		stackConfig(
			generateHCL(
				labels("testnew.tf"),
				terraform(
					str("required_version", "1.11"),
				),
			),
			generateHCL(
				labels("another.tf"),
				backend(
					labels("type"),
				),
			),
		).String())

	assertOutdatedFiles([]string{"another.tf", "test.tf", "testnew.tf"})

	s.Generate()

	assertOutdatedFiles([]string{})

	// Detects configurations that have been removed.
	stackEntry.CreateConfig(stackConfig().String())

	assertOutdatedFiles([]string{"another.tf", "testnew.tf"})

	s.Generate()

	assertOutdatedFiles([]string{})
}

func TestCheckOutdatedIgnoresEmptyGenerateHCLBlocks(t *testing.T) {
	s := sandbox.New(t)

	stackEntry := s.CreateStack("stacks/stack")
	stack := stackEntry.Load()

	assertOutdatedFiles := func(want []string) {
		t.Helper()

		got, err := generate.CheckStack(s.RootDir(), stack)
		assert.NoError(t, err)
		assertEqualStringList(t, got, want)
	}

	// Checking detection when the config is empty at first
	stackEntry.CreateConfig(
		stackConfig(
			generateHCL(
				labels("test.tf"),
			),
		).String())

	assertOutdatedFiles([]string{})

	// Checking detection when the config isnt empty, is generated and then becomes empty
	stackEntry.CreateConfig(
		stackConfig(
			generateHCL(
				labels("test.tf"),
				block("whatever"),
			),
		).String())

	assertOutdatedFiles([]string{"test.tf"})

	s.Generate()

	assertOutdatedFiles([]string{})

	stackEntry.CreateConfig(
		stackConfig(
			generateHCL(
				labels("test.tf"),
			),
		).String())

	assertOutdatedFiles([]string{"test.tf"})

	s.Generate()

	assertOutdatedFiles([]string{})
}

func TestCheckReturnsOutdatedStackFilenamesForBackendAndLocals(t *testing.T) {
	s := sandbox.New(t)

	stack1 := s.CreateStack("stacks/stack-1")
	stack2 := s.CreateStack("stacks/stack-2")

	stack1val, err := tmstack.Load(s.RootDir(), stack1.Path())
	assert.NoError(t, err)
	stack2val, err := tmstack.Load(s.RootDir(), stack2.Path())
	assert.NoError(t, err)

	assertAllStacksAreUpdated := func() {
		t.Helper()

		for _, stack := range []tmstack.S{stack1val, stack2val} {
			got, err := generate.CheckStack(s.RootDir(), stack)
			assert.NoError(t, err)
			assertEqualStringList(t, got, []string{})
		}
	}

	assertAllStacksAreUpdated()

	// Checking detection when there is no config generated yet
	// for both locals and backend config
	stack1.CreateConfig(
		hcldoc(
			stack(),
			exportAsLocals(
				expr("test", "terramate.path"),
			),
		).String())

	got, err := generate.CheckStack(s.RootDir(), stack1val)
	assert.NoError(t, err)
	assertEqualStringList(t, got, []string{generate.LocalsFilename})

	stack2.CreateConfig(
		hcldoc(
			terramate(
				backend(labels("test")),
			),
			stack(),
		).String())

	got, err = generate.CheckStack(s.RootDir(), stack2val)
	assert.NoError(t, err)
	assertEqualStringList(t, got, []string{generate.BackendCfgFilename})

	s.Generate()

	assertAllStacksAreUpdated()

	// Now checking when we have code + it gets outdated for both stacks.
	stack1.CreateConfig(
		hcldoc(
			stack(),
			exportAsLocals(
				expr("changed", "terramate.name"),
			),
		).String())

	got, err = generate.CheckStack(s.RootDir(), stack1val)
	assert.NoError(t, err)
	assertEqualStringList(t, got, []string{generate.LocalsFilename})

	stack2.CreateConfig(
		hcldoc(
			terramate(
				backend(labels("changed")),
			),
			stack(),
		).String())

	got, err = generate.CheckStack(s.RootDir(), stack2val)
	assert.NoError(t, err)
	assertEqualStringList(t, got, []string{generate.BackendCfgFilename})

	s.Generate()

	assertAllStacksAreUpdated()

	stack2.CreateConfig(
		hcldoc(
			terramate(
				backend(labels("anotherChange")),
			),
			exportAsLocals(
				expr("changed", "terramate.name"),
			),
			stack(),
		).String())

	got, err = generate.CheckStack(s.RootDir(), stack2val)
	assert.NoError(t, err)
	assertEqualStringList(t, got, []string{
		generate.BackendCfgFilename,
		generate.LocalsFilename,
	})

	s.Generate()

	assertAllStacksAreUpdated()

	// Changing generated filenames will trigger detection, with new filenames
	// Here we detect both new files missing and also current files that should
	// not exist, since no configuration generates them anymore

	const (
		backendFilename = "backend.tf"
		localsFilename  = "locals.tf"
	)

	codegenConfig := hcldoc(
		terramate(
			block("config",
				block("generate",
					str("backend_config_filename", backendFilename),
					str("locals_filename", localsFilename),
				),
			),
		),
	)

	rootEntry := s.DirEntry(".")
	rootEntry.CreateConfig(codegenConfig.String())

	got, err = generate.CheckStack(s.RootDir(), stack1val)
	assert.NoError(t, err)
	assertEqualStringList(t, got, []string{generate.LocalsFilename, localsFilename})

	got, err = generate.CheckStack(s.RootDir(), stack2val)
	assert.NoError(t, err)
	assertEqualStringList(t, got, []string{
		generate.BackendCfgFilename,
		generate.LocalsFilename,
		backendFilename,
		localsFilename,
	})

	s.Generate()
	assertAllStacksAreUpdated()
}

func TestCheckFailsWithInvalidConfig(t *testing.T) {
	invalidConfigs := []string{
		hcldoc(
			terramate(
				backend(
					labels("test"),
					expr("undefined", "terramate.undefined"),
				),
			),
			stack(),
		).String(),
		hcldoc(
			exportAsLocals(
				expr("undefined", "terramate.undefined"),
			),
			stack(),
		).String(),
		hcldoc(
			generateHCL(
				expr("undefined", "terramate.undefined"),
			),
			stack(),
		).String(),
	}

	for _, invalidConfig := range invalidConfigs {
		s := sandbox.New(t)

		stackEntry := s.CreateStack("stack")
		stackEntry.CreateConfig(invalidConfig)

		stack, err := tmstack.Load(s.RootDir(), stackEntry.Path())
		assert.NoError(t, err)

		_, err = generate.CheckStack(s.RootDir(), stack)
		assert.Error(t, err, "should fail for configuration:\n%s", invalidConfig)
	}
}