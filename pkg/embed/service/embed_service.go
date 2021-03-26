package service

import (
	"bytes"
	"compress/gzip"
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"

	"github.com/thecxx/embed/pkg/embed/asset/config"
)

var Embed = NewEmbedService()

type EmbedService struct {
}

// NewEmbedService returns a embed service.
func NewEmbedService() *EmbedService {
	return new(EmbedService)
}

func (e *EmbedService) Build(ctx context.Context) error {

	buffer := bytes.NewBuffer(nil)

	e.emitPackage(buffer, config.Embed.Package)
	e.emitImport(buffer, []string{
		"compress/gzip",
		"io",
		"io/ioutil",
	})

	if len(config.Embed.Items) > 0 {

		storage := bytes.NewBuffer(nil)

		fmt.Fprintf(buffer, "var(\n")
		fmt.Fprintf(storage, "var(\n")
		for _, f := range config.Embed.Items {
			buf := bytes.NewBuffer(nil)
			com := gzip.NewWriter(buf)

			data, err := ioutil.ReadFile(f.File)
			if err != nil {
				return err
			}
			com.Write(data)
			com.Flush()

			fmt.Fprintf(storage, "    %s_data = []byte{", f.Name)

			for i, c := range buf.Bytes() {
				if i+1 >= buf.Len() {
					fmt.Fprintf(storage, "0x%x", c)
				} else {
					fmt.Fprintf(storage, "0x%x, ", c)
				}
			}
			fmt.Fprintf(storage, "}\n")

			fmt.Fprintf(buffer, "    %s = file{\"%s\", gzipReader(&buffer{data: %s_data, index: 0})}\n", f.Name, f.File, f.Name)

		}
		fmt.Fprintf(buffer, ")\n\n")
		fmt.Fprintf(storage, ")\n\n")

		fmt.Fprintf(buffer, "%s\n", storage.String())
	}

	e.emitTemplate(buffer)

	return ioutil.WriteFile(config.Embed.Path, buffer.Bytes(), 0777)

}

func (e *EmbedService) emitPackage(buffer io.Writer, name string) {
	fmt.Fprintf(buffer, "package %s\n\n", name)
}

func (e *EmbedService) emitImport(buffer io.Writer, packages []string) {
	if len(packages) > 0 {
		fmt.Fprintf(buffer, "import(\n")
		for _, pkg := range packages {
			fmt.Fprintf(buffer, "    \"%s\"\n", pkg)
		}
		fmt.Fprintf(buffer, ")\n\n")
	}
}

func (e *EmbedService) emitTemplate(buffer io.Writer) {
	tpl := `type buffer struct {
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

type emptyBuffer struct {
}

func (b *emptyBuffer) Read([]byte) (int, error) {
	return 0, io.EOF
}

var invalidBuffer = new(emptyBuffer)

func gzipReader(in io.Reader) io.Reader {
	r, err := gzip.NewReader(in)
	if err != nil {
		return invalidBuffer
	}
	return r
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
`
	fmt.Fprintf(buffer, "%s\n", tpl)
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
