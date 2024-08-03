package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"imantz/daily_it_meeting_helper_ppp/internal/models"
	"net/http"
	"os"
)

const chatGPTAPIURL = "https://api.openai.com/v1/engines/davinci-codex/completions"

type ChatGPTRequest struct {
	Prompt    string `json:"prompt"`
	MaxTokens int    `json:"max_tokens"`
}

type ChatGPTResponse struct {
	Choices []struct {
		Text string `json:"text"`
	} `json:"choices"`
}

func CallChatGPT(msg models.Message) (string, error) {
	prompt := fmt.Sprintf("Progress: %s\nPlans: %s\nProblems: %s\nMake it simple, fix grammar.",
		msg.Progress, msg.Plans, msg.Problems)

	reqBody, err := json.Marshal(ChatGPTRequest{
		Prompt:    prompt,
		MaxTokens: 150,
	})
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", chatGPTAPIURL, bytes.NewBuffer(reqBody))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+os.Getenv("CHATGPT_API_KEY"))

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var chatGPTResp ChatGPTResponse
	if err := json.NewDecoder(resp.Body).Decode(&chatGPTResp); err != nil {
		return "", err
	}

	if len(chatGPTResp.Choices) > 0 {
		return chatGPTResp.Choices[0].Text, nil
	}

	return "", fmt.Errorf("no response from ChatGPT")
}
