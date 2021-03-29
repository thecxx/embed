package service

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
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
		vars    = make(map[string]string)
		bufs    = make(map[string][]byte)
		buffer  = bytes.NewBuffer(nil)
		imports = []string{"io", "io/ioutil"}
	)

	var (
		sf       = pack.NewSourceFile(config.Embed.Package)
		compress = ""
	)

	sf.Import("io")
	sf.Import("io/ioutil")

	switch config.Embed.Compress {
	// gzip
	case "gzip", "gz":
		compress = "gzip"
		sf.Import("compress/gzip")
		sf.DeclareGzipReader()

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
		if len(data) <= 0 || sign == "" {
			sign, data = "empty", make([]byte, 0)
		}
		if _, ok := bufs[sign]; !ok {
			bufs[sign] = data
			// Save bytes
			sf.DeclareVar(pack.Variable{Name: fmt.Sprintf("_%s", sign), Assign: "[]byte{}"})
		}

		b := fmt.Sprintf("&buffer{data: _%s, index: 0}", sign)

		if len(data) > 0 {
			switch config.Embed.Compress {
			// gzip
			case "gzip", "gz":
				b = fmt.Sprintf("gzipReader(%s)", b)
			}
		}

		vars[f.Name] = fmt.Sprintf("file{\"%s\", %s}", f.File, b)
	}

	e.emitPackage(buffer, config.Embed.Package)
	e.emitImport(buffer, imports)

	if len(vars) > 0 {
		fmt.Fprintf(buffer, "var (\n")
		for n, assign := range vars {
			fmt.Fprintf(buffer, "    %s = %s\n", n, assign)
		}
		fmt.Fprintf(buffer, ")\n\n")
	}

	if len(bufs) > 0 {
		fmt.Fprintf(buffer, "var (\n")
		for sign, data := range bufs {
			if len(data) > 0 {
				buf := bytes.NewBuffer(nil)
				com := gzip.NewWriter(buf)

				com.Write(data)
				com.Flush()

				fmt.Fprintf(buffer, "    _%s = []byte{", sign)

				for i, c := range buf.Bytes() {
					if i+1 >= buf.Len() {
						fmt.Fprintf(buffer, "0x%x", c)
					} else {
						fmt.Fprintf(buffer, "0x%x, ", c)
					}
				}

				fmt.Fprintf(buffer, "}\n")
			} else {
				fmt.Fprintf(buffer, "    _%s = []byte{}\n", sign)
			}
		}
		fmt.Fprintf(buffer, ")\n\n")
	}

	e.emitTemplate(buffer)

	if err := os.MkdirAll(path.Dir(config.Embed.Path), 0644); err != nil {
		return err
	}
	return ioutil.WriteFile(config.Embed.Path, buffer.Bytes(), 0644)

}

func (e *EmbedService) emitPackage(buffer io.Writer, name string) {
	fmt.Fprintf(buffer, "package %s\n\n", name)
}

func (e *EmbedService) emitImport(buffer io.Writer, packages []string) {
	if len(packages) > 0 {
		fmt.Fprintf(buffer, "import (\n")
		for _, pkg := range packages {
			fmt.Fprintf(buffer, "    \"%s\"\n", pkg)
		}
		fmt.Fprintf(buffer, ")\n\n")
	}
}

func (e *EmbedService) emitVariable(buffer io.Writer, name, filename, buf string, cast string) {
	buf = fmt.Sprintf("&buffer{data: %s, index: 0}", buf)
	if len(cast) >= 0 {
		buf = fmt.Sprintf("%s(%s)", cast, buf)
	}
	fmt.Fprintf(buffer, "    %s = file{\"%s\", %s}\n", name, filename, buf)
}

func (e *EmbedService) emitBuffer(buffer io.Writer, name string, data []byte) {

}

func (e *EmbedService) emitTemplate(buffer io.Writer) {
	fmt.Fprintf(buffer, "%s\n", `type buffer struct {
	data  []byte
	index int64
}

func (b *buffer) Read(p []byte) (int, error) {
	if len(b.data) <= 0 {
		return 0, io.EOF
	}
	if b.index >= int64(len(b.data)) {
		return 0, io.EOF
	}
	// Copy
	n := copy(p, b.data[b.index:])
	b.index += int64(n)

	return n, nil
}

type file struct {
	file   string
	reader io.Reader
}

func (r *file) Read(buffer []byte) (int, error) {
	return r.reader.Read(buffer)
}

func (r *file) ReadAll() ([]byte, error) {
	return ioutil.ReadAll(r.reader)
}

func (r *file) Bytes() []byte {
	if b, err := ioutil.ReadAll(r.reader); err == nil {
		return b
	}
	return make([]byte, 0)
}
`)
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
