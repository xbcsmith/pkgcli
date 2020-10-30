// Copyright © 2020 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package source

import (
	"encoding/json"
	"fmt"

	yaml "gopkg.in/yaml.v3"
)

// Source struct for source
type Source struct {
	Archive     string `json:"archive" yaml:"archive"`
	Destination string `json:"destination,omitempty" yaml:"destination,omitempty"`
	SHA256      string `json:"sha256,omitempty" yaml:"sha256,omitempty"`
	MD5         string `json:"md5,omitempty" yaml:"md5,omitempty"`
	OnDisk      string `json:"ondisk,omitempty" yaml:"ondisk,omitempty"`
	Size        string `json:"size,omitempty" yaml:"size,omitempty"`
}

// Sources struct for sources
type Sources []struct {
	Sources []Source `json:"sources" yaml:"sources"`
}

// ToYAML func takes no input and returns []byte, error
func (s *Sources) ToYAML() ([]byte, error) {
	content, err := yaml.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to yaml : %v", err)
	}
	return content, nil
}

// ToJSON func takes no input and returns []byte, error
func (s *Sources) ToJSON() ([]byte, error) {
	content, err := json.Marshal(s)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to json : %v", err)
	}
	return content, nil
}
