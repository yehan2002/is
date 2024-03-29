package is

import (
	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
)

// Option an options to be applied to the test.
type Option func(*options)

// CmpAllUnexported enables comparing all unexported fields using [Is.Equal].
// Using this options is not recommended since this will compare the unexported fields of structs from other
// packages. Use [CmpUnexported] instead.
func CmpAllUnexported() Option {
	return func(o *options) { o.cmpAllUnexported = true }
}

// CmpUnexported allows comparing unexported fields of the given struct types.
func CmpUnexported(types ...interface{}) Option {
	return func(o *options) {
		if o.cmpUnexportedMap == nil {
			o.cmpUnexportedMap = make(map[reflect.Type]struct{})
		}

		for _, i := range types {
			o.cmpUnexportedMap[reflect.TypeOf(i)] = struct{}{}
		}
		o.cmpUnexported = append(o.cmpUnexported, types...)
	}
}

// EquateNaN sets if two NaN values are considered equal
// Default: true
func EquateNaN(eq bool) Option {
	return func(o *options) { o.equateNaN = eq }
}

// EquateErrors sets if the equality of errors is checked using [errors.Is].
// Default: false
func EquateErrors(eq bool) Option {
	return func(o *options) { o.equateErrors = eq }
}

// EquateEmpty sets if empty maps/slices and nil maps/slices are considered to be equal.
// Default: true
func EquateEmpty(eq bool) Option {
	return func(o *options) { o.equateEmpty = eq }
}

// CmpOpt adds the given [cmp.Option]s to the list of options passed to [cmp.Diff].
func CmpOpt(opts ...cmp.Option) Option {
	return func(o *options) { o.userOpts = append(o.userOpts, opts...) }
}

type options struct {
	cmpUnexported    []interface{}
	cmpUnexportedMap map[reflect.Type]struct{}

	cmpAllUnexported bool

	equateNaN    bool
	equateErrors bool
	equateEmpty  bool

	userOpts []cmp.Option

	cmpOpts []cmp.Option
}

func newOptions(opts []Option) *options {
	o := &options{equateEmpty: true, equateNaN: true}
	for _, opt := range opts {
		if opt != nil {
			opt(o)
		}
	}
	return o
}

func (o *options) CmpOpts() []cmp.Option {
	if o.cmpOpts == nil {
		o.cmpOpts = append(o.cmpOpts, cmp.FilterPath(func(p cmp.Path) bool {
			sf, ok := p.Index(-1).(cmp.StructField)
			if !ok {
				return false
			}

			parent := p.Index(-2)
			field := parent.Type().Field(sf.Index())
			ignoreTag := field.Tag != "" && (field.Tag.Get("deep") == "-" || field.Tag.Get("cmp") == "-")

			isExported := field.PkgPath == ""

			isUnexported := !isExported
			if isUnexported {
				if _, ok := o.cmpUnexportedMap[parent.Type()]; ok || o.cmpAllUnexported {
					isUnexported = false
				}
			}

			return isUnexported || ignoreTag
		}, cmp.Ignore()))

		if o.cmpAllUnexported {
			o.cmpOpts = append(o.cmpOpts, cmp.Exporter(func(t reflect.Type) bool {
				return o.cmpAllUnexported
			}))
		} else if len(o.cmpUnexported) != 0 {
			o.cmpOpts = append(o.cmpOpts, cmp.AllowUnexported(o.cmpUnexported...))
		}

		if o.equateEmpty {
			o.cmpOpts = append(o.cmpOpts, cmpopts.EquateEmpty())
		}

		if o.equateErrors {
			o.cmpOpts = append(o.cmpOpts, cmpopts.EquateErrors())
		}

		if o.equateNaN {
			o.cmpOpts = append(o.cmpOpts, cmpopts.EquateNaNs())
		}

		o.cmpOpts = append(o.cmpOpts, o.userOpts...)
	}

	cmpopts.EquateErrors()
	return o.cmpOpts
}
