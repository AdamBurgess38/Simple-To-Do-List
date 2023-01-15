package main

import (
	"fmt"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type ToDo struct{
	ID int `json:id`
	Title string `json:title`
	Body string `json:body`
	Done bool `json:done`
}

func main(){
	fmt.Printf("Hello world!")

	app := fiber.New()

	todos := []ToDo{}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	//fmt.Println(todos[0])

	app.Get("/healthcheck", func(c *fiber.Ctx) error{
		return c.SendString("OK")
	})

	app.Post("/api/todos", func(c *fiber.Ctx) error{
		todo := &ToDo{}

		if err := c.BodyParser(todo); err != nil{
			return err;
		}

		todo.ID = len(todos)+1

		todos = append(todos, *todo)


		
		return c.JSON(todos)

	})

	app.Patch("/api/todos/:id/done", func(c *fiber.Ctx) error {
		var ID int
		ID, err := c.ParamsInt("id");
		
		if err != nil{
			return c.Status(401).SendString("Invalid String")
		}
		

		for i, t := range todos{
			if (t.ID == ID){
				todos[i].Done = true
				break;
			}
		}
		
		return c.JSON(todos)
	})


	app.Get("/api/todos/", func(c *fiber.Ctx) error{

		fmt.Println(todos)
		fmt.Println(c.JSON(todos))
		return c.JSON(todos)
	})
	log.Fatal(app.Listen(":4000"))
}