package memory

import (
	"github.com/poonman/entry-task/dora/misc/lion/loader"
	"github.com/poonman/entry-task/dora/misc/lion/reader"
	"github.com/poonman/entry-task/dora/misc/lion/source"
)

// WithSource appends a source to list of sources
func WithSource(s source.Source) loader.Option {
	return func(o *loader.Options) {
		o.Source = append(o.Source, s)
	}
}

// WithReader sets the config reader
func WithReader(r reader.Reader) loader.Option {
	return func(o *loader.Options) {
		o.Reader = r
	}
}
