package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path"
	"path/filepath"
)

type StringBuffer struct {
	buf *bytes.Buffer
}

func (sb *StringBuffer) String() string {
	result := "[]byte{"
	for i, b := range sb.buf.Bytes() {
		if i != sb.buf.Len()-1 {
			result += fmt.Sprintf("%d, ", b)
		} else {
			result += fmt.Sprintf("%d", b)
		}
	}
	result += "}"
	return result
}

func main() {
	println("generating assets:")
	filesDir := path.Join("assets", "files")
	dir := http.Dir(filesDir)

	root, err := dir.Open("/")
	must(err)

	info, err := root.Stat()
	must(err)

	if !info.IsDir() {
		must(errors.New("given path is not a directory"))
	}

	assets, err := root.Readdir(-1)
	must(err)

	generated, err := os.Create(path.Join("assets", "generated.go"))
	must(err)

	io.WriteString(generated, fmt.Sprintf(`// Code generated by go generate; DO NOT EDIT.

package assets

func init() {
	Assets = New()
`))

	for _, asset := range assets {
		raw, err := ioutil.ReadFile(path.Join(filesDir, asset.Name()))
		must(err)

		ct := mime.TypeByExtension(filepath.Ext(asset.Name()))

		println("\t" + asset.Name())

		io.WriteString(generated, fmt.Sprintf(`	Assets.files["%s"] = &AssetFile{bytes: %v, ct: "%s", size: "%d"}
`, asset.Name(), &StringBuffer{bytes.NewBuffer(raw)}, ct, len(raw)))
	}

	io.WriteString(generated, "}\n")
}

func must(err error) {
	if err != nil {
		println(err.Error())
		os.Exit(1)
	}
}