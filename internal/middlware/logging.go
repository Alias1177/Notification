package middlware

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"github.com/fatih/color"
)

// CustomColorHandler определяет обработчик для кастомного логгера с цветным выводом
type CustomColorHandler struct{}

// NewCustomColorHandler создаёт новый экземпляр кастомного обработчика логов
func NewCustomColorHandler() slog.Handler {
	return &CustomColorHandler{}
}

// Enabled проверяет, включён ли логгер для определённого уровня (в данном случае всегда true)
func (h *CustomColorHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return true
}

// Handle обрабатывает каждый лог-запись и форматирует её с цветным выводом
func (h *CustomColorHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level
	var levelColorFunc func(a ...interface{}) string

	// Определяем цвет для каждого уровня логирования
	switch level {
	case slog.LevelDebug:
		levelColorFunc = func(a ...interface{}) string {
			return color.MagentaString("%v", a...)
		}
	case slog.LevelInfo:
		levelColorFunc = func(a ...interface{}) string {
			return color.BlueString("%v", a...)
		}
	case slog.LevelWarn:
		levelColorFunc = func(a ...interface{}) string {
			return color.YellowString("%v", a...)
		}
	case slog.LevelError:
		levelColorFunc = func(a ...interface{}) string {
			return color.RedString("%v", a...)
		}
	default:
		levelColorFunc = fmt.Sprint
	}

	// Форматируем и выводим сообщение лога
	msg := r.Message
	timestamp := r.Time.Format(time.RFC3339)
	levelStr := levelColorFunc(level.String())
	fmt.Fprintf(os.Stdout, "%s [%s]: %s\n", timestamp, levelStr, msg)

	return nil
}

// WithAttrs добавляет атрибуты (в данном случае возвращает тот же обработчик без изменений)
func (h *CustomColorHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return h
}

// WithGroup добавляет группу (в данном случае возвращает тот же обработчик без изменений)
func (h *CustomColorHandler) WithGroup(name string) slog.Handler {
	return h
}

// ColorLogger настраивает глобальный логгер на использование кастомного цветного обработчика
func ColorLogger() {
	colorHandler := NewCustomColorHandler()
	slog.SetDefault(slog.New(colorHandler))
	color.NoColor = false
}
