package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/scionproto/scion/pkg/addr"
)

// TODO
// https://docs.scion.org/en/latest/command/scion/scion_ping.html#scion-ping
// When the --healthy-only option is set, ping first determines healthy paths through probing and chooses amongst them.

var (
	scionPingMutex          sync.Mutex
	currentScionPing        CommandState
	scionPingResultLocation = "/var/lib/scion-node-manager/scion-ping-results/"
	scionPingManager        = &ProcessManager{
		Mutex:     &scionPingMutex,
		State:     &currentScionPing,
		LogDir:    scionPingResultLocation,
		LogPrefix: "scion-ping_",
	}
)

func initScionPingState() {
	currentScionPing = CommandState{InProgress: false}
}

func startScionPing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Invalid Method - use POST",
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
			Message: "Invalid body",
		})
		return
	}

	_, err := addr.ParseAddr(request.Dst)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Invalid Scion Address %s. error: %s", request.Dst, err.Error()),
		})
		return
	}

	scionPingMutex.Lock()
	defer scionPingMutex.Unlock()

	if currentScionPing.InProgress {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "A scion ping is already in progress",
		})
		return
	}

	args := []string{"ping", request.Dst}
	if request.Count != nil {
		countStr := fmt.Sprintf("%d", *request.Count)
		args = append([]string{"-c", countStr}, args...)
	}
	logFileName := fmt.Sprintf("scion-ping_%d.log", time.Now().Unix())
	err = scionPingManager.StartProcess(args, logFileName, "scion")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to start scion ping: %v", err),
		})
		return
	}
	log.Printf("Started scion ping (PID: %d), dst: %s", currentScionPing.PID, request.Dst)
	json.NewEncoder(w).Encode(APIResponse{
		Status:  "success",
		Message: fmt.Sprintf("Scion Pinging started (PID: %d)", currentScionPing.PID),
		Data:    currentScionPing,
	})
}

func stopScionPing(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Invalid Method - use POST",
		})
		return
	}

	scionPingMutex.Lock()
	defer scionPingMutex.Unlock()

	if !currentScionPing.InProgress {
		w.WriteHeader(http.StatusConflict)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "No scion ping in progress",
		})
		return
	}

	stoppedPID := currentScionPing.PID

	err := scionPingManager.StopProcess()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: fmt.Sprintf("Failed to stop scion ping: %v", err),
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Status:  "success",
		Message: fmt.Sprintf("Stop signal sent to scion ping process (PID: %d)", stoppedPID),
	})
}

func getScionPingStatus(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Invalid Method - use GET",
		})
		return
	}

	scionPingMutex.Lock()
	defer scionPingMutex.Unlock()

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(APIResponse{
		Status: "success",
		Data:   currentScionPing,
	})
}
func getAvailableScionPingResults(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		json.NewEncoder(w).Encode(APIResponse{
			Status:  "error",
			Message: "Method not allowed - use GET",
		})
		return
	}

	fileInfos, err := getFilesFromDirectory(scionPingResultLocation)
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
