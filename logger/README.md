# Enhanced Coloured Logger

This is a utility logger that print the logs in with colour of the log level.

Example:

<span style='color:green'>[YourApp] 2024-01-30T10:41:25+09:00 LOG [SomeFunction] Hello, world!</span>

<span style='color:yellow'>[YourApp] 2024-01-30T10:41:25+09:00 WARN [OtherFunction] Hello, world!</span>

<span style='color:red'>[YourApp] 2024-01-30T10:41:25+09:00 ERROR [MaybeMain] Hello, world!</span>

## Log structure

The log is consisted of the following components

`[App name] Time LOGLEVEL [Logger Name] Message`

- App name
  - The name of the app. This will help identify the source app of the log
- Time
  - The time of the log creation. ISO format
- Log Level
  - The log level. (TRACE, DEBUG, LOG, WARN, ERROR, FATAL, PANIC)
  - The Debug level is only available when env `RUNTIME=development`
- Logger Name
  - The name of the logger. Pass in the name of the component (struct name, method name) to help identify the source of the log
- Message
  - The message of the log

## Basic Usage

### Initialization

The Logger shares a global context for the App's name. Set this by calling the `SetAppName` method from the logger. If the `SetAppName` is not called, then the default app name `GoApp` will be used.

The logger can be instantiated using `logger.NewLogger`. At least the name of the logger must be given through the init options

```golang
func main() {
  // Instantiate a new Logger
  l := logger.NewLogger(logger.LoggerOption {
    Name: "ThisLogger"
  })

  // Default text log
  l.Log("Hello world!")

  // Formatted version
  l.Logf("%s", "Hello World!")
}
```

All Levels of the logger provide a formatting version `~f` thus allows a formatted string to be used in the log.

## Advanced usage

### Extra streams

The logger provides an interface spec that the user can provide to acquire the logs and manipulate in their own manner.

```golang
type (
  LogMessage struct {
    AppName string
    Time    string
    Name    string
    Color   string
    Level   string
    Msg     string
  }

  ILogStream interface {
    Write(msg LogMessage)
  }
)
```

Any struct that satisfies the `ILogStream` interface can be injected with the logger

For example, the default extrastream `FileLogStream` util can be initialized like below to write the same logs of the stdout to a rollover filestream

```golang
l := logger.NewLogger(logger.LoggerOption{
  Name: "test",
  ExtraStreams: []logger.ILogStream{
    logger.NewFileLogStream(logger.FileLogStreamOption{
      LogDirectory: "./logs",
      FileName:     "app",
    }),
  },
})
```

If you want to enable a specific extra stream global so every new logger has the custom extra stream.

```golang
logger.AddGlobalExtraStream([]logger.ILogStream{
  logger.NewFileLogStream(logger.FileLogStreamOption{
    LogDirectory: "./logs",
    FileName:     "app",
  }),
})

// Both loggers below will write a file log
l := logger.NewLogger(logger.LoggerOption {
  Name: "Logger 1",
})

l2 := logger.NewLogger(logger.LoggerOption {
  Name: "Logger 2",
})
```

### File Log Stream

This is a default log stream that writes the logs to a file for persistance.

The initializer options are as follows

```golang
FileLogStreamOption struct {
  LogDirectory  string
  FileName      string
  FileRollover  bool
  MaxFileSizeKb int
}
```

- LogDirectory
  - The directory that the logs are saved to. The directory is created automatically if not exists
- FileName
  - The file name prefix of the log file. If the `RollOver` options is not enabled, then the file name would be `${FileName}.log`
  - If the Rollover is enabled, `${FileName}.${Date}.log` will be used
- FileRollover
  - If `true`, the logfile will move on when the date changes
- MaxFileSizeKb (Not supported Yet)
  - If the Log's size reaches this size in KB, a new log file is created
