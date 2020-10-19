// Copyright © 2019 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package models

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

// Instruction struct for instruction
type Instruction struct {
	Build     string `json:"build" yaml:"build"`
	Configure string `json:"configure,omitempty" yaml:"configure,omitempty"`
	Install   string `json:"install" yaml:"install"`
	Post      string `json:"post,omitempty" yaml:"post,omitempty"`
	Pre       string `json:"pre,omitempty" yaml:"pre,omitempty"`
	Test      string `json:"test,omitempty" yaml:"test,omitempty"`
	Unpack    string `json:"unpack" yaml:"unpack"`
}

// Instructions struct for instructions
type Instructions struct {
	Instructions []Instruction `json:"instructions" yaml:"instructions"`
}

// ToYAML func takes no input and returns []byte, error
func (i *Instructions) ToYAML() ([]byte, error) {
	content, err := yaml.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to yaml : %v", err)
	}
	return content, nil
}

// ToJSON func takes no input and returns []byte, error
func (i *Instructions) ToJSON() ([]byte, error) {
	content, err := json.Marshal(i)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to json : %v", err)
	}
	return content, nil
}
