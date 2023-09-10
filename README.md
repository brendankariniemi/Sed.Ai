# Sed.Ai

#### sed 's/Human Text/Command/g'

## Introduction

**Sed.Ai** is a Go-based CLI tool that harnesses OpenAI's ChatGPT to convert plain English input into executable command-line instructions. 
It simplifies tasks by transforming "Find files accessed in the last day" into the actionable command "find /dir -mtime 0" enhancing productivity through natural language interaction.

### Features

- **Intelligent Command Generation**: Chat with ChatGPT to generate valid command-line instructions.
- **Command History**: Users can specify the number of recent commands they wish to save, and navigate those commands using the up/down keys at launch. This feature allows for easy command re-execution while retaining the context provided to ChatGPT which allows users to add on extra details to the commands instead of re-describing your goals from scratch.

## Installation

1. Clone the project and navigate to project directory
2. Install dependencies and build the project:

```
go get -d -v ./...

go build
```
3. Specify OPENAI_API_KEY and other configuration settings in /internal/config/config.yaml
4. Create a softlink to the executable from somwhere in your PATH (ex: /usr/local/bin) to make easily call it from anywhere
```
ln -s /path_to_built_executable/sai /usr/local/bin/sai
```
Replace /path_to_built_executable/sai with the actual path to the built executable. 

5. Run the application
```
sai
```
