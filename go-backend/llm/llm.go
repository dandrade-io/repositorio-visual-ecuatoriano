package llm

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"repository/config"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

var client openai.Client

func init() {
	client = openai.NewClient(
		option.WithBaseURL(config.LLMBaseUrl),
		option.WithAPIKey(config.LLMApiKey),
	)
}

func AskWithJsonSchema[T any](
	ctx context.Context,
	schema *jsonschema.Schema,
	systemPrompt string,
	userPrompts []openai.ChatCompletionContentPartUnionParam,
) (T, error) {
	var zero T

	data, err := json.Marshal(schema)
	if err != nil {
		return zero, fmt.Errorf("marshal error: %w", err)
	}

	var schemaMap map[string]any
	if err := json.Unmarshal(data, &schemaMap); err != nil {
		return zero, fmt.Errorf("unmarshal error: %w", err)
	}

	// pretty, err := json.MarshalIndent(schema, "", "  ")
	// if err != nil {
	// 	return zero, err
	// }
	// fmt.Println(string(pretty))

	completion, err := client.Chat.Completions.New(ctx, openai.ChatCompletionNewParams{
		Model: config.LLMModel,
		Messages: []openai.ChatCompletionMessageParamUnion{
			openai.SystemMessage(systemPrompt),
			openai.UserMessage(userPrompts),
		},
		ResponseFormat: openai.ChatCompletionNewParamsResponseFormatUnion{
			OfJSONSchema: &openai.ResponseFormatJSONSchemaParam{
				JSONSchema: openai.ResponseFormatJSONSchemaJSONSchemaParam{
					Name:   reflect.TypeOf(zero).Name(),
					Schema: schemaMap,
					Strict: openai.Bool(true),
				},
			},
		},
	})
	if err != nil {
		return zero, fmt.Errorf("chat/completions: %w", err)
	}

	var v T
	if err := json.Unmarshal([]byte(completion.Choices[0].Message.Content), &v); err != nil {
		return zero, fmt.Errorf("decode json: %w", err)
	}

	return v, nil
}
