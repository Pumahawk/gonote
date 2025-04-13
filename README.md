# 📝 GoNote

**GoNote** is a simple CLI tool written in Go to manage notes stored in either YAML or Markdown format. It lets you list, inspect, and filter your notes efficiently with a few commands.

---

## 📦 Features

- Supports notes written in **YAML** or **Markdown** with frontmatter.
- Powerful filtering by tags and regex on titles or IDs.
- Simple CLI interface with table output format.
- Lightweight and fast — perfect for dev notes, bug tracking, or TODOs.

---

## 🗂 Note Formats

### YAML Example

```yaml
notes:
- id: BUG-001
  title: "Login issue on dashboard"
  tags: [ "bug", "todo" ]
  note: |
    The dashboard crashes when a user with no profile logs in.
```

### Markdown Example

```markdown
---
id: MD-HELLO
title: "Welcome Note"
tags: ["md"]
---

# Hello, GoNote!

This is a sample markdown note powered by GoNote.
```

---

## 🚀 Usage

```bash
gonote -root ./notes
```

`-root` specifies the directory where notes are stored (defaults to current `.` directory).

### 📄 Commands

#### `ls` – List all notes

```bash
gonote ls [options]
```

**Options:**

| Flag               | Description                            |
|--------------------|----------------------------------------|
| `-o`               | Output format. Default: `table`        |
| `-t`               | Filter notes by tags (AND logic)       |
| `-tor`             | Filter notes by tags (OR logic)        |
| `-xid`             | Regex match for ID                     |
| `-xtitle`          | Regex match for title                  |
| `-tableIdWidth`    | Set ID column width                    |
| `-tableTitleWidth` | Set Title column width                 |
| `-tableTagsWidth`  | Set Tags column width                  |

**Example:**

```bash
gonote ls -t bug -xtitle "auth"
```

#### `info` – Display summary info about notes

```bash
gonote info [options]
```

**Options are the same as in `ls`**.

**Example Output:**

```
Notes: 9
Tags: 
    bug: 7
    fixed: 3
    frontend: 3
    md: 1
    remind: 1
    skip: 2
    todo: 2
```

#### `show` – Show full note content

```bash
gonote show <NOTE_ID>
```

Displays the complete content of a note identified by its `id`.

**Example:**

```bash
gonote show MD-HELLO
```

---

## 📁 Example Output

### `ls`

```
ID                        TITLE                                                     TAGS                      PATH
--------------------------------------------------------------------------------------------------------------
MD-HELLO                  Welcome Note                                              md                        md/note.md:1
BUG-001                   Login issue on dashboard                                  bug, fixed                note.yaml:2
BUG-002                   Crash on payment success callback                         bug, remind, todo         note.yaml:26
BUG-003                   Notification system not triggering alerts properly        bug, fixed                note.yaml:35
```

### `info`

```
Notes: 9
Tags: 
    bug: 7
    fixed: 3
    frontend: 3
    md: 1
    remind: 1
    skip: 2
    todo: 2
```

### `show`

```markdown
---
id: MD-HELLO
title: "Welcome Note"
tags: ["md"]
---

# Hello, GoNote!

This is a sample markdown note powered by GoNote.
```

---

## 🔧 Installation

Build it from source:

```bash
go build -o gonote
```

Or add it to your `$GOPATH/bin` for global access.
