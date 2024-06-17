package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

func getRequestBody(imageUrl string) string {
	// replace newlines with spaces
	return fmt.Sprintf(`{
    "model": "gpt-4o",
    "messages": [
        {
            "role": "system",
            "content": [
                {
                    "type": "text",
                    "text": "What dish or food is displayed on this picture?"
                }
            ]
        },
        {
            "role": "user",
            "content": [
                {
                    "type": "image_url",
                    "image_url": {
                        "url": "%s"
                    }
                }
            ]
        }
    ],
    "tools": [{
        "type": "function",
        "function": {
            "name": "DecomposeLaunch",
            "description": "Name dish and ingredients of food displayed on photo",
            "parameters": {
                "type": "object",
                "properties": {
					"spam": {
						"type": "boolean",
						"description": "True if photo is not related to food"
					},
                    "dish": {
                        "type": "string",
                        "description": "Short-name of dish displayed on photo"
                    },
                    "ingredients": {
                        "type": "array",
                        "items": {
                            "type": "string"
                        },
                        "description": "Name of all edible ingridient displayed on photo"
                    }
                },
                "required": [
                    "ingredients",
                    "dish",
					"spam"
                ]
            }
        }
    }],
    "temperature": 0.7,
    "max_tokens": 150,
    "top_p": 1,
    "frequency_penalty": 0,
    "presence_penalty": 0
}`, imageUrl)
}

type OpenAIResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role      string `json:"role"`
			Content   string `json:"content"`
			ToolCalls []struct {
				ID       string `json:"id"`
				Type     string `json:"type"`
				Function struct {
					Name      string `json:"name"`
					Arguments string `json:"arguments"`
				} `json:"function"`
			} `json:"tool_calls"`
		} `json:"message"`
		FinishReason string `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

func sendOpenAIRequest(reqBody string, token string) (*OpenAIResponse, error) {

	client := &http.Client{}

	req, err := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", strings.NewReader(reqBody))

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))

	resp, err := client.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var openAIResponse OpenAIResponse

	if err := json.Unmarshal(body, &openAIResponse); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &openAIResponse, nil
}

type FunctionResponse struct {
	DishName    string   `json:"dish"`
	Ingredients []string `json:"ingredients"`
	IsSpam      bool     `json:"spam"`
}

func GetFoodPictureInfo(imgUrl string, openAIKey string) (*FunctionResponse, error) {
	log.Printf("Getting food picture info for %s\n", imgUrl)

	reqBody := getRequestBody(imgUrl)

	resp, err := sendOpenAIRequest(reqBody, openAIKey)

	if err != nil {
		return nil, fmt.Errorf("failed to send OpenAI request: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in OpenAI response")
	}

	choice := resp.Choices[0]

	if choice.FinishReason != "tool_calls" {
		return nil, fmt.Errorf("unexpected finish reason: %s", choice.FinishReason)
	}

	var functionResponse FunctionResponse

	for _, toolCall := range choice.Message.ToolCalls {
		if toolCall.Function.Name == "DecomposeLaunch" {
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &functionResponse); err != nil {
				return nil, fmt.Errorf("failed to unmarshal function response: %w", err)
			}
		}
	}

	log.Printf("DishName: %s\n", functionResponse.DishName)
	log.Printf("Ingredients: %v\n", functionResponse.Ingredients)
	log.Printf("Is Spam: %v\n", functionResponse.IsSpam)

	return &functionResponse, nil
}
