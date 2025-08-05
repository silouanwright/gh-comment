package cmd

import (
	"bytes"
	"embed"
	"fmt"
	"io/fs"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

var (
	promptCategory string
	listPrompts    bool
)

//go:embed prompts/**/*.md
var promptFiles embed.FS

var promptsCmd = &cobra.Command{
	Use:   "prompts [prompt-name]",
	Short: "Get AI-powered code review prompts and best practices",
	Long: heredoc.Doc(`
		Access a curated collection of professional code review prompts optimized
		for AI assistants like Claude, ChatGPT, and other automated workflows.

		These prompts follow research-backed best practices for effective code
		review communication, including psychological safety principles and
		the CREG emoji system (🔧🤔♻️📝😃📌).

		Perfect for: Senior developers, code review leads, QA teams, and AI
		assistants who need sophisticated review templates beyond basic feedback.
	`),
	Example: heredoc.Doc(`
		# List all available prompts
		$ gh comment prompts --list

		# Get specific security review prompt
		$ gh comment prompts security-audit

		# List prompts by category
		$ gh comment prompts --category performance --list

		# Get performance optimization prompt
		$ gh comment prompts performance-optimization

		# Security-focused review prompt
		$ gh comment prompts security-audit

		# Architecture review prompt
		$ gh comment prompts architecture-review
	`),
	Args: cobra.MaximumNArgs(1),
	RunE: runPrompts,
}

func init() {
	rootCmd.AddCommand(promptsCmd)

	promptsCmd.Flags().StringVar(&promptCategory, "category", "", "Filter prompts by category: security, performance, architecture, quality, ai")
	promptsCmd.Flags().BoolVar(&listPrompts, "list", false, "List available prompts")
}

func runPrompts(cmd *cobra.Command, args []string) error {
	if listPrompts {
		return listAvailablePrompts()
	}

	if len(args) == 0 {
		return listAvailablePrompts()
	}

	promptName := args[0]

	// Handle "list" as an argument (user-friendly alias for --list flag)
	if promptName == "list" {
		return listAvailablePrompts()
	}

	prompt, exists := getPrompt(promptName)
	if !exists {
		fmt.Printf("❌ Prompt '%s' not found.\n\n", promptName)
		fmt.Println("Available prompts:")
		return listAvailablePrompts()
	}

	fmt.Printf("📋 **%s**\n\n", prompt.Title)
	fmt.Printf("🎯 **Category**: %s\n", prompt.Category)
	fmt.Printf("⏱️  **Estimated Time**: %s\n\n", prompt.EstimatedTime)

	fmt.Println("📝 **Prompt:**")
	fmt.Println("```")
	fmt.Println(prompt.Content)
	fmt.Println("```")

	if len(prompt.Examples) > 0 {
		fmt.Println("\n💡 **Example Usage:**")
		for _, example := range prompt.Examples {
			fmt.Printf("• %s\n", example)
		}
	}

	if len(prompt.Tags) > 0 {
		fmt.Printf("\n🏷️  **Tags**: %s\n", strings.Join(prompt.Tags, ", "))
	}

	return nil
}

type Prompt struct {
	Name          string
	Title         string
	Category      string
	Content       string
	Examples      []string
	Tags          []string
	EstimatedTime string
}

type PromptMetadata struct {
	Name          string   `yaml:"name"`
	Title         string   `yaml:"title"`
	Category      string   `yaml:"category"`
	EstimatedTime string   `yaml:"estimated_time"`
	Tags          []string `yaml:"tags"`
	Examples      []string `yaml:"examples"`
}

func getPrompt(name string) (Prompt, bool) {
	prompts := getAllPrompts()
	prompt, exists := prompts[name]
	return prompt, exists
}

func listAvailablePrompts() error {
	prompts := getAllPrompts()

	if promptCategory != "" {
		fmt.Printf("📋 **Code Review Prompts - %s Category**\n\n", strings.Title(promptCategory))
	} else {
		fmt.Print("📋 **Available Code Review Prompts**\n\n")
	}

	categories := make(map[string][]Prompt)
	for _, prompt := range prompts {
		if promptCategory != "" && prompt.Category != promptCategory {
			continue
		}
		categories[prompt.Category] = append(categories[prompt.Category], prompt)
	}

	for category, categoryPrompts := range categories {
		fmt.Printf("## %s\n", strings.Title(category))
		for _, prompt := range categoryPrompts {
			fmt.Printf("  **%s** - %s (%s)\n", prompt.Name, prompt.Title, prompt.EstimatedTime)
			fmt.Printf("    %s\n", strings.Join(prompt.Tags, " • "))
		}
		fmt.Print("\n")
	}

	fmt.Println("💡 **Usage**: `gh comment prompts <prompt-name>` to get the full prompt")
	return nil
}

func getAllPrompts() map[string]Prompt {
	prompts := make(map[string]Prompt)

	// Walk through all markdown files in the prompts directory
	fs.WalkDir(promptFiles, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		// Only process .md files
		if !strings.HasSuffix(path, ".md") {
			return nil
		}

		// Read the markdown file
		content, err := promptFiles.ReadFile(path)
		if err != nil {
			return err
		}

		// Parse the markdown file with frontmatter
		prompt, err := parseMarkdownPrompt(content)
		if err != nil {
			// Skip files that can't be parsed, but continue processing others
			return nil
		}

		prompts[prompt.Name] = prompt
		return nil
	})

	return prompts
}

func parseMarkdownPrompt(content []byte) (Prompt, error) {
	// Split frontmatter and content
	parts := bytes.SplitN(content, []byte("---"), 3)
	if len(parts) < 3 {
		return Prompt{}, fmt.Errorf("invalid markdown format: missing frontmatter")
	}

	// Parse YAML frontmatter
	var metadata PromptMetadata
	if err := yaml.Unmarshal(parts[1], &metadata); err != nil {
		return Prompt{}, fmt.Errorf("failed to parse frontmatter: %w", err)
	}

	// Extract markdown content (remove leading/trailing whitespace)
	markdownContent := strings.TrimSpace(string(parts[2]))

	return Prompt{
		Name:          metadata.Name,
		Title:         metadata.Title,
		Category:      metadata.Category,
		Content:       markdownContent,
		Examples:      metadata.Examples,
		Tags:          metadata.Tags,
		EstimatedTime: metadata.EstimatedTime,
	}, nil
}
