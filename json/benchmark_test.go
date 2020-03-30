package json

import (
	"bytes"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/klauspost/compress/zstd"
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

func benchmarkFromFile(b *testing.B, filename string) {
	infile, err := os.Open(filepath.Join("testdata", filename+".json.zst"))
	if err != nil {
		b.Fatal(err)
	}

	msg, err := ioutil.ReadAll(infile)
	infile.Close()
	if err != nil {
		b.Fatal(err)
	}

	dec, err := zstd.NewReader(nil)
	if err != nil {
		b.Fatal(err)
	}

	msg, err = dec.DecodeAll(msg, nil)
	dec.Close()
	if err != nil {
		b.Fatal(err)
	}

	// fmt.Println(string(msg))

	b.SetBytes(int64(len(msg)))
	b.ReportAllocs()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		sax := new(EmptySAX)
		rMsg := bytes.NewReader(msg)
		if err := Parse(rMsg, sax, nil); err != nil {
			b.Fatal(err)
		}
	}
}
