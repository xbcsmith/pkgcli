package utils

import (
	"testing"

	"gotest.tools/assert"
	is "gotest.tools/assert/cmp"
)

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

func TestConverterJson(t *testing.T) {
	tests := NewTests()
	expected := `- key:`
	actual, err := converter(tests.a, false)
	assert.Assert(t, is.Nil(err))
	assert.Assert(t, !IsJSON(actual))
	assert.Assert(t, is.Contains(string(actual), expected))
}

func TestConverterAlsoJson(t *testing.T) {
	tests := NewTests()
	expected := `brackets:`
	actual, err := converter(tests.e, false)
	assert.Assert(t, is.Nil(err))
	assert.Assert(t, !IsJSON(actual))
	assert.Assert(t, is.Contains(string(actual), expected))
}

func TestConverterYaml(t *testing.T) {
	tests := NewTests()
	expected := `"bar": [`
	actual, err := converter(tests.b, false)
	assert.Assert(t, is.Nil(err))
	assert.Assert(t, IsJSON(actual))
	assert.Assert(t, is.Contains(string(actual), expected))
}

func TestConverterYamlNoIndent(t *testing.T) {
	tests := NewTests()
	expected := `bar":[`
	actual, err := converter(tests.b, true)
	assert.Assert(t, is.Nil(err))
	assert.Assert(t, IsJSON(actual))
	assert.Assert(t, is.Contains(string(actual), expected))
}

func TestGetEnv(t *testing.T) {
	os.Setenv("FOO", "1")
	foo := GetEnv("FOO", "2")
	assert.Assert(t, foo == "1")
	bar := GetEnv("BAR", "42")
	assert.Assert(t, bar == "42")
}

func TestNewULID(t *testing.T) {
	u := NewULID()
	assert.Assert(t, len(u) == 26)
}
