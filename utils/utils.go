package utils

import (
  	"encoding/json"
	"unicode"
  	yaml "gopkg.in/yaml.v3"
)

/ IsJSON try to guess if file is json or yaml
func IsJSON(buf []byte) bool {
	trim := bytes.TrimLeftFunc(buf, unicode.IsSpace)
	if bytes.HasPrefix(trim, []byte("{")) {
		return true
	}
	if bytes.HasPrefix(trim, []byte("[")) {
		return true
	}
	return false
}

func json2yaml(raw []byte) ([]byte, error) {
	var output interface{}
	if err := json.Unmarshal([]byte(raw), &output); err != nil {
		return nil, err
	}

	content, err := yaml.Marshal(output)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to yaml : %v", err)
	}
	return content, nil
}

func yaml2json(raw []byte, noindent bool) ([]byte, error) {
	// ms := yaml.MapSlice{}
	var output interface{}
	if err := yaml.Unmarshal(raw, &output); err != nil {
		return nil, err
	}

	if noindent {
		content, err := json.Marshal(output)
		if err != nil {
			return nil, fmt.Errorf("failed to convert to yaml : %v", err)
		}
		return content, nil
	}

	content, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to convert to yaml : %v", err)
	}
	return content, nil

}

// Convert yaml to json or json to yaml
func Convert(raw []byte, noindent bool) ([]byte, error) {
	isjson := IsJSON(raw)
	if isjson != true {
		output, err := yaml2json(raw, noindent)
		if err != nil {
			fmt.Printf("decode data: %v", err)
			return nil, err
		}
		return output, nil
	}
	output, err := json2yaml(raw)
	if err != nil {
		fmt.Printf("decode data: %v", err)
		return nil, err
	}
	return output, nil
}
