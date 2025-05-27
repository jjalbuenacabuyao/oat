package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

const todoFile = "todos.csv"

type Todo struct {
	ID   int
	Task string
}

func loadTodos() ([]Todo, error) {
	file, err := os.OpenFile(todoFile, os.O_CREATE | os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader :=csv.NewReader(file)
	lines, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var todos []Todo
	for _, line := range lines {
		id, err := strconv.Atoi(line[0])
		if err != nil {
			continue
		}
		todos = append(todos, Todo{ID: id, Task: line[1]})
	}

	return todos, nil
}

func saveTodos(todos []Todo) error {
	file, err := os.OpenFile(todoFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	for _, todo := range todos {
		writer.Write([]string{strconv.Itoa(todo.ID), todo.Task})
	}
	writer.Flush()
	return writer.Error()
}

func addTask(task string) error {
	todos, err := loadTodos()
	if err != nil {
		return err
	}
	id := 1
	if len(todos) > 0 {
		id = todos[len(todos)-1].ID + 1
	}
	todos = append(todos, Todo{ID: id, Task: task})
	return saveTodos(todos)
}

func listTasks() error {
	todos, err := loadTodos()
	if err != nil {
		return err
	}
	fmt.Println("Your TODOs:")
	for _, todo := range todos {
		fmt.Printf("%d: %s\n", todo.ID, todo.Task)
	}
	return nil
}

func deleteTask(id int) error {
	todos, err := loadTodos()
	if err != nil {
		return err
	}
	var updated []Todo
	for _, todo := range todos {
		if todo.ID != id {
			updated = append(updated, todo)
		}
	}
	return saveTodos(updated)
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: todo [add|list|delete] [task|id]")
		return
	}

	switch strings.ToLower(os.Args[1]) {
	case "add":
		if len(os.Args) < 3 {
			fmt.Println("Please provide a task to add.")
			return
		}
		task := strings.Join(os.Args[2:], " ")
		if err := addTask(task); err != nil {
			fmt.Println("Error adding task:", err)
		}
		fmt.Println("Task added.")
	case "list":
		if err := listTasks(); err != nil {
			fmt.Println("Error listing tasks:", err)
		}
	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Please provide the ID of the task to delete.")
			return
		}
		id, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Println("Invalid ID provided.")
			return
		}
		if err := deleteTask(id); err != nil {
			fmt.Println("Error deleting task:", err)
		} else {
			fmt.Println("Task deleted.")
		}
	default:
		fmt.Println("Unknown command.")
	}
}
