package main

import (
	"bytes"
	"errors"
	"os/exec"
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

var jobs []job

func (s *job) start() {
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

func ClearJobs() error {
	for _, JobEl := range jobs {
		if JobEl.Status == "pending" {
			return errors.New("there are pending jobs")
		}
	}
	jobs = []job{}
	return nil
}

func AddJob(WebRoot, Plugin, Params string) int {
	jobs = append(jobs, job{WebRoot, Plugin, Params, "", "", "pending", 0})
	go jobs[len(jobs)-1].start()
	return len(jobs) - 1
}
