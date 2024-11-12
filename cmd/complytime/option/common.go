package option

import (
	"io"

	"github.com/spf13/pflag"
)

type Common struct {
	Verbose bool
	IOStreams
}

type IOStreams struct {
	// In think, os.Stdin
	In io.Reader
	// Out think, os.Stdout
	Out io.Writer
	// ErrOut think, os.Stderr
	ErrOut io.Writer
}

func (o *Common) BindFlags(fs *pflag.FlagSet) {
	fs.BoolVarP(&o.Verbose, "verbose", "v", false, "verbose output")
}
