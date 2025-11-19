package rpeat

// https://stackoverflow.com/questions/10571182/how-to-disable-a-log-logger

import (
	"io"
	"log"
	"os"
)

var (
	ServerLogger     *log.Logger
	ConnectionLogger *log.Logger
	RequestLogger    *log.Logger
	UpdatesLogger    *log.Logger
)

func initServerLogging(out io.Writer) {
	ServerLogger = log.New(out, "Server     |", log.Ldate|log.Ltime|log.LUTC|log.Lmicroseconds)
}
func initConnectionLogging(out io.Writer) {
	ConnectionLogger = log.New(out, "Connection |", log.Ldate|log.Ltime|log.LUTC|log.Lmicroseconds)
}
func initRequestLogging(out io.Writer) {
	RequestLogger = log.New(out, "Request    |", log.Ldate|log.Ltime|log.LUTC|log.Lmicroseconds)
}
func initUpdatesLogging(out io.Writer) {
	UpdatesLogger = log.New(out, "Updates    |", log.Ldate|log.Ltime|log.LUTC|log.Lmicroseconds)
}

func StartLogs() {
	initServerLogging(os.Stderr)
	initConnectionLogging(os.Stderr)
	initRequestLogging(os.Stderr)
	initUpdatesLogging(os.Stderr)
}
