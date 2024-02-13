package htmlmatch

import (
	"bytes"
	"fmt"
	"testing"

	_ "embed"

	"golang.org/x/net/html"
)

//go:embed testdata/example.html
var example []byte

func TestContainsTree(t *testing.T) {
	exampleNode, err := html.Parse(bytes.NewBuffer(example))
	if err != nil {
		t.Fatal(err)
	}

	tests := map[string]struct {
		partial string
		full    *html.Node
		match   bool
	}{
		"basic match":    {partial: "<td>a1</td>", full: exampleNode, match: true},
		"basic mismatch": {partial: "<td>noes</td>", full: exampleNode, match: false},
		"should match text inside <td>, not present": {
			partial: "<tr>a1</tr>",
			full:    exampleNode,
			match:   true,
		},
		"should not match text inside <td>": {
			partial: "<tr>nodes</tr>",
			full:    exampleNode,
			match:   false,
		},
		"should match multiple <b> elements": {
			partial: `<h2>
				<b></b>
				<b></b>
			</h2>`,
			full:  exampleNode,
			match: true,
		},
		"should not match too many <b> elements": {
			partial: `<h2>
				<b></b>
				<b></b>
				<b></b>
			</h2>`,
			full:  exampleNode,
			match: false,
		},
		"should match ordered elements with attributes": {
			partial: `<h2>
				<b class="b1"></b>
				<b class="b2"></b>
			</h2>`,
			full:  exampleNode,
			match: true,
		},
		"should not match elements with attributes in wrong order": {
			partial: `<h2>
				<b class="b2"></b>
				<b class="b1"></b>
			</h2>`,
			full:  exampleNode,
			match: false,
		},
		"should match table contents": {
			partial: `<table>
				<tr></tr>
				<tr></tr>
			</table>`,
			full:  exampleNode,
			match: true,
		},
		"should match substring": {
			partial: `<div>substring:ipsum</div>`,
			full:    exampleNode,
			match:   true,
		},
		"should not verbatim": {
			partial: `<div>verbatim: world </div>`,
			full:    exampleNode,
			match:   false,
		},
		"should not match verbatim": {
			partial: `<div>verbatim:world</div>`,
			full:    exampleNode,
			match:   false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			partial, err := ParseVerbatim(tc.partial)
			if err != nil {
				t.Fatal(err)
			}

			if tc.match != ContainsTree(tc.full, partial) {
				t.Errorf("expected %v", tc.match)
			}
		})
	}
}

func ExampleContainsTree() {
	full := MustParseVerbatim(`<html>
		<body>
			<div hx-get="/foo" hx-target="#mytarget">Click here</div>
			<div>lorem ipsum</div>
			<div>lorem <span id="mytarget">ipsum</span></div>
		</body>
		</html>`)
	pattern := MustParseVerbatim(`
		<div hx-target="#mytarget"></div>
		<span id="mytarget" />`)
	fmt.Println(ContainsTree(full, pattern))
	// Output: true
}
