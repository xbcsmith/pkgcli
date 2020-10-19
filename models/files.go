// Copyright © 2019 Brett Smith <xbcsmith@gmail.com>, . All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package models

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

// File struct for file
type File struct {
	Path   string `json:"path,omitempty" yaml:"path,omitempty"`
	Name   string `json:"name,omitempty" yaml:"name,omitempty"`
	Mode   string `json:"mode,omitempty" yaml:"mode,omitempty"`
	SHA256 string `json:"sha256,omitempty" yaml:"sha256,omitempty"`
	MD5    string `json:"md5sum,omitempty" yaml:"md5sum,omitempty"`
}

// Files struct for files
type Files struct {
	Files []File `json:"files" yaml:"files"`
}

// for sorting by length
type byLength []string

func (s byLength) Len() int {
	return len(s)
}
func (s byLength) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byLength) Less(i, j int) bool {
	return len(s[i]) < len(s[j])
}

// splitPath returns []string path
func splitPath(path string) []string {
	path = strings.TrimPrefix(path, "/")
	if path == "" {
		return nil
	}
	return strings.Split(path, "/")
}

// FindDir finds a given directory
func FindDir(base, dirname string, depth int) ([]string, error) {
	var dirs []string
	err := filepath.Walk(base,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if info.Mode().IsDir() {
				relpath, err := filepath.Rel(base, path)
				if err != nil {
					return err
				}
				if len(splitPath(relpath)) < depth {
					if filepath.Base(path) == dirname {
						dirs = append(dirs, path)
					}
				}
			}
			return nil
		})
	if err != nil {
		log.Println(err)
	}
	sort.Sort(byLength(dirs))
	return dirs, err
}

// FindFile finds a file in a given directory
func FindFile(base, filename string) ([]string, error) {
	var files []string
	abspath, err := filepath.Abs(base)
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

// DownloadFiles downloads a map[string]string of files
func DownloadFiles(filemap map[string]string) error {
	return nil
}

// DownloadFile downloads a binary from an http location
// DownloadFile func takes no input and returns filepath string, url string error
func DownloadFile(filepath string, url string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	resp, err := http.Get(url) // nolint:gosec
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
