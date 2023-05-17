package clap

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
)

var ErrInvalidTag = errors.New("invalid tag")

const (
	trailing  string = "trailing"
	mandatory string = "mandatory"
)

func getTrailingFieldDescription(tags []string, field reflect.StructField) (*fieldDescription, error) {
	fieldDesc := &fieldDescription{Type: field.Type}
	if len(tags) != 1 {
		return nil, fmt.Errorf("field '%s': %w (got '%s', expected 'trailing')", field.Name,
			ErrInvalidTag, field.Tag.Get("clap"))
	}
	if fieldDesc.Type.Kind() != reflect.String && fieldDesc.Type.Elem().Kind() != reflect.String {
		return nil, fmt.Errorf("field '%s' should be a []string: %w", field.Name, ErrInvalidTag)
	}
	return fieldDesc, nil
}

func getShortNameFieldDescription(tags []string, field reflect.StructField) (*fieldDescription, error) {
	fieldDesc := &fieldDescription{Type: field.Type}
	if len(tags) > 3 || len(tags) < 2 {
		return nil, fmt.Errorf("field '%s': %w (got '%s', expected two or three values)", field.Name,
			ErrInvalidTag, field.Tag.Get("clap"))
	}
	fieldDesc.ShortName = strings.Trim(tags[1], " -")
	if len(fieldDesc.ShortName) != 1 {
		return nil, fmt.Errorf("field '%s': %w (got '%s', expected a single char value)", field.Name,
			ErrInvalidTag, field.Tag.Get("clap"))
	}
	if len(tags) == 3 {
		tag := strings.Trim(tags[2], " ")
		if tag != mandatory {
			return nil, fmt.Errorf("field '%s': %w (got '%s', expected '%s,mandatory')", field.Name,
				ErrInvalidTag, field.Tag.Get("clap"), field.Name)
		}
		fieldDesc.Mandatory = true
	}
	return fieldDesc, nil
}

func getLongNameFieldDescription(tags []string, field reflect.StructField) (*fieldDescription, error) {
	fieldDesc := &fieldDescription{Type: field.Type}
	if len(tags) > 1 {
		tag := strings.Trim(tags[1], " -")
		if tag == mandatory {
			fieldDesc.Mandatory = true
		} else {
			fieldDesc.ShortName = tag
			if len(fieldDesc.ShortName) > 1 {
				return nil, fmt.Errorf("field '%s': %w (got '%s', expected a single char value)", field.Name,
					ErrInvalidTag, field.Tag.Get("clap"))
			}
			if len(tags) == 3 {
				tag := strings.Trim(tags[2], " ")
				if tag != mandatory {
					return nil, fmt.Errorf("field '%s': %w (got '%s', expected '%s,mandatory')", field.Name,
						ErrInvalidTag, field.Tag.Get("clap"), field.Name)
				}
				fieldDesc.Mandatory = true
			}
		}
	}
	return fieldDesc, nil
}

func computeFieldDescriptions(t reflect.Type) (map[string]*fieldDescription, error) {
	fieldDescs := make(map[string]*fieldDescription)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		tagString := field.Tag.Get("clap")
		if tagString != "" {
			var err error
			var fieldDesc *fieldDescription
			tags := strings.Split(tagString, ",")
			tag := strings.Trim(tags[0], " ")
			switch tag {
			case trailing:
				fieldDesc, err = getTrailingFieldDescription(tags, field)
				if err != nil {
					return nil, err
				}
				fieldDescs[trailing] = fieldDesc
			case "":
				fieldDesc, err = getShortNameFieldDescription(tags, field)
				if err != nil {
					return nil, err
				}
				fieldDescs["-"+fieldDesc.ShortName] = fieldDesc
			default:
				fieldDesc, err = getLongNameFieldDescription(tags, field)
				if err != nil {
					return nil, err
				}
				fieldDesc.LongName = strings.Trim(tag, "-")
				if fieldDesc.LongName != "" {
					fieldDescs["--"+fieldDesc.LongName] = fieldDesc
					if field.Type.Kind() == reflect.Bool {
						fieldDescs["--no-"+fieldDesc.LongName] = fieldDesc
					}
				}
				if fieldDesc.ShortName != "" {
					fieldDescs["-"+fieldDesc.ShortName] = fieldDesc
				}
			}
			fieldDesc.Field = i
		}
	}
	return fieldDescs, nil
}
