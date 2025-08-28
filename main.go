package main

import (

	//"bufio"
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
	var char_status rune
	if !t.status {
		//Unicode Character “✗” (U+2717)
		char_status = '✗'
	} else {
		//Unicode Character “✓” (U+2713)
		char_status = '✓'
	}

	fmt.Printf("%s\t%c\n", t.content, char_status)
}

//-------------------------------------------------------------------------------------

type List struct {
	s []Task
}

func (l *List) appendTask(t Task) {
	if len(l.s) >= cap(l.s) {
		fmt.Printf("Max amount of Tasks reached (%d). Please remove some Task to add more", capacity)
	}
	l.s = append(l.s, t)
}

func (l *List) markAsDone(n int) {
	// before calling this func check if 0 <= n && n < len(l.s)
	if 0 < n && len(l.s) < n {
		fmt.Printf("len of l.s %d\n", len(l.s))
		fmt.Printf("Given index %d outside range of todolist\n", n)
		return
	}
	l.s[n].status = true
	fmt.Println("status changed")
}

func (l *List) printTasks() {
	for idx, task := range l.s {
		fmt.Printf("%d.  ", idx+1)
		task.print()
	}
}

//-------------------------------------------------------------------------------------

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadData(l *List) {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_RDWR, 0666)
	check(err)
	defer file.Close()

	// get file size in bytes and allocate buffer
	finfo, err := file.Stat()
	check(err)
	if finfo.Size() == 0 {
		fmt.Print("File Size is 0")
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

func saveData(l *List) {
	return
}

func main() {
	// underlying array
	task_array := [capacity]Task{}
	// create a slice of array
	task_list := List{s: task_array[:0]}
	loadData(&task_list)

	// Main Loop
	var flag bool = true
	for flag {

		// Show Task list
		task_list.printTasks()

		fmt.Println()
		fmt.Print("(1) Add Task; (2) Mark as Done; (3) Remove Task; (0) Exit")
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
		case 1:
			fmt.Println("Adding task...")
		case 2:
			task_list.printTasks()
			fmt.Print(": ")
			_, err := fmt.Scanf("%d", &selection)
			if err != nil {
				continue
			}
			task_list.markAsDone(selection - 1)

		case 3:
			fmt.Println("removing task...")
		default:
			fmt.Println("please enter a number from 1-4 or use ctrl-c")

		}

	}

}
