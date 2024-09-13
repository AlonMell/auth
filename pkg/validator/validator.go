package validator

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"sync"
)

// All tags:
// - required
// - alpha
// - alphanum
// - email
// - phone
// - password
// - uuid

type ValidationErrors struct {
	Errors []error
}

func (ve ValidationErrors) Error() string {
	var sb strings.Builder
	sb.WriteString("validation errors:")
	for _, err := range ve.Errors {
		sb.WriteString(fmt.Sprintf(" - %s", err))
	}
	return sb.String()
}

func (ve *ValidationErrors) Add(err error) {
	ve.Errors = append(ve.Errors, err)
}

func (ve ValidationErrors) HasErrors() bool {
	return len(ve.Errors) > 0
}

func (ve ValidationErrors) UnWrap() []error {
	return ve.Errors
}

var (
	ErrNotValidValue    = errors.New("field is not valid")
	ErrNotValidPassword = errors.New("password is not valid")
	ErrUnknownTag       = errors.New("unknown tag value")
	ErrRequiredField    = errors.New("field required")
)

var regexMap = map[string]*regexp.Regexp{
	"alpha":    regexp.MustCompile(`^[a-zA-Z]+$`),
	"alphanum": regexp.MustCompile(`^[a-zA-Z0-9]+$`),
	"email":    regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`),
	"phone":    regexp.MustCompile(`^\+?[0-9]{1,15}$`),
	"uuid":     regexp.MustCompile(`^[a-f0-9]{8}-[a-f0-9]{4}-4[a-f0-9]{3}-[89ab][a-f0-9]{3}-[a-f0-9]{12}$`),
}

type Field struct {
	Name  string
	Value reflect.Value
	Tags  []string
}

// Struct validates a struct based on the tags provided.
// It returns an error if the validation fails.
func Struct(s any) error {
	v := reflect.ValueOf(s)
	t := reflect.TypeOf(s)

	nums := t.NumField()
	errChan := make(chan error, nums)
	var wg sync.WaitGroup

	wg.Add(nums)
	for num := range nums {
		go func(n int) {
			defer wg.Done()

			f := Field{
				Name:  t.Field(n).Name,
				Value: v.Field(n),
				Tags:  strings.Split(t.Field(n).Tag.Get("validate"), ","),
			}

			if err := f.validate(); err != nil {
				errChan <- err
			}
		}(num)
	}

	go func() {
		wg.Wait()
		close(errChan)
	}()

	var errs ValidationErrors
	for err := range errChan {
		errs.Add(err)
	}

	if errs.HasErrors() {
		return errs
	}

	return nil
}

func (f *Field) validate() error {
	for _, tag := range f.Tags {
		if tag == "" {
			continue
		}

		if err := f.validateTag(tag); err != nil {
			return err
		}
	}
	return nil
}

func (f *Field) validateTag(tag string) error {
	if pattern, ok := regexMap[tag]; ok {
		return f.validatePattern(tag, pattern)
	}
	return f.switcher(tag)
}

func (f *Field) validatePattern(tag string, pattern *regexp.Regexp) error {
	if str := f.Value.String(); !pattern.MatchString(str) {
		return formatError(f.Name, ErrNotValidValue)
	}
	return nil
}

func (f *Field) switcher(tag string) error {
	switch tag {
	case "required":
		return f.validateRequired()
	case "password":
		return f.validatePassword()
	default:
		return formatError(f.Name, ErrUnknownTag)
	}
}

func (f *Field) validateRequired() error {
	if f.Value.IsZero() {
		return formatError(f.Name, ErrRequiredField)
	}
	return nil
}

func (f *Field) validatePassword() error {
	s := f.Value.String()
	pass := regexp.MustCompile(`^[a-zA-Z\d]{8,}$`).MatchString(s)
	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(s)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(s)
	hasDigit := regexp.MustCompile(`\d`).MatchString(s)

	if !(pass && hasUpper && hasLower && hasDigit) {
		return formatError(f.Name, ErrNotValidPassword)
	}

	return nil
}

func formatError(name string, err error) error {
	return fmt.Errorf("%s: %w", name, err)
}
