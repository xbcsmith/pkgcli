// Copyright © 2019 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package utils

import (
	"archive/tar"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"unicode"

	"github.com/ulikunitz/xz"
	yaml "gopkg.in/yaml.v3"
)

// GetEnv returns an env variable value or a default
// GetEnv func takes no input and returns key, fallback string string
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// IsJSON try to guess if file is json or yaml
// IsJSON func takes no input and returns buf []byte bool
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

// json2yaml func takes raw []byte input and returns []byte, error
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

// yaml2json func takes raw []byte, noindent bool input and returns []byte, error
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
// Convert func takes raw []byte, noindent bool input and returns []byte, error
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
// DownloadFile func takes no input and returns filepath string, url string error
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

// FindFile finds a file in a given directory
func FindFile(dirpath, filename string) ([]string, error) {
	var files []string
	abspath, err := filepath.Abs(dirpath)
	if err != nil {
		return files, err
	}
	err = filepath.Walk(abspath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsDir() {
				return nil
			}
			if filepath.Base(path) == filename {
				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return files, err
}

// FindArtifacts finds artifacts in a given directory
func FindArtifacts(dirpath, suffix string) ([]string, error) {
	var files []string
	abspath, err := filepath.Abs(dirpath)
	if err != nil {
		return files, err
	}
	err = filepath.Walk(abspath,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsDir() {
				return nil
			}
			if strings.HasSuffix(path, suffix) {
				files = append(files, path)
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	return files, err
}

// GetEnvsByPrefix finds all ENV vars that start with prefix
// GetEnvsByPrefix func takes no input and returns prefix string, strip bool map[string]string
func GetEnvsByPrefix(prefix string, strip bool) map[string]string {
	envs := make(map[string]string)
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		if strings.HasPrefix(pair[0], prefix) {
			if len(pair[1]) > 0 {
				k := pair[0]
				if strip {
					k = strings.Split(pair[0], prefix+"_")[1]
				}
				envs[k] = pair[1]
			}
		}
	}
	return envs
}

// StringInSlice is used to find a string in a list
// StringInSlice func takes no input and returns a string, list []string bool
func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}

// DeepEqualStringArray compares []string to []string and returns bool
// DeepEqualStringArray func takes no input and returns first []string, second []string bool
func DeepEqualStringArray(first []string, second []string) bool {
	if len(first) != len(second) {
		return false
	}
	for _, v := range first {
		if !StringInSlice(v, second) {
			return false
		}

	}
	return true
}

// GetTokens func takes prefix string as input and returns *Tokens
func GetTokens(prefix string) *Tokens {
	tokens := make(map[string]string)
	envs := GetEnvsByPrefix(prefix, true)
	if len(envs) > 0 {
		for k, v := range envs {
			key := "@" + k + "@"
			tokens[key] = v
		}
	}
	return &Tokens{
		Tokens: tokens,
	}
}

// Tokens struct for tokens
type Tokens struct {
	Tokens map[string]string `json:"tokens"`
}

// Keys func takes no input and returns []string
func (t *Tokens) Keys() []string {
	var keys []string
	for k := range t.Tokens {
		keys = append(keys, k)
	}
	return keys
}

// Get func takes token string input and returns string
func (t *Tokens) Get(token string) string {
	for k, v := range t.Tokens {
		if k == token {
			return v
		}
	}
	return ""
}

// Replacer replaces values for tokens
// Replacer func takes *Tokens and text as input and returns a string
func Replacer(tokens *Tokens, text string) string {
	for _, token := range tokens.Keys() {
		text = strings.Replace(text, token, tokens.Get(token), -1)
	}
	return text
}

// Compress func takes no input and returns src string, excludes []string, writers ...io.Writer error
func Compress(src string, excludes []string, writers ...io.Writer) error {

	fmt.Printf("Creating tar xz : %s\n", src)

	if _, err := os.Stat(src); err != nil {
		return fmt.Errorf("Unable to tar files - %v", err.Error())
	}

	mw := io.MultiWriter(writers...)

	xzw, err := xz.NewWriter(mw)
	if err != nil {
		log.Fatalf("xz.NewWriter error %s", err)
	}
	defer xzw.Close()

	tw := tar.NewWriter(xzw)
	defer tw.Close()

	return filepath.Walk(src, func(file string, fi os.FileInfo, err error) error {

		if err != nil {
			return err
		}

		if StringInSlice(fi.Name(), excludes) {
			fmt.Printf("skipping file : %s\n", fi.Name())
			return nil
		}

		if !fi.Mode().IsRegular() {
			return nil
		}

		header, err := tar.FileInfoHeader(fi, fi.Name())
		if err != nil {
			return err
		}

		header.Name = strings.TrimPrefix(strings.Replace(file, src, "", -1), string(filepath.Separator))

		if err := tw.WriteHeader(header); err != nil {
			return err
		}

		f, err := os.Open(file)
		if err != nil {
			return err
		}

		if _, err := io.Copy(tw, f); err != nil {
			return err
		}

		f.Close()

		return nil
	})
}
