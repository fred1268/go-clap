package clap

import "reflect"

type fieldDescription struct {
	Field     int
	ShortName string
	LongName  string
	Type      reflect.Type
	Args      []string
	Mandatory bool
	Found     bool
	Visited   bool
}
