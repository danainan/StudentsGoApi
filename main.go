package main

import (
	"encoding/json"

	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type Student struct {
	ID    int    `json:"id"`
	Fname string `json:"fname"`
	Lname string `json:"lname"`
}

var students = []Student{
	{ID: 63103600, Fname: "Danainan", Lname: "Chamnanpaison"},
	{ID: 63103666, Fname: "Danainan2", Lname: "Chamnanpaison2"},
}

func getStudents(ctx *fiber.Ctx) error {
	return ctx.JSON(students)
}

func getStudent(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	for _, student := range students {
		if strconv.Itoa(student.ID) == id {
			return ctx.JSON(student)
		}
	}
	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
		"massage": "Student not found",
	})
}

func createStudent(ctx *fiber.Ctx) error {
	student := new(Student)
	if err := ctx.BodyParser(student); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Error 400"})
	}

	students = append(students, *student)
	return ctx.JSON(student)
}

func updateStudent(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	for i, student := range students {
		if strconv.Itoa(student.ID) == id {
			updatedStudent := new(Student)
			if err := ctx.BodyParser(updatedStudent); err != nil {
				return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Invalid form field"})
			}
			updatedStudent.ID = student.ID
			students[i] = *updatedStudent
			return ctx.JSON(updatedStudent)
		}
	}
	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Student Not Found"})
}

func deleteStudent(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	for i, student := range students {
		if strconv.Itoa(student.ID) == id {
			students = append(students[:i], students[i+1:]...)
			return ctx.SendStatus(fiber.StatusNoContent)
		}
	}
	return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Student Not Found"})
}

func main() {
	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
	})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowMethods: "GET, POST, PUT, DELETE",
	}))
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("hello World")
	})

	app.Get("/getstudents", getStudents)
	app.Get("/getstudent/:id", getStudent)
	app.Post("/createStudent", createStudent)
	app.Put("/updateStudent/:id", updateStudent)
	app.Delete("/deleteStudent/:id", deleteStudent)

	app.Use(logger.New())
	app.Use(requestid.New())

	app.Listen(":8080")
}
