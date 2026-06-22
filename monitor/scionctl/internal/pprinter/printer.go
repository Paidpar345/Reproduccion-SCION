package pprinter

// Using go pretty print
// https://github.com/jedib0t/go-pretty

// Functionalities for printing to stdout in a nice way

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"scionctl/internal/api"
	"strings"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

func PrintTable(input TableInput) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(input.header)
	t.AppendRows(input.content)
	t.Render()
}

func FileInfosTable(fileInfos []api.FileInfo) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(table.Row{"idx", "fn", "size"})

	for _, fi := range fileInfos {
		t.AppendRow(table.Row{fi.Index, fi.Name, fi.Size})
	}

	t.Render()
}

func PrintStatus(cs api.CommandState) {
	status := text.FgRed.Sprint("No")
	if cs.InProgress {
		status = text.FgGreen.Sprint("Yes")
	}

	var formattedString = fmt.Sprintf(`%s
%s: %s
%s: %d
%s: %s
%s: %s`,
		text.Bold.Sprint("Command Status"),
		text.FgBlue.Sprint("In Progress"), status,
		text.FgBlue.Sprint("PID"), cs.PID,
		text.FgBlue.Sprint("Start Time"), cs.StartTime.Format("2006-01-02 15:04:05"),
		text.FgBlue.Sprint("Output File"), cs.OutputFile,
	)

	fmt.Println(formattedString)
}

// Simple HTTP Response to Stdout
func HTTPResponseToStdout(resp *http.Response, err error) {
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Check if response is nil
	if resp == nil {
		fmt.Println("Error: nil response received")
		return
	}

	defer resp.Body.Close()

	fmt.Println("Response Status:", resp.Status)
	fmt.Println("Response Body:")

	buf := new(strings.Builder)
	_, copyErr := io.Copy(buf, resp.Body)
	if copyErr != nil {
		PrintError(copyErr)
		return
	}

	jsonTransformer := text.NewJSONTransformer("", "    ")
	prettyJSON := jsonTransformer(buf.String())
	fmt.Println(prettyJSON)

	// Add a newline for better formatting
	fmt.Println()
}

func PrintError(err error) {
	if err != nil {
		fmt.Println("Error:", err)
	}
}

// HTTPResponseRawToStdout prints the HTTP response body as raw text without JSON formatting
// Useful for YAML files, plain text, or other non-JSON content
func HTTPResponseRawToStdout(resp *http.Response, err error) {
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Check if response is nil
	if resp == nil {
		fmt.Println("Error: nil response received")
		return
	}

	defer resp.Body.Close()

	// Read the response body
	buf := new(strings.Builder)
	_, copyErr := io.Copy(buf, resp.Body)
	if copyErr != nil {
		PrintError(copyErr)
		return
	}

	// Print status and raw content
	fmt.Printf("%s: %s\n", text.FgGreen.Sprint("Response Status"), resp.Status)
	fmt.Printf("\n%s:\n\n", text.FgCyan.Sprint("Path Policy Configuration"))
	fmt.Println(buf.String())
}
