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
	wg   *sync.WaitGroup
	list []*job
}

func (j *jobs) init() {
	j.wg = new(sync.WaitGroup)
	j.list = []*job{}
}

func (j *jobs) ClearJobs() error {
	var jobs_busy []string
	for i := 0; i < len(j.list); i++ {
		if j.list[i].Status == "pending" {
			jobs_busy = append(jobs_busy, strconv.Itoa(i))
		}
	}
	if len(jobs_busy) != 0 {
		return errors.New("there are pending jobs: " + strings.Join(jobs_busy, ", "))
	}
	j.list = []*job{}
	return nil
}

func (j *jobs) AddJob(WebRoot, Plugin, Params string) int {
	j.wg.Add(1)
	j.list = append(j.list, &job{WebRoot, Plugin, Params, "", "", "pending", 0})
	id := len(j.list) - 1
	go j.list[id].start(j.wg)
	return id
}
