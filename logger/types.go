package logger

// LogConfig the config of core module
type LogConfig struct {
	ConfigFile string    `yaml:"config_file"`
	LoggerIns  LoggerIns `yaml:"log"`
}

type LoggerIns struct {
	RequestLog LogModuleConfig `yaml:"request"`
	SystemLog  LogModuleConfig `yaml:"system"`
}

// LogModuleConfig 设置不同模块的日志
type LogModuleConfig struct {
	LogLevelDefault string            `yaml:"log_level_default"`
	LogLevels       map[string]string `yaml:"log_levels"`
	FilePath        string            `yaml:"file_path"`
	MaxAge          int               `yaml:"max_age"`
	RotationTime    int               `yaml:"rotation_time"`
	RotationSize    int64             `yaml:"rotation_size"`
	LogInConsole    bool              `yaml:"log_in_console"`
	ShowColor       bool              `yaml:"show_color"`
	StackTraceLevel string            `yaml:"stack_trace_level"`
}
