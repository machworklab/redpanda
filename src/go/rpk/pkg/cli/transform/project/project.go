// Copyright 2023 Redpanda Data, Inc.
//
// Use of this software is governed by the Business Source License
// included in the file licenses/BSL.md
//
// As of the Change Date specified in that file, in accordance with
// the Business Source License, use of this software will be governed
// by the Apache License, Version 2.0

package project

import (
	"fmt"
	"strings"

	"github.com/spf13/afero"
	"gopkg.in/yaml.v3"
)

type WasmLang string

const (
	WasmLangTinygo WasmLang = "tinygo"
)

var AllWasmLangs = []string{"tinygo"}

func (l *WasmLang) Set(str string) error {
	lower := strings.ToLower(str)
	for _, s := range AllWasmLangs {
		if lower == s {
			*l = WasmLang(s)
			return nil
		}
	}
	return fmt.Errorf("unknown language: %q", str)
}

func (l WasmLang) String() string {
	return string(l)
}

func (WasmLang) Type() string {
	return "string"
}

type Config struct {
	Name        string            `yaml:"name"`
	Description string            `yaml:"description,omitempty"`
	InputTopic  string            `yaml:"input-topic"`
	OutputTopic string            `yaml:"output-topic"`
	Language    WasmLang          `yaml:"language"`
	Env         map[string]string `yaml:"env,omitempty"`
}

var ConfigFileName = "transform.yaml"

func MarshalConfig(c Config) ([]byte, error) {
	return yaml.Marshal(c)
}

func UnmarshalConfig(b []byte, c *Config) error {
	return yaml.Unmarshal(b, c)
}

func LoadCfg(fs afero.Fs) (c Config, err error) {
	b, err := afero.ReadFile(fs, ConfigFileName)
	if err != nil {
		return c, err
	}
	err = UnmarshalConfig(b, &c)
	return c, err
}
