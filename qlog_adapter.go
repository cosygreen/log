package log

import "github.com/rs/zerolog"

type QlogAdapter struct {
	zerolog.Logger
}

func NewQlogAdapter(logger zerolog.Logger) *QlogAdapter {
	return &QlogAdapter{
		Logger: logger,
	}
}

func (q QlogAdapter) Errorf(format string, args ...any) {
	q.Logger.Error().Msgf(format, args...)
}
