package logger

import (
	"context"
	"fmt"

	"github.com/natefinch/lumberjack"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Claves de contexto predeterminadas
type contextKey string

const (
	CtxUserIDKey contextKey = "userID"
	CtxReqIDKey  contextKey = "reqID"
)

// Permite ejecutar acciones extra en cada log
type HookFunc func(level zapcore.Level, msg string, fields map[string]interface{})

type Logger struct {
	sugar *zap.SugaredLogger
	hooks []HookFunc
}

// Crea un logger con rotación y nivel configurable
func NewLogger(logFile string, level zapcore.Level, hooks ...HookFunc) (*Logger, error) {
	writer := zapcore.AddSync(&lumberjack.Logger{
		Filename:   logFile,
		MaxSize:    10, // MB
		MaxBackups: 5,
		MaxAge:     30, // días
		Compress:   true,
	})

	encoderCfg := zap.NewProductionEncoderConfig()
	encoderCfg.TimeKey = "timestamp"
	encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder

	core := zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		writer,
		level,
	)

	logger := zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))
	sugar := logger.Sugar()

	return &Logger{
		sugar: sugar,
		hooks: hooks,
	}, nil
}

// Añade userID al contexto
func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, CtxUserIDKey, userID)
}

// Añade reqID al contexto
func WithReqID(ctx context.Context, reqID string) context.Context {
	return context.WithValue(ctx, CtxReqIDKey, reqID)
}

// Obtiene campos del contexto para log estructurado
func getFieldsFromCtx(ctx context.Context) map[string]interface{} {
	fields := map[string]interface{}{}
	if ctx == nil {
		return fields
	}
	if v := ctx.Value(CtxUserIDKey); v != nil {
		fields["userID"] = v
	}
	if v := ctx.Value(CtxReqIDKey); v != nil {
		fields["reqID"] = v
	}
	return fields
}

func (l *Logger) runHooks(level zapcore.Level, msg string, fields map[string]interface{}) {
	for _, h := range l.hooks {
		h(level, msg, fields)
	}
}

func (l *Logger) Info(ctx context.Context, msg string, args ...any) {
	fields := getFieldsFromCtx(ctx)
	l.sugar.Infow(fmt.Sprintf(msg, args...), fields)
	l.runHooks(zapcore.InfoLevel, msg, fields)
}

func (l *Logger) Debug(ctx context.Context, msg string, args ...any) {
	fields := getFieldsFromCtx(ctx)
	l.sugar.Debugw(fmt.Sprintf(msg, args...), fields)
	l.runHooks(zapcore.DebugLevel, msg, fields)
}

func (l *Logger) Warn(ctx context.Context, msg string, args ...any) {
	fields := getFieldsFromCtx(ctx)
	l.sugar.Warnw(fmt.Sprintf(msg, args...), fields)
	l.runHooks(zapcore.WarnLevel, msg, fields)
}

func (l *Logger) Error(ctx context.Context, msg string, args ...any) {
	fields := getFieldsFromCtx(ctx)
	l.sugar.Errorw(fmt.Sprintf(msg, args...), fields)
	l.runHooks(zapcore.ErrorLevel, msg, fields)
}

// Devuelve un contexto válido si el recibido es nil
func EnsureContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.TODO()
	}
	return ctx
}

// Sync asegura que todos los logs se escriban antes de cerrar
func (l *Logger) Sync() {
	_ = l.sugar.Sync()
}
