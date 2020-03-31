package parse

import (
	"bytes"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/klauspost/compress/zstd"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func BenchmarkApache_builds(b *testing.B)  { benchmarkFromFile(b, "apache_builds") }
func BenchmarkCanada(b *testing.B)         { benchmarkFromFile(b, "canada") }
func BenchmarkCitm_catalog(b *testing.B)   { benchmarkFromFile(b, "citm_catalog") }
func BenchmarkGithub_events(b *testing.B)  { benchmarkFromFile(b, "github_events") }
func BenchmarkGsoc_2018(b *testing.B)      { benchmarkFromFile(b, "gsoc-2018") }
func BenchmarkInstruments(b *testing.B)    { benchmarkFromFile(b, "instruments") }
func BenchmarkMarine_ik(b *testing.B)      { benchmarkFromFile(b, "marine_ik") }
func BenchmarkMesh(b *testing.B)           { benchmarkFromFile(b, "mesh") }
func BenchmarkMesh_pretty(b *testing.B)    { benchmarkFromFile(b, "mesh.pretty") }
func BenchmarkNumbers(b *testing.B)        { benchmarkFromFile(b, "numbers") }
func BenchmarkRandom(b *testing.B)         { benchmarkFromFile(b, "random") }
func BenchmarkTwitter(b *testing.B)        { benchmarkFromFile(b, "twitter") }
func BenchmarkTwitterescaped(b *testing.B) { benchmarkFromFile(b, "twitterescaped") }
func BenchmarkUpdate_center(b *testing.B)  { benchmarkFromFile(b, "update-center") }

var (
	win16be  = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	utf16bom = unicode.BOMOverride(win16be.NewDecoder())
)

func benchmarkFromFile(b *testing.B, filename string) {
	infile, err := os.Open(filepath.Join("testdata", "json", filename+".json.zst"))
	if err != nil {
		b.Fatal(err)
	}

	dec, err := zstd.NewReader(infile)
	infile.Close()
	if err != nil {
		b.Fatal(err)
	}

	// TODO Input is UNICDOE - either handle this in the parser OR convert prior to parsing.
	r := transform.NewReader(dec, utf16bom)
	var buf bytes.Buffer
	if _, err = io.Copy(&buf, r); err != nil {
		b.Fatal(err)
	}

	src := buf.Bytes()

	b.SetBytes(int64(len(src)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		if v, n := JSON(src, nil); v == nil {
			b.Fatal("nil was returned")
		} else if n != len(src) {
			b.Fatal("Length mismatch", n, len(src), v)
		}
	}
}
