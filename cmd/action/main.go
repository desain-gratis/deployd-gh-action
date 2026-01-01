package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

func env(key string, required bool) string {
	v := os.Getenv(key)
	if v == "" && required {
		fmt.Fprintf(os.Stderr, "missing env %s\n", key)
		os.Exit(1)
	}
	return v
}

func main() {
	url := env("INPUT_URL", true)
	namespace := env("INPUT_NAMESPACE", true)
	name := env("INPUT_NAME", true)

	eventPath := env("GITHUB_EVENT_PATH", true)
	eventData, err := os.ReadFile(eventPath)
	if err != nil {
		panic(err)
	}

	var eventPayload map[string]any
	_ = json.Unmarshal(eventData, &eventPayload)

	// workflow dispatch
	payload := map[string]any{
		"namespace":    namespace,
		"name":         name,
		"commit_id":    env("GITHUB_SHA", true),
		"branch":       "",
		"actor":        env("GITHUB_ACTOR", true),
		"tag":          "",
		"data":         json.RawMessage(eventData),
		"messsage":     "",
		"published_at": time.Now(),
		"source":       "github",
		"os_arch": []string{
			"linux/amd64",
		},
	}

	if env("GITHUB_REF_TYPE", true) == "branch" {
		payload["branch"] = env("GITHUB_REF_NAME", true)
	} else if env("GITHUB_REF_TYPE", true) == "tag" {
		payload["tag"] = env("GITHUB_REF_NAME", true)
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(body))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Artifactd-Namespace", namespace)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "github-action-go")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		b, _ := io.ReadAll(resp.Body)
		panic(fmt.Sprintf("failed (%d): %s", resp.StatusCode, b))
	}

	msg, _ := io.ReadAll(resp.Body)
	fmt.Println("Metadata sent successfully: %v", string(msg))
}
