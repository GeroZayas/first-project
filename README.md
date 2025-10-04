# ✅ FocusFlow CLI

A colorful, human-friendly todo manager for your terminal. Built in Go to keep your focus sharp and your list tidy. ✨

## 🚀 Features
- 🎨 Rich ANSI color theme for instant status cues
- ⚡ Snappy keyboard-driven workflow with single-key shortcuts
- 📋 Clear table layout showing status, task, and timestamps
- 🧹 Batch-clear completed tasks without touching the rest
- 🛡️ No config, no setup—just run and start typing

## 🧭 Usage
```bash
go run main.go
```

### Menu shortcuts
- `1 / a` – Add a task
- `2 / t` – Toggle completion
- `3 / d` – Delete a task
- `4 / c` – Clear all completed
- `q`     – Quit gracefully

## 🧪 Build & Run
Want to build the binary?
```bash
go build -o todos
./todos
```

> ⚠️ Sandboxed environments might block Go's build cache. If `go build` fails with a permission error, rerun it outside the sandbox or clear the Go cache directory.

## 💡 Tips
- Keep task titles short; FocusFlow trims long ones with an ellipsis.
- Completed items show their finish time to celebrate progress. 🎉
- Pair it with your favorite terminal font for max vibes.

Happy shipping! 🚢
