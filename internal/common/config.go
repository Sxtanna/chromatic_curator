package common

import "reflect"

type Configuration interface{}

type ConfigurationWithPostProcessing interface {
	Configuration

	Process() error
}

type ConfigurationWithValidation interface {
	Configuration

	Validate() error
}

func OptProcess(config Configuration) error {
	if post, ok := config.(ConfigurationWithPostProcessing); ok {
		return post.Process()
	}

	return nil
}

func OptValidate(config Configuration) error {
	if val, ok := config.(ConfigurationWithValidation); ok {
		return val.Validate()
	}

	return nil
}

func FindConfiguration[C Configuration](config Configuration) *C {

	resolveTargetFromValue := func(value any) *C {
		if casted, ok := value.(*C); ok {
			return casted
		}
		if casted, ok := value.(C); ok {
			return &casted
		}

		return nil
	}

	configType := reflect.TypeOf(config)

	if configType.Kind() == reflect.Ptr {
		configType = configType.Elem()
	}

	if target := resolveTargetFromValue(config); target != nil {
		return target
	}

	if configType.Kind() != reflect.Struct {
		return nil
	}

	configValueElement := reflect.ValueOf(config)

	if configValueElement.Kind() == reflect.Ptr {
		configValueElement = configValueElement.Elem()
	}

	values := make([]any, 0)

	for _, field := range reflect.VisibleFields(configType) {
		if field.IsExported() {
			values = append(values, configValueElement.FieldByName(field.Name).Interface())
		}
	}

	for _, value := range values {
		if target := resolveTargetFromValue(value); target != nil {
			return target
		}
	}

	for _, value := range values {
		if conf := FindConfiguration[C](value.(Configuration)); conf != nil {
			return conf
		}
	}

	return nil
}
