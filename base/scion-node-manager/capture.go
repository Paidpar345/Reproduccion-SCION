package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

// Could consider changing to this:
// https://github.com/ghedo/go.pkt

var (
	captureCmd                   *exec.Cmd
	captureMutex                 sync.Mutex
	currentCapture               CommandState
	packetCapturesResultLocation = "/var/lib/scion-node-manager/packet-captures/"
)

func initCaptureState() {
	currentCapture = CommandState{InProgress: false}
}

// TODO add duration parameter maybe?
func startCapture(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Method not allowed - use POST",
		})
		return
	}

	var request struct {
		Interface string `json:"interface,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		request.Interface = "any"
	}

	if request.Interface == "" {
		request.Interface = "any"
	}

	captureMutex.Lock()
	defer captureMutex.Unlock()

	if currentCapture.InProgress {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Capture already in progress",
		})
		return
	}

	outputFile := fmt.Sprintf("%scapture_%d.pcap", packetCapturesResultLocation, time.Now().Unix())
	args := []string{"-i", request.Interface, "-w", outputFile}

	captureCmd = exec.Command("tcpdump", args...)

	if err := captureCmd.Start(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to start capture: %v", err),
		})
		return
	}

	currentCapture = CommandState{
		InProgress: true,
		PID:        captureCmd.Process.Pid,
		StartTime:  time.Now(),
		OutputFile: outputFile,
	}

	log.Printf("Started tcpdump (PID: %d) on %s", currentCapture.PID, request.Interface)

	go func() {
		err := captureCmd.Wait()
		captureMutex.Lock()
		defer captureMutex.Unlock()
		if err != nil {
			log.Printf("tcpdump process (PID: %d) finished with error: %v", currentCapture.PID, err)
		} else {
			log.Printf("tcpdump process (PID: %d) finished successfully", currentCapture.PID)
		}
		currentCapture.InProgress = false
		currentCapture.PID = 0
	}()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Status:  "success",
		Message: fmt.Sprintf("Capture started on %s (PID: %d)", request.Interface, currentCapture.PID),
		Data:    currentCapture,
	})
}

func stopCapture(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Method not allowed - use POST",
		})
		return
	}

	captureMutex.Lock()
	defer captureMutex.Unlock()

	if !currentCapture.InProgress || captureCmd == nil || captureCmd.Process == nil {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "No capture is currently running",
		})
		return
	}

	if err := captureCmd.Process.Signal(os.Interrupt); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to stop capture: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Status:  "success",
		Message: fmt.Sprintf("Stop signal sent to capture process (PID: %d)", currentCapture.PID),
	})
}

func getCaptureStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Method not allowed - use GET",
		})
		return
	}

	captureMutex.Lock()
	defer captureMutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Status: "success",
		Data:   currentCapture,
	})
}

func getAvailableCaptures(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Method not allowed - use GET",
		})
		return
	}

	fileInfos, err := getFilesFromDirectory(packetCapturesResultLocation)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: err.Error(),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Status: "success",
		Data:   fileInfos,
	})
}
