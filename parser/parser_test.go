package parser

import (
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	const src = `import "test";
import "something/else" as something;

let v = 3;`

	script, err := Parse(strings.NewReader(src))
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%+v", script)
}
