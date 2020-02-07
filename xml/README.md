# xml
A fast XML tokenizer and SAX parser for GO.

While the standard GO xml reader is (probably) more complete, this tokenizer is measured to be approximately 5 to 6 times faster.

Compared to other SAX parsers this implementation is somewhat simpler and consider namespace as an integral part of tag and value identifiers. Separation of namespace and identifier must be done on the recieving/client side.

## Example
```
package main

import (
	"bytes"
	"fmt"
	"github.com/xianhammer/format/xml"
)

type receiver struct {
	xml.Partial
}

func (r *receiver) Tag(name []byte) {
	fmt.Printf("Tag [%s]\n", name)
}

func (r *receiver) TagEnd(autoclose bool) {
	fmt.Printf("Tagend\n")
}

func (r *receiver) Text(value []byte) {
	fmt.Printf("Text [%s]\n", value)
}

var sample = `<!DOCTYPE cafProductFeed SYSTEM "http://www.affiliatewindow.com/DTD/affiliate/datafeed.1.5.dtd">
<tag attr1="value1" attr2="2" hidden="1">'MixedChars'!$A$1:$AE$10754</tag>
`

func main() {
	br := bytes.NewBufferString(sample)
	t := NewTokenizer(e)
	t.ReadFrom(br)
}
```
