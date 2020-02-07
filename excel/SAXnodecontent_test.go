package excel

import (
	"bytes"
	"testing"
)

const sampleXMLsharedstrings = `<?xml version="1.0" encoding="UTF-8" standalone="yes" ?>
<sst count="6" uniqueCount="6" xmlns="http://schemas.openxmlformats.org/spreadsheetml/2006/main">
	<si>
		<t>First</t>
	</si>
	<si>
		<t>2</t>
	</si>
	<si>
		<t></t>
	</si>
	<si>
		<t>Fourth</t>
	</si>
	<si>
		<t/>
	</si>
	<si>
		<t>Sixth</t>
	</si>
</sst>`

func TestNodeContentEmpty(t *testing.T) {
	r := bytes.NewBufferString(sampleXMLsharedstrings)
	content, err := NodeContent(r, []byte("UNKNOWNTAG"))
	if err != nil {
		t.Errorf("Expected error [%v], got [%v]", nil, err)
	}
	if len(content) != 0 {
		t.Errorf("Expected content count [%v], got [%v]", 0, len(content))
	}
}

func TestNodeContent(t *testing.T) {
	expectedCount := 6
	r := bytes.NewBufferString(sampleXMLsharedstrings)
	content, err := NodeContent(r, []byte("t"))
	if err != nil {
		t.Errorf("Expected error [%v], got [%v]", nil, err)
	}
	if len(content) != expectedCount {
		t.Errorf("Expected content count [%v], got [%v]", expectedCount, len(content))
	}
	expect := []string{
		"First",
		"2",
		"",
		"Fourth",
		"",
		"Sixth",
	}

	for i, e := range expect {
		if content[i] != e {
			t.Errorf("Expected element at [%v] [%v], got [%v]", i, e, content[i])
		}
	}
}
