package validator

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/aws-cloudformation/cloudformation-cli-go-plugin/cfn/handler"
	"github.com/aws/aws-sdk-go/service/cloudformation"
	"github.com/mongodb/mongodbatlas-cloudformation-resources/util/constants"
	progress_events "github.com/mongodb/mongodbatlas-cloudformation-resources/util/progress_event"
)

type ValidatorDefinition interface {
	GetCreateFields() []string
	GetReadFields() []string
	GetUpdateFields() []string
	GetDeleteFields() []string
	GetListFields() []string
}

func ValidateModel(event constants.Event, def ValidatorDefinition, model interface{}) *handler.ProgressEvent {

	fields := []string{}

	switch event {
	case constants.Create:
		fields = def.GetCreateFields()
	case constants.Read:
		fields = def.GetReadFields()
	case constants.Update:
		fields = def.GetUpdateFields()
	case constants.Delete:
		fields = def.GetDeleteFields()
	default:
		fields = def.GetListFields()
	}

	requiredFields := ""

	for _, field := range fields {
		if fieldIsEmpty(model, field) {
			requiredFields = fmt.Sprintf("%v %v", requiredFields, field)
		}
	}

	if requiredFields == "" {
		return nil
	}

	progressEvent := progress_events.GetFailedEventByCode(fmt.Sprintf("The next fields are required%v", requiredFields),
		cloudformation.HandlerErrorCodeInvalidRequest)

	return &progressEvent
}

func fieldIsEmpty(model interface{}, field string) bool {

	var f reflect.Value
	if strings.Contains(field, ".") {
		fields := strings.Split(field, ".")
		r := reflect.ValueOf(model)

		for _, f := range fields {
			fmt.Printf("%v", f)
			baseProperty := reflect.Indirect(r).FieldByName(f)

			if baseProperty.IsNil() {
				return true
			}

			r = baseProperty
		}
		return false
	} else {
		r := reflect.ValueOf(model)
		f = reflect.Indirect(r).FieldByName(field)
	}

	return f.IsNil()
}