package main

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/anthropics/anthropic-sdk-go"
)

type llmMessage struct {
	role string `json:"role"`
	data string `json:"user"`
}

func main() {
	sentence := ""
	results := make([]string, 0)
	fmt.Println("-> Request:")
	scanner := bufio.NewScanner(os.Stdin)
	if scanner.Scan() {
		sentence = scanner.Text()
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading input:", err)
	}

	// Initialize configs value
	initializeConfigs()
	clientLLM := initClientLLM()

	maxLoop := getMaxLoop()
	if maxLoop <= 0 {
		panic("maxLoop cannot be less than 0")
	}

	messages := make([]llmMessage, 0)
	idx := 1
	messages = append(messages, llmMessage{"user", addContext(sentence)})
	for idx <= maxLoop {
		fmt.Println(fmt.Sprintf("<- %d ->", idx))

		response := LLMRequester(clientLLM, messages)
		if strings.Contains(response, "FINAL:") {
			fmt.Println(fmt.Sprintf("-> Response: %s", strings.TrimPrefix(response, "FINAL:")))
			idx = 100
			break
		}
		auditResponse(response, idx)
		messages = append(messages, llmMessage{"assistant", response})

		result, err := LLMCoder(response)
		if err != nil {
			results = append(results, result)
		}
		results = append(results, result)

		fmt.Println(fmt.Sprintf("<- %d -> Result: %s", idx, result))
		messages = append(messages, llmMessage{"user", addObservation(addContext(sentence), results)})
		idx++

	}
}

func addObservation(sentence string, results []string) string {
	resultSummary := summaryList(results)
	if len(results) == 0 {
		replacer := strings.NewReplacer("{observations}", "Sin resultados obtenidos")
		sentence = replacer.Replace(sentence)
	}
	replacer := strings.NewReplacer("{observations}", resultSummary)
	return replacer.Replace(sentence)
}

func summaryList(list []string) string {
	summary := ""
	if len(list) > 0 {
		for i, obj := range list {
			summary = fmt.Sprintf("%s - Step %d: %s \n", summary, i+1, obj)
		}
		return fmt.Sprintf("\n%s", summary)
	}
	return ""
}

func LLMRequester(client anthropic.Client, messages []llmMessage) string {
	messagesClaude := make([]anthropic.MessageParam, 0)

	for i := 0; i < len(messages); i++ {
		msg := messages[i]
		switch msg.role {
		case "user":
			messagesClaude = append(messagesClaude, anthropic.NewUserMessage(anthropic.NewTextBlock(msg.data)))
		case "assistant":
			messagesClaude = append(messagesClaude, anthropic.NewAssistantMessage(anthropic.NewTextBlock(msg.data)))
		}
	}
	message, err := client.Messages.New(context.TODO(), anthropic.MessageNewParams{
		Model:     anthropic.ModelClaudeHaiku4_5,
		MaxTokens: int64(1024),
		Messages:  messagesClaude,
	})

	if err != nil {
		panic(err)
	}
	msgClient := ""
	if message != nil && len(message.Content) > 0 {
		msgClient = message.Content[0].Text
	}
	return msgClient
}

func LLMCoder(code string) (string, error) {
	defer deleteTempFile()
	var stout, stderr bytes.Buffer
	fileName := fmt.Sprintf("./%s.go", "tmp")
	err := os.WriteFile(fileName, []byte(extractGoCode(code)), 0644)
	if err != nil {
		fmt.Printf("Error creating: %s\n", err)
	}
	err = exec.Command("go", "mod", "tidy").Run()
	if err != nil {
		fmt.Printf("Error ejecutando tidy: %s\n", err)
		return "", err
	}
	cmd := exec.Command("go", "run", fileName)
	cmd.Stderr = &stderr
	cmd.Stdout = &stout
	err = cmd.Run()
	if err != nil {
		fmt.Printf("Error executing: %s\n", err)
		return stderr.String(), err
	}
	return stout.String(), nil
}

func deleteTempFile() {
	err := exec.Command("rm", fmt.Sprintf("./%s.go", "tmp")).Run()
	if err != nil {
		fmt.Printf("Error deleting: %s\n", err)
		panic(err)
	}
}

func auditResponse(code string, step int) {
	if os.Getenv("AUDIT-RESPONSE") == "true" {
		fileName, _ := namedTempFile()
		err := os.MkdirAll("./agent/tmp", 0755)
		if err != nil {
			fmt.Printf("Error audit: %s\n", err)
		}
		err = os.WriteFile(fmt.Sprintf("./agent/tmp/%d_%s.txt", step, fileName), []byte(code), 0644)
		if err != nil {
			fmt.Printf("Error creating: %s\n", err)
		}
	}
}
