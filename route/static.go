package route

import (
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"regexp"

	"github.com/theTardigrade/fbdServer-v2/options"
)

const (
	staticFilepathSeparator  = string(filepath.Separator)
	staticPrevDirReplacement = staticFilepathSeparator + "." + staticFilepathSeparator
)

var (
	staticPrevDirRegexp = regexp.MustCompile(`[/\\][.]{2,}[/\\]`)
)

func staticMimeType(filePath string) (mimeType string) {
	switch ext := filepath.Ext(filePath); ext[1:] {
	case "js":
		mimeType = "application/javascript"
	default:
		mimeType = mime.TypeByExtension(ext)

		if mimeType == "" {
			mimeType = "application/octet-stream"
		}
	}

	return
}

var (
	staticGetHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var success bool

		defer func() {
			if err := recover(); err != nil {
				serverErrorHandler(w, r)
			} else if !success {
				notFoundHandler(w, r)
			}
		}()

		localPath := filepath.Join(
			".",
			staticPrevDirRegexp.ReplaceAllLiteralString(r.URL.Path[1:], staticPrevDirReplacement),
		)

		fileInfo, err := os.Stat(localPath)
		if err != nil {
			if !os.IsNotExist(err) {
				panic(err)
			}

			return
		}
		if !fileInfo.Mode().IsRegular() {
			return
		}
		// fileModTime := fileInfo.ModTime()

		fileContents, err := os.ReadFile(localPath)
		if err != nil {
			panic(err)
		}

		header := w.Header()

		header.Set("Content-Type", staticMimeType(localPath))

		// useEtag := true

		// if fileInfo.Size() < staticEtagFileSizeMin {
		// 	useEtag = false
		// }

		success = true

		// if useEtag {
		// 	etag := datum.etag

		// 	if r.Header.Get("If-None-Match") == etag {
		// 		w.WriteHeader(http.StatusNotModified)
		// 		w.Write([]byte{})
		// 		return
		// 	}

		// 	header.Set("Etag", etag)
		// }

		w.WriteHeader(http.StatusOK)
		w.Write(fileContents)
	})
)

const (
	staticPath = "/static/*"
)

func init() {
	options.Options.Routes.Get[staticPath] = staticGetHandler
}
