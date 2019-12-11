package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"unicode"

	yaml "gopkg.in/yaml.v3"
)

// GetEnv returns an env variable value or a default
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// IsJSON try to guess if file is json or yaml
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
	if !isjson {
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

// DownloadFile downloads a binary from an http location
func DownloadFile(filepath string, url string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	_, err = io.Copy(f, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

func FindFile(dirpath string, filename string) ([]string, error) {
	filelist := []string{}
	err := filepath.Walk(dirpath, func(path string, f os.FileInfo, err error) error {
		if f.Name() == filename {
			filelist = append(filelist, path)
		}
		return nil
	})
	return filelist, err
}
