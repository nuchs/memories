# Memories

This repository holds **project-specific notes** ("memories") that provide long-lived context when working with AI assistants.  
Each project gets its own Markdown file (e.g. `how_to_commit_fraud.md`, `embarassing_medicals.md`) thatcaptures background information, design notes, or anything else you want to carry across sessions.

Temporary **pointer files** are generated locally and uploaded into a project chat. These pointers tell the assistant which memory file to look at in this repo, so you don’t have to paste notes or maintain long-term context manually.

*N.B.* This is very much tailored for my own use

---

## Usage

1. Install the `memories` tool:

   ```bash
   go install github.com/nuchs/memories@latest
    ```
Make sure ~/go/bin (or your Go bin path) is on your $PATH.

2. Set the environment variable MYMEMORIES to point at this repo’s local clone:

    ```bash
    export MYMEMORIES=~/dev/memories
    ```

3. Generate a pointer and open the memory file:

    ```bash
    memories alibi
    ```
This will:

    * Create alibi.pointer.txt in the current directory.
    * Ensure alibi.md exists in $MYMEMORIES (creating a stub if needed).
    * Open alibi.md in Neovim for editing.

4. Upload the generated snaked.pointer.txt into the ChatGPT project chat.
That tells the assistant: “For context, look in snaked.md in the
memories repo.”

## Conventions

To keep memories easy to navigate and version:

* File naming: each project has one memory file named <project>.md.

* Pointer naming: temporary pointers are <project>.pointer.txt and should
be .gitignored.

* Structure inside memory files:
    * Begin with a top-level heading: # <Project> Memory.
    * Optionally include a version marker and date, e.g.:
        ```Text
        # Alibi Memory
        _v3 — 2025-09-19_

        - I am allegric to grass and books
        - I don't even know what a knoll is
        ```
    * Use bullet points or short sections for clarity.
    * Add links or references to code, issues, or external resources as needed.

* Versioning: use Git history for tracking. Bump the inline version number
when making significant updates.

* Scope: keep each file focused on one project. Create new files rather
than mixing multiple projects.

## Notes

Pointer files are temporary — add *.pointer.txt to your .gitignore so
they don’t get committed.

Memory files (*.md) are versioned in Git, so you can track changes,
branch them, or roll back like any other source file.
