package main

import (
	"bytes"
	"errors"
	"os/exec"
	"strconv"
	"strings"
	"sync"
)

type job struct {
	WebRoot  string
	Plugin   string
	Params   string
	Stdout   string
	Stderr   string
	Status   string
	ExitCode int
}

func (s *job) start(wg *sync.WaitGroup) {
	defer wg.Done()
	cmd := exec.Command("cmd", "/C", s.WebRoot+"\\plugins\\"+s.Plugin+"\\exec.bat "+s.Params)
	var OutBuffer, ErrBuffer bytes.Buffer
	cmd.Stdout = &OutBuffer
	cmd.Stderr = &ErrBuffer
	err := cmd.Run()
	s.Stdout = OutBuffer.String()
	s.Stderr = ErrBuffer.String()
	if err != nil {
		s.Status = "failed"
		if exitError, ok := err.(*exec.ExitError); ok {
			s.ExitCode = exitError.ExitCode()
		} else {
			s.ExitCode = -1
		}
	} else {
		s.ExitCode = 0
		s.Status = "complete"
	}
}

type jobs struct {
	wg    *sync.WaitGroup
	cfg   *ConfigServer
	list  []*job
	limit int
}

func (j *jobs) init() {
	j.wg = new(sync.WaitGroup)
	j.list = []*job{}
	j.limit = j.cfg.JobsLimit
}

func (j *jobs) ClearJobs() error {
	var jobs_busy []string
	i := 0
	for i < len(j.list) {
		if j.list[i].Status == "complete" {
			start := i
			for j.list[i].Status == "complete" && i < len(j.list)-1 {
				i++
			}
			end := i
			j.list = append(j.list[:start], j.list[end+1:]...)
			i = i - (end - start)
		} else {
			jobs_busy = append(jobs_busy, strconv.Itoa(i))
			i++
		}
	}
	if len(jobs_busy) != 0 {
		return errors.New("there are pending jobs: " + strings.Join(jobs_busy, ", "))
	}
	j.list = []*job{}
	return nil
}

func (j *jobs) AddJob(WebRoot, Plugin, Params string) int {
	if (len(j.list) + 1) > j.limit {
		j.ClearJobs()
	}
	j.wg.Add(1)
	j.list = append(j.list, &job{WebRoot, Plugin, Params, "", "", "pending", 0})
	id := len(j.list) - 1
	go j.list[id].start(j.wg)
	return id
}
