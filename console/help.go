package console

import (
	"fmt"
	"os"
	"sort"
	"strings"
	"text/template"
)

var helpTemplate = `{{if len .Description | ne 0}}{{.Description}}{{end}}
{{if .ShowCommand}}
Current Command: 
  {{.Name}} {{.Curr}} {{.Args}} {{range $ex := .FlagsEx}} {{$ex}}{{end}}

Flags:
{{range $ex := .Flags}}  {{$ex}}
{{end}}
Examples:
{{range $ex := .Examples}}  {{$ex}}
{{end}}
_____________________________________________________{{end}}
{{if len .Next | ne 0}}
Usage: 
  {{.Name}} {{.Curr}} [command] [args]

Available Commands:
{{range $ex := .Next}}  {{$ex}}
{{end}}{{end}}
_____________________________________________________
Use flag --help for more information about a command.

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

func help(tool string, desc string, c CommandGetter, args []string) {
	model := &helpModel{
		ShowCommand: c != nil && c.Call() != nil,
		Name:        tool,
		Description: desc,

		Curr: strings.Join(args, " "),
		Next: func() (out []string) {
			if c == nil {
				return
			}
			var max int
			next := c.List()
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

	if c != nil {
		model.Examples = func() (out []string) {
			for _, v := range c.Examples() {
				out = append(out, tool+" "+v)
			}
			return
		}()
		model.Args = "[arg]"
		model.Flags = func() (out []string) {
			max := 0
			c.Flags().Info(func(r bool, n string, v interface{}, u string) {
				if len(n) > max {
					max = len(n)
				}
			})
			c.Flags().Info(func(r bool, n string, v interface{}, u string) {
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
			c.Flags().Info(func(r bool, n string, v interface{}, u string) {
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
