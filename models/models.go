package models

type Pkg struct {
	Description  string        `json:"description"`
	Instructions []Instruction `json:"instructions"`
	Name         string        `json:"name"`
	Package      string        `json:"package"`
	PlatformID   string        `json:"platform_id"`
	Provides     []string      `json:"provides"`
	Release      string        `json:"release"`
	Requires     []string      `json:"requires"`
	Sources      []Source      `json:"sources"`
	Summary      string        `json:"summary"`
	Version      string        `json:"version"`
}

type Pkgs struct {
	Packages []Pkg `json:"packages" yaml:"packages"`
}

type Instruction struct {
	Build     string `json:"build"`
	Configure string `json:"configure"`
	Install   string `json:"install"`
	Post      string `json:"post"`
	Pre       string `json:"pre"`
	Test      string `json:"test"`
	Unpack    string `json:"unpack"`
}

type Instructions struct {
	Instructions []Instruction `json:"instructions" yaml:"instructions"`
}

type Source struct {
	Archive     string `json:"archive" yaml:"archive"`
	Destination string `json:"destination" yaml:"destination"`
	SHA256      string `json:"sha256" yaml:"sha256"`
	MD5         string `json:"md5" yaml:"md5"`
}

type Sources []struct {
	Sources []Source `json:"sources" yaml:"sources"`
}

type File struct {
	Path   string `json:"path" yaml:"path"`
	Name   string `json:"name" yaml:"name"`
	Mode   string `json:"mode" yaml:"mode"`
	SHA256 string `json:"sha256" yaml:"sha256"`
}

type Files struct {
	Files []File `json:"files" yaml:"files"`
}
