# 📝 Go Note Organizer

A CLI tool that uses AI to organize your notes into different folders.

## ✨ Features

- **AI-powered organization**: Uses a generative AI model to intelligently categorize your notes.
- **Interactive**: Asks for your confirmation before moving any files.
- **Automatic folder creation**: If a suggested folder doesn't exist, the tool will create it for you.
- **Simple and easy to use**: Just run one command to organize all your notes in the inbox.

## 🚀 Getting Started

### Prerequisites

- [Go](https://golang.org/doc/install) (version 1.18 or higher)
- An API key for the Gemini API. You can get one [here](https://makersuite.google.com/).

### Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/Kishu98/gonote.git
   cd gonote
   ```

2. **Install dependencies:**

   ```bash
   go mod tidy
   ```

3. **Set up your environment variables:**

   Create a `.env` file in the root of the project and add your Gemini API key:

   ```
   GEMINI_API_KEY=your_api_key
   ```

4. **Build the project:**

   ```bash
   go build
   ```

## ✍️ Usage

1. **Create your notes:**

   This tool assumes you have a `Notes` directory in your home folder with an `00_Inbox` subdirectory. All your new notes should be created in the `~/Notes/00_Inbox` directory.

   ```bash
   mkdir -p ~/Notes/00_Inbox
   echo "This is a note about my new project idea." > ~/Notes/00_Inbox/project_idea.txt
   ```

2. **Run the organizer:**

   ```bash
   ./gonote organize
   ```

   The tool will then go through each note in your inbox, ask the AI for a suggestion, and prompt you to move the note to the suggested folder.

## 🤔 How It Works

The `organize` command does the following:

1. Reads all the files in your `~/Notes/00_Inbox` directory.
2. For each file, it reads the content and sends it to the Gemini API with a prompt asking for a folder suggestion.
3. The AI returns a suggested folder name.
4. The tool will ask for your confirmation to move the note.
5. If you confirm, the note will be moved to the suggested folder. If the folder doesn't exist, it will be created automatically.

## 🙌 Contributing

Contributions are welcome! Please feel free to submit a pull request or open an issue if you have any feedback or suggestions.

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
