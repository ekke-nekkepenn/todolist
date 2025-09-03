package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

const (
	filename string = "data.txt"
	capacity int    = 100
)

// -------------------------------------------------------------------------------------
type Task struct {
	content string
	notes   string
	status  bool
}

func (t *Task) print() {
	fmt.Printf("%s\t%c\n", t.content, t.getStatusSymbol())
}

func (t *Task) getFormattedString() string {
	return fmt.Sprintf("%s;%s;%c\n", t.content, t.notes, t.getStatusSymbol())
}

func (t *Task) getStatusSymbol() rune {
	if t.status {
		//Unicode Character “✗” (U+2717)
		//char_status = '✗'
		return 'o'
	}
	//Unicode Character “✓” (U+2713)
	//char_status = '✓'
	return 'x'
}

//-------------------------------------------------------------------------------------

type List struct {
	items []Task
}

func (l *List) appendTask(t Task) {
	if len(l.items) >= cap(l.items) {
		fmt.Printf("Max amount of Tasks reached (%d). Please remove some Task to add more", capacity)
	}
	l.items = append(l.items, t)
}

func (l *List) isIndexValid(n int) bool {
	if 0 < n && len(l.items) <= n {
		return false
	}
	return true
}

func (l *List) markAsDone(n int) {
	// check n is positive AND enough room in l
	if !(l.isIndexValid(n)) {
		fmt.Printf("TThere is no task with index %d", n+1)
		return
	}
	l.items[n].status = true
	fmt.Println("status changed")
}

func (l *List) printTasks() {
	for i, task := range l.items {
		fmt.Printf("%d.  ", i+1)
		task.print()
	}
}

//-------------------------------------------------------------------------------------

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadData(l *List, file *os.File) {
	// get file size in bytes and allocate buffer
	finfo, err := file.Stat()
	check(err)

	if finfo.Size() == 0 {
		return
	}

	buffer := make([]byte, finfo.Size())

	// read file data into buffer
	_, err = file.Read(buffer)
	check(err)

	lines := strings.Lines(string(buffer))
	for line := range lines {
		line = strings.TrimSuffix(line, "\n")
		parts := strings.Split(line, ";")

		var status bool
		// squiggly lines?
		if strings.ToLower(parts[2]) == "x" {
			status = false
		} else if strings.ToLower(parts[2]) == "o" {
			status = true
		} else {
			fmt.Println("Status data in file is wrong")
			panic(status)
		}

		l.appendTask(Task{
			content: parts[0],
			notes:   parts[1],
			status:  status,
		})
	}
}

func saveData(l *List, file *os.File) {
	// deletes content of file & and writes TaskList items to file
	file.Truncate(0)
	file.Seek(0, 0)
	for _, task := range l.items {
		file.WriteString(task.getFormattedString())
	}
}

func addTask(l *List) {
	// https://pkg.go.dev/bufio#Scanner
	// from the examples
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Task: ")
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	content := scanner.Text()

	fmt.Print("Notes: ")
	scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
	notes := scanner.Text()

	l.appendTask(Task{content: content, notes: notes, status: false})
}

func removeTask(l *List) {
	l.printTasks()

	var selection int
	_, err := fmt.Scanf("%d", &selection)
	check(err)
	selection--

	if selection < 0 || selection >= len(l.items) {
		fmt.Printf("%d is not a valid index\n", selection+1)
	}

	// delete Task
	l.items[selection] = Task{}

	// close up "hole" in array by copying later task into it
	for i := selection + 1; i < len(l.items); i++ {
		l.items[i-1] = l.items[i]
	}

	// resize the slice
	l.items = l.items[:len(l.items)-1]
}

func main() {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR, 0666)

	check(err)
	defer file.Close()

	// Task Storage
	task_array := [capacity]Task{}
	task_list := List{items: task_array[:0]}
	loadData(&task_list, file)

	// Main Loop
	var flag bool = true
	for flag {
		// Show Task list
		task_list.printTasks()

		fmt.Println()
		fmt.Print("(1) Add Task; ")
		fmt.Print("(2) Mark as Done; ")
		fmt.Print("(3) Remove Task; ")
		fmt.Print("(0) Exit\n")
		fmt.Print(": ")

		var selection int
		_, err := fmt.Scanf("%d", &selection)

		// if wrong input is given, ignore it because idk how to catch errors in golang
		if err != nil {
			continue
		}

		switch selection {

		case 0:
			fmt.Println("Exiting...")
			saveData(&task_list, file)
			return
		case 1:
			fmt.Println("Adding task...")
			addTask(&task_list)

		case 2:
			task_list.printTasks()
			fmt.Print(": ")
			_, err := fmt.Scanf("%d", &selection)
			if err != nil {
				continue
			}
			task_list.markAsDone(selection - 1)

		case 3:
			fmt.Println("Removing task...")
			removeTask(&task_list)

		default:
			fmt.Println("please enter a number from 1-4 or use ctrl-c")

		}
	}
}
