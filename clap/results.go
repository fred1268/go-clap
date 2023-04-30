package clap

/*
Represents the results of the command line parsing.

Unexpected: contains unexpected parameters type (for instance expecting an integer
and getting a string)

Missing: contains non-boolean parameters with missing value(s)

Ignored: contains parameters presents on the command line, but not recognized by the
program

Mandatory: contains mandatory (non optional) parameters that are not present in the
command line

Duplicated: contains parameters that are duplicated on the command line
*/
type Results struct {
	Unexpected []string
	Missing    []string
	Ignored    []string
	Mandatory  []string
	Duplicated []string
}

/*
Returns true if the parser returns any error. Errors are:

- Unexpected parameter

- Missing int or string parameter

- Mandatory parameters not present

- Duplicated parameters
*/
func (r *Results) HasErrors() bool {
	return len(r.Unexpected) != 0 || len(r.Missing) != 0 || len(r.Mandatory) != 0 || len(r.Duplicated) != 0
}

/*
Returns true if the parser returns any Warning. Warnings are:

- Ignored parameters
*/
func (r *Results) HasWarnings() bool {
	return len(r.Ignored) != 0
}
