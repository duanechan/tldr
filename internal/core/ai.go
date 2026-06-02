package core

import (
	"context"

	"google.golang.org/genai"
)

type AIModel interface {
	List(
		context.Context,
		*genai.ListModelsConfig,
	) (genai.Page[genai.Model], error)
	GenerateContent(
		context.Context,
		string,
		[]*genai.Content,
		*genai.GenerateContentConfig,
	) (*genai.GenerateContentResponse, error)
}

type AIClient struct {
	models *genai.Models
}

func NewAIClient(apiKey string) (*AIClient, error) {
	client, err := genai.NewClient(context.Background(), &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, err
	}

	return &AIClient{client.Models}, nil
}

func (a *AIClient) List(
	ctx context.Context,
	config *genai.ListModelsConfig,
) (genai.Page[genai.Model], error) {
	return a.models.List(ctx, config)
}

func (a *AIClient) GenerateContent(
	ctx context.Context,
	model string,
	contents []*genai.Content,
	config *genai.GenerateContentConfig,
) (*genai.GenerateContentResponse, error) {
	return a.models.GenerateContent(ctx, model, contents, config)
}
