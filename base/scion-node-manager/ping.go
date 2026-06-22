package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

var (
	pingMutex          sync.Mutex
	currentPing        CommandState
	pingResultLocation = "/var/lib/scion-node-manager/ping-results/"
	pingManager        = &ProcessManager{
		Mutex:     &pingMutex,
		State:     &currentPing,
		LogDir:    pingResultLocation,
		LogPrefix: "ping_",
	}
)

func initPingState() {
	currentPing = CommandState{InProgress: false}
}

func startPing(w http.ResponseWriter, r *http.Request) {
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
		Dst   string `json:"dst"`
		Count *int   `json:"count,omitempty"`
	}

	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Invalid Body",
		})
		return
	}

	if ip := net.ParseIP(request.Dst); ip == nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Invalid dst IP",
		})
		return
	}

	pingMutex.Lock()
	defer pingMutex.Unlock()

	if currentPing.InProgress {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "A ping is already in progress",
		})
		return
	}

	args := []string{request.Dst}
	if request.Count != nil {
		countStr := fmt.Sprintf("%d", *request.Count)
		args = append([]string{"-c", countStr}, args...)
	}
	logFileName := fmt.Sprintf("ping_%d.log", time.Now().Unix())
	err := pingManager.StartProcess(args, logFileName, "ping")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to start ping: %v", err),
		})
		return
	}
	log.Printf("Started ping (PID: %d), dst: %s", currentPing.PID, request.Dst)

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Status:  "success",
		Message: fmt.Sprintf("Pinging started (PID: %d)", currentPing.PID),
		Data:    currentPing,
	})
}

func stopPing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Method not allowed - use POST",
		})
		return
	}

	pingMutex.Lock()
	defer pingMutex.Unlock()

	if !currentPing.InProgress {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "No ping in progress",
		})
		return
	}

	stoppedPID := currentPing.PID

	err := pingManager.StopProcess()
	if err != nil {
		log.Printf("Failed to stop ping (PID: %d)", stoppedPID)

		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to stop ping: %v", err),
		})
		return
	}

	log.Printf("Stopped ping (PID: %d)", stoppedPID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Status:  "success",
		Message: fmt.Sprintf("Stop signal sent to ping process (PID: %d)", stoppedPID),
	})
}

func getAvailablePingResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Method not allowed - use GET",
		})
		return
	}

	fileInfos, err := getFilesFromDirectory(pingResultLocation)
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

func getPingStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Invalid Method - use GET",
		})
		return
	}

	pingMutex.Lock()
	defer pingMutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Status: "success",
		Data:   currentPing,
	})
}
