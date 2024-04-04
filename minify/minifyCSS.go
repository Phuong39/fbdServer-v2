package minify

const (
	minifyCSSDirPath       = "static/styles/"
	minifyCSSInputFileExt  = ".css"
	minifyCSSOutputFileExt = ".min" + minifyCSSInputFileExt
)

func init() {
	walkDirAndModifyFiles(
		minifyCSSDirPath,
		minifyCSSInputFileExt,
		minifyCSSOutputFileExt,
		func(contentInput []byte) (contentOutput []byte, err error) {
			contentOutput, err = minifyEngine.Bytes("text/css", contentInput)
			if err != nil {
				return
			}

			return
		},
	)
}
