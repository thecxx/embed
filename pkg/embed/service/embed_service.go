package service

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"

	"github.com/thecxx/embed/pkg/pack"

	"github.com/thecxx/embed/pkg/embed/asset/config"
)

var Embed = NewEmbedService()

type EmbedService struct {
}

// NewEmbedService returns a embed service.
func NewEmbedService() *EmbedService {
	return new(EmbedService)
}

func (e *EmbedService) Init(ctx context.Context, file string) error {
	return ioutil.WriteFile(file, ([]byte)(`# embed build [-f embed.yaml]

---
pkg: embed

path: embed/embed.go

compress: gz

items:
  - name: TestFile
    file: test.txt

`), 0644)
}

func (e *EmbedService) Build(ctx context.Context) error {

	var (
		compress = ""
		bufs     = make(map[string][]byte)
		vars1    = make([]pack.Variable, 0)
		vars2    = make([]pack.Variable, 0)
		sf       = pack.NewSourceFile(config.Embed.Package)
		arch     = pack.NewSourceFile(config.Embed.Package)
	)

	sf.Import("bytes")
	sf.Import("io")
	sf.Import("io/ioutil")

	switch config.Embed.Compress {
	// gzip
	case "gzip", "gz":
		compress = "gzip"
		sf.Import("compress/gzip")
		sf.DeclareGzipReader()
	// Empty
	case "":
		// Nothing to do
	// Not supported
	default:
		return errors.New("compress method not supported")
	}

	for _, f := range config.Embed.Items {
		// Load file
		data, sign, err := e.loadFile(f.File)
		if err != nil {
			return err
		}
		if len(data) <= 0 {
			sign, data = "ffffffffffffffffffffffffffffffff", make([]byte, 0)
		} else {
			if _, ok := bufs[sign]; !ok {
				bufs[sign] = data

				var (
					raw []byte
					tmp bytes.Buffer
				)
				// Save bytes
				if len(compress) > 0 {
					writer := gzip.NewWriter(&tmp)
					writer.Write(data)
					writer.Flush()
					writer.Close()

					raw = tmp.Bytes()
				} else {
					raw = data
				}
				vars1 = append(vars1, pack.Variable{
					Name:   fmt.Sprintf("_%s", sign),
					Assign: e.formatBytes(raw),
				})
			}
		}

		var buf string
		if len(data) > 0 && len(compress) > 0 {
			buf = fmt.Sprintf("%sReader(_%s)", compress, sign)
		} else {
			buf = fmt.Sprintf("bufferReader(_%s)", sign)
		}
		// Save buffer
		vars2 = append(vars2, pack.Variable{
			Name:   f.Name,
			Assign: fmt.Sprintf("file{\"%s\", %s}", f.File, buf),
			// Comment: fmt.Sprintf("File: %s", f.File),
		})
	}

	if len(vars1) > 0 {
		if config.Embed.Archive {
			arch.DeclareVar(vars1...)
		} else {
			sf.DeclareVar(vars1...)
		}
	}

	if len(vars2) > 0 {
		sf.DeclareVar(vars2...)
	}

	sf.DeclareFileReader()

	dir := path.Dir(config.Embed.Path)

	err := os.MkdirAll(dir, 0644)
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(dir+"/archive.go", arch.Bytes(), 0644)
	if err != nil {
		return err
	}

	return ioutil.WriteFile(config.Embed.Path, sf.Bytes(), 0644)

}

func (e *EmbedService) loadFile(filename string) ([]byte, string, error) {
	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, "", err
	}
	if len(buf) <= 0 {
		return nil, "", nil
	}
	md5sum := md5.New()
	md5sum.Write(buf)
	return buf, hex.EncodeToString(md5sum.Sum(nil)), nil
}

func (e *EmbedService) formatBytes(data []byte) string {
	if len(data) <= 0 {
		return "[]byte{}"
	}
	buf := bytes.NewBufferString(fmt.Sprintf("0x%02X", data[0]))
	for i := 1; i < len(data); i++ {
		fmt.Fprintf(buf, ", 0x%02X", data[i])
	}
	return fmt.Sprintf("[]byte{%s}", buf.String())
}
