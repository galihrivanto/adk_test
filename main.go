package main

import (
	"context"
	"log"
	"os"

	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
	"google.golang.org/adk/agent"
	"google.golang.org/adk/agent/llmagent"
	"google.golang.org/adk/cmd/launcher"
	"google.golang.org/adk/cmd/launcher/full"
	openaiModel "google.golang.org/adk/model/openai"
)

func main() {
	ctx := context.Background()

	client := openai.NewClient(
		option.WithBaseURL("https://openrouter.ai/api/v1"),
		option.WithAPIKey(os.Getenv("OPENROUTER_API_KEY")),
	)

	log.Println("model:", os.Getenv("OPENROUTER_MODEL"))

	model, err := openaiModel.NewModel(ctx, os.Getenv("OPENROUTER_MODEL"), &client)
	if err != nil {
		log.Fatalf("Failed to create model: %v", err)
	}

	a, err := llmagent.New(llmagent.Config{
		Name:        "helpful agent",
		Model:       model,
		Description: "Agent to answer questions about Malang City in Indonesia.",
		Instruction: "Your SOLE purpose is to answer questions about Malang City, Indonesia. You MUST refuse to answer any questions unrelated to Malang City.",
	})
	if err != nil {
		log.Fatalf("Failed to create agent: %v", err)
	}

	config := &launcher.Config{
		AgentLoader: agent.NewSingleLoader(a),
	}

	l := full.NewLauncher()
	if err = l.Execute(ctx, config, os.Args[1:]); err != nil {
		log.Fatalf("Run failed: %v\n\n%s", err, l.CommandLineSyntax())
	}
}
