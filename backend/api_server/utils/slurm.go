package utils

import (
	"bytes"
	"errors"
	"io"
	"os/exec"
	"regexp"

	"api_server/logger"
)

type Slurm struct{}

var slurm *Slurm

func NewSlurm() *Slurm {
	if slurm == nil {
		slurm = &Slurm{}
	}

	return slurm
}

func (slurm *Slurm) Sbatch(scriptPath string, scriptArgs []string) (string, error) {
	var out bytes.Buffer

	sbatchArgs := append([]string{scriptPath}, scriptArgs...)
	cmd := exec.Command("sbatch", sbatchArgs...)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		return "", err
	}
	logger.Debug(cmd.String())

	output := out.String()
	reg := regexp.MustCompile(`Submitted batch job (\d+)`)
	matches := reg.FindStringSubmatch(output)
	if len(matches) != 2 {
		return "", errors.New(output)
	}

	return matches[1], nil
}

func (slurm *Slurm) Scancel(jobID string) error {
	cmd := exec.Command("scancel", jobID)
	err := cmd.Run()
	if err != nil {
		return err
	}

	return nil
}

func (slurm *Slurm) IsTerminatedJob(jobID string) bool {
	r, w := io.Pipe()
	var out bytes.Buffer

	squeueCmd := exec.Command("squeue")
	squeueCmd.Stdout = w

	grepCmd := exec.Command("grep", jobID)
	grepCmd.Stdin = r
	grepCmd.Stdout = &out

	errSQ := squeueCmd.Start()
	if errSQ != nil {
		logger.Debug(errSQ)
	}

	errCmd := grepCmd.Start()
	if errCmd != nil {
		logger.Debug(errCmd)
	}

	go func() {
		defer w.Close()
		errWait := squeueCmd.Wait()
		if errWait != nil {
			logger.Debug(errWait)
		}
	}()
	errWait := grepCmd.Wait()
	if errWait != nil {
		logger.Debug(errWait)
	}

	output := out.String()

	if output == "" {
		return true
	} else {
		return false
	}
}
