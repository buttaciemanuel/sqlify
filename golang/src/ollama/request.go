package ollama

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

func Generate(model, prompt string) (string, error) {
	params, _ := json.Marshal(map[string]any{
		"model":  model,
		"prompt": prompt,
		"stream": false,
	})
	response, err := http.Post("http://localhost:11434/api/generate", "application/json", bytes.NewBuffer(params))

	if err != nil {
		return "", sendPromptRequestError(err)
	}

	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)

	if err != nil {
		return "", readResponseError(err)
	}

	var content struct {
		Response string `json:"response"`
	}

	if err := json.Unmarshal(body, &content); err != nil {
		return "", unpackResponseError(err)
	}

	startIndex := strings.Index(content.Response, "```sql")

	if startIndex >= 0 {
		content.Response = content.Response[startIndex+len("```sql"):]
	}

	endIndex := strings.Index(content.Response, "```")

	if endIndex >= 0 {
		content.Response = content.Response[:endIndex]
	}

	return content.Response, nil
}
