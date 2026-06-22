package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"scionctl/internal/api"
	"scionctl/internal/config"
	"scionctl/internal/pprinter"
	"time"
)

func HandleGetFile(args []string, src string) {
	node, exist := config.CmdNodeManager.GetNode(args[0])

	if !exist {
		pprinter.PrintError(fmt.Errorf("node [%s] does not exist", args[0]))
		return
	}

	c := api.NewClient(api.ClientConfig{
		BaseURL: "http://" + node.Addr + ":" + fmt.Sprint(node.Port), Timeout: time.Second * 5,
	})

	// Todo Check args[1] too
	resp, err := c.GetFile(args[1], src)

	if err != nil {
		pprinter.PrintError(fmt.Errorf("failed to get file: %v", err))
		return
	}
	defer resp.Body.Close()

	// Check for successful response
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		pprinter.PrintError(fmt.Errorf("server error: %s - %s", resp.Status, string(body)))
		return
	}

	fileName := args[1]
	if src == "capture" {
		fileName += ".pcap"
	} else {
		fileName += ".log"
	}

	// copy
	_, err = io.Copy(os.Stdout, resp.Body)
	if err != nil {
		pprinter.PrintError(fmt.Errorf("failed to write file: %v", err))
		return
	}
}
