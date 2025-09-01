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
	items []Task
}

func (l *List) appendTask(t Task) {
	if len(l.items) >= cap(l.items) {
		fmt.Printf("Max amount of Tasks reached (%d). Please remove some Task to add more", capacity)
	}
	l.items = append(l.items, t)
}

func (l *List) markAsDone(n int) {
	// check n is positive AND enough room in l 
	if 0 < n && len(l.items) < n {
		fmt.Printf("len of l.items %d\n", len(l.items))
		fmt.Printf("Given index %d outside range of todolist\n", n)
		return
	}
	l.items[n].status = true
	fmt.Println("status changed")
}

func (l *List) printTasks() {
	for idx, task := range l.items {
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

func loadData(l *List, file *os.File) {
	// get file size in bytes and allocate buffer
	finfo, err := file.Stat()
	check(err)
	
	fmt.Println(finfo.Size())
	if finfo.Size() == 0 {
		fmt.Println(finfo.Size())
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
	file.WriteString("Content;Notes;x\n")
}

func main() {
	// Open File os.O_APPEND|
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
