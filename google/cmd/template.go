package main

import (
	"io"
	"text/template"

	"github.com/pkg/errors"
)

const (
	// packageTmpl it's the package definition
	packageTmpl = `
	package google
	// Code generated by 'go generate'; DO NOT EDIT
	import (
		"context"

		"github.com/pkg/errors"

		"google.golang.org/api/compute/v1"
	)
	`

	// functionTmpl it's the implementation of a reader function
	functionTmpl = `
	// List{{ .Resource }}s will returns a list of {{ .Resource }} within a project {{ if .Zone }}and a zone {{ end }}
	func (r *GCPReader) List{{ .Resource}}s(ctx context.Context, filter string) ({{ if .Zone }}map[string]{{end}}[]compute.{{ .Resource }}, error) {
		service := compute.New{{ .Resource}}sService(r.compute)
		{{ if .Zone }}
		list := make(map[string][]compute.{{ .Resource }})
		zones, err := r.getZones()
		if err != nil {
			return nil, errors.Wrap(err, "unable to get zones in region")
		}
		for _, zone := range zones {
		{{ end }}
		resources := make([]compute.{{ .Resource }}, 0)
		{{ if .Zone }}
		if err := service.List(r.project, zone).
		{{ else }}
		if err := service.List(r.project).
		{{ end }}
			Filter(filter).
			MaxResults(int64(r.maxResults)).
			Pages(ctx, func(list *compute.{{ .Resource }}List) error {
				for _, res := range list.Items {
					resources = append(resources, *res)
				}
				return nil
			}); err != nil {
			return nil, errors.Wrap(err, "unable to list compute {{ .Resource }} from google APIs")
		}
		{{ if .Zone }}
		list[zone] = resources
		}
		return list, nil
		{{ else }}
		return resources, nil
		{{ end }}
	}
	`
)

var (
	fnTmpl  *template.Template
	pkgTmpl *template.Template
)

func init() {
	var err error
	fnTmpl, err = template.New("template").Parse(functionTmpl)
	if err != nil {
		panic(err)
	}
	pkgTmpl, err = template.New("template").Parse(packageTmpl)
	if err != nil {
		panic(err)
	}
}

// Function is the definition of one of the functions
type Function struct {
	// Resource is the Google name of the entity, like
	// Firewall, Instance, etc.
	// https://godoc.org/google.golang.org/api/compute/v1
	Resource string

	// Zone is used to determine whether the resource is located within google zones or not
	Zone bool
}

// Execute uses the fnTmpl to interpolate f
// and write the result to w
func (f Function) Execute(w io.Writer) error {
	if err := fnTmpl.Execute(w, f); err != nil {
		return errors.Wrapf(err, "failed to Execute with Function %+v", f)
	}
	return nil
}