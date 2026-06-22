package main

import (
	"log"
	"os"
	"os/exec"
	"sync"
	"time"
)

type ProcessManager struct {
	Cmd       *exec.Cmd
	Mutex     *sync.Mutex
	State     *CommandState
	LogDir    string
	LogPrefix string
}

func (pm *ProcessManager) StartProcess(args []string, logFileName string, command string) error {
	outputFile := pm.LogDir + logFileName
	pm.Cmd = exec.Command(command, args...)
	file, err := os.OpenFile(outputFile, os.O_CREATE|os.O_WRONLY, 0644)
	log.Printf("preparing logfile: %s for process (PID: %d)", file.Name(), pm.State.PID)
	if err != nil {
		log.Printf("couldn't open log file: %v", err)
		return err
	}
	defer file.Close()
	pm.Cmd.Stdout = file
	pm.Cmd.Stderr = file
	if err := pm.Cmd.Start(); err != nil {
		log.Printf("couldn't start process: (PID: %d). Error: %v", pm.State.PID, err)

		return err
	}
	*pm.State = CommandState{
		InProgress: true,
		PID:        pm.Cmd.Process.Pid,
		StartTime:  time.Now(),
		OutputFile: outputFile,
	}
	go func() {
		err := pm.Cmd.Wait()
		pm.Mutex.Lock()
		defer pm.Mutex.Unlock()
		if err != nil {
			// log the error
			log.Printf("process (PID: %d) finished with error: %v", pm.State.PID, err)
		}
		log.Printf("process (PID: %d) finished, logfile: %s", pm.State.PID, outputFile)

		pm.State.InProgress = false
		pm.State.PID = 0
	}()
	return nil
}

func (pm *ProcessManager) StopProcess() error {
	if pm.Cmd == nil || pm.Cmd.Process == nil {
		return os.ErrInvalid
	}
	return pm.Cmd.Process.Signal(os.Interrupt)
}
