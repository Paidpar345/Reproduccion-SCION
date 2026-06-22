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

func HandleScionPingStart(args []string, count int) {
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
	// Format to correct scion ping addr e.g. 16-ffaa:1:1,127.0.0.1
	resp, err := c.StartScionPing(api.NodeToScionAddress(&pingedNode), count)
	pprinter.HTTPResponseToStdout(resp, err)
}

func HandleScionPingStop(args []string) {
	node, exist := config.CmdNodeManager.GetNode(args[0])

	if !exist {
		pprinter.PrintError(fmt.Errorf("node [%s] does not exist", args[0]))
		return
	}

	c := api.NewClient(api.ClientConfig{
		BaseURL: "http://" + node.Addr + ":" + fmt.Sprint(node.Port), Timeout: time.Second * 5,
	})

	resp, err := c.StopScionPing()
	if err != nil {
		pprinter.PrintError(err)
		return
	}

	pprinter.HTTPResponseToStdout(resp, err)
}

// Lots of duplicate code i guess
func HandleScionPingStatus(args []string) {
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

func HandleScionPingList(args []string) {
	node, exist := config.CmdNodeManager.GetNode(args[0])

	if !exist {
		pprinter.PrintError(fmt.Errorf("node [%s] does not exist", args[0]))
		return
	}

	c := api.NewClient(api.ClientConfig{
		BaseURL: "http://" + node.Addr + ":" + fmt.Sprint(node.Port), Timeout: time.Second * 5,
	})

	resp, err := c.ScionGetResultsPing()
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
