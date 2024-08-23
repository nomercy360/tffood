package handler

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"rednit/db"
	"regexp"
	"strings"
	"time"
)

type LanguageContent struct {
	AnalyzePrompt               string
	AnalyzeDescription          string
	SpamDescription             string
	DishDescription             string
	TagsDescription             string
	IngredientsDescription      string
	NutritionalPrompt           string
	NutritionalDescription      string
	IngredientNameDescription   string
	IngredientAmountDescription string
}

func getLanguageContent(language string) LanguageContent {
	if language == "ru" {
		return LanguageContent{
			AnalyzePrompt:               "Какое блюдо или продукт изображены на этой картинке?",
			AnalyzeDescription:          "Анализ изображения с едой для определения, является ли оно спамом, идентификации блюда, маркировки изображения по содержанию и перечисления ингредиентов вместе с их приблизительным количеством.",
			SpamDescription:             "Указывает, считается ли изображение спамом или не относящимся к задаче",
			DishDescription:             "Определенное основное блюдо на изображении.",
			TagsDescription:             "Теги, описывающие блюдо на фото с учетом диетических предпочтений, ингредиентов или вкусовых профилей.",
			IngredientsDescription:      "Перечислите все видимые ингредиенты и оцените приблизительное количество каждого в граммах, используя стандартные объекты на фото, такие как столовые приборы или посуда для масштабирования.",
			NutritionalPrompt:           "Анализируем информацию о питательности продукта и предоставляем данные о калориях, макронутриентах и диетической информации.",
			NutritionalDescription:      "Форматирует ответ анализа питательности в структурированный, читаемый формат для отображения или дальнейшей обработки. Эта функция организует полученные данные анализа питательности в разделы калорийности, макро- и микронутриентов, а также пригодности блюда для различных диет.",
			IngredientNameDescription:   "Название ингредиента",
			IngredientAmountDescription: "Приблизительное количество ингредиента в граммах",
		}
	}
	return LanguageContent{
		AnalyzePrompt:               "What dish or food is displayed on this picture?",
		AnalyzeDescription:          "Analyzes an image of food to determine if it's spam, identify the dish, tag the image based on its contents, and list the ingredients along with their approximate amounts.",
		SpamDescription:             "Indicates whether the image is considered spam or irrelevant to the task",
		DishDescription:             "The identified main dish in the image.",
		TagsDescription:             "Tags that describe the dish displayed in the photo based on dietary preferences, ingredients, or taste profiles.",
		IngredientsDescription:      "List all visible ingredients and estimate the approximate amount of each in grams, using standard objects in the photo such as utensils or dishware for scale.",
		NutritionalPrompt:           "Analyzing the nutritional information of the food and provide insights on the calories, macronutrients, and dietary information.",
		NutritionalDescription:      "Formats the nutritional analysis response into a structured, readable format for display or further processing. This function takes the raw data from a nutritional analysis and organizes it into sections for calories, macronutrients, micronutrients, and dietary suitability.",
		IngredientNameDescription:   "Name of the ingredient",
		IngredientAmountDescription: "Approximate amount of the ingredient in grams",
	}
}

