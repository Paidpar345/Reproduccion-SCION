package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"scionctl/internal/api"
	"scionctl/internal/config"
	"scionctl/internal/pprinter"
	"time"
)

// Todo Docs
func HandleCaptureStart(args []string) {
	node, exists := config.CmdNodeManager.GetNode(args[0])

	if !exists {
		pprinter.PrintError(fmt.Errorf("node [%s] does not exist", args[0]))
		pprinter.PrintError(fmt.Errorf("node [%s] does not exist", args[0]))
		return
	}

	// Timeout currently after 5s, arbitrary value
	c := api.NewClient(api.ClientConfig{
		BaseURL: "http://" + node.Addr + ":" + fmt.Sprint(node.Port),
		Timeout: time.Second * 5,
	})

	// Need to make sure pingCount is set
	resp, err := c.StartCapture()
	pprinter.HTTPResponseToStdout(resp, err)
}

// Todo Docs
func HandleCaptureStop(args []string) {
	node, exist := config.CmdNodeManager.GetNode(args[0])

	if !exist {
		pprinter.PrintError(fmt.Errorf("node [%s] does not exist", args[0]))
		return
	}

	c := api.NewClient(api.ClientConfig{
		BaseURL: "http://" + node.Addr + ":" + fmt.Sprint(node.Port), Timeout: time.Second * 5,
	})

	resp, err := c.StopCapture()
	if err != nil {
		pprinter.PrintError(err)
		return
	}

	pprinter.HTTPResponseToStdout(resp, err)
}

// Needs documentation
func HandleCaptureList(args []string) {
	node, exist := config.CmdNodeManager.GetNode(args[0])

	if !exist {
		pprinter.PrintError(fmt.Errorf("node [%s] does not exist", args[0]))
		return
	}

	c := api.NewClient(api.ClientConfig{
		BaseURL: "http://" + node.Addr + ":" + fmt.Sprint(node.Port), Timeout: time.Second * 5,
	})

	resp, err := c.GetResultsCapture()
	if err != nil {
		pprinter.PrintError(err)
		return
	}

	defer resp.Body.Close()
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		pprinter.PrintError(err)
		return
	}

	var apiResp = api.FileInfosAPIResponse{}
	err = json.Unmarshal(bodyBytes, &apiResp)
	if err != nil {
		pprinter.PrintError(err)
		return
	}
	// We only care about the Data part
	if resp.StatusCode != http.StatusOK {
		pprinter.PrintError(fmt.Errorf("error: %s", apiResp.Message))
		return
	}

	fmt.Printf("Available ping logfiles on %s:\n\n", node.Name)

	fileInfos := apiResp.Data

	pprinter.FileInfosTable(fileInfos)
}

// Needs DOC
func HandleCaptureStatus(args []string) {
	node, exist := config.CmdNodeManager.GetNode(args[0])

	if !exist {
		pprinter.PrintError(fmt.Errorf("node [%s] does not exist", args[0]))
		return
	}

	c := api.NewClient(api.ClientConfig{
		BaseURL: "http://" + node.Addr + ":" + fmt.Sprint(node.Port), Timeout: time.Second * 5,
	})

	resp, err := c.StatusCapture()
	if err != nil {
		pprinter.PrintError(err)
		return
	}

	defer resp.Body.Close()

	// Unpack body into PingStatusAPIResponse
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		pprinter.PrintError(err)
		return
	}

	var apiResp = api.StatusAPIResponse{}
	err = json.Unmarshal(bodyBytes, &apiResp)
	if err != nil {
		pprinter.PrintError(err)
		return
	}
	// We only care about the Data part
	if resp.StatusCode != http.StatusOK {
		pprinter.PrintError(fmt.Errorf("error: %s", apiResp.Message))
		return
	}

	pprinter.PrintStatus(apiResp.Data)

}
