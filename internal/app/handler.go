package app

import (
	"errors"
	"fmt"

	"github.com/AlekseyAytov/skillrock-todo/internal/models/master"
	"github.com/AlekseyAytov/skillrock-todo/internal/models/task"
	"github.com/gofiber/fiber/v2"
)

// ToDoAPI some comment
type ToDoAPI struct {
	api        *fiber.App
	taskmaster *master.TaskMaster
}

// NewToDoAPI some comment
func NewToDoAPI(taskmaster *master.TaskMaster) *ToDoAPI {
	td := ToDoAPI{
		api:        fiber.New(),
		taskmaster: taskmaster,
	}
	td.endpoints()
	return &td
}

// StartServer some comment
func (td *ToDoAPI) StartServer(socket string) error {
	return td.api.Listen(socket)
}

func (td *ToDoAPI) endpoints() {
	g1 := td.api.Group("tasks")
	g1.Post("/", td.addTask)
	g1.Get("/", td.getAllTasks)
	g1.Put("/:id", td.changeTask)
	g1.Delete("/:id", td.deleteTask)
}

func (td *ToDoAPI) addTask(c *fiber.Ctx) error {
	var t task.TaskHeads
	if err := c.BodyParser(&t); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}
	err := td.taskmaster.Add(t)
	if err != nil {
		if errors.Is(err, master.ErrBadStatus) || errors.Is(err, master.ErrEmptyTitle) {
			return c.Status(fiber.StatusBadRequest).SendString(err.Error())
		}
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusCreated).SendString("task was successfully created")
}

func (td *ToDoAPI) getAllTasks(c *fiber.Ctx) error {
	list, err := td.taskmaster.GetAll()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.JSON(list)
}

func (td *ToDoAPI) changeTask(c *fiber.Ctx) error {
	id := c.Params("id")

	var t task.TaskHeads
	if err := c.BodyParser(&t); err != nil {
		return c.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	err := td.taskmaster.UpdateBy(id, t)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("task with id %s successfully updated", id))
}

func (td *ToDoAPI) deleteTask(c *fiber.Ctx) error {
	id := c.Params("id")
	err := td.taskmaster.Delete(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}
	return c.Status(fiber.StatusOK).SendString(fmt.Sprintf("task with id %s deleted", id))
}
