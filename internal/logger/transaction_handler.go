package logger

import (
	middlewares "auth-service/internal/middleware"
	"context"
	"log/slog"
)

type trxHandler struct {
	next slog.Handler
}

func (h *trxHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.next.Enabled(ctx, level)
}

func (h *trxHandler) Handle(ctx context.Context, r slog.Record) error {
	if trxID, ok := ctx.Value(middlewares.TrxIDKey).(string); ok {
		r.AddAttrs(slog.String("masterTransactionId", trxID))
	}
	return h.next.Handle(ctx, r)
}

func (h *trxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &trxHandler{next: h.next.WithAttrs(attrs)}
}

func (h *trxHandler) WithGroup(name string) slog.Handler {
	return &trxHandler{next: h.next.WithGroup(name)}
}
