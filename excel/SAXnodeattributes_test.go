package excel

import (
	"bytes"
	"testing"
)

const sampleXMLrels = `<?xml version="1.0" encoding="UTF-8" standalone="yes"?>
<Relationships xmlns="http://schemas.openxmlformats.org/package/2006/relationships">
<Relationship Id="rId3" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties" Target="docProps/app.xml"/>
<Relationship Id="rId2" Type="http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties" Target="docProps/core.xml"/>
<Relationship Id="rId1" Type="http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument" Target="xl/workbook.xml"/>
</Relationships>`

func TestNodeAttributesEmpty(t *testing.T) {
	// NodeAttributes(r io.Reader, tag []byte) (attributes []map[string]string, err error)
	r := bytes.NewBufferString(sampleXMLrels)
	attr, err := NodeAttributes(r, []byte("UNKNOWNTAG"))
	if err != nil {
		t.Errorf("Expected error [%v], got [%v]", nil, err)
	}
	if len(attr) != 0 {
		t.Errorf("Expected attribute count [%v], got [%v]", 0, len(attr))
	}
}

func TestNodeAttributesOK(t *testing.T) {
	// NodeAttributes(r io.Reader, tag []byte) (attributes []map[string]string, err error)
	r := bytes.NewBufferString(sampleXMLrels)
	attr, err := NodeAttributes(r, []byte("Relationship"))
	if err != nil {
		t.Errorf("Expected error [%v], got [%v]", nil, err)
	}
	if len(attr) != 3 {
		t.Errorf("Expected attribute count [%v], got [%v]", 3, len(attr))
	}

	tests := [][]string{
		{"Id", "rId3", "Type", "http://schemas.openxmlformats.org/officeDocument/2006/relationships/extended-properties", "Target", "docProps/app.xml"},
		{"Id", "rId2", "Type", "http://schemas.openxmlformats.org/package/2006/relationships/metadata/core-properties", "Target", "docProps/core.xml"},
		{"Id", "rId1", "Type", "http://schemas.openxmlformats.org/officeDocument/2006/relationships/officeDocument", "Target", "xl/workbook.xml"},
	}

	for i, element := range attr {
		test := tests[i]
		if len(test) != len(element)*2 { // Elements are key-value pairs
			t.Errorf("Expected element len [%v], got [%v]", len(test), len(element)*2)
		}

		for j := 0; j < len(tests); j += 2 {
			key, expect := test[j], test[j+1]
			got := element[key]
			if expect != got {
				t.Errorf("Expected value [%v], got [%v], for key [%v]", expect, got, key)
			}
		}
	}
}
