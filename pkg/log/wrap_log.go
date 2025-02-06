package log

import (
	"fmt"
	charmlog "github.com/charmbracelet/log"
	"github.com/hashicorp/go-hclog"
	"io"
	"os"
)

var ErrMissingValue = fmt.Errorf("missing value")

// Need to initialize function

// Wrap the function
func WrapLog(charmlog *charmlog.Logger) hclog.Logger { return &CharmHclog{charmlog} }

// CharmHclog will be a structure that accesses the attributes of charm log
type CharmHclog struct {
	logger *charmlog.Logger
}

// LoggerOption is an option for a logger (from charm logger)
//type LoggerOption = func(*Logger)

// CharmHclog will implement the hclog.Logger

var _ hclog.Logger = &CharmHclog{}

var hclogCharmLevels = map[hclog.Level]charmlog.Level{
	hclog.NoLevel: charmlog.InfoLevel,  // There is no "NoLevel" equivalent in charm, use info
	hclog.Trace:   charmlog.DebugLevel, // There is no "Trace" equivalent in charm, use debug
	hclog.Debug:   charmlog.DebugLevel,
	hclog.Info:    charmlog.InfoLevel,
	hclog.Warn:    charmlog.WarnLevel,
	hclog.Error:   charmlog.ErrorLevel,
	hclog.Off:     charmlog.FatalLevel, // There is no "Off" level equivalent in charm
}

var charmHclogLevels = map[charmlog.Level]hclog.Level{
	charmlog.DebugLevel: hclog.Debug,
	charmlog.InfoLevel:  hclog.Info,
	charmlog.WarnLevel:  hclog.Warn,
	charmlog.ErrorLevel: hclog.Error,
	charmlog.FatalLevel: hclog.Error, // There is no "fatal" equivalent in hclog
}

func (c *CharmHclog) Log(level hclog.Level, msg string, args ...interface{}) {
	c.logger.Log(hclogCharmLevels[level], fmt.Sprintf(msg, args...))
}
func (c *CharmHclog) Trace(msg string, args ...interface{}) {
	c.logger.Debug(msg, args...)
}
func (c *CharmHclog) Debug(msg string, args ...interface{}) {
	c.logger.Debug(msg, args...)
}
func (c *CharmHclog) Info(msg string, args ...interface{}) {
	c.logger.Info(msg, args...)
}
func (c *CharmHclog) Warn(msg string, args ...interface{}) {
	c.logger.Warn(msg, args...)
}
func (c *CharmHclog) Error(msg string, args ...interface{}) {
	c.logger.Error(msg, args...)
}

// Functions from go-hc-log

func (c *CharmHclog) IsTrace() bool     { return false }
func (c *CharmHclog) IsDebug() bool     { return false }
func (c *CharmHclog) IsInfo() bool      { return false }
func (c *CharmHclog) IsWarn() bool      { return false }
func (c *CharmHclog) IsError() bool     { return false }
func (c *CharmHclog) ImpliedArgs() bool { return false }
func (c *CharmHclog) With(args ...interface{}) hclog.Logger {
	return &CharmHclog{c.logger.With(args...)}
}

// Need to configure a Name function
func (c *CharmHclog) Name() string { return c.Name() }

// Take input and then prepend name string
func (c *CharmHclog) Named(name string) hclog.Logger {
	return &CharmHclog{c.logger.With()}
}

// go-hclog logger resetnamed function to implement
//
//func (c *CharmHclog) ResetNamed(name string) hclog.Logger {
//	logger, err := charmlog.NewWithOptions()
//	if err != nil { panic(err) }
//	return &CharmHclog{logger.WithAttrs(name)}
//}

// Enables setting log level
func (c *CharmHclog) SetLevel(level hclog.Level) {
	charmlog.SetLevel(hclogCharmLevels[level])
}

// GetLevel using charm logger GetLevel
func (c *CharmHclog) GetLevel() hclog.Level {
	return charmHclogLevels[hclog.Level(charmlog.GetLevel())]
}

// GetLevel using charm logger Level method to extract
//func (c *CharmHclog) GetLevel() hclog.Level { return charmHclogLevels[hclog.Level(charmlog.Level(level))]}

// Look at stdlog.go for forcing level
// The standard logger needs to be implemented to return a standard logger of the charm type
func (c *CharmHclog) StandardLogger(opts *hclog.StandardLoggerOptions) *charmlog.Logger {
	return charmlog.NewWithOptions(c.logger.StandardLog())
}

func (c *CharmHclog) StandardWriter(opts *hclog.StandardLoggerOptions) io.Writer { return os.Stdout }

func Logger() *charmlog.Logger { return charmlog.StandardLog() }
