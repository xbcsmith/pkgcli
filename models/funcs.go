package models

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"crypto/md5"
	"crypto/sha256"

	"github.com/xbcsmith/pkgcli/utils"
	yaml "gopkg.in/yaml.v3"
)

// NewRelease func takes no input and returns a string
func NewRelease() string {
	t := time.Now()
	return t.Format("20060102.1573068157933")
}

// NewPkg func takes name and version as input and returns *Pkg
func NewPkg(name, version string) *Pkg {
	release := NewRelease()
	description := strings.Title(name) + " " + version + " " + release
	return &Pkg{
		Name:        name,
		Version:     version,
		Release:     release,
		Description: description,
		Summary:     description,
		Package:     "tar.xz",
		PlatformID:  "x86_64-gnu-linux-9",
	}
}

// GetNVRA returns the Name Version Release Arch of a package.
// foo-1.0.0-20191209.1573068157933.x86_64
func (p *Pkg) GetNVRA() string {
	release := "None"
	arch := "None"
	if len(p.Release) > 0 {
		release = p.Release
	}
	if len(strings.Split(p.PlatformID, "-")) > 0 {
		arch = strings.Split(p.PlatformID, "-")[0]
	}
	return p.Name + "-" + p.Version + "-" + release + "." + arch
}

// FetchSources fetches the sources from a pkg
func (p *Pkg) FetchSources(destdir string) ([]string, error) {
	var filelist []string
	err := os.MkdirAll(destdir, 0755)
	if err != nil {
		return nil, err
	}
	for _, src := range p.Sources {
		filename := path.Base(src.Archive)
		filepath := path.Join(destdir, filename)
		err := utils.DownloadFile(filepath, src.Archive)
		if err != nil {
			return nil, err
		}
		raw, err := ioutil.ReadFile(filepath)
		if err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
		md5sum := fmt.Sprintf("%x", md5.Sum(raw))
		if src.MD5 != "" {
			if md5sum != src.MD5 {
				return nil, fmt.Errorf("%s : MD5 sums do not match %s != %s", src.Archive, src.MD5, md5sum)
			}
		}
		sha256sum := fmt.Sprintf("%x", sha256.Sum256(raw))
		if src.SHA256 != "" {
			if sha256sum != src.SHA256 {
				return nil, fmt.Errorf("%s : SHA256 sums do not match %s != %s", src.Archive, src.SHA256, sha256sum)
			}
		}
		filelist = append(filelist, sha256sum+" "+filepath)
	}
	return filelist, nil
}

// ToYAML func takes no input and returns []byte, error
func (p *Pkg) ToYAML() ([]byte, error) {
	content, err := yaml.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to yaml : %v", err)
	}
	return content, nil
}

// ToJSON func takes no input and returns []byte, error
func (p *Pkg) ToJSON() ([]byte, error) {
	content, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to json : %v", err)
	}
	return content, nil
}

// ToPrettyJSON func takes no input and returns []byte, error
func (p *Pkg) ToPrettyJSON() ([]byte, error) {
	content, err := json.MarshalIndent(p, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to convert to json : %v", err)
	}
	return content, nil
}

// ToYAML func takes no input and returns []byte, error
func (p *Pkgs) ToYAML() ([]byte, error) {
	content, err := yaml.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to yaml : %v", err)
	}
	return content, nil
}

// ToJSON func takes no input and returns []byte, error
func (p *Pkgs) ToJSON() ([]byte, error) {
	content, err := json.Marshal(p)
	if err != nil {
		return nil, fmt.Errorf("failed to convert to json : %v", err)
	}
	return content, nil
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

// DecodePkgFromJSON func takes reader io.Reader as input and returns *Pkg, error
func DecodePkgFromJSON(reader io.Reader) (*Pkg, error) {
	pkg := &Pkg{}
	err := json.NewDecoder(reader).Decode(pkg)
	if err != nil {
		return nil, fmt.Errorf("error decoding json : %v", err)
	}
	return pkg, nil
}

// DecodePkgFromYAML func takes reader io.Reader as input and returns *Pkg, error
func DecodePkgFromYAML(reader io.Reader) (*Pkg, error) {
	pkg := &Pkg{}
	err := yaml.NewDecoder(reader).Decode(pkg)
	if err != nil {
		return nil, fmt.Errorf("error decoding yaml : %v", err)
	}
	return pkg, nil
}
