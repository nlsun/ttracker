package tracker

import (
	"fmt"
	"sort"
	"strconv"
	"strings"

	log "github.com/Sirupsen/logrus"
	"github.com/shirou/gopsutil/process"
)

type byName []Target

type Target struct {
	Name string
	Pid  int32

	proc *process.Process
}

func NewTarget(name string, pid int32) (Target, error) {
	proc, err := process.NewProcess(pid)
	if err != nil {
		return Target{}, err
	}
	return Target{
		Name: name,
		Pid:  pid,
		proc: proc,
	}, nil
}

func RunOnce(target Target) string {
	cpuPercent, err := target.proc.Percent(0)
	if err != nil {
		log.WithFields(log.Fields{
			"error": err,
			"name":  target.Name,
			"pid":   target.Pid,
		}).Info("tracker failed")
		return "error"
	}
	return fmt.Sprintf("%f", cpuPercent)
}

func ParseSpec(spec string) ([]Target, error) {
	// Format: "name1,pid1;name2,pid2;..."
	//
	// Returns in sorted order!

	var targets []Target

	chunks := strings.Split(spec, ";")
	for _, chk := range chunks {
		fields := strings.Split(chk, ",")
		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			return nil, err
		}
		tgt, err := NewTarget(fields[0], int32(pid))
		if err != nil {
			return nil, err
		}
		targets = append(targets, tgt)
	}
	sort.Sort(byName(targets))
	return targets, nil
}

func (s byName) Len() int {
	return len(s)
}

func (s byName) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s byName) Less(i, j int) bool {
	return s[i].Name < s[j].Name
}
