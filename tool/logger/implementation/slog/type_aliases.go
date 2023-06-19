package slog

import "golang.org/x/exp/slog"

// Logger is the upstream type, semantically analogous to logger.Logger.
type Logger = slog.Logger

// Handler is the upstream type, semantically analogous to logger.Emitter.
type Handler = slog.Handler

// Level is the upstream type, semantically analogous to logger.Level.
type Level = slog.Level

// Attr is the upstream type, semantically analogous to field.Field.
type Attr = slog.Attr

// Record is the upstream type, semantically analogous to logger.Entry.
type Record = slog.Record
