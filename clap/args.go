package clap

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"strings"
)

var (
	ErrUnexpectedArgument   = errors.New("unexpected argument")
	ErrMissingArgumentValue = errors.New("missing argument value")
	ErrIgnoredArgument      = errors.New("ignored argument")
	ErrMandatoryArgument    = errors.New("mandatory argument")
	ErrDuplicatedArgument   = errors.New("duplicated argument")
)

func consumeArguments(start int, args []string, count int) (int, []string) {
	var values []string
	for ; start < count; start++ {
		if strings.HasPrefix(args[start], "-") {
			break
		}
		values = append(values, args[start])
	}
	return start - 1, values
}

func stringsToInts(strs []string) ([]int, error) {
	var ints []int
	for _, str := range strs {
		var val int64
		var err error
		if val, err = strconv.ParseInt(str, 10, 64); err != nil {
			return nil, fmt.Errorf("%w (got '%s', expected integer)", err, str)
		}
		ints = append(ints, int(val))
	}
	return ints, nil
}

func argsToFields(args []string, fieldDescs map[string]*fieldDescription, cfg any) (*Results, error) {
	results := &Results{}
	reflectValue := reflect.ValueOf(cfg).Elem()
	for i := 0; i < len(args); i++ {
		arg := args[i]
		if desc, ok := fieldDescs[arg]; ok {
			if desc.Found {
				results.Duplicated = append(results.Duplicated, arg)
				return results, fmt.Errorf("argument '%s': %w (duplicated argument)", arg, ErrDuplicatedArgument)
			}
			desc.Found = true
			field := reflectValue.Field(desc.Field)
			if !field.CanSet() {
				continue
			}
			switch desc.Type.Kind() {
			case reflect.Int, reflect.String, reflect.Float64:
				i++
				if i >= len(args) || strings.HasPrefix(args[i], "-") {
					results.Missing = append(results.Missing, arg)
					return results, fmt.Errorf("argument '%s': %w (missing argument)", arg, ErrMissingArgumentValue)
				}
				desc.Args = append(desc.Args, args[i])
			case reflect.Bool:
				_, values := consumeArguments(i+1, args, len(args))
				if len(values) == 0 {
					values = append(values, "true")
				}
				desc.Args = append(desc.Args, values...)
			case reflect.Slice, reflect.Array:
				var values []string
				count := len(args)
				if desc.Type.Kind() == reflect.Array {
					count = i + 1 + desc.Type.Len()
				}
				i, values = consumeArguments(i+1, args, count)
				if len(values) == 0 {
					results.Missing = append(results.Missing, arg)
					return results, fmt.Errorf("argument '%s': %w (missing argument)", arg, ErrMissingArgumentValue)
				}
				desc.Args = append(desc.Args, values...)
			}
		} else {
			found := false
			for j := i; j < len(args); j++ {
				if strings.HasPrefix(args[j], "-") {
					results.Ignored = append(results.Ignored, arg)
					// does not generate an error
					found = true
					break
				}
			}
			if !found {
				if desc, ok := fieldDescs[trailing]; ok {
					field := reflectValue.Field(desc.Field)
					if field.CanSet() {
						var values []string
						for j := i; j < len(args); j++ {
							values = append(values, args[j])
						}
						desc.Args = append(desc.Args, values...)
						break
					}
				}
			}
		}
	}
	for _, desc := range fieldDescs {
		if !desc.Found && desc.Mandatory {
			name := desc.LongName
			if name == "" {
				name = desc.ShortName
			}
			results.Mandatory = append(results.Mandatory, name)
			return results, fmt.Errorf("mandatory argument '%s' not found: %w", name, ErrMandatoryArgument)
		}
	}
	return results, nil
}

func fillStruct(args []string, fieldDescs map[string]*fieldDescription, cfg any) (*Results, error) {
	results, err := argsToFields(args, fieldDescs, cfg)
	if err != nil {
		return results, err
	}
	reflectValue := reflect.ValueOf(cfg).Elem()
	for name, desc := range fieldDescs {
		field := reflectValue.Field(desc.Field)
		if !field.CanSet() || len(desc.Args) == 0 || desc.Visited {
			continue
		}
		desc.Visited = true
		switch desc.Type.Kind() {
		case reflect.String:
			field.SetString(desc.Args[0])
		case reflect.Int:
			val, err := strconv.ParseInt(desc.Args[0], 10, 64)
			if err != nil {
				results.Unexpected = append(results.Unexpected, name)
				return results, fmt.Errorf("argument '%s': %w (got '%s', expected integer)", name,
					ErrUnexpectedArgument, desc.Args[0])
			}
			field.SetInt(val)
		case reflect.Float64:
			val, err := strconv.ParseFloat(desc.Args[0], 64)
			if err != nil {
				results.Unexpected = append(results.Unexpected, name)
				return results, fmt.Errorf("argument '%s': %w (got '%s', expected float)", name,
					ErrUnexpectedArgument, desc.Args[0])
			}
			field.SetFloat(val)
		case reflect.Bool:
			field.SetBool(true)
			if len(desc.Args) >= 1 {
				value := strings.ToLower(desc.Args[0])
				if value == "0" || value == "false" || value == "no" {
					field.SetBool(false)
				}
			}
		case reflect.Slice:
			if desc.Type.Elem().Kind() == reflect.String {
				field.Set(reflect.ValueOf(desc.Args))
			} else if desc.Type.Elem().Kind() == reflect.Int {
				ints, err := stringsToInts(desc.Args)
				if err != nil {
					results.Unexpected = append(results.Unexpected, name)
					return results, fmt.Errorf("argument '%s': %w", name, err)
				}
				field.Set(reflect.ValueOf(ints))
			}
		case reflect.Array:
			if desc.Type.Elem().Kind() == reflect.String {
				arrayType := reflect.ArrayOf(desc.Type.Len(), desc.Type.Elem())
				v := reflect.New(arrayType).Elem()
				reflect.Copy(v, reflect.ValueOf(desc.Args))
				field.Set(v)
			} else if desc.Type.Elem().Kind() == reflect.Int {
				ints, err := stringsToInts(desc.Args)
				if err != nil {
					results.Unexpected = append(results.Unexpected, name)
					return results, fmt.Errorf("argument '%s': %w", name, err)
				}
				arrayType := reflect.ArrayOf(desc.Type.Len(), desc.Type.Elem())
				v := reflect.New(arrayType).Elem()
				reflect.Copy(v, reflect.ValueOf(ints))
				field.Set(v)
			}
		}
	}
	return results, nil
}
