package console

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
)

var helpTemplate = `{{if len .Description | ne 0}}{{.Description}}

{{end}}Usage: {{if .ShowCommand}}
  {{.Name}} {{.Curr}} {{.Args}} {{range $ex := .FlagsEx}} {{$ex}}{{end}}

Flags:
{{range $ex := .Flags}}  {{$ex}}
{{end}}

Examples:
{{range $ex := .Examples}}  {{$ex}}
{{end}}
{{else}}
  {{.Name}} {{.Curr}} [command] [args]

Available Commands:
{{range $ex := .Next}}  {{$ex}}
{{end}}
_____________________________________________________
Use flag --help for more information about a command.
{{end}}
`

type helpModel struct {
	Name        string
	Description string
	ShowCommand bool

	Args     string
	Examples []string
	FlagsEx  []string
	Flags    []string

	Curr string
	Next []string
}

func help(tool string, desc string, next []CommandGetter, command CommandGetter, cur []string) {
	model := &helpModel{
		ShowCommand: len(next) == 0,
		Name:        tool,
		Description: desc,

		Curr: strings.Join(cur, " "),
		Next: func() (out []string) {
			var max int
			for _, v := range next {
				if max < len(v.Name()) {
					max = len(v.Name())
				}
			}
			sort.Slice(next, func(i, j int) bool {
				return next[i].Name() < next[j].Name()
			})
			for _, v := range next {
				out = append(out, v.Name()+strings.Repeat(" ", max-len(v.Name()))+"    "+v.Description())
			}
			return
		}(),
	}

	if command != nil {
		model.Examples = func() (out []string) {
			for _, v := range command.Examples() {
				out = append(out, tool+" "+v)
			}
			return
		}()
		model.Args = strings.TrimSpace(strings.Repeat("[arg] ", command.ArgCount()))
		model.Flags = func() (out []string) {
			max := 0
			command.Flags().Info(func(r bool, n string, v interface{}, u string) {
				if len(n) > max {
					max = len(n)
				}
			})
			command.Flags().Info(func(r bool, n string, v interface{}, u string) {
				ex, i := "", 1
				if !r {
					ex = fmt.Sprintf("(default: %+v)", v)
				}
				if len(n) > 1 {
					i = 2
				}
				out = append(out, fmt.Sprintf(
					"%s%s%s    %s %s",
					strings.Repeat("-", i), n, strings.Repeat(" ", max-len(n)), u, ex))
			})
			return
		}()
		model.FlagsEx = func() (out []string) {
			command.Flags().Info(func(r bool, n string, v interface{}, u string) {
				i, ex := 1, ""
				if len(n) > 1 {
					i = 2
				}
				switch v.(type) {
				case bool:
				default:
					ex = fmt.Sprintf("=%+v", v)
				}
				out = append(out, fmt.Sprintf(
					"%s%s%s",
					strings.Repeat("-", i), n, ex))
			})
			return
		}()
	}

	if err := template.Must(template.New("").Parse(helpTemplate)).Execute(os.Stdout, model); err != nil {
		Fatalf(err.Error())
	}
}
