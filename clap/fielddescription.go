package clap

import "reflect"

type fieldDescription struct {
	Field     int
	ShortName string
	LongName  string
	Mandatory bool
	Found     bool
	Type      reflect.Type
	Args      []string
}
