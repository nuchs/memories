package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"text/template"
	"time"
)

// Default fallback template (used if TEMPLATE.md is absent).
const defaultTemplate = `# {{ .Project }} Memory
_v{{ .Version }} — {{ .Date }}_

## Context
- 

## Goals / Next Up
- 

## Decisions
- 

## Links
- Repo:
- Issues / Boards:
- Docs:

## Notes
- 
`

type FilePaths struct {
	pointer  string
	memory   string
	template string
}

type Details struct {
	Project string
	Version int
	Date    string
}

func main() {
	if len(os.Args) != 2 {
		log.Fatalf("Usage: %s <project-name>", os.Args[0])
	}
	project := os.Args[1]
	paths := GetFilePaths(project)

	if err := EnsureMemory(paths, project); err != nil {
		log.Fatalf("failed to prepare memory: %v", err)
	}

	if err := OpenInNvim(paths.memory); err != nil {
		log.Fatalf("failed to open editor: %v", err)
	}
}

// GetFilePaths returns the local pointer path, the memory path, and TEMPLATE.md path.
func GetFilePaths(project string) FilePaths {
	// Pointer always in current working dir
	cwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get current directory: %v", err)
	}
	pointer := filepath.Join(cwd, project+".pointer.txt")

	// Memory in $MYMEMORIES (or cwd if unset)
	memDir := os.Getenv("MYMEMORIES")
	if memDir == "" {
		memDir = cwd
	}
	memory := filepath.Join(memDir, project+".md")
	template := filepath.Join(memDir, "TEMPLATE.md")

	return FilePaths{pointer, memory, template}
}

// EnsurePointer creates/overwrites the pointer file that references the GitHub memory path.
func EnsurePointer(path, project string) error {
	content := fmt.Sprintf("See: github.com/nuchs/memories/blob/main/%s.md\n", project)
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return fmt.Errorf("failed to write pointer: %w", err)
	}
	log.Printf("Created/updated pointer %s", path)
	return nil
}

// EnsureMemory ensures the memory file exists. If missing, it delegates creation
// to CreateNewMemory (which uses TEMPLATE.md if present).
func EnsureMemory(paths FilePaths, project string) error {
	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(paths.memory), 0755); err != nil {
		return fmt.Errorf("failed to ensure memory dir: %w", err)
	}

	// If it already exists, nothing to do
	if _, err := os.Stat(paths.memory); err == nil {
		log.Printf("Memory file exists: %s", paths.memory)
		return nil
	} else if !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat memory file: %w", err)
	}

	// only create the pointer if the memory does not alread exist
	if err := EnsurePointer(paths.pointer, project); err != nil {
		log.Fatalf("failed to prepare pointer: %v", err)
	}

	return CreateNewMemory(paths, project)
}

// CreateNewMemory creates a fresh memory file at memPath. If tmplPath exists,
// it renders it as a text/template with the provided fields; otherwise writes a stub.
func CreateNewMemory(paths FilePaths, project string) error {
	tbytes := []byte(defaultTemplate)
	if fi, err := os.Stat(paths.template); err == nil && !fi.IsDir() {
		b, rerr := os.ReadFile(paths.template)
		if rerr != nil {
			return fmt.Errorf("failed to read TEMPLATE.md: %w", rerr)
		}
		tbytes = b
	} else if err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to stat TEMPLATE.md: %w", err)
	}

	tmpl, perr := template.New("mem").Parse(string(tbytes))
	if perr != nil {
		return fmt.Errorf("failed to parse template: %w", perr)
	}

	data := Details{
		Project: project,
		Version: 1,
		Date:    time.Now().Format("1979-01-01"),
	}

	var buf bytes.Buffer
	if exerr := tmpl.Execute(&buf, data); exerr != nil {
		return fmt.Errorf("failed to render template: %w", exerr)
	}

	if werr := os.WriteFile(paths.memory, buf.Bytes(), 0644); werr != nil {
		return fmt.Errorf("failed to write memory file: %w", werr)
	}

	log.Printf("Created memory file %s", paths.memory)
	return nil
}

// OpenInNvim launches $EDITOR (or nvim) with the given file.
func OpenInNvim(path string) error {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nvim"
	}
	cmd := exec.Command(editor, path)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
