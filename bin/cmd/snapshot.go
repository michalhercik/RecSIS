package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/spf13/cobra"
)

var snapshotCmd = &cobra.Command{
	Use:     "snap",
	Aliases: []string{"snapshot"},
	Short:   "Create meilisearch snapshot",
	Long:    "Create meilisearch snapshot",
	Run: func(cmd *cobra.Command, args []string) {
		// fmt.Printf("Addition of %s and %s = %s.\n\n", args[0], args[1], Add(args[0], args[1]))
		fmt.Println("Creating meilisearch snapshot...")
		res, err := snapshot()
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		fmt.Println("Response:", string(res))
	},
}

const noID = -1

var tasksCmd = &cobra.Command{
	Use:   "tasks",
	Short: "List meilisearch tasks",
	Long:  "List meilisearch tasks",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		id := noID
		if len(args) == 1 {
			var err error
			id, err = strconv.Atoi(args[0])
			if err != nil {
				fmt.Println("Invalid task ID:", args[0], err)
				return
			}
		}
		res, err := tasks(id)
		if err != nil {
			fmt.Println("Error:", err)
			return
		}
		prettyPrintJSON(res)
	},
}

func init() {
	rootCmd.AddCommand(snapshotCmd)
	rootCmd.AddCommand(tasksCmd)
}

func snapshot() ([]byte, error) {
	url := "http://localhost:7700/snapshots"
	token := "MASTER_KEY"
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("status: %s", resp.Status)
	}
	return body, nil
}

func tasks(id int) ([]byte, error) {
	url := "http://localhost:7700/tasks"
	if id > 0 {
		url += fmt.Sprintf("/%d", id)
	} else {
		url += "?reverse=true&afterEnqueuedAt=2025-10-03T00:00:00Z"
	}
	token := "MASTER_KEY"
	// reqBody := `{
	// 	"reverse": true
	// }`
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", "Bearer "+token)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		return body, fmt.Errorf("status: %s", resp.Status)
	}
	return body, nil
}

func prettyPrintJSON(data []byte) {
	var out bytes.Buffer
	err := json.Indent(&out, data, "", "  ")
	if err != nil {
		fmt.Println("Invalid JSON:", err)
		fmt.Println(string(data))
		return
	}
	fmt.Println(out.String())
}
