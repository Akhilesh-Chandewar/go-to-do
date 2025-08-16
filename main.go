package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"gopkg.in/mgo.v2/bson"
)

type Todo struct {
	ID        primitive.ObjectID `json:"id,omitempty" bson:"_id"`
	Completed bool               `json:"completed" bson:"completed"`
	Body      string             `json:"body" bson:"body"`
}

var todoCollection *mongo.Collection

func main() {
	evnError := godotenv.Load()
	if evnError != nil {
		log.Fatal("Error loading .env file", evnError)
	}

	mongodb_uri := os.Getenv("MONGO_URI")
	clientOptions := options.Client().ApplyURI(mongodb_uri)
	client, mongoError := mongo.Connect(context.Background(), clientOptions)
	if mongoError != nil {
		log.Fatal("Error connecting to MongoDB", mongoError)
	}

	defer client.Disconnect(context.Background())

	connetionError := client.Ping(context.Background(), nil)
	if connetionError != nil {
		log.Fatal("Error connecting to MongoDB", connetionError)
	}

	fmt.Println("Connected to MongoDB")

	todoCollection = client.Database("todo").Collection("todos")

	app := fiber.New()
	port := os.Getenv("PORT")
	if port == "" {
		port = "8000"
	}

	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Todo API is running!"})
	})
	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodo)
	app.Put("/api/todos/:id", updateTodo)
	app.Delete("/api/todos/:id", deleteTodo)

	fmt.Println("Listening on port " + port)
	log.Fatal(app.Listen("0.0.0.0:" + port))
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo

	cursor, err := todoCollection.Find(context.Background(), bson.M{})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error fetching todos from MongoDB",
		})
	}
	defer cursor.Close(context.Background())

	if err := cursor.All(context.Background(), &todos); err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error decoding todos",
		})
	}

	return c.Status(http.StatusOK).JSON(fiber.Map{"data": todos})
}

func createTodo(c *fiber.Ctx) error {
	todo := new(Todo)
	if parsingError := c.BodyParser(todo); parsingError != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Error parsing request body",
		})
	}

	if todo.Body == "" {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Todo body cannot be empty",
		})
	}

	insertResult, insertError := todoCollection.InsertOne(context.Background(), todo)
	if insertError != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{
			"error": "Error inserting todo into MongoDB",
		})
	}

	todo.ID = insertResult.InsertedID.(primitive.ObjectID)

	return c.Status(http.StatusCreated).JSON(fiber.Map{
		"data": insertResult,
		"todo": todo,
	})
}

func updateTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId , error :=  primitive.ObjectIDFromHex(id)
	if error != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}
	filter := bson.M{"_id": objectId}
	update := bson.M{"$set": bson.M{"completed": true}}
	result, err := todoCollection.UpdateOne(context.Background(), filter, update)
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

func deleteTodo(c *fiber.Ctx) error {
	id := c.Params("id")
	objectId , error :=  primitive.ObjectIDFromHex(id)
	if error != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}
	filter := bson.M{"_id": objectId}
	result, err := todoCollection.DeleteOne(context.Background(), filter)
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
