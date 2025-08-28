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
	if t.status {
		char_status = 'X'
	} else {
		char_status = 'O'
	}

	fmt.Printf("%s\t%c\n", t.content, char_status)
}

//-------------------------------------------------------------------------------------

// stores Tasks
type TaskContainer struct {
	// slice for Tasks
	slc []Task
}

func (tc *TaskContainer) addTask(t Task) {
	if len(tc.slc) >= cap(tc.slc) {
		fmt.Printf("Cannot store more task. Max amount: %d", capacity)
	}
	// append can resize the underlying array so we check before
	tc.slc = append(tc.slc, t)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func loadData(tc *TaskContainer) {
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
		if parts[2] == "x" {
			status = false
		} else if parts[2] == "o" {
			status = true
		} else {
			fmt.Println("Status data in file is wrong")
			panic(status)
		}

		t := Task{
			content: parts[0],
			notes:   parts[1],
			status:  status,
		}

		tc.addTask(t)
	}
}

func main() {
	task_array := [capacity]Task{}
	// create a slice of array
	task_list := TaskContainer{slc: task_array[:0]}
	loadData(&task_list)

}
