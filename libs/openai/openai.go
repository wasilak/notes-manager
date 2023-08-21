package openai

import (
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

func GetAIResponse(note db.Note) (db.Note, error) {
	content := `
	title: %s
	content: %s
	`

	c := openai.NewClient(os.Getenv("OPENAI_API_KEY"))

	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		// MaxTokens: 20,
		Messages: []openai.ChatCompletionMessage{
			{Role: "system", Content: "You are devops or cloud engineer"},
			{Role: "assistant", Content: `You will be presented with JSON document consisting of 'title', 'tags' and 'content' fields. Result has to be a only JSON document (no other text, either before or after JSON) with keys: 'title' and 'content' and 'tags'. Always preserve language from request. Response 'content' field should be an enriched, better described or simply rewritten 'content' using Markdown format. Response 'title' field should be improved as well but not in Markdown. Response 'tags' field should be a list of tags describing content and title, use current tags or propose new ones. Tags need to be lowercased, replace spaces with hyphens. Preserve links to images. Add comments to code blocks/ Go over each code block and add inline comments to relevant code lines or blocks, etc.: loops or conditions.`},
			{Role: "user", Content: fmt.Sprintf(content, note.Title, note.Content)},
		},
		Stream: false,
	}

	response, err := c.CreateChatCompletion(common.CTX, req)
	if err != nil {
		slog.ErrorContext(common.CTX, "ChatCompletion error: %v\n", err)
		return note, err
	}

	chatResponse := response.Choices[0].Message.Content

	slog.DebugContext(common.CTX, "AI response", chatResponse)
	prefix := "```json"
	suffix := "```"

	if len(chatResponse) > len(prefix) && chatResponse[:len(prefix)] == prefix {
		chatResponse = chatResponse[len(prefix):]
	}

	if len(chatResponse) > len(suffix) && chatResponse[len(chatResponse)-len(suffix):] == suffix {
		chatResponse = chatResponse[:len(chatResponse)-len(suffix)]
	}

	var AIResponse NoteAIResponse
	err = json.Unmarshal([]byte(chatResponse), &AIResponse)
	if err != nil {
		slog.ErrorContext(common.CTX, "Error decoding OpenAI response.", err)
		return note, err
	}

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

	return note, nil
}