func getRequestBody(lang, imageUrl string, caption *string) string {
	content := getLanguageContent(lang)

	var captionText string
	if caption != nil {
		captionText = fmt.Sprintf(`{
          "type": "text",
          "text": "%s"
        },`, *caption)
	}

	// replace newlines with spaces
	return fmt.Sprintf(`{
    "model": "gpt-4o-2024-08-06",
    "messages": [
        {
            "role": "system",
            "content": [
                {
                    "type": "text",
                    "text": "%s"
                }
            ]
        },
        {
            "role": "user",
            "content": [
				%s
                {
                    "type": "image_url",
                    "image_url": {
                        "url": "%s"
                    }
                }
            ]
        }
    ],
    "response_format": {
        "type": "json_schema",
        "json_schema": {
            "name": "analyze_food_image",
            "description": "%s",
            "strict": true,
            "schema": {
                "type": "object",
                "properties": {
                    "spam": {
                        "type": "boolean",
                        "description": "%s"
                    },
                    "dish": {
                        "type": "string",
                        "description": "%s"
                    },
                    "tags": {
                        "type": "array",
                        "items": {
                            "type": "string",
                            "enum": [
                                "веган",
                                "без глютена",
								"без лактозы",
								"кето",
								"палео",
								"вегетарианец",
								"белковая диета"
                            ]
                        },
                        "description": "%s"
                    },
                    "ingredients": {
                        "type": "array",
                        "items": {
                            "type": "object",
                            "properties": {
                                "name": {
                                    "type": "string",
                                    "description": "%s"
                                },
                                "amount": {
                                    "type": "number",
                                    "description": "%s"
                                }
                            },
                            "additionalProperties": false,
                            "required": ["name", "amount"]
                        },
                        "description": "%s"
                    }
                },
                 "additionalProperties": false,
                "required": [
                    "ingredients",
                    "dish",
                    "spam",
                    "tags"
                ]
            }
        }
    },
    "temperature": 0.7,
    "max_tokens": 150,
    "top_p": 1,
    "frequency_penalty": 0,
    "presence_penalty": 0
}`, content.AnalyzePrompt, captionText, imageUrl, content.AnalyzeDescription, content.SpamDescription, content.DishDescription, content.TagsDescription, content.IngredientNameDescription, content.IngredientAmountDescription, content.IngredientsDescription)
}

func nutritionRequestBody(lang, foodInfo string) string {
	content := getLanguageContent(lang)

	return fmt.Sprintf(`{
    "model": "gpt-4o-2024-08-06",
    "messages": [
        {
            "role": "system",
            "content": [
                {
                    "type": "text",
                    "text": "%s"
                }
            ]
        },
        {
            "role": "user",
            "content": [
                {
                    "type": "text",
                    "text": "%s: %s"
                }
            ]
        }
    ],
    "response_format": 
        {
            "type": "json_schema",
            "json_schema": {
                "name": "formatNutritionalResponse",
                "description": "%s",
                "strict": true,
                "schema": {
                    "type": "object",
                    "properties": {
                        "calories": {
                            "type": "number",
                            "description": "Total calories of the dish."
                        },
                        "macronutrients": {
                            "type": "object",
                            "description": "Breakdown of macronutrients in grams.",
                            "properties": {
                                "carbohydrates": {
                                    "type": "number"
                                },
                                "proteins": {
                                    "type": "number"
                                },
                                "fats": {
                                    "type": "number"
                                }
                            },
                            "additionalProperties": false,
                            "required": [
                                "carbohydrates",
                                "proteins",
                                "fats"
                            ]
                        },
                        "dietaryInformation": {
                            "type": "array",
                            "items": {
                                "type": "string"
                            },
                            "description": "Information on the suitability of the dish for various diets."
                        }
                    },
                    "additionalProperties": false,
                    "required": [
                        "calories",
                        "macronutrients",
                        "dietaryInformation"
                    ]
                }
            }
        }
    ,
    "temperature": 0.7,
    "max_tokens": 150,
    "top_p": 1,
    "frequency_penalty": 0,
    "presence_penalty": 0
}`, content.NutritionalPrompt, content.NutritionalPrompt, foodInfo, content.NutritionalDescription)
}

type OpenAIResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role    string      `json:"role"`
			Content string      `json:"content"`
			Refusal interface{} `json:"refusal"`
		} `json:"message"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	SystemFingerprint string `json:"system_fingerprint"`
}

func formatIngredients(ingredients []db.Ingredient) string {
	var formattedIngredients string

	for _, ingredient := range ingredients {
		formattedIngredients += fmt.Sprintf("Ingredient: %s, Amount: %d grams\n", ingredient.Name, int(ingredient.Amount))
	}

	formattedIngredients = strings.ReplaceAll(formattedIngredients, `"`, `\"`)

	re := regexp.MustCompile(`\r?\n`)
	formattedIngredients = re.ReplaceAllString(formattedIngredients, " ")

	return formattedIngredients
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

type ImageRecognitionResponse struct {
	DishName    string `json:"dish"`
	Ingredients []db.Ingredient
	IsSpam      bool     `json:"spam"`
	Tags        []string `json:"tags"`
}

