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
            "name": "analyzeFoodImage",
            "description": "Analyzes an image of food to determine if it's spam, identify the dish, tag the image based on its contents, and list the ingredients along with their approximate amounts.",
            "parameters": {
                "type": "object",
                "properties": {
                    "spam": {
                        "type": "boolean",
                        "description": "Indicates whether the image is considered spam or irrelevant to the task"
                    },
                    "dish": {
                        "type": "string",
                        "description": "The identified main dish in the image."
                    },
                    "tags": {
                        "type": "array",
                        "items": {
                            "type": "string",
                            "enum": [
                                "vegetarian",
                                "vegan",
                                "gluten-free",
                                "meat",
                                "seafood",
                                "dessert",
                                "spicy",
                                "sweet",
                                "salty"
                            ]
                        },
                        "description": "Tags that describe the dish displayed in the photo based on dietary preferences, ingredients, or taste profiles."
                    },
                    "ingredients": {
                        "type": "array",
                        "items": {
                            "type": "object",
                            "properties": {
                                "name": {
                                    "type": "string",
                                    "description": "Name of the ingredient"
                                },
                                "amount": {
                                    "type": "number",
                                    "description": "Approximate amount of the ingredient in grams"
                                }
                            },
                            "required": ["name", "amount"]
                        },
                        "description": "List all visible ingredients and estimate the approximate amount of each in grams, using standard objects in the photo such as utensils or dishware for scale."
                    }
                },
                "required": [
                    "ingredients",
                    "dish",
                    "spam",
                    "tags"
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

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, body)
	}

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
	DishName    string `json:"dish"`
	Ingredients []struct {
		Name   string  `json:"name"`
		Amount float64 `json:"amount"`
	}
	IsSpam bool     `json:"spam"`
	Tags   []string `json:"tags"`
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
		if toolCall.Function.Name == "analyzeFoodImage" {
			if err := json.Unmarshal([]byte(toolCall.Function.Arguments), &functionResponse); err != nil {
				return nil, fmt.Errorf("failed to unmarshal function response: %w", err)
			}
		}
	}

	log.Printf("DishName: %s\n", functionResponse.DishName)
	log.Printf("Ingredients: %v\n", functionResponse.Ingredients)
	log.Printf("Is Spam: %v\n", functionResponse.IsSpam)
	log.Printf("Tags: %v\n", functionResponse.Tags)

	return &functionResponse, nil
}
