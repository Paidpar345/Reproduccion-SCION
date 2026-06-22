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

// Parse input to IPs
//
// args[0] is the source node as ScionNode type
//
// args[1] is the destination node as ScionNode type
//
// Example usage: scionctl ping start scion11 scion31 -c 5
func HandlePingStart(args []string, count int) {
	pingerNode, pingedNode, err := validateNodes(args)
	if err != nil {
		pprinter.PrintError(err)
		return
	}

	// Timeout currently after 5s, arbitrary value
	c := api.NewClient(api.ClientConfig{
		BaseURL: "http://" + pingerNode.Addr + ":" + fmt.Sprint(pingerNode.Port),
		Timeout: time.Second * 5,
	})

	// Need to make sure pingCount is set
	resp, err := c.StartPing(pingedNode.Addr, count)
	pprinter.HTTPResponseToStdout(resp, err)
}

func validateNodes(args []string) (api.ScionNode, api.ScionNode, error) {
	pingerNode, exists := config.CmdNodeManager.GetNode(args[0])
	if !exists {
		return api.ScionNode{}, api.ScionNode{}, fmt.Errorf("sender node [%s] does not exist", args[0])
	}
	pingedNode, exists := config.CmdNodeManager.GetNode(args[1])
	if !exists {
		return api.ScionNode{}, api.ScionNode{}, fmt.Errorf("receiver node [%s] does not exist", args[1])
	}

	return *pingerNode, *pingedNode, nil
}

// HandlePingStop handles the logic for stopping a ping operation
//
// args: [node]
func HandlePingStop(args []string) {
	node, exist := config.CmdNodeManager.GetNode(args[0])

	if !exist {
		pprinter.PrintError(fmt.Errorf("node [%s] does not exist", args[0]))
		return
	}

	c := api.NewClient(api.ClientConfig{
		BaseURL: "http://" + node.Addr + ":" + fmt.Sprint(node.Port), Timeout: time.Second * 5,
	})

	resp, err := c.StopPing()
	if err != nil {
		pprinter.PrintError(err)
		return
	}

	pprinter.HTTPResponseToStdout(resp, err)
}

// Needs documentation
func HandlePingList(args []string) {
	node, exist := config.CmdNodeManager.GetNode(args[0])

	if !exist {
		pprinter.PrintError(fmt.Errorf("node [%s] does not exist", args[0]))
		return
	}

	c := api.NewClient(api.ClientConfig{
		BaseURL: "http://" + node.Addr + ":" + fmt.Sprint(node.Port), Timeout: time.Second * 5,
	})

	resp, err := c.GetResultsPing()
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
func HandlePingStatus(args []string) {
	node, exist := config.CmdNodeManager.GetNode(args[0])

	if !exist {
		pprinter.PrintError(fmt.Errorf("node [%s] does not exist", args[0]))
		return
	}

	c := api.NewClient(api.ClientConfig{
		BaseURL: "http://" + node.Addr + ":" + fmt.Sprint(node.Port), Timeout: time.Second * 5,
	})

	resp, err := c.StatusPing()
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
