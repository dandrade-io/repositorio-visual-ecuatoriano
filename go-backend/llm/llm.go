package llm

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"reflect"
	"repository/config"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

var client openai.Client

func init() {
	if config.CAPath == "" {
		client = openai.NewClient(
			option.WithBaseURL(config.LLMBaseUrl),
			option.WithAPIKey(config.LLMApiKey),
		)

		return
	}

	pem, err := os.ReadFile(config.CAPath)
	if err != nil {
		panic(err)
	}

	pool, err := x509.SystemCertPool()
	if err != nil || pool == nil {
		pool = x509.NewCertPool()
	}

	if !pool.AppendCertsFromPEM(pem) {
		panic(fmt.Errorf("no valid certificate found in CA file"))
	}

	tlsConfig := &tls.Config{
		RootCAs: pool,
	}
	tr := &http.Transport{
		TLSClientConfig: tlsConfig,
	}

	client = openai.NewClient(
		option.WithBaseURL(config.LLMBaseUrl),
		option.WithAPIKey(config.LLMApiKey),
		option.WithHTTPClient(&http.Client{Transport: tr}),
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
