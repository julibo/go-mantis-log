package conf

type Option func(*Options)

type Options struct {
	LogPath     string // 日志保存路径
	LogName     string // 日志保存的名称，不写随机生成
	LogLevel    string // 日志记录级别
	MaxSize     int    // 日志分割的尺寸 MB
	MaxAge      int    // 分割日志保存的时间 day
	MaxBackups  int    // 日志文件最多保存多少个备份
	Stacktrace  string // 记录堆栈的级别
	IsStdOut    string // 是否标准输出 console输出
	ProjectName string // 项目名称
}

func WithLogPath(logpath string) Option {
	return func(o *Options) {
		o.LogPath = logpath
	}
}

func WithLogName(logname string) Option {
	return func(o *Options) {
		o.LogName = logname
	}
}

func WithLogLevel(loglevel string) Option {
	return func(o *Options) {
		o.LogLevel = loglevel
	}
}

func WithMaxSize(maxsize int) Option {
	return func(o *Options) {
		o.MaxSize = maxsize
	}
}

func WithMaxAge(maxage int) Option {
	return func(o *Options) {
		o.MaxAge = maxage
	}
}

func WithMaxBackups(maxBackups int) Option {
	return func(o *Options) {
		o.MaxBackups = maxBackups
	}
}

func WithStacktrace(stacktrace string) Option {
	return func(o *Options) {
		o.Stacktrace = stacktrace
	}
}

func WithIsStdOut(isstdout string) Option {
	return func(o *Options) {
		o.IsStdOut = isstdout
	}
}

func WithProjectName(projectname string) Option {
	return func(o *Options) {
		o.ProjectName = projectname
	}
}
