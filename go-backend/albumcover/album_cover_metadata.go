package albumcover

import (
	"context"
	"fmt"
	"regexp"
	"repository/llm"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/openai/openai-go/v3"
)

type Metadata struct {
	Title          *string       `json:"title" jsonschema:"Album title exactly as printed on the cover, in its original language. Null if not legible."`
	Year           *int          `json:"year" jsonschema:"Release year only if printed on the cover. Null otherwise."`
	Artist         *string       `json:"artist" jsonschema:"Artist or group name as printed. Null if not legible."`
	Rhythms        []string      `json:"rhythms" jsonschema:"Musical rhythms or genres, only if printed on the cover or clearly implied. Empty array if unknown."`
	ReprTechnique  ReprTechnique `json:"representationTechnique"`
	Predominance   Predominance  `json:"graphicPredominance" jsonschema:"What dominates the cover."`
	DominantColors []string      `json:"dominantColors" jsonschema:"Only fill this when graphicPredominance is 'color' or 'image_and_color'. For 'image' or 'typography' it MUST be null. When filled: up to 3 colors covering the largest areas, ordered by area, largest first, as lowercase hex codes like #1a2b3c."`
	TitleSyntax    TitleSyntax   `json:"titleSyntax" jsonschema:"What the title relates to."`
}

var hexColor = regexp.MustCompile(`^#[0-9a-f]{6}$`)

func (m *Metadata) Normalize() error {
	if m.Predominance != Color && m.Predominance != ImageAndColor {
		m.DominantColors = nil
		return nil
	}

	for _, c := range m.DominantColors {
		if !hexColor.MatchString(c) {
			return fmt.Errorf("invalid hex color %q", c)
		}
	}

	return nil
}

type ReprTechnique string

const (
	Photography            ReprTechnique = "photography"
	Illustration           ReprTechnique = "illustration"
	Collage                ReprTechnique = "collage"
	Abstract               ReprTechnique = "abstract"
	Symbolic               ReprTechnique = "symbolic"
	TypographicComposition ReprTechnique = "typographic_composition"
	UnknownTechnique       ReprTechnique = "unknown"
)

type Predominance string

const (
	Typography          Predominance = "typography"
	ImageProdominance   Predominance = "image"
	Color               Predominance = "color"
	ImageAndColor       Predominance = "image_and_color"
	UnknownPredominance Predominance = "unknown"
)

type TitleSyntax string

const (
	Song          TitleSyntax = "song"
	Image         TitleSyntax = "image"
	Artist        TitleSyntax = "artist"
	MusicType     TitleSyntax = "music_type"
	Company       TitleSyntax = "company"
	Sentiment     TitleSyntax = "sentiment"
	UnknownSyntax TitleSyntax = "unknown"
)

func Extract(ctx context.Context, imgDataURL string) (Metadata, error) {
	schema, err := jsonschema.For[Metadata](nil)
	if err != nil {
		return Metadata{}, err
	}

	schema.Properties["representationTechnique"].Enum = []any{Photography, Illustration, Collage, Abstract, Symbolic, TypographicComposition, UnknownTechnique}
	schema.Properties["graphicPredominance"].Enum = []any{Typography, ImageProdominance, Color, ImageAndColor, UnknownPredominance}
	schema.Properties["titleSyntax"].Enum = []any{Song, Image, Artist, MusicType, Company, Sentiment, UnknownSyntax}

	userPrompts := []openai.ChatCompletionContentPartUnionParam{
		openai.TextContentPart("Analyze this album cover."),
		openai.ImageContentPart(openai.ChatCompletionContentPartImageImageURLParam{
			URL: imgDataURL,
		}),
	}

	const systemPrompt = `You analyze album cover images and fill the provided schema.
	graphicPredominance: what dominates the cover.
	- "typography": text is the main element, with no dominant image.
	- "image": any illustration, drawing, painting, collage or photograph is the main element. Small text over an image does not change this.
	- "color": flat color or gradients dominate, with little or no text.
	- "image_and_color": image and flat color or gradients together dominate, with no dominant text.
	dominantColors: only for "color" or "image_and_color"; otherwise null. Give up to 3 colors covering the largest areas, largest first, as lowercase hex codes. Ignore small details, shadows and noise.`

	metadata, err := llm.AskWithJsonSchema[Metadata](ctx, schema, systemPrompt, userPrompts)
	if err != nil {
		return Metadata{}, err
	}

	if err := metadata.Normalize(); err != nil {
		return Metadata{}, err
	}

	return metadata, nil
}
