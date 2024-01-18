package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"

	openai "github.com/sashabaranov/go-openai"
	"github.com/wasilak/notes-manager/libs/common"
	"github.com/wasilak/notes-manager/libs/providers/db"
)

type NoteAIResponse struct {
	Content string   `bson:"content" json:"content"`
	Title   string   `bson:"title" json:"title"`
	Tags    []string `bson:"tags" json:"tags,omitempty"`
}

func GetAIResponseInstruct(ctx context.Context, note db.Note) (db.Note, error) {
	ctx, span := common.Tracer.Start(ctx, "GetAIResponseInstruct")
	chatRequest := fmt.Sprintf(`Rewrite this article in more descriptive and human friendly way with examples using markdown: %s. Content must be in markdown.
	Write title and tags to generated article. Do not use markdown for title and tags.
	Format response as valid RFC8259 compliant JSON document with 'content', 'title' and 'tags' fields.
	Do not include any explanations, only provide a  RFC8259 compliant JSON response  following this format without deviation.
	
	Example of response:

	{
		"content": "Article **content**",
		"title": "This is awesome title",
		"tags": ["tag1", "tag2"]
	}
	`, note.Content)

	req := openai.CompletionRequest{
		Model:     openai.GPT3Dot5TurboInstruct,
		Prompt:    chatRequest,
		MaxTokens: 3000,
		Echo:      false,
		Stream:    false,
		// Temperature: 1,
	}

	c := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	ctx, spanCreateCompletion := common.Tracer.Start(ctx, "CreateCompletion")
	response, err := c.CreateCompletion(context.TODO(), req)
	if err != nil {
		return db.Note{}, err
	}
	spanCreateCompletion.End()

	chatResponse := response.Choices[0].Text

	ctx, spanUnmarshal := common.Tracer.Start(ctx, "Unmarshal")
	var AIResponse NoteAIResponse
	err = json.Unmarshal([]byte(chatResponse), &AIResponse)
	if err != nil {
		slog.ErrorContext(ctx, "Error decoding OpenAI response.", err)
		return note, err
	}
	spanUnmarshal.End()

	containsAIGenerated := false
	for _, tag := range AIResponse.Tags {
		if tag == "ai-generated" {
			containsAIGenerated = true
			break
		}
	}

	if !containsAIGenerated {
		AIResponse.Tags = append(AIResponse.Tags, "ai-generated")
	}

	note.Title = AIResponse.Title
	note.Content = AIResponse.Content
	note.Tags = AIResponse.Tags

	span.End()
	return note, nil
}

func GetAIResponse(ctx context.Context, note db.Note) (db.Note, error) {
	ctx, span := common.Tracer.Start(ctx, "GetAIResponse")
	b, err := json.MarshalIndent(note, "", "  ")
	if err != nil {
		return note, err
	}

	c := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	req := openai.ChatCompletionRequest{
		Model:     openai.GPT3Dot5Turbo,
		MaxTokens: 3000,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: `I want you to act as an API. 
			I will send JSON documents having title, content and tags fields and you will reply with what JSON having only following fields: content, title and tags.
			Response 'content' field should be an enriched, better described or rewritten in more descriptive and human friendly way with examples in Markdown format.
			Response 'title' field should be improved as well but not in Markdown. 
			Response 'tags' field should be a list of tags describing content and title, use current tags or propose new ones. Tags need to be lowercased, replace spaces with hyphens.
			Preserve links to images. Add comments to code blocks. Go over each code block and add inline comments to relevant code lines or blocks, etc.: loops or conditions.
			Move title from content to title field. You will generate tags describing new content and title and place them as an array in tags field. 
			Format response as valid RFC8259 compliant JSON document with 'content', 'title' and 'tags' fields.
			Do not include any explanations, only provide a  RFC8259 compliant JSON response  following this format without deviation.

			Example of response:

			{
				"content": "Article **content**",
				"title": "This is awesome title",
				"tags": ["tag1", "tag2"]
			}
			`},
			{Role: "user", Content: string(b)},
		},
		Stream: false,
	}

	ctx, spanCreateChatCompletion := common.Tracer.Start(ctx, "CreateChatCompletion")
	response, err := c.CreateChatCompletion(ctx, req)
	if err != nil {
		slog.ErrorContext(ctx, "ChatCompletion error: %v\n", err)
		return note, err
	}
	spanCreateChatCompletion.End()

	chatResponse := response.Choices[0].Message.Content

	slog.DebugContext(ctx, "AI response", "chatResponse", chatResponse)
	// prefix := "```json"
	// suffix := "```"

	// if len(chatResponse) > len(prefix) && chatResponse[:len(prefix)] == prefix {
	// 	chatResponse = chatResponse[len(prefix):]
	// }

	// if len(chatResponse) > len(suffix) && chatResponse[len(chatResponse)-len(suffix):] == suffix {
	// 	chatResponse = chatResponse[:len(chatResponse)-len(suffix)]
	// }

	ctx, spanUnmarshal := common.Tracer.Start(ctx, "Unmarshal")
	var AIResponse NoteAIResponse
	err = json.Unmarshal([]byte(chatResponse), &AIResponse)
	if err != nil {
		slog.ErrorContext(ctx, "Error decoding OpenAI response.", err)
		return note, err
	}
	spanUnmarshal.End()

	containsAIGenerated := false
	for _, tag := range AIResponse.Tags {
		if tag == "ai-generated" {
			containsAIGenerated = true
			break
		}
	}

	if !containsAIGenerated {
		AIResponse.Tags = append(AIResponse.Tags, "ai-generated")
	}

	note.Title = AIResponse.Title
	note.Content = AIResponse.Content
	note.Tags = AIResponse.Tags

	span.End()
	return note, nil
}
