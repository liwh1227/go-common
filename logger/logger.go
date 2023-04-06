package logger

import (
	"strings"
	"sync"

	"gitee.com/liwh1227/common/logger/core"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Log module
const (
	MODULE_REQUEST = "[Request]"
	MODULE_SYSTEM  = "[System]"

	MODULE_BRIEF = "[Brief]"
	MODULE_EVENT = "[Event]"

	DefaultStackTraceLevel = "PANIC"
)

var (
	// map[module-name]map[module-name+serviceName]zap.AtomicLevel
	loggerLevels = make(map[string]map[string]zap.AtomicLevel)
	loggerMutex  sync.Mutex
	logConfig    *LogConfig

	// map[moduleName+serviceName]*Logger
	serLoggers = sync.Map{}
)

// Logger is an implementation of service logger.
type Logger struct {
	zlog        *zap.SugaredLogger
	name        string
	serviceName string
	lock        sync.RWMutex
	logLevel    core.LOG_LEVEL
}

func (l *Logger) Logger() *zap.SugaredLogger {
	l.lock.RLock()
	defer l.lock.RUnlock()
	return l.zlog
}

func (l *Logger) Debug(args ...interface{}) {
	l.zlog.Debug(args...)
}
func (l *Logger) Debugf(format string, args ...interface{}) {
	l.zlog.Debugf(format, args...)
}
func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.zlog.Debugw(msg, keysAndValues...)
}
func (l *Logger) Error(args ...interface{}) {
	l.zlog.Error(args...)
}
func (l *Logger) Errorf(format string, args ...interface{}) {
	l.zlog.Errorf(format, args...)
}
func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.zlog.Errorw(msg, keysAndValues...)
}
func (l *Logger) Fatal(args ...interface{}) {
	l.zlog.Fatal(args...)
}
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.zlog.Fatalf(format, args...)
}
func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.zlog.Fatalw(msg, keysAndValues...)
}
func (l *Logger) Info(args ...interface{}) {
	l.zlog.Info(args...)
}
func (l *Logger) Infof(format string, args ...interface{}) {
	l.zlog.Infof(format, args...)
}
func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.zlog.Infow(msg, keysAndValues...)
}
func (l *Logger) Panic(args ...interface{}) {
	l.zlog.Panic(args...)
}
func (l *Logger) Panicf(format string, args ...interface{}) {
	l.zlog.Panicf(format, args...)
}
func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.zlog.Panicw(msg, keysAndValues...)
}
func (l *Logger) Warn(args ...interface{}) {
	l.zlog.Warn(args...)
}
func (l *Logger) Warnf(format string, args ...interface{}) {
	l.zlog.Warnf(format, args...)
}
func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.zlog.Warnw(msg, keysAndValues...)
}

func (l *Logger) DebugDynamic(getStr func() string) {
	if l.logLevel == core.LEVEL_DEBUG {
		str := getStr()
		l.zlog.Debug(str)
	}
}
func (l *Logger) InfoDynamic(getStr func() string) {
	if l.logLevel == core.LEVEL_DEBUG || l.logLevel == core.LEVEL_INFO {
		l.zlog.Info(getStr())
	}
}

// SetLogger set logger.
func (l *Logger) SetLogger(logger *zap.SugaredLogger) {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.zlog = logger
}

// newServiceLogger create a new Logger.
func newServiceLogger(name string, serviceName string, logger *zap.SugaredLogger, logLevel core.LOG_LEVEL) *Logger {
	return &Logger{name: name, serviceName: serviceName, zlog: logger, logLevel: logLevel}
}

// SetLogConfig set the config of logger module, called in initialization of config module
func SetLogConfig(config *LogConfig) {
	logConfig = config
	RefreshLogConfig(logConfig)
}

// GetLogger find or create a Logger with module name, usually called in initialization of all module.
// After one module get the logger, the module can use it forever until the program terminate.
func GetLogger(name string) *Logger {
	return GetLoggerByService(name, "")
}

// GetLoggerByChain find the Logger object with module name and serviceName,
// usually called in initialization of all module.
// One module can get a logger for each chain, then logger can be use forever until the program terminate.
func GetLoggerByService(name, serviceName string) *Logger {
	logHeader := name + serviceName
	var logger *Logger
	loggerVal, ok := serLoggers.Load(logHeader)
	if ok {
		logger, _ = loggerVal.(*Logger)
		return logger
	}
	zapLogger, logLevel := createLoggerByService(name, serviceName)

	logger = newServiceLogger(name, serviceName, zapLogger, logLevel)
	loggerVal, ok = serLoggers.LoadOrStore(logHeader, logger)
	if ok {
		logger, _ = loggerVal.(*Logger)
	}
	return logger

}

