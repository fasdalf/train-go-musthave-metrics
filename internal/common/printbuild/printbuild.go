// Package printbuild - Print Build
package printbuild

import (
	"os"
	"text/template"
)

type Data struct {
	BuildVersion string
	BuildDate    string
	BuildCommit  string
}

const Template = `Build version: {{ if .BuildVersion }}{{ .BuildVersion }}{{ else }}N/A{{ end }}
Build date: {{ if .BuildDate }}{{ .BuildDate }}{{ else }}N/A{{ end }}
Build commit: {{ if .BuildCommit }}{{ .BuildCommit }}{{ else }}N/A{{ end }} 
`

func (d *Data) Print() {
	t := template.Must(template.New("-").Parse(Template))

	_ = t.Execute(os.Stdout, d)
}
