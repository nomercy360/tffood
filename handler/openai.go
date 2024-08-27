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
	IngredientNameDescription   string
	IngredientListDescription   string
	IngredientAmountDescription string
	MacroNutrientsDescription   string
	CaloriesDescription         string
	NutritionIngredientWeight   string
	Tags                        string
}

func getLanguageContent(language string) LanguageContent {
	if language == "ru" {
		return LanguageContent{
			AnalyzePrompt:               "Какое блюдо или продукт изображены на этой картинке?",
			AnalyzeDescription:          "Анализ изображения с едой для определения, является ли оно спамом, идентификации блюда, маркировки изображения по содержанию и перечисления ингредиентов вместе с их приблизительным количеством.",
			SpamDescription:             "Указывает, считается ли изображение спамом или не относящимся к задаче",
			DishDescription:             "Определенное основное блюдо на изображении.",
			TagsDescription:             "Теги, описывающие блюдо на фото с учетом диетических предпочтений, ингредиентов или вкусовых профилей.",
			IngredientsDescription:      "Перечисли все видимые ингредиенты и оцените приблизительное количество каждого в граммах, используя стандартные объекты на фото, такие как столовые приборы или посуда для масштабирования.",
			NutritionalPrompt:           "Проанализируй информацию о питательности продукта и предоставь данные о калориях, макронутриентах и диетической информации.",
			IngredientNameDescription:   "Название ингредиента",
			IngredientAmountDescription: "Приблизительное количество ингредиента в граммах",
			IngredientListDescription:   "Список ингредиентов с их питательной информацией.",
			MacroNutrientsDescription:   "Разбивка макронутриентов в граммах для этого ингредиента.",
			CaloriesDescription:         "Калории для этого ингредиента.",
			NutritionIngredientWeight:   "Вес ингредиента в граммах",
			Tags:                        "\"веган\", \"без глютена\", \"богатый белком\", \"низкоуглеводный\", \"палео\", \"без лактозы\", \"вегетарианский\", \"без сахара\", \"низкожирный\", \"средиземноморский\", \"богатый клетчаткой\"",
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
		IngredientNameDescription:   "Name of the ingredient",
		IngredientAmountDescription: "Approximate amount of the ingredient in grams",
		MacroNutrientsDescription:   "Breakdown of macronutrients in grams for this ingredient.",
		IngredientListDescription:   "List of ingredients with their nutritional information.",
		CaloriesDescription:         "Calories for this ingredient.",
		NutritionIngredientWeight:   "Weight of the ingredient in grams",
		Tags:                        "\"vegan\", \"gluten-free\", \"high-protein\", \"low-carb\", \"paleo\", \"dairy-free\", \"vegetarian\", \"sugar-free\", \"low-fat\", \"mediterranean\", \"high-fiber\"",
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
            "name": "food_image_analysis",
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
                            "enum": [%s]
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
    "max_tokens": 200,
    "top_p": 1,
    "frequency_penalty": 0,
    "presence_penalty": 0
}`,
		content.AnalyzePrompt,
		captionText, imageUrl,
		content.AnalyzeDescription,
		content.SpamDescription,
		content.DishDescription,
		content.Tags,
		content.TagsDescription,
		content.IngredientNameDescription,
		content.IngredientAmountDescription,
		content.IngredientsDescription,
	)
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
                    "text": "%s"
                }
            ]
        }
    ],
    "response_format": 
        {
            "type": "json_schema",
            "json_schema": {
				"name": "nutrition_info",
                "strict": true,
                "schema": {
                    "type": "object",
                    "properties": {
                        "ingredients": {
                            "type": "array",
                            "description": "%s",
                            "items": {
                                "type": "object",
                                "properties": {
                                    "name": {
                                        "type": "string",
                                        "description": "%s"
                                    },
                                    "calories": {
                                        "type": "number",
                                        "description": "%s"
                                    },
									"weight": {
										"type": "number",
										"description": "%s"	
									},
                                    "macronutrients": {
                                        "type": "object",
                                        "description": "%s",
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
                                    }
                                },
                                "additionalProperties": false,
                                "required": [
                                    "name",
                                    "calories",
                                    "macronutrients",
									"weight"
                                ]
                            }
                        }
                    },
                    "additionalProperties": false,
                    "required": [
                        "ingredients"
                    ]
                }
            }
        }
    ,
    "temperature": 0.7,
    "max_tokens": 300,
    "top_p": 1,
    "frequency_penalty": 0,
    "presence_penalty": 0
}`,
		content.NutritionalPrompt,
		foodInfo,
		content.IngredientListDescription,
		content.IngredientNameDescription,
		content.CaloriesDescription,
		content.NutritionIngredientWeight,
		content.MacroNutrientsDescription,
	)
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

type Ingredient struct {
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

func formatIngredients(lang string, ingredients []Ingredient) string {
	var formattedIngredients string

	for _, ingredient := range ingredients {
		if lang == "ru" {
			formattedIngredients += fmt.Sprintf("Ингредиент: %s, Количество: %d грамм.", ingredient.Name, int(ingredient.Amount))
			continue
		} else if lang == "en" {
			formattedIngredients += fmt.Sprintf("Ingredient: %s, Amount: %d grams.", ingredient.Name, int(ingredient.Amount))
			continue
		}
	}

	formattedIngredients = strings.ReplaceAll(formattedIngredients, `"`, `\"`)

	re := regexp.MustCompile(`\r?\n`)
	formattedIngredients = re.ReplaceAllString(formattedIngredients, " ")
	// last space is not needed
	formattedIngredients = formattedIngredients[:len(formattedIngredients)-1]

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
	DishName    string       `json:"dish"`
	Ingredients []Ingredient `json:"ingredients"`
	IsSpam      bool         `json:"spam"`
	Tags        []string     `json:"tags"`
}

type NutritionResponse struct {
	Ingredients []db.Ingredient `json:"ingredients"`
}

func getNutritionInfo(lang, foodInfo string, openAIKey string) (*NutritionResponse, error) {
	log.Printf("Getting nutrition info for %s\n", foodInfo)

	reqBody := nutritionRequestBody(lang, foodInfo)

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

	log.Printf("Nutrition info: %v\n", functionResponse)

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

func getFoodPictureInfo(lang, imgUrl string, caption *string, openAIKey string) (*ImageRecognitionResponse, error) {
	log.Printf("Getting food picture info for %s\n", imgUrl)

	reqBody := getRequestBody(lang, imgUrl, caption)

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
