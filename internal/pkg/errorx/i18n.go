package errorx

import (
	"github.com/epkgs/i18n"
	i18nerrx "github.com/epkgs/i18n/errorx"
)

func Define[Args any](i18n *i18n.I18n, code int, format string, httpStatus int) *i18nerrx.Definition[*AppError, Args] {
	return i18nerrx.Define[Args](i18n, format, wrapAppError(code, httpStatus))
}

func DefineSimple(i18n *i18n.I18n, code int, format string, httpStatus int) *i18nerrx.DefinitionSimple[*AppError] {
	return i18nerrx.DefineSimple(i18n, format, wrapAppError(code, httpStatus))
}

func wrapAppError(code, httpStatus int) i18nerrx.Wrapper[*AppError] {
	return func(err *i18nerrx.Error) *AppError {
		return NewAppError(code, err, httpStatus)
	}
}
