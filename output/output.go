package output

import (
	"fmt"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
)

type LogConfig struct {
	cpuLog  *os.File
	hookLog *os.File
	// some kind of socket that is to be listened on for api calls
}

type Logger interface {
	LogStats(stats string) error
	LogHook(key, val string) error // a arbitrary key/value pair passed in through the socket
	Flush() error                  // force all writes through to disk
}

func InitLogger(prefix string) (Logger, error) {
	cpufile, err := os.Create(fmt.Sprintf("%s_cpu.csv", prefix))
	if err != nil {
		return nil, err
	}
	hookfile, err := os.Create(fmt.Sprintf("%s_hook.csv", prefix))
	if err != nil {
		return nil, err
	}
	return LogConfig{cpuLog: cpufile, hookLog: hookfile}, nil
}

// Prints to stdout and logs
func (c LogConfig) LogStats(stats string) error {
	log.WithField("stats", stats).Info("log stats")
	_, err := c.cpuLog.WriteString(fmt.Sprintf("%s\n", stats))
	return err
}

// Prints to stdout and logs
func (c LogConfig) LogHook(rawKey, rawVal string) error {
	key := strings.TrimSpace(rawKey)
	val := strings.TrimSpace(rawVal)
	log.WithFields(log.Fields{
		"key":   key,
		"value": val,
	}).Info("log hook")
	_, err := c.hookLog.WriteString(fmt.Sprintf("%s,%s\n", key, val))
	return err
}

func (c LogConfig) Flush() error {
	return nil
}
