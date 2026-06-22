package main

import "time"

type APIResponse struct {
	Status  string      `json:"status"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data,omitempty"`
}

type CommandState struct {
	InProgress bool      `json:"in_progress"`
	PID        int       `json:"pid,omitempty"`
	StartTime  time.Time `json:"start_time,omitempty"`
	OutputFile string    `json:"output_file,omitempty"`
}

type FileInfo struct {
	Index int32  `json:"index"`
	Name  string `json:"name"`
	Size  int64  `json:"size"`
}
