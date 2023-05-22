/*
Package clap provides a lightweight, non intrusive library to
parse the command line arguments.

Very easy to start with:

 1. declare a struct containing your configuration

 2. add struct tags as hints for clap

 3. call `clap.Parse(args, &config)`

```

	type config struct {
		Cookie      string   `clap:"--cookie"`
		HTTPOnly    bool     `clap:"--httpOnly"`
		Secure      bool     `clap:"--secure"`
		Origins     [4]string `clap:"--origins,-O,mandatory"`
		Port        int      `clap:",-P,mandatory"`
		ConfigFiles []string `clap:"trailing"`
	}

```

Please note that this automatically generates a `--no-secure` and `--no-httpOnly`
flags that you can use on the command line to set the corresponding booleans
to the `false` value. This is useful when you want to give a boolean a `true`
default value.

A clap struct tag has the following structure:

	Name        Type    `clap:"longName[,shortName][,mandatory]"`

longName is a... well... long name, like `--recursive` or `--credentials`

shortName is a single letter name, like `-R` or `-c`

mandatory can be added to make the non-optional parameters

In your main, just make a call to `clap.Parse()`:

	func main() {
	    var err error
	    var results *clap.Results
	    // define your defaults
		cfg := &config{Secure: true}
	    // note you may want to skip the first few
	    // parameters (like command and subcommand)
	    // by passing args[2:] instead of args
	    if results, err = clap.Parse(args, cfg); err != nil {
	        // results contains a list of arguments in error
	        // can be used for user friendly error handling
	        return err
	    }
	    // results contains the list of arguments being ignored
	    // can be used for user friendly error handling
	}

The following parameter types are supported by clap:

	bool: `--param` or `--no-param`
	string: `--param test`
	int: `--param 10`
	float: `--param 12.3`
	string array of any size (here 3): `--param a b c`
	int array of any size (here 2): `--param 80 443`
	string slice: `--param a b c`
	int slice: `--param 80 443`
*/
package clap

import "reflect"

/*
Parses the command line arguments into the given struct
Struct tag looks like one of the following:

	`clap:"longname[,shortname][,mandatory]"`

longname: represents the long name of the command line
parameter, without the --. For instance:

	recursive for --recursive

shortname: represents a one letter command line argument
without the -. For instance:

	R for -R

Examples:

	`clap:"recursive"`
	`clap:",R,mandatory"`
	`clap:"recursive,,mandatory"`
	`clap:"recursive,R,mandatory"`

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
