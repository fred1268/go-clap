package clap

import "reflect"

/*
Parses the command line arguments into the given struct
Struct tag looks like one of the following:

	`clap:"longname[, shortname][, optional]"`

longname: represents the long name of the command line
parameter, without the --. For instance:

	recursive for --recursive

shortname: represents a one letter command line argument
without the -. For instance:

	R for -R

Examples:

	`clap:"recursive"`
	`clap:",R,optional"`
	`clap:"recursive,,optional"`
	`clap:"recursive,R,optional"`

There is a special longname that you can use to retrieve
all trailing parameters on your command line: trailing.
It is used like this:

	`clap:"trailing"`

Supported field types:

	int
	bool
	string
	[]int
	[]string
*/
func Parse[T any](args []string, cfg *T) (*Results, error) {
	var err error
	var results *Results
	var fieldDescs map[string]*fieldDescription
	fieldDescs, err = computeFieldDescriptions(reflect.TypeOf(*cfg))
	if err != nil {
		return nil, err
	}
	if results, err = fillStruct(args, fieldDescs, cfg); err != nil {
		return results, err
	}
	return results, nil
}
