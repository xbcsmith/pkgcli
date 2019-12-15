package utils

import (
	"encoding/json"
	"os"
	"testing"

	yaml "gopkg.in/yaml.v3"
	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

// Tests struct for tests
type Tests struct {
	a []byte
	b []byte
	c map[string]interface{}
	d interface{}
	e []byte
}

const ystr = `
bar:
  - buz
  - cuz
  - duz
baz:
    caz: fuz
flag: true
foo: show_value_of_foo
fuzzy:
    complicated-it-is:
        could_be-but:
            not-really_possible:
                until_it_is: true
yyy:
  - one
  - 2
  - true
  - "4"
  - key: value
  - - 1
    - "2"
    - things:
        - complicated: true
          couldbe: maybe
          notreally: false
zzz-zzz:
    buz:
      - 1
      - 2
      - 3
`

const jstr = `{
  "bar": [
    "buz",
    "cuz",
    "duz"
  ],
  "baz": {
    "caz": "fuz"
  },
  "flag": true,
  "foo": "show_value_of_foo",
  "fuzzy": {
    "complicated-it-is": {
      "could_be-but": {
        "not-really_possible": {
          "until_it_is": true
        }
      }
    }
  },
  "yyy": [
    "one",
    2,
    true,
    "4",
    {
      "key": "value"
    },
    [
      1,
      "2",
      {
        "things": [
          {
            "complicated": true,
            "couldbe": "maybe",
            "notreally": false
          }
        ]
      }
    ]
  ],
  "zzz-zzz": {
    "buz": [
      1,
      2,
      3
    ]
  }
}
`

// NewTests func takes no input and returns *Tests
func NewTests() *Tests {
	var j map[string]interface{}
	if err := json.Unmarshal([]byte(jstr), &j); err != nil {
		panic(err)
	}
	var y map[string]interface{}
	if err := yaml.Unmarshal([]byte(ystr), &y); err != nil {
		panic(err)
	}
	tests := &Tests{
		a: []byte(jstr),
		b: []byte(ystr),
		c: j,
		d: y,
		e: []byte(`[{"brackets" : "json"}]`),
	}
	return tests
}

// TestConvertJson func takes no input and returns t *testing.T
func TestConvertJson(t *testing.T) {
	tests := NewTests()
	expected := `- key:`
	actual, err := Convert(tests.a, false)
	assert.Assert(t, is.Nil(err))
	assert.Assert(t, !IsJSON(actual))
	assert.Assert(t, is.Contains(string(actual), expected))
}

// TestConvertAlsoJson func takes no input and returns t *testing.T
func TestConvertAlsoJson(t *testing.T) {
	tests := NewTests()
	expected := `brackets:`
	actual, err := Convert(tests.e, false)
	assert.Assert(t, is.Nil(err))
	assert.Assert(t, !IsJSON(actual))
	assert.Assert(t, is.Contains(string(actual), expected))
}

// TestConvertYaml func takes no input and returns t *testing.T
func TestConvertYaml(t *testing.T) {
	tests := NewTests()
	expected := `"bar": [`
	actual, err := Convert(tests.b, false)
	assert.Assert(t, is.Nil(err))
	assert.Assert(t, IsJSON(actual))
	assert.Assert(t, is.Contains(string(actual), expected))
}

// TestConvertYamlNoIndent func takes no input and returns t *testing.T
func TestConvertYamlNoIndent(t *testing.T) {
	tests := NewTests()
	expected := `bar":[`
	actual, err := Convert(tests.b, true)
	assert.Assert(t, is.Nil(err))
	assert.Assert(t, IsJSON(actual))
	assert.Assert(t, is.Contains(string(actual), expected))
}

// TestGetEnv func takes no input and returns t *testing.T
func TestGetEnv(t *testing.T) {
	os.Setenv("FOO", "1")
	foo := GetEnv("FOO", "2")
	assert.Assert(t, foo == "1")
	bar := GetEnv("BAR", "42")
	assert.Assert(t, bar == "42")
}

// TestStringInSlice validates StringInSlice works
// TestStringInSlice func takes no input and returns t *testing.T
func TestStringInSlice(t *testing.T) {
	list := []string{"foo", "bar", "caz"}
	result := StringInSlice("foo", list)
	assert.Equal(t, result, true, "Error in StringInSlice")
}

// TestDeepEqualStringArray func takes no input and returns t *testing.T
func TestDeepEqualStringArray(t *testing.T) {
	a := []string{"01DQDHG7GK3NQDREC3P47DJT48", "01DQWKGE76ABCYDZ07E6DZ2QC7"}
	b := []string{"01DQWKGE76ABCYDZ07E6DZ2QC7", "01DQDHG7GK3NQDREC3P47DJT48"}
	c := []string{"01DQWKGE76ABCYDZ07E6DZ2QC7", "foobar"}
	d := []string{"01DQWKGE76ABCYDZ07E6DZ2QC7"}
	assert.Assert(t, DeepEqualStringArray(a, b), "Error DeepEqualArray")
	assert.Assert(t, !DeepEqualStringArray(b, c), "Error DeepEqualArray")
	assert.Assert(t, !DeepEqualStringArray(b, d), "Error DeepEqualArray")
}

// TestGetEnvsByPrefix func takes no input and returns t *testing.T
func TestGetEnvsByPrefix(t *testing.T) {
	// 46 and 2
	os.Setenv("GET_ENV_PREFIX_FOO", "46")
	os.Setenv("GET_ENV_PREFIX_BAR", "2")
	prefix := "GET_ENV_PREFIX"
	tokens := GetEnvsByPrefix(prefix, true)
	assert.Assert(t, tokens["FOO"] == "46")
	assert.Assert(t, tokens["BAR"] == "2")
}

// TestGetTokens func takes no input and returns t *testing.T
func TestGetTokens(t *testing.T) {
	// 46 and 2
	os.Setenv("TOKEN_PREFIX_FOO", "46")
	os.Setenv("TOKEN_PREFIX_BAR", "2")
	prefix := "TOKEN_PREFIX"
	tokens := GetTokens(prefix)
	assert.Assert(t, tokens.Get("@FOO@") == "46")
	assert.Assert(t, tokens.Get("@BAR@") == "2")
	assert.Assert(t, DeepEqualStringArray(tokens.Keys(), []string{"@FOO@", "@BAR@"}))
}

// TestReplacer func takes no input and returns t *testing.T
func TestReplacer(t *testing.T) {
	data := `{
  "description": "The Sharutils package contains utilities that can create 'shell' archives",
  "instructions": [
    {
      "build": "make",
      "configure": "./configure --prefix=/usr",
      "install": "make install",
      "post": "",
      "pre": "sed -i 's/IO_ftrylockfile/IO_EOF_SEEN/' lib/*.c \u0026\u0026 echo '#define _IO_IN_BACKUP 0x100' \u003e\u003e lib/stdio-impl.h",
      "test": "make check",
      "unpack": "tar -xvf $name-$version.tar.xz \u0026\u0026 cd $name-$version"
    }
  ],
  "name": "@NAME@",
  "package": "@PACKAGE@",
  "platform_id": "@PLATFORM_ID@",
  "provides": [],
  "release": null,
  "requires": [],
  "sources": [
    {
      "archives": "http://ftp.gnu.org/gnu/sharutils/sharutils-4.15.2.tar.xz",
      "destination": "/usr",
      "md5": "5975ce21da36491d7aa6dc2b0d9788e0",
      "sha256": "2b05cff7de5d7b646dc1669bc36c35fdac02ac6ae4b6c19cb3340d87ec553a9a"
    }
  ],
  "summary": "sharutils utilities for shell archives",
  "version": "@VERSION@"
}`
	expected := `{
  "description": "The Sharutils package contains utilities that can create 'shell' archives",
  "instructions": [
    {
      "build": "make",
      "configure": "./configure --prefix=/usr",
      "install": "make install",
      "post": "",
      "pre": "sed -i 's/IO_ftrylockfile/IO_EOF_SEEN/' lib/*.c \u0026\u0026 echo '#define _IO_IN_BACKUP 0x100' \u003e\u003e lib/stdio-impl.h",
      "test": "make check",
      "unpack": "tar -xvf $name-$version.tar.xz \u0026\u0026 cd $name-$version"
    }
  ],
  "name": "sharutils",
  "package": "tar.xz",
  "platform_id": "x86_64-lfs-linux-9",
  "provides": [],
  "release": null,
  "requires": [],
  "sources": [
    {
      "archives": "http://ftp.gnu.org/gnu/sharutils/sharutils-4.15.2.tar.xz",
      "destination": "/usr",
      "md5": "5975ce21da36491d7aa6dc2b0d9788e0",
      "sha256": "2b05cff7de5d7b646dc1669bc36c35fdac02ac6ae4b6c19cb3340d87ec553a9a"
    }
  ],
  "summary": "sharutils utilities for shell archives",
  "version": "4.15.2"
}`
	tokens := &Tokens{
		Tokens: map[string]string{
			"@NAME@":        "sharutils",
			"@PLATFORM_ID@": "x86_64-lfs-linux-9",
			"@VERSION@":     "4.15.2",
			"@PACKAGE@":     "tar.xz",
		},
	}

	new := Replacer(tokens, data)
	var m map[string]interface{}
	err := json.Unmarshal([]byte(new), &m)
	assert.Assert(t, is.Nil(err))
	assert.Equal(t, m["name"], "sharutils", "Replace name failed")
	assert.Equal(t, new, expected, "Expected text failed")
}

func TestFindFile(t *testing.T) {
	filename := "sharutils.yml"
	filepath := "tests/pkg/"
	expected := []string{"tests/pkg/sharutils.yml"}
	paths, err := FindFile(filepath, filename)
	assert.Assert(t, is.Nil(err))
	assert.Equal(t, paths, expected, "Expected text failed")
}
