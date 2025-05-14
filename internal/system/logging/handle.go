package logging

import (
	"context"
	"emperror.dev/emperror"
	"emperror.dev/errors"
	"log/slog"
)

type slogHandler struct {
	logger *slog.Logger
}

func NewSlogHandler(logger *slog.Logger) emperror.ErrorHandler {
	return &slogHandler{
		logger: logger,
	}
}

func (h *slogHandler) Handle(err error) {
	if err == nil {
		return
	}

	logger := h.logger

	// Extract details from the error and attach it to the log
	if details := errors.GetDetails(err); len(details) > 0 {
		logger = h.logger.With(details)
	}

	type errorCollection interface {
		Errors() []error
	}

	if errs, ok := err.(errorCollection); ok {
		for _, e := range errs.Errors() {
			logger.With("parent", err.Error()).Error(e.Error())
		}
	} else {
		logger.Error(err.Error())
	}
}

func (h *slogHandler) HandleContext(_ context.Context, err error) {
	if err == nil {
		return
	}

	h.Handle(err)
}
