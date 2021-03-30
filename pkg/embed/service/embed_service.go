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

	"github.com/thecxx/embed/pkg/embed/asset/config"
	"github.com/thecxx/embed/pkg/pack"
)

const (
	hextable  = "0123456789ABCDEF"
	emptySign = "ffffffffffffffffffffffffffffffff"
)

var Embed = NewEmbedService()

type EmbedService struct {
}

// NewEmbedService returns a embed service.
func NewEmbedService() *EmbedService {
	return new(EmbedService)
}

func (e *EmbedService) Init(ctx context.Context, file string) error {

	stat, err := os.Stat(file)
	if err != nil {
		if !os.IsNotExist(err) {
			return err
		}
	} else {
		if stat.IsDir() {
			return errors.New("invalid config file")
		} else {
			return errors.New("the config file already exists")
		}
	}

	return ioutil.WriteFile(file, ([]byte)(`# embed build [-f embed.yaml]
---

# The package name
pkg: "embed"

# Source file path
path: "embed/embed.go"

# The compression method for embedded resource, it can be set to [no] [gz|gzip]
compress: "gz"

# Archive the data to a separate file
archive: false

# Resource list
items:
  - name: "TestFile"
    file: "tests/test.txt"
    comment: "A test file, \njust for test."

  - name: "EmptyFile"
    file: "tests/empty.txt"
    comment: "An empty file"

`), 0755)
}

func (e *EmbedService) Build(ctx context.Context) error {

	var (
		compressor = ""
		datas      = make(map[string][]byte)
		exports    = make([]pack.Variable, 0)
		archives   = make([]pack.Variable, 0)
		sf         = pack.NewSourceFile(config.Embed.Package)
		af         = pack.NewSourceFile(config.Embed.Package)
	)

	sf.Import("bytes")
	sf.Import("io")
	sf.Import("io/ioutil")
	sf.DeclareFileReader()

	switch config.Embed.Compress {
	// No compressor
	case "no":
		compressor = ""
	// gzip
	case "gzip", "gz":
		compressor = "gzip"
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
		if len(data) <= 0 {
			sign, data = emptySign, make([]byte, 0)
		}

		compress := len(data) > 0 && len(compressor) > 0

		if _, ok := datas[sign]; !ok {
			datas[sign] = data

			var (
				err error
				raw []byte
			)
			if compress {
				// Compress with gzip
				if raw, err = e.compressWithGzip(data); err != nil {
					return err
				}
			} else {
				raw = data
			}
			archives = append(archives, pack.Variable{
				Name:   fmt.Sprintf("_%s", sign),
				Assign: e.formatBytes(raw),
			})
		}

		var buf string

		if compress {
			buf = fmt.Sprintf("%sReader(_%s)", compressor, sign)
		} else {
			buf = fmt.Sprintf("directReader(_%s)", sign)
		}

		exports = append(exports, pack.Variable{
			Name:    f.Name,
			Comment: f.Comment,
			Assign:  fmt.Sprintf("file{\"%s\", %s}", f.File, buf),
		})
	}

	directory := path.Dir(config.Embed.Path)
	archfile := directory + "/archive.go"

	// Create the package directory
	if err := os.MkdirAll(directory, 0755); err != nil {
		return err
	}
	// Clean
	os.Remove(config.Embed.Path)
	os.Remove(archfile)

	// Save all data
	if len(archives) > 0 {
		if config.Embed.Archive {
			for _, arch := range archives {
				af.DeclareVar(arch)
			}
			if err := ioutil.WriteFile(archfile, af.Bytes(), 0755); err != nil {
				return err
			}
		} else {
			sf.DeclareVar(archives...)
		}
	}

	// Export variables
	if len(exports) > 0 {
		sf.DeclareVar(exports...)
	}

	return ioutil.WriteFile(config.Embed.Path, sf.Bytes(), 0755)

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

func (e *EmbedService) compressWithGzip(uncompressed []byte) ([]byte, error) {
	var (
		out bytes.Buffer
		com = gzip.NewWriter(&out)
	)
	if _, err := com.Write(uncompressed); err != nil {
		return nil, err
	}
	if err := com.Flush(); err != nil {
		return nil, err
	}
	if err := com.Close(); err != nil {
		return nil, err
	}
	return out.Bytes(), nil
}

func (e *EmbedService) formatBytes(data []byte) string {
	if len(data) <= 0 {
		return "make([]byte, 0)"
	}
	buf := bytes.NewBuffer([]byte{'0', 'x', hextable[data[0]>>4], hextable[data[0]&0x0F]})
	for i := 1; i < len(data); i++ {
		buf.Write([]byte{',', ' ', '0', 'x', hextable[data[i]>>4], hextable[data[i]&0x0F]})
	}
	return fmt.Sprintf("[]byte{%s}", buf.String())
}
