package main

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"os/exec"
	"strconv"
	"strings"
)

type Pid struct {
	Shell string
}

type Service struct {
	Pid     Pid
	Command string
}

func (s *Service) getPid() (int, error) {
	cmd := exec.Command("bash", "-c", s.Pid.Shell)
	var out bytes.Buffer
	cmd.Stdout = &out

	if err := cmd.Run(); err != nil {
		return 0, err
	}
	outS := out.String()
	fmt.Printf("cmd output: %s\n", outS)
	outS = strings.Trim(outS, " \n")
	if outS == "" {
		return 0, nil
	}
	i, err := strconv.Atoi(outS)
	if err != nil {
		// TODO(vad): wrap exception, make explicit that the error is due to a bad command
		return 0, err
	}

	return i, nil
}

func (s *Service) Status() (bool, error) {
	pid, err := s.getPid()
	if err != nil {
		return false, err
	}
	if pid > 0 {
		return true, nil
	}
	return false, nil
}

func (s *Service) Start() error {
	cmd := exec.Command("spotify")

	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

type Conf struct {
	Services map[string]*Service
}

func parseConf(f string) *Conf {
	c := Conf{}

	data, err := ioutil.ReadFile(f)
	if err != nil {
		log.Fatalln("Error reading configuration file:", err)
	}

	if err := yaml.Unmarshal([]byte(data), &c); err != nil {
		log.Fatalln("Error parsing configuration file:", err)
	}

	return &c
}
