package output

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/nlsun/ttracker/util"
)

type HookReq struct {
	Value string `json:"value"`
}

type logConfig struct {
	cpuLog  *os.File
	hookLog *os.File
	addrLog *os.File
	stamp   util.Stamp
}

type Logger interface {
	LogStats(stats string) error
	LogHook(val string) error
	LogAddr(addr string) error

	// http.Handler interface
	ServeHTTP(http.ResponseWriter, *http.Request)
}

func InitLogger(prefix string, stamp util.Stamp) (Logger, error) {
	cpufile, err := os.Create(fmt.Sprintf("%s_cpu.csv", prefix))
	if err != nil {
		return nil, err
	}
	hookfile, err := os.Create(fmt.Sprintf("%s_hook.csv", prefix))
	if err != nil {
		return nil, err
	}
	addrfile, err := os.Create(fmt.Sprintf("%s_addr.log", prefix))
	if err != nil {
		return nil, err
	}
	return logConfig{
		cpuLog:  cpufile,
		hookLog: hookfile,
		addrLog: addrfile,
		stamp:   stamp,
	}, nil
}

// Prints to stdout and logs
func (c logConfig) LogStats(stats string) error {
	log.WithField("stats", stats).Debug("log stats")
	_, err := c.cpuLog.WriteString(fmt.Sprintf("%s\n", stats))
	return err
}

// Prints to stdout and logs
func (c logConfig) LogHook(rawVal string) error {
	stamp := c.stamp()
	val := strings.TrimSpace(rawVal)
	log.WithFields(log.Fields{
		"stamp": stamp,
		"value": val,
	}).Debug("log hook")
	logline := fmt.Sprintf("%s,%s\n", stamp, val)
	_, err := c.hookLog.WriteString(logline)
	return err
}

func (c logConfig) LogAddr(addr string) error {
	log.WithField("addr", addr).Debug("log addr")
	_, err := c.addrLog.Seek(0, 0)
	if err != nil {
		return err
	}
	err = c.addrLog.Truncate(0)
	if err != nil {
		return err
	}
	_, err = c.addrLog.WriteString(addr)
	return err
}

func (c logConfig) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != "PUT" {
		http.Error(w, "error: not a PUT", http.StatusMethodNotAllowed)
		return
	}

	decoder := json.NewDecoder(r.Body)
	var hreq HookReq
	err := decoder.Decode(&hreq)
	if err != nil {
		http.Error(w, "could not decode hook request", http.StatusBadRequest)
		log.WithError(err).Info("could not decode hook request")
	}
	err = c.LogHook(hreq.Value)
	if err != nil {
		log.WithError(err).Info()
	}
}
