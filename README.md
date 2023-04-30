# clap :clap:

The lightweight, non-intrusive **Command Line Argument Parser**

---

## Why?

The most famous command line parser is undoubtly Cobra. Cobra is very nice and has tons of features,
however, I was not quite happy with it. First, it does way too many things for my liking, and also,
it is pretty intrusive, i.e. you need to build your code around Cobra, rather than having Cobra
silently collaborating with your code.

---

### Enters clap :clap:

clap is a **non intrusive** Command Line Argument Parser, that will use struct tags to understand
what you want to parse. You don't have to implement interface or use a CLI to scaffold your project:
you just call `clap.Parse(args, &YourConfig)` and you are done. You are notetheless responsible for
handling potential commands and subcommands, but clap will fill up your CLI configuration struct
with the values passed on the command line for those commands / subcommands.

---

### Benefits of clap :clap:

- lightweight
- non intrusive
- obvious defaults
- no dependencies
- easy to configure
- no CLI nor scaffolding

---

## Installation

```go
go get github.com/fred1268/go-clap
```

---

## Quick start

Very easy to start with:

1. declare a struct containing your configuration
2. add struct tags as hints for clap
3. call `clap.Parse(args, &config)`

```go
    type config struct {
    	Cookie      string   `clap:"--cookie"`
    	HTTPOnly    bool     `clap:"--httpOnly"`
    	Secure      bool     `clap:"--secure"`
    	Origins     [4]string `clap:"--origins,-O,mandatory"`
    	Port        int      `clap:",-P,mandatory"`
    	ConfigFiles []string `clap:"trailing"`
    }
```

A clap struct tag has the following structure:

```go
    Name        Type    `clap:"longName[,shortName][,mandatory]"`
```

longName is a... well... long name, like --recursive or --credentials

shortName is a single letter name, like -R or -c

mandatory can be added to make the non-optional parameters

In your main, just make a call to clap.Parse():

```go
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
```

Assuming the command line looks like:

```shell
    -P 8080 --cookie clapcookie --httpOnly --origins http://localhost:5137 \
    https://localhost:5173 http://localhost:3000 https://localhost:3000 \
    config-db.json config-log.json
```

You will get the following struct:

```go
    config{
    	Cookie:   "clapcookie",
    	HTTPOnly: true,
    	Secure:   true, // comes from your default (not in the command line)
    	Origins: []string{
    		"http://localhost:5137", "https://localhost:5173",
    		"http://localhost:3000", "https://localhost:3000",
    	},
    	Port:        8080,
    	// trailing parameters
    	ConfigFiles: []string{"config-db.json", "config-log.json"},
    }
```

> Note that it is important to use arrays rather than slices when you can,
> since arrays will consume, at maximum the requested number of
> parameters, whereas slices will consume all possible parameters.
> Thus, a slice parameter should not be used just before the trailing
> otherwise it will consume the trailing parameters.

---

## Supported parameter types

The following parameter types are supported by clap:

- bool: `--param`
- string: `--param test`
- int: `--param 10`
- float: `--param 12.3`
- string array of any size (here 3): `--param a b c`
- int array of any size (here 2): `--param 80 443`
- string slice: `--param a b c`
- int slice: `--param 80 443`

---

## License & contribution

clap is licensed under the MIT license (see [LICENSE](LICENSE)).

Issues or PR are welcome.
