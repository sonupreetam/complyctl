// SPDX-License-Identifier: Apache-2.0

package log

import (
	"fmt"
	"io"
	"log"
	"os"

	"github.com/hashicorp/go-hclog"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Enable changing the log level
var atom = zap.NewAtomicLevel()

// Pre-configure the logger
func init() {
	zap.ReplaceGlobals(zap.Must(zap.NewProduction()))
}

// Adapt zap to hclog
type ZapHclog struct {
	logger *zap.SugaredLogger
}

// Ensure ZapHclog implements hclog.Logger
var _ hclog.Logger = (*ZapHclog)(nil)

var hclogZapLevels = map[hclog.Level]zapcore.Level{
	// NoLevel is a special level used to indicate that no level has been
	// set and allow for a default to be used.
	//hclog.NoLevel
	hclog.Trace: zapcore.DebugLevel, // No "Trace" equivalent
	hclog.Debug: zapcore.DebugLevel,
	hclog.Info:  zapcore.InfoLevel,
	hclog.Warn:  zapcore.WarnLevel,
	hclog.Error: zapcore.ErrorLevel,
	hclog.Off:   zapcore.FatalLevel, // No "Off" equivalent
}

var zapHclogLevels = map[zapcore.Level]hclog.Level{
	zapcore.DebugLevel:   hclog.Debug,
	zapcore.InfoLevel:    hclog.Info,
	zapcore.WarnLevel:    hclog.Warn,
	zapcore.ErrorLevel:   hclog.Error,
	zapcore.PanicLevel:   hclog.Error, // No "Panic" equivalent
	zapcore.FatalLevel:   hclog.Error, // No "Fatal" equivalent
	zapcore.InvalidLevel: hclog.Error, // No "Invalid" equivalent
}

func (z ZapHclog) Log(level hclog.Level, msg string, args ...interface{}) {
	z.logger.Logw(hclogZapLevels[level], fmt.Sprintf(msg, args...))
}

func (z ZapHclog) Trace(msg string, args ...interface{}) {
	z.logger.Debugw(msg, args...)
}

func (z ZapHclog) Debug(msg string, args ...interface{}) {
	z.logger.Debugw(msg, args...)
}

func (z ZapHclog) Info(msg string, args ...interface{}) {
	z.logger.Infow(msg, args...)
}

func (z ZapHclog) Warn(msg string, args ...interface{}) {
	z.logger.Warnw(msg, args...)
}

func (z ZapHclog) Error(msg string, args ...interface{}) {
	z.logger.Errorw(msg, args...)
}

func (z ZapHclog) IsTrace() bool {
	return false
}

func (z ZapHclog) IsDebug() bool {
	return false
}

func (z ZapHclog) IsInfo() bool {
	return false
}

func (z ZapHclog) IsWarn() bool {
	return false
}

func (z ZapHclog) IsError() bool {
	return false
}

func (z ZapHclog) ImpliedArgs() []interface{} {
	//do nothing
	return nil
}

func (z ZapHclog) With(args ...interface{}) hclog.Logger {
	return &ZapHclog{z.logger.With(args...)}
}

func (z ZapHclog) Name() string {
	return z.logger.Desugar().Name()
}

func (z ZapHclog) Named(name string) hclog.Logger {
	return &ZapHclog{z.logger.Named(name)}
}

func (z ZapHclog) ResetNamed(name string) hclog.Logger {
	logger, err := zap.NewProduction()
	if err != nil {
		// TODO decide what to do
		panic(err)
	}
	return &ZapHclog{logger.Sugar().Named(name)}
}

func (z ZapHclog) SetLevel(level hclog.Level) {
	atom.SetLevel(hclogZapLevels[level])
}

func (z ZapHclog) GetLevel() hclog.Level {
	return zapHclogLevels[atom.Level()]
}

func (z ZapHclog) StandardLogger(opts *hclog.StandardLoggerOptions) *log.Logger {
	//TODO decide action
	panic("not implemented")
}

func (z ZapHclog) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer {
	return os.Stdout
}

func Logger() *zap.SugaredLogger {
	return zap.S()
}
