// Copyright 2024 Terramate GmbH
// SPDX-License-Identifier: MPL-2.0

stack {
  name        = "package hcl // import \"github.com/terramate-io/terramate/hcl\""
  description = "package hcl // import \"github.com/terramate-io/terramate/hcl\"\n\nPackage hcl provides parsing functionality for Terramate HCL configuration.\nIt also provides printing and formatting for Terramate configuration.\n\nconst ErrHCLSyntax errors.Kind = \"HCL syntax error\" ...\nconst ErrScriptNoLabels errors.Kind = \"terramate schema error: (script): must provide at least one label\" ...\nconst SharingIsCaringExperimentName = \"outputs-sharing\"\nconst StackBlockType = \"stack\"\nfunc IsRootConfig(rootdir string) (bool, error)\nfunc MatchAnyGlob(globs []glob.Glob, s string) bool\nfunc PrintConfig(w io.Writer, cfg Config) error\nfunc PrintImports(w io.Writer, imports []string) error\nfunc ValueAsStringList(val cty.Value) ([]string, error)\ntype AssertConfig struct{ ... }\ntype ChangeDetectionConfig struct{ ... }\ntype CloudConfig struct{ ... }\ntype Command ast.Attribute\n    func NewScriptCommand(attr ast.Attribute) *Command\ntype Commands ast.Attribute\n    func NewScriptCommands(attr ast.Attribute) *Commands\ntype Config struct{ ... }\n    func NewConfig(dir string) (Config, error)\n    func ParseDir(root string, dir string, experiments ...string) (Config, error)\ntype Evaluator interface{ ... }\ntype GenFileBlock struct{ ... }\ntype GenHCLBlock struct{ ... }\ntype GenerateConfig struct{ ... }\ntype GenerateRootConfig struct{ ... }\ntype GitConfig struct{ ... }\n    func NewGitConfig() *GitConfig\ntype Input struct{ ... }\ntype Inputs []Input\ntype ManifestConfig struct{ ... }\ntype ManifestDesc struct{ ... }\ntype OptionalCheck int\n    const CheckIsUnset OptionalCheck = iota ...\n    func ToOptionalCheck(v bool) OptionalCheck\ntype Output struct{ ... }\ntype Outputs []Output\ntype RawConfig struct{ ... }\n    func NewCustomRawConfig(handlers map[string]mergeHandler) RawConfig\n    func NewTopLevelRawConfig() RawConfig\ntype RootConfig struct{ ... }\ntype RunConfig struct{ ... }\n    func NewRunConfig() *RunConfig\ntype RunEnv struct{ ... }\ntype Script struct{ ... }\ntype ScriptJob struct{ ... }\ntype SharingBackend struct{ ... }\ntype SharingBackendType int\n    const TerraformSharingBackend SharingBackendType = iota + 1\ntype SharingBackends []SharingBackend\ntype Stack struct{ ... }\ntype StackFilterConfig struct{ ... }\ntype TargetsConfig struct{ ... }\ntype TerragruntChangeDetectionEnabledOption int\n    const TerragruntAutoOption TerragruntChangeDetectionEnabledOption = iota ...\ntype TerragruntConfig struct{ ... }\ntype Terramate struct{ ... }\ntype TerramateParser struct{ ... }\n    func NewStrictTerramateParser(rootdir string, dir string, experiments ...string) (*TerramateParser, error)\n    func NewTerramateParser(rootdir string, dir string, experiments ...string) (*TerramateParser, error)\ntype VendorConfig struct{ ... }"
  tags        = ["golang", "hcl"]
  id          = "dd11014a-a1bb-4f1c-99da-8f6d188d36d1"
}
