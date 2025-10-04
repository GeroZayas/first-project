package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const (
	colorReset   = "\033[0m"
	colorBold    = "\033[1m"
	colorFaint   = "\033[2m"
	colorRed     = "\033[31m"
	colorGreen   = "\033[32m"
	colorYellow  = "\033[33m"
	colorBlue    = "\033[34m"
	colorMagenta = "\033[35m"
	colorCyan    = "\033[36m"
	colorGray    = "\033[90m"
)

type todo struct {
	title       string
	completed   bool
	createdAt   time.Time
	completedAt *time.Time
}

type feedback struct {
	message string
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	todos := make([]todo, 0, 16)
	var note *feedback

	for {
		clearScreen()
		printBanner()
		renderTodos(todos)

		if note != nil {
			fmt.Println(note.message)
		}

		printMenu()

		fmt.Print(colorCyan + "Choose an option: " + colorReset)
		input, err := readLine(reader)
		if err != nil {
			fmt.Println(colorRed + "Failed to read input:" + err.Error() + colorReset)
			return
		}

		choice := strings.ToLower(input)
		note = nil

		switch choice {
		case "1", "a", "add":
			fmt.Print(colorCyan + "Task title: " + colorReset)
			title, err := readLine(reader)
			if err != nil {
				note = &feedback{message: colorRed + "Could not read the task title." + colorReset}
				continue
			}

			title = strings.TrimSpace(title)
			if title == "" {
				note = &feedback{message: colorYellow + "Empty task discarded." + colorReset}
				continue
			}

			todos = append(todos, todo{title: title, createdAt: time.Now()})
			note = &feedback{message: colorGreen + "Added \"" + title + "\"" + colorReset}

		case "2", "t", "toggle":
			if len(todos) == 0 {
				note = &feedback{message: colorYellow + "Nothing to toggle yet." + colorReset}
				continue
			}

			fmt.Print(colorCyan + "Toggle task #: " + colorReset)
			idx, err := parseIndex(reader, len(todos))
			if err != nil {
				note = &feedback{message: colorRed + err.Error() + colorReset}
				continue
			}

			toggled := todos[idx]
			todos[idx].completed = !todos[idx].completed
			if todos[idx].completed {
				now := time.Now()
				todos[idx].completedAt = &now
				note = &feedback{message: colorGreen + "Completed \"" + toggled.title + "\"" + colorReset}
			} else {
				todos[idx].completedAt = nil
				note = &feedback{message: colorBlue + "Marked \"" + toggled.title + "\" as pending" + colorReset}
			}

		case "3", "d", "delete":
			if len(todos) == 0 {
				note = &feedback{message: colorYellow + "Nothing to delete yet." + colorReset}
				continue
			}

			fmt.Print(colorCyan + "Delete task #: " + colorReset)
			idx, err := parseIndex(reader, len(todos))
			if err != nil {
				note = &feedback{message: colorRed + err.Error() + colorReset}
				continue
			}

			removed := todos[idx].title
			todos = append(todos[:idx], todos[idx+1:]...)
			note = &feedback{message: colorMagenta + "Deleted \"" + removed + "\"" + colorReset}

		case "4", "c", "clear":
			if len(todos) == 0 {
				note = &feedback{message: colorYellow + "List already empty." + colorReset}
				continue
			}

			kept := todos[:0]
			removed := 0
			for _, td := range todos {
				if td.completed {
					removed++
					continue
				}
				kept = append(kept, td)
			}
			todos = kept

			if removed == 0 {
				note = &feedback{message: colorYellow + "No completed tasks to clear." + colorReset}
			} else {
				note = &feedback{message: colorMagenta + fmt.Sprintf("Cleared %d completed task(s).", removed) + colorReset}
			}

		case "q", "quit", "exit":
			clearScreen()
			fmt.Println(colorGreen + "Bye! Go crush those goals." + colorReset)
			return

		default:
			note = &feedback{message: colorYellow + "Unknown option. Try again." + colorReset}
		}
	}
}

func renderTodos(todos []todo) {
	fmt.Println()
	if len(todos) == 0 {
		fmt.Println(colorGray + "No todos yet. Add your first one!" + colorReset)
		fmt.Println()
		return
	}

	fmt.Printf(colorBold+"%-4s %-7s %-45s %-16s"+colorReset+"\n", "#", "Status", "Task", "Created")
	fmt.Println(colorGray + strings.Repeat("─", 80) + colorReset)

	for idx, td := range todos {
		status := colorYellow + "[ ]" + colorReset
		when := td.createdAt.Format("Jan 02 15:04")
		title := truncate(td.title, 42)

		if td.completed {
			status = colorGreen + "[✔]" + colorReset
			if td.completedAt != nil {
				when = td.completedAt.Format("Done 15:04")
			} else {
				when = "Done"
			}
			title = colorFaint + title + colorReset
		}

		fmt.Printf("%s%2d%s   %-7s %-45s %-16s\n",
			colorBlue, idx+1, colorReset, status, title, when,
		)
	}
	fmt.Println()
}

func printBanner() {
	fmt.Println(colorMagenta + strings.Repeat("═", 60) + colorReset)
	fmt.Println(colorBold + colorMagenta + "  ✨  FocusFlow CLI — manage tasks at the speed of thought" + colorReset)
	fmt.Println(colorMagenta + strings.Repeat("═", 60) + colorReset)
}

func printMenu() {
	fmt.Println(colorGray + strings.Repeat("-", 60) + colorReset)
	fmt.Printf("%s[1]%s Add    %s[2]%s Toggle    %s[3]%s Delete    %s[4]%s Clear Completed    %s[Q]%s Quit\n",
		colorMagenta, colorReset,
		colorBlue, colorReset,
		colorRed, colorReset,
		colorGreen, colorReset,
		colorYellow, colorReset,
	)
	fmt.Println(colorGray + strings.Repeat("-", 60) + colorReset)
}

func parseIndex(reader *bufio.Reader, total int) (int, error) {
	raw, err := readLine(reader)
	if err != nil {
		return 0, errors.New("failed to read the number")
	}

	num, err := strconv.Atoi(strings.TrimSpace(raw))
	if err != nil {
		return 0, errors.New("please type a valid number")
	}

	num--
	if num < 0 || num >= total {
		return 0, fmt.Errorf("pick between 1 and %d", total)
	}

	return num, nil
}

func readLine(reader *bufio.Reader) (string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(line), nil
}

func truncate(text string, limit int) string {
	if len([]rune(text)) <= limit {
		return text
	}

	runes := []rune(text)
	return string(runes[:limit-1]) + "…"
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}
