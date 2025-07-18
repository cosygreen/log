package log

import (
	"github.com/rs/zerolog"
)

type QlogAdapter struct {
	logger zerolog.Logger
}

func NewQlogAdapter(logger zerolog.Logger) *QlogAdapter {
	return &QlogAdapter{
		logger: logger,
	}
}

func (q *QlogAdapter) Errorf(format string, args ...any) {
	q.logger.Error().Msgf(format, args...)
}