func createLoggerByService(name, serviceName string) (*zap.SugaredLogger, core.LOG_LEVEL) {
	var config core.LogConfig
	var pureName string

	if logConfig == nil {
		logConfig = DefaultLogConfig()
	}

	if logConfig.LoggerIns.SystemLog.LogLevelDefault == "" {
		defaultLogNode := GetDefaultLogModuleConfig()
		config = core.LogConfig{
			Module:          "[DEFAULT]",
			ServiceName:     serviceName,
			LogPath:         defaultLogNode.FilePath,
			LogLevel:        core.GetLogLevel(defaultLogNode.LogLevelDefault),
			MaxAge:          defaultLogNode.MaxAge,
			RotationTime:    defaultLogNode.RotationTime,
			JsonFormat:      false,
			ShowLine:        true,
			LogInConsole:    defaultLogNode.LogInConsole,
			ShowColor:       defaultLogNode.ShowColor,
			IsBrief:         false,
			StackTraceLevel: defaultLogNode.StackTraceLevel,
		}
	} else {
		if name == MODULE_REQUEST {
			config = core.LogConfig{
				Module:          name,
				ServiceName:     serviceName,
				LogPath:         logConfig.LoggerIns.RequestLog.FilePath,
				LogLevel:        core.GetLogLevel(logConfig.LoggerIns.RequestLog.LogLevelDefault),
				MaxAge:          logConfig.LoggerIns.RequestLog.MaxAge,
				RotationTime:    logConfig.LoggerIns.RequestLog.RotationTime,
				RotationSize:    logConfig.LoggerIns.RequestLog.RotationSize,
				JsonFormat:      true,
				ShowLine:        true,
				LogInConsole:    logConfig.LoggerIns.RequestLog.LogInConsole,
				ShowColor:       logConfig.LoggerIns.RequestLog.ShowColor,
				IsBrief:         false,
				StackTraceLevel: logConfig.LoggerIns.RequestLog.StackTraceLevel,
			}
		} else {
			pureName = strings.ToLower(strings.Trim(name, "[]"))
			value, exists := logConfig.LoggerIns.SystemLog.LogLevels[pureName]
			if !exists {
				value = logConfig.LoggerIns.SystemLog.LogLevelDefault
			}

			config = core.LogConfig{
				Module:          name,
				ServiceName:     serviceName,
				LogPath:         logConfig.LoggerIns.SystemLog.FilePath,
				LogLevel:        core.GetLogLevel(value),
				MaxAge:          logConfig.LoggerIns.SystemLog.MaxAge,
				RotationTime:    logConfig.LoggerIns.SystemLog.RotationTime,
				RotationSize:    logConfig.LoggerIns.SystemLog.RotationSize,
				JsonFormat:      true,
				ShowLine:        true,
				LogInConsole:    logConfig.LoggerIns.SystemLog.LogInConsole,
				ShowColor:       logConfig.LoggerIns.SystemLog.ShowColor,
				IsBrief:         false,
				StackTraceLevel: logConfig.LoggerIns.SystemLog.StackTraceLevel,
			}
		}
	}
	logger, level := core.InitSugarLogger(&config)
	if pureName != "" {
		if _, exist := loggerLevels[pureName]; !exist {
			loggerLevels[pureName] = make(map[string]zap.AtomicLevel)
		}
		logHeader := name + serviceName
		loggerLevels[pureName][logHeader] = level
	}

	return logger, config.LogLevel
}

func refreshAllLoggerOfCmLoggers() {
	serLoggers.Range(func(_, value interface{}) bool {
		serLogger, _ := value.(*Logger)
		newLogger, logLevel := createLoggerByService(serLogger.name, serLogger.serviceName)
		serLogger.SetLogger(newLogger)
		serLogger.logLevel = logLevel
		return true
	})
}

// RefreshLogConfig refresh core levels of modules at initiation time of core module
// or refresh core levels of modules dynamically at running time.
func RefreshLogConfig(config *LogConfig) {
	loggerMutex.Lock()
	defer loggerMutex.Unlock()
	// scan loggerLevels and find the level from config, if can't find level, set it to default
	for name, loggers := range loggerLevels {
		var logLevevl zapcore.Level
		var strlevel string
		var exist bool
		if strlevel, exist = config.LoggerIns.SystemLog.LogLevels[name]; !exist {
			strlevel = config.LoggerIns.SystemLog.LogLevelDefault
		}
		switch core.GetLogLevel(strlevel) {
		case core.LEVEL_DEBUG:
			logLevevl = zap.DebugLevel
		case core.LEVEL_INFO:
			logLevevl = zap.InfoLevel
		case core.LEVEL_WARN:
			logLevevl = zap.WarnLevel
		case core.LEVEL_ERROR:
			logLevevl = zap.ErrorLevel
		default:
			logLevevl = zap.InfoLevel
		}
		for _, aLevel := range loggers {
			aLevel.SetLevel(logLevevl)
		}
	}

	refreshAllLoggerOfCmLoggers()
}

// DefaultLogConfig create default config for core module
func DefaultLogConfig() *LogConfig {
	defaultLogNode := GetDefaultLogModuleConfig()
	config := &LogConfig{
		ConfigFile: "",
		LoggerIns: LoggerIns{
			SystemLog: LogModuleConfig{
				LogLevelDefault: defaultLogNode.LogLevelDefault,
				FilePath:        defaultLogNode.FilePath,
				MaxAge:          defaultLogNode.MaxAge,
				RotationTime:    defaultLogNode.RotationTime,
				RotationSize:    defaultLogNode.RotationSize,
				LogInConsole:    defaultLogNode.LogInConsole,
				StackTraceLevel: defaultLogNode.StackTraceLevel,
			},
			RequestLog: LogModuleConfig{
				LogLevelDefault: defaultLogNode.LogLevelDefault,
				FilePath:        defaultLogNode.FilePath,
				MaxAge:          defaultLogNode.MaxAge,
				RotationTime:    defaultLogNode.RotationTime,
				RotationSize:    defaultLogNode.RotationSize,
				LogInConsole:    defaultLogNode.LogInConsole,
				StackTraceLevel: defaultLogNode.StackTraceLevel,
			},
		},
	}
	return config
}

// GetDefaultLogModuleConfig create a default core config of node
func GetDefaultLogModuleConfig() LogModuleConfig {
	return LogModuleConfig{
		LogLevelDefault: core.DEBUG,
		FilePath:        "./default.core",
		MaxAge:          core.DEFAULT_MAX_AGE,
		RotationTime:    core.DEFAULT_ROTATION_TIME,
		RotationSize:    core.DEFAULT_ROTATION_SIZE,
		LogInConsole:    true,
		ShowColor:       true,
		StackTraceLevel: DefaultStackTraceLevel,
	}
}
