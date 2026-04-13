package validator

import (
	"regexp"
	"strings"

	ut "github.com/go-playground/universal-translator"
	validator "github.com/go-playground/validator/v10"
)

var bareGitVersionPattern = regexp.MustCompile(`^\d+\.\d+\.\d+([-.][0-9A-Za-z.-]+)?$`)

type GitRevisionValidator struct{}

func (GitRevisionValidator) Tag() string {
	return "git_revision"
}

func (GitRevisionValidator) Func() validator.Func {
	return func(fl validator.FieldLevel) bool {
		value := strings.TrimSpace(fl.Field().String())
		if value == "" {
			return false
		}

		// Reject bare semver like 1.0.0 so callers must use the exact Git ref,
		// for example main, a commit SHA, or a v-prefixed tag like v1.0.0.
		return !bareGitVersionPattern.MatchString(value)
	}
}

func (GitRevisionValidator) Translation() (string, validator.TranslationFunc) {
	return "{0} must be an exact git ref such as main, a commit SHA, or a v-prefixed tag like v1.0.0",
		func(ut ut.Translator, fe validator.FieldError) string {
			msg, err := ut.T(fe.Tag(), fe.Field())
			if err != nil {
				return fe.Field() + " must be an exact git ref such as main, a commit SHA, or a v-prefixed tag like v1.0.0"
			}
			return msg
		}
}
