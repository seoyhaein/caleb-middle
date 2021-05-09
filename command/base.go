package command

import (
	"bytes"
	"flag"
	"fmt"
	"github.com/mitchellh/cli"
)

type flagSet struct {
	*flag.FlagSet
}

func (fs *flagSet) help() string {
	var buf bytes.Buffer
	fs.VisitAll(func(f *flag.Flag) {
		fmt.Fprintf(&buf, "  -%s\n", f.Name)
		fmt.Fprintf(&buf, "      %s", f.Usage)
		if f.DefValue != "" {
			fmt.Fprintf(&buf, " (default: %s)", f.DefValue)
		}
		fmt.Fprint(&buf, "\n\n")
	})
	return buf.String()
}

type BaseCommand struct {
	Ui cli.Ui
}

// Printf outputs formatted message.
func (c *BaseCommand) Printf(format string, a ...interface{}) {
	c.Ui.Output(fmt.Sprintf(format, a...))
}