type NutritionResponse struct {
	Calories float32 `json:"calories"`
	Macros   struct {
		Proteins float32 `json:"proteins"`
		Fats     float32 `json:"fats"`
		Carbs    float32 `json:"carbohydrates"`
	} `json:"macronutrients"`
	DietaryInfo []string `json:"dietaryInformation"`
}

func GetNutritionInfo(foodInfo string, openAIKey string) (*NutritionResponse, error) {
	log.Printf("Getting nutrition info for %s\n", foodInfo)

	reqBody := nutritionRequestBody("ru", foodInfo)

	resp, err := sendOpenAIRequest(reqBody, openAIKey)

	if err != nil {
		return nil, fmt.Errorf("failed to send OpenAI request: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in OpenAI response")
	}

	choice := resp.Choices[0]

	if choice.FinishReason == "length" {
		return nil, fmt.Errorf("unexpected finish reason: %s", choice.FinishReason)
	}

	var functionResponse NutritionResponse

	if choice.Message.Refusal != nil {
		return nil, fmt.Errorf("OpenAI refused to process the request. Here's why: %s", choice.Message.Content)
	}

	if choice.FinishReason == "content_filter" {
		return nil, fmt.Errorf("OpenAI content filter triggered")
	}

	if choice.FinishReason == "stop" {
		if err := json.Unmarshal([]byte(choice.Message.Content), &functionResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal function response: %w", err)
		}
	}

	log.Printf("Calories: %f\n", functionResponse.Calories)
	log.Printf("Proteins: %f\n", functionResponse.Macros.Proteins)
	log.Printf("Fats: %f\n", functionResponse.Macros.Fats)
	log.Printf("Carbs: %f\n", functionResponse.Macros.Carbs)
	log.Printf("Dietary Info: %v\n", functionResponse.DietaryInfo)

	return &functionResponse, nil
}

func checkImageAvailable(imgUrl string) error {
	check := func(url string) bool {
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("Failed to fetch image: %v\n", err)
			return false
		}
		defer resp.Body.Close()
		return resp.StatusCode == http.StatusOK
	}

	delays := []time.Duration{0, 1 * time.Second, 3 * time.Second}
	for _, delay := range delays {
		time.Sleep(delay)
		if check(imgUrl) {
			return nil
		}
		log.Printf("Retry after %v\n", delay)
	}

	log.Printf("Image not available after retries: %s\n", imgUrl)
	return fmt.Errorf("image not available: %s", imgUrl)
}

func GetFoodPictureInfo(imgUrl string, caption *string, openAIKey string) (*ImageRecognitionResponse, error) {
	log.Printf("Getting food picture info for %s\n", imgUrl)

	reqBody := getRequestBody("ru", imgUrl, caption)

	if err := checkImageAvailable(imgUrl); err != nil {
		return nil, err
	}

	resp, err := sendOpenAIRequest(reqBody, openAIKey)

	if err != nil {
		return nil, fmt.Errorf("failed to send OpenAI request: %w", err)
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("no choices in OpenAI response")
	}

	choice := resp.Choices[0]

	if choice.FinishReason == "length" {
		return nil, fmt.Errorf("unexpected finish reason: %s", choice.FinishReason)
	}

	var imageResponse ImageRecognitionResponse

	if choice.Message.Refusal != nil {
		return nil, fmt.Errorf("OpenAI refused to process the request. Here's why: %s", choice.Message.Content)
	}

	if choice.FinishReason == "content_filter" {
		return nil, fmt.Errorf("OpenAI content filter triggered")
	}

	if choice.FinishReason == "stop" {
		if err := json.Unmarshal([]byte(choice.Message.Content), &imageResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal function response: %w", err)
		}
	}

	log.Printf("DishName: %s\n", imageResponse.DishName)
	log.Printf("Ingredients: %v\n", imageResponse.Ingredients)
	log.Printf("Is Spam: %v\n", imageResponse.IsSpam)
	log.Printf("Tags: %v\n", imageResponse.Tags)

	return &imageResponse, nil
}
