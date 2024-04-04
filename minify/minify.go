package minify

import (
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/tdewolff/minify"
	"github.com/tdewolff/minify/css"
	"github.com/tdewolff/minify/js"
	globalFilepath "github.com/theTardigrade/golang-globalFilepath"
)

var (
	minifyEngine *minify.M
)

func init() {
	globalFilepath.Init("..")

	minifyEngine = minify.New()

	minifyEngine.AddFunc("text/css", css.Minify)
	minifyEngine.AddFunc("text/javascript", js.Minify)
}

type WalkDirAndModifyFilesCallback func(contentInput []byte) (contentOutput []byte, err error)

func walkDirAndModifyFiles(dirPath, inputFileExt, outputFileExt string, callback WalkDirAndModifyFilesCallback) (err error) {
	err = filepath.WalkDir(
		globalFilepath.Join(dirPath),
		func(path string, d fs.DirEntry, err error) (err2 error) {
			if err != nil {
				err2 = err
				return
			}

			if d.IsDir() {
				return
			}

			name := strings.TrimSpace(d.Name())
			ext := filepath.Ext(name)

			if ext != inputFileExt || strings.HasSuffix(name, outputFileExt) {
				return
			}

			nameWithoutExt := name[:len(name)-len(ext)]
			outputName := nameWithoutExt + outputFileExt
			pathDir := filepath.Dir(path)
			outputPath := filepath.Join(pathDir, outputName)

			content, err2 := os.ReadFile(path)
			if err2 != nil {
				return
			}

			content, err2 = callback(content)
			if err2 != nil {
				return
			}

			err2 = os.WriteFile(outputPath, content, os.ModePerm)
			if err2 != nil {
				return
			}

			return
		},
	)
	if err != nil {
		return
	}

	return
}
