package main

import "fmt"

// "log"
// "os"
// "strconv"

// "github.com/gofiber/fiber/v2"
// "github.com/joho/godotenv"

// type Todo struct {
// 	ID        int    `json:"id"`
// 	Completed bool   `json:"completed"`
// 	Body      string `json:"body"`
// }

func mainWithoutdb() {
	fmt.Println("without db apis")
	// 	app := fiber.New()

	// 	err := godotenv.Load()
	// 	if err != nil {
	// 		log.Fatal("Error loading .env file" , err)
	// 	}

	// 	port := os.Getenv("PORT")

	// 	todos := []Todo{}

	// 	app.Get("/", func(c *fiber.Ctx) error {
	// 		return c.JSON(fiber.Map{"message": "Todo API is running!"})
	// 	})

	// 	app.Get("/api/todos", func(c *fiber.Ctx) error {
	// 		return c.JSON(fiber.Map{"data": todos})
	// 	})

	// 	app.Get("/api/todos/:id", func(c *fiber.Ctx) error {
	// 		id, err := strconv.Atoi(c.Params("id"))
	// 		if err != nil {
	// 			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	// 		}

	// 		for _, todo := range todos {
	// 			if todo.ID == id {
	// 				return c.JSON(fiber.Map{"data": todo})
	// 			}
	// 		}

	// 		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	// 	})

	// 	app.Post("/api/todos", func(c *fiber.Ctx) error {
	// 		newTodo := Todo{}
	// 		if err := c.BodyParser(&newTodo); err != nil {
	// 			return err
	// 		}
	// 		if newTodo.Body == "" {
	// 			return c.Status(400).JSON(fiber.Map{"error": "Body is required"})
	// 		}

	// 		newTodo.ID = len(todos) + 1
	// 		newTodo.Completed = false
	// 		todos = append(todos, newTodo)

	// 		return c.Status(201).JSON(fiber.Map{"data": newTodo})
	// 	})

	// 	app.Put("/api/todos/:id", func(c *fiber.Ctx) error {
	// 		id, err := strconv.Atoi(c.Params("id"))
	// 		if err != nil {
	// 			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	// 		}

	// 		updateData := Todo{}
	// 		if err := c.BodyParser(&updateData); err != nil {
	// 			return err
	// 		}

	// 		for i, todo := range todos {
	// 			if todo.ID == id {
	// 				if updateData.Body != "" {
	// 					todo.Body = updateData.Body
	// 				}
	// 				todo.Completed = updateData.Completed
	// 				todos[i] = todo
	// 				return c.JSON(fiber.Map{"data": todo})
	// 			}
	// 		}

	// 		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	// 	})

	// 	app.Delete("/api/todos/:id", func(c *fiber.Ctx) error {
	// 		id, err := strconv.Atoi(c.Params("id"))
	// 		if err != nil {
	// 			return c.Status(400).JSON(fiber.Map{"error": "Invalid ID"})
	// 		}

	// 		for i, todo := range todos {
	// 			if todo.ID == id {
	// 				todos = append(todos[:i], todos[i+1:]...)
	// 				return c.SendStatus(204)
	// 			}
	// 		}

	// 		return c.Status(404).JSON(fiber.Map{"error": "Todo not found"})
	// 	})

	// log.Fatal(app.Listen(":" + port))
}
