package json

import (
	"bytes"
	"encoding/json"
	"io"
	"io/ioutil"
	"path/filepath"
	"testing"

	"golang.org/x/text/encoding/unicode"
)

func BenchmarkApache_builds(b *testing.B)      { benchmarkFromFile(b, "apache_builds") }
func BenchmarkApache_builds_std(b *testing.B)  { stdBenchmarkFromFile(b, "apache_builds") }
func BenchmarkCanada(b *testing.B)             { benchmarkFromFile(b, "canada") }
func BenchmarkCanada_std(b *testing.B)         { stdBenchmarkFromFile(b, "canada") }
func BenchmarkCitm_catalog(b *testing.B)       { benchmarkFromFile(b, "citm_catalog") }
func BenchmarkCitm_catalog_std(b *testing.B)   { stdBenchmarkFromFile(b, "citm_catalog") }
func BenchmarkGithub_events(b *testing.B)      { benchmarkFromFile(b, "github_events") }
func BenchmarkGithub_events_std(b *testing.B)  { stdBenchmarkFromFile(b, "github_events") }
func BenchmarkGsoc_2018(b *testing.B)          { benchmarkFromFile(b, "gsoc-2018") }
func BenchmarkGsoc_2018_std(b *testing.B)      { stdBenchmarkFromFile(b, "gsoc-2018") }
func BenchmarkInstruments(b *testing.B)        { benchmarkFromFile(b, "instruments") }
func BenchmarkInstruments_std(b *testing.B)    { stdBenchmarkFromFile(b, "instruments") }
func BenchmarkMarine_ik(b *testing.B)          { benchmarkFromFile(b, "marine_ik") }
func BenchmarkMarine_ik_std(b *testing.B)      { stdBenchmarkFromFile(b, "marine_ik") }
func BenchmarkMesh(b *testing.B)               { benchmarkFromFile(b, "mesh") }
func BenchmarkMesh_std(b *testing.B)           { stdBenchmarkFromFile(b, "mesh") }
func BenchmarkMesh_pretty(b *testing.B)        { benchmarkFromFile(b, "mesh.pretty") }
func BenchmarkMesh_pretty_std(b *testing.B)    { stdBenchmarkFromFile(b, "mesh.pretty") }
func BenchmarkNumbers(b *testing.B)            { benchmarkFromFile(b, "numbers") }
func BenchmarkNumbers_std(b *testing.B)        { stdBenchmarkFromFile(b, "numbers") }
func BenchmarkRandom(b *testing.B)             { benchmarkFromFile(b, "random") }
func BenchmarkRandom_std(b *testing.B)         { stdBenchmarkFromFile(b, "random") }
func BenchmarkTwitter(b *testing.B)            { benchmarkFromFile(b, "twitter") }
func BenchmarkTwitter_std(b *testing.B)        { stdBenchmarkFromFile(b, "twitter") }
func BenchmarkTwitterescaped(b *testing.B)     { benchmarkFromFile(b, "twitterescaped") }
func BenchmarkTwitterescaped_std(b *testing.B) { stdBenchmarkFromFile(b, "twitterescaped") }
func BenchmarkUpdate_center(b *testing.B)      { benchmarkFromFile(b, "update-center") }
func BenchmarkUpdate_center_std(b *testing.B)  { stdBenchmarkFromFile(b, "update-center") }

var (
	win16le  = unicode.UTF16(unicode.LittleEndian, unicode.IgnoreBOM)
	win16be  = unicode.UTF16(unicode.BigEndian, unicode.IgnoreBOM)
	utf16bom = unicode.BOMOverride(win16le.NewDecoder())
	utf8bom  = unicode.BOMOverride(unicode.UTF8.NewDecoder())
)

func readFile(b *testing.B, filename string) (src []byte, srcLen int) {
	// infile, err := os.Open(filepath.Join("testdata", "json", filename+".json.zst"))
	// if err != nil {
	// 	b.Fatal(err)
	// }
	// defer infile.Close()

	// dec, err := zstd.NewReader(infile)
	// if err != nil {
	// 	b.Fatal(err)
	// }
	// defer dec.Close()

	// r := transform.NewReader(dec, utf8bom)
	// src, err := ioutil.ReadAll(r)
	// if err != nil {
	// 	b.Fatal(err)
	// }

	src, err := ioutil.ReadFile(filepath.Join("testdata", "json0", filename+"-utf8.json"))
	if err != nil {
		b.Fatal(err)
	}
	srcLen = len(bytes.TrimSpace(src))

	b.SetBytes(int64(srcLen))

	return
}

func benchmarkFromFile(b *testing.B, filename string) {
	src, srcLen := readFile(b, filename)
	buffer := make([]byte, srcLen) // Needed, since giving JSON a nil argument will corrupt the input after first loop.
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		got, _, err := Parse(src, buffer)
		if err != nil && err != io.EOF {
			b.Fatal(err)
		}
		if got == nil {
			b.Fatalf("Expected non-nil (%v, %v)\n", got, err)
			// b.Fatal("Expected non-nil (%v, %v, %v)\n", v, n, err)
		}
		// if gotSize != srcLen {
		// 	b.Fatal("Length mismatch", gotSize, srcLen)
		// }
	}
}

func stdBenchmarkFromFile(b *testing.B, filename string) {
	src, _ := readFile(b, filename)
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		var got interface{}
		err := json.Unmarshal(src, &got)
		if err != nil {
			b.Fatal(err)
		}
		if got == nil {
			b.Fatalf("Expected non-nil (%v, %v)\n", got, err)
			// b.Fatal("Expected non-nil (%v, %v, %v)\n", v, n, err)
		}
	}
}
