package field

// Prefixer adds Prefix to Key-s of all fields.
//
// TODO: redesign this, see: https://github.com/facebookincubator/go-belt/issues/6
type Prefixer struct {
	Prefix string
	Fields AbstractFields
}

var _ AbstractFields = (*Prefixer)(nil)

// Prefix adds Prefix to Key-s of all fields.
func Prefix(prefix string, fields AbstractFields) Prefixer {
	return Prefixer{
		Prefix: prefix,
		Fields: fields,
	}
}

// ForEachField implements AbstractFields.
func (p Prefixer) ForEachField(callback func(f *Field) bool) bool {
	var f Field
	return p.Fields.ForEachField(func(in *Field) bool {
		f = *in
		f.Key = p.Prefix + in.Key
		return callback(&f)
	})
}

// Len implements AbstractFields.
func (p Prefixer) Len() int {
	return p.Fields.Len()
}
