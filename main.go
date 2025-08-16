package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Completed bool               `json:"completed" bson:"completed"`
	Body      string             `json:"body" bson:"body"`
}

var todoCollection *mongo.Collection

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Connect to MongoDB
	mongoURI := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal("Error connecting to MongoDB:", err)
	}
	defer client.Disconnect(context.Background())

	if err := client.Ping(context.Background(), nil); err != nil {
		log.Fatal("MongoDB ping failed:", err)
	}
	fmt.Println("✅ Connected to MongoDB")

	// Collection reference
	todoCollection = client.Database("todo").Collection("todos")

	// Fiber app
	app := fiber.New()

	// Middleware
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
	}))

	// Routes
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Todo API is running!"})
	})
	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Put("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}
	fmt.Println("🚀 Listening on port " + port)
	log.Fatal(app.Listen("0.0.0.0:" + port))
}

// GET /api/todos
func getTodos(c *fiber.Ctx) error {
	var todos []Todo
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	cursor, err := todoCollection.Find(ctx, bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching todos from MongoDB",
		})
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &todos); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error decoding todos",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"data": todos})
}

// POST /api/todos
func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)
	if err := c.BodyParser(todo); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	todo.Body = strings.TrimSpace(todo.Body)
	if todo.Body == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Todo body cannot be empty",
		})
	}

	todo.Completed = false
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	insertResult, err := todoCollection.InsertOne(ctx, todo)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error inserting todo into MongoDB",
		})
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(http.StatusCreated).JSON(fiber.Map{"todo": todo})
}

// PUT /api/todos/:id
func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	var updateData struct {
		Completed bool `json:"completed"`
	}
	if err := c.BodyParser(&updateData); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"completed": updateData.Completed}}

	result, err := todoCollection.UpdateOne(ctx, filter, update)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error updating todo in MongoDB",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Todo updated successfully",
		"result":  result,
	})
}

// DELETE /api/todos/:id
func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": objectId}
	result, err := todoCollection.DeleteOne(ctx, filter)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error deleting todo from MongoDB",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{
		"message": "Todo deleted successfully",
		"result":  result,
	})
}
