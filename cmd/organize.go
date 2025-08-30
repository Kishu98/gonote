// TODO
// Make the note into a struct
package cmd

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"google.golang.org/genai"
)

func init() {
	rootCMD.AddCommand(organizeCMD)
}

var organizeCMD = &cobra.Command{
	Use:   "organize",
	Short: "Organizes the quick notes into specific folders.",
	Run: func(cmd *cobra.Command, args []string) {
		err := godotenv.Load()
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			fmt.Println("Error: GEMINI_API_KEY is missing")
			return
		}

		ctx := context.TODO()
		client, err := genai.NewClient(ctx, &genai.ClientConfig{
			APIKey:  apiKey,
			Backend: genai.BackendGeminiAPI,
		})

		home, err := os.UserHomeDir()
		if err != nil {
			fmt.Println("Error getting home directory:", err)
			return
		}

		inboxPath := filepath.Join(home, "Notes", "00_Inbox")
		notesPath := filepath.Join(home, "Notes")

		files, err := os.ReadDir(inboxPath)
		if err != nil {
			fmt.Println("Error reading inbox directory:", err)
			return
		}

		folders, err := getFolders(notesPath)
		if err != nil {
			fmt.Println("Error getting folders:", err)
			return
		}

		for _, file := range files {
			if file.IsDir() {
				continue
			}

			filePath := filepath.Join(inboxPath, file.Name())
			content, err := os.ReadFile(filePath)
			if err != nil {
				fmt.Printf("Error reading file %s: %v\n", file.Name(), err)
				continue
			}

			fmt.Printf("\nProcessing Note: %s\n", file.Name())

			suggestedFolder, err := getAISuggestion(ctx, client, string(content), folders)
			if err != nil {
				fmt.Printf("Error getting AI Suggestion for %s: %v\n", file.Name(), err)
				continue
			}

			// fmt.Printf("Suggested Folder for note\n\n%s\n is --> %s\n", string(content), suggestedFolder)
			// #TODO: Note Handling

			handleNote(suggestedFolder, folders, filePath)
		}
	},
}

func handleNote(suggestedFolder string, folders []string, notePath string) {
	contains := false
	for _, folder := range folders {
		if folder == suggestedFolder {
			contains = true
		}
	}
	if contains {
		fmt.Printf("AI wants to move the note to folder: %s\n", suggestedFolder)
		fmt.Println("Do you want me to continue?")
		if checkUserInput() {
			fmt.Printf("Moving note to %s", suggestedFolder)
			moveNote(notePath, suggestedFolder)
		} else {
			fmt.Println("Continueing....")
		}
	} else {
		fmt.Printf("AI wants to move the folder to %s\n", suggestedFolder)
		fmt.Println("Do you want me to continue?")
		if checkUserInput() {
			home, err := os.UserHomeDir()
			if err != nil {
				log.Fatal(err)
			}
			os.Mkdir(filepath.Join(home, "Notes", suggestedFolder), 0755)
			moveNote(notePath, suggestedFolder)
		} else {
			fmt.Println("Continueing....")
		}
	}

	fmt.Printf("\n**********************************************\n")
}

func checkUserInput() bool {
	var input string
	fmt.Scan(&input)
	if input == "y" {
		return true
	}
	return false
}

func moveNote(notePath string, suggestedFolder string) {
	home, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		return
	}
	if err := os.Rename(notePath, filepath.Join(home, "Notes", suggestedFolder, filepath.Base(notePath))); err != nil {
		log.Fatal(err)
	}
	fmt.Printf("Moved note to %s\n", suggestedFolder)
}

func getAISuggestion(ctx context.Context, client *genai.Client, content string, folders []string) (string, error) {
	prompt := fmt.Sprintf(`You are helping organize notes into folders. 

EXISTING CATEGORIES: %s

TASK: Analyze the following note content and determine the best category for it.

RULES:
1. If the note clearly fits and strongly aligns into one of the existing categories, respond with just the category name
2. If you're unsure or the note doesn't fit or does not strongly aligns existing categories , respond with "suggested_category" where suggested_category is a single word
3. Keep suggested categories simple and descriptive (e.g., "Work", "Personal", "Ideas", "Recipes", "Finance")

NOTE CONTENT:
%s

RESPONSE:`, strings.Join(folders, ", "), content)

	// var budget int32 = 0
	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
		// &genai.GenerateContentConfig{
		// 	ThinkingConfig: &genai.ThinkingConfig{
		// 		IncludeThoughts: true,
		// 		// ThinkingBudget:  new(int32),
		// 	},
		// },
	)
	// for _, part := range result.Candidates[0].Content.Parts {
	// 	if part.Text != "" {
	// 		if part.Thought {
	// 			fmt.Println("Thoughts Summary:")
	// 			fmt.Println(part.Text)
	// 		} else {
	// 			fmt.Println("Answer:")
	// 			fmt.Println(part.Text)
	// 		}
	// 	}
	// 	fmt.Printf("\n*****************************************************\n\n")
	// }
	if err != nil {
		log.Fatal(err)
	}

	return result.Text(), nil
}

func getFolders(path string) ([]string, error) {
	var folders []string

	entries, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}

	for _, entry := range entries {
		if entry.IsDir() && entry.Name() != "00_Inbox" && entry.Name() != "Projects" {
			folders = append(folders, entry.Name())
		}
	}

	return folders, nil
}
