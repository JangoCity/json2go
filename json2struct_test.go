package json2struct

import (
	"bufio"
	"bytes"
	"fmt"
	"testing"
)

var basic = []byte(`{
	"foo": "fooer",
	"bar": "bars",
	"biz": 1,
	"baz": 42.1,
	"foo_bar": "frood"
}`)
var expectedBasic = "type Basic struct {\n\tBar string `json:\"bar\"`\n\tBaz float64 `json:\"baz\"`\n\tBiz int `json:\"biz\"`\n\tFoo string `json:\"foo\"`\n\tFooBar string `json:\"foo_bar\"`\n}\n"

func TestBasicStruct(t *testing.T) {
	def, err := Gen("Basic", basic)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if string(def) != expectedBasic {
		t.Errorf("expected %q got %q", expectedBasic, string(def))
	}
}

var intermediate = []byte(`{
	"name": "Marvin",
	"id": 42,
	"bot": true,
	"quotes": [
		"I think you ought to know I'm feeling very depressed",
		"Life? Don't talk to me about life!"
	],
	"ints": [
		1,
		2,
		3,
		4
	],
	"floats": [
		1.01,
		2.01,
		3.01
	],
	"bools": [
		true,
		false
	],
	"date": "Fri Jan 23 13:02:46 +0000 2015"
}`)
var expectedIntermediate = "type Intermediate struct {\n\tBools []bool `json:\"bools\"`\n\tBot bool `json:\"bot\"`\n\tDate string `json:\"date\"`\n\tFloats []float64 `json:\"floats\"`\n\tId int `json:\"id\"`\n\tInts []int `json:\"ints\"`\n\tName string `json:\"name\"`\n\tQuotes []string `json:\"quotes\"`\n}\n"

func TestIntermediateStruct(t *testing.T) {
	def, err := Gen("Intermediate", intermediate)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if string(def) != expectedIntermediate {
		t.Errorf("expected %q got %q", expectedIntermediate, string(def))
	}
}

// widget example from json.org/example.html
var widget = []byte(`{
	"widget": {
		"debug": "on",
		"window": {
			"title": "Sample Konfabulator Widget",
        		"name": "main_window",
        		"width": 500,
        		"height": 500
		},
		"image": {
        		"src": "Images/Sun.png",
        		"name": "sun1",
        		"hOffset": 250,
        		"vOffset": 250,
        		"alignment": "center"
		},
		"text": {
			"data": "Click Here",
        		"size": 36,
        		"style": "bold",
        		"name": "text1",
        		"hOffset": 250,
        		"vOffset": 100,
        		"alignment": "center",
        		"onMouseUp": "sun1.opacity = (sun1.opacity / 100) * 90;"
		}
	}
}`)

var expectedWidget = "type TestW struct {\n\tWidget `json:\"widget\"`\n}\ntype Widget struct {\n\tDebug string `json:\"debug\"`\n\tImage `json:\"image\"`\n\tText `json:\"text\"`\n\tWindow `json:\"window\"`\n}\ntype Image struct {\n\tAlignment string `json:\"alignment\"`\n\tHOffset int `json:\"hOffset\"`\n\tName string `json:\"name\"`\n\tSrc string `json:\"src\"`\n\tVOffset int `json:\"vOffset\"`\n}\ntype Text struct {\n\tAlignment string `json:\"alignment\"`\n\tData string `json:\"data\"`\n\tHOffset int `json:\"hOffset\"`\n\tName string `json:\"name\"`\n\tOnMouseUp string `json:\"onMouseUp\"`\n\tSize int `json:\"size\"`\n\tStyle string `json:\"style\"`\n\tVOffset int `json:\"vOffset\"`\n}\ntype Window struct {\n\tHeight int `json:\"height\"`\n\tName string `json:\"name\"`\n\tTitle string `json:\"title\"`\n\tWidth int `json:\"width\"`\n}\n"

func TestWidget(t *testing.T) {
	s, err := Gen("TestW", widget)
	if err != nil {
		t.Errorf("unexpected error %q", err)
	}
	if string(s) != expectedWidget {
		t.Errorf("expected:\n%q, got:\n%q", expectedWidget, string(s))
	}
}

var wnull = []byte(`{
	"foo": "fooer",
	"bar": null,
	"biz": 1,
	"baz": 42.1,
	"foo_bar": "frood"
}`)

var expectedWNull = "type WNull struct {\n\tBar interface{} `json:\"bar\"`\n\tBaz float64 `json:\"baz\"`\n\tBiz int `json:\"biz\"`\n\tFoo string `json:\"foo\"`\n\tFooBar string `json:\"foo_bar\"`\n}\n"

func TestWNull(t *testing.T) {
	def, err := Gen("WNull", wnull)
	if err != nil {
		t.Errorf("unexpected error: %s", err)
	}
	if string(def) != expectedWNull {
		t.Errorf("expected %q got %q", expectedWNull, string(def))
	}
}

func TestTransmogrify(t *testing.T) {
	tests := []struct {
		pkg        string
		importJSON bool
		expected   string
	}{
		{"", false, fmt.Sprintf("package main\n\n%s", expectedBasic)},
		{"test", false, fmt.Sprintf("package test\n\n%s", expectedBasic)},
		{"", true, fmt.Sprintf("package main\n\nimport (\n\t\"encoding/json\"\n)\n\n%s", expectedBasic)},
		{"test", true, fmt.Sprintf("package test\n\nimport (\n\t\"encoding/json\"\n)\n\n%s", expectedBasic)},
	}

	var b bytes.Buffer
	for i, test := range tests {
		// create reader
		r := bytes.NewReader(basic)
		// create Writer
		w := bufio.NewWriter(&b)

		calvin := NewTransmogrifier("Basic", r, w)
		if test.pkg != "" {
			calvin.SetPkg(test.pkg)
		}
		calvin.SetImportJSON(test.importJSON)
		err := calvin.Gen()
		if err != nil {
			t.Errorf("%d: unexpected error %q", i, err)
			continue
		}
		w.Flush()
		if b.String() != test.expected {
			t.Errorf("%d: expected %q, got %q", i, test.expected, b.String())
		}
		b.Reset()
	}
}
