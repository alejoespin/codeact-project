package main

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
)

func initializeConfigs() {
	file, err := os.Open("./agent/configs.env")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		register := strings.Split(scanner.Text(), "=")
		if len(register) == 2 {
			os.Setenv(register[0], register[1])
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func initClientLLM() anthropic.Client {
	client := anthropic.NewClient(
		option.WithAPIKey(os.Getenv("ANTHROPIC_KEY")),
	)
	return client
}

func getMaxLoop() int {
	var err error
	maxLoop := -1

	maxLoopStr := os.Getenv("LOOP-MAX")
	if maxLoopStr != "" {
		maxLoop, err = strconv.Atoi(maxLoopStr)
		if maxLoop <= 0 && err != nil {
			maxLoop = 0
		}
	}
	return maxLoop
}

func addContext(prompt string) string {
	basePrompt, err := os.ReadFile("agent/base_prompt.md")
	if err != nil {
		panic(err)
	}
	contextData, err := os.ReadFile("agent/context.md")
	if err != nil {
		panic(err)
	}
	replacer := strings.NewReplacer(
		"{user_request}", prompt,
		"{context}", string(contextData))
	return replacer.Replace(string(basePrompt))
}

func extractGoCode(input string) string {
	lines := strings.Split(input, "\n")
	start := -1
	end := -1

	for i, line := range lines {
		trimmed := strings.TrimSpace(line)
		if start == -1 && trimmed == "package main" {
			start = i
		}
		if start != -1 && trimmed == "```" {
			end = i
			break
		}
	}

	if start == -1 || end == -1 {
		return ""
	}

	return strings.Join(lines[start:end], "\n")
}

func namedTempFile() (string, error) {
	bytes := make([]byte, 4)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
