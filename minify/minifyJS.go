package minify

const (
	minifyJSDirPath       = "static/scripts/"
	minifyJSInputFileExt  = ".js"
	minifyJSOutputFileExt = ".min" + minifyJSInputFileExt
)

func init() {
	walkDirAndModifyFiles(
		minifyJSDirPath,
		minifyJSInputFileExt,
		minifyJSOutputFileExt,
		func(contentInput []byte) (contentOutput []byte, err error) {
			contentOutput, err = minifyEngine.Bytes("text/javascript", contentInput)
			if err != nil {
				return
			}

			return
		},
	)
}
