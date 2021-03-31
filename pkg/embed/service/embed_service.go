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
	"path/filepath"

	"github.com/thecxx/embed/pkg/embed/asset/embed"

	"github.com/thecxx/embed/pkg/embed/asset/options"

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

func (e *EmbedService) Init(ctx context.Context) error {
	stat, err := os.Stat(options.InitCmd.File)
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
	return ioutil.WriteFile(options.InitCmd.File, embed.EmbedYaml.Bytes(), 0755)
}

func (e *EmbedService) Build(ctx context.Context) error {

	var (
		sf         = pack.NewSourceFile(config.Embed.Package)
		af         = pack.NewSourceFile(config.Embed.Package)
		exports    = make([]pack.Variable, 0)
		archives   = make([]pack.Variable, 0)
		compressor = ""
	)

	sf.Import("io")
	sf.Import("bytes")
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

	datas := make(map[string][]byte)

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
			buf = fmt.Sprintf("%sRead(_%s)", compressor, sign)
		} else {
			buf = fmt.Sprintf("directRead(_%s)", sign)
		}

		exports = append(exports, pack.Variable{
			Name:    f.Name,
			Comment: f.Comment,
			Assign:  fmt.Sprintf("item{\"%s\", %s}", f.File, buf),
		})
	}

	// Clean
	e.remove(config.Embed.Path, config.Embed.Archive)

	// Save all data
	if len(archives) > 0 {
		if len(config.Embed.Archive) > 0 {
			for _, arch := range archives {
				af.DeclareVar(arch)
			}
			if err := e.save(config.Embed.Archive, af, 0755); err != nil {
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

	return e.save(config.Embed.Path, sf, 0755)

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

func (e *EmbedService) remove(files ...string) error {
	for _, f := range files {
		_, err := os.Stat(f)
		if err != nil {
			if os.IsNotExist(err) {
				continue
			}
			return err
		}
		if err := os.Remove(f); err != nil {
			return err
		}
	}
	return nil
}

func (e *EmbedService) save(file string, src *pack.SourceFile, perm os.FileMode) error {
	dir := filepath.Dir(file)
	// Create the package directory
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	// Create file
	return ioutil.WriteFile(file, src.Bytes(), perm)
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
