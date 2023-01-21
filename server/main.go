package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"packages/exercise"
	"packages/sort"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type ToDo struct{
	ID int `json:id`
	Title string `json:title`
	Body string `json:body`
	Done bool `json:done`
}

type GetRequest struct{
	StatsFormat int `json:statsFormat`
	ExceriseName string `json:exceriseName`
}

type DeleteRequest struct{
	WholeExcerise bool `json:wholeExcerise`
	ID int `json:ID`
	ExceriseName string `json:exceriseName`
}

func initaliseListOfKeys(){
	var kys []string
	for key, _ := range userInfo.Exercises{
		kys = append(kys, key)
	}
	kys = sort.Sort(kys)
	keys = kys
	initialiseExceriseNamesAndIDs()
}

func printAllExceriseNames(){
	for v, d := range keys{
		fmt.Printf("%d. %s\n", v, d)
	}
}

func initialiseExceriseNamesAndIDs(){
	for v, d := range keys{
		if entry, ok := userInfo.Exercises[d]; ok {
			entry.ID = v
			userInfo.Exercises[d] = entry
		}
	}
}

var userInfo *exercise.UsersExercise;
var currentExercise string;
var keys []string

/*
Returns false is the current excerise doesn't exist
*/
func updateCurrentExercise(newName string) bool{
	return findExercise(newName);
}

func findExercise(exerciseRequested string) bool {
	for v, d := range keys{
		fmt.Println(v, exerciseRequested, d)
		if(strconv.Itoa(v) == exerciseRequested || d == exerciseRequested){
			currentExercise = d;
			return true;
		}
	}
	currentExercise = exerciseRequested;
	return false;
}

func addNewExerciseInstant(c *fiber.Ctx) error{
	userInput := &exercise.UserInput{}

	if err := c.BodyParser(userInput); err != nil{
			return err;
	}
	exists := updateCurrentExercise(userInput.ExceriseName);
	
	if(exercise.UserRequestNewIteration(userInfo, currentExercise, *exercise.UserTempIteration(userInput.Reps, userInput.Weights, userInput.Sets, userInput.Weight, 
		userInput.DaysAgo ,userInput.Note))){
			fmt.Println("Instance of ", currentExercise ," has been successfully added")
			if(!exists){
				initaliseListOfKeys()
			}
			return c.SendString("OK")
	}
	return c.Status(401).SendString("Error creating exercise") ;
	
}

func deleteEntireExercise(c *fiber.Ctx) error{
	deleteRequest := &DeleteRequest{}
	if err := c.BodyParser(deleteRequest); err != nil{
		return c.Status(400).SendString("Invalid body")
	}
	exists := findExercise(deleteRequest.ExceriseName)
	if(!exists){
		return c.Status(400).SendString("Invalid exercise")
	}
	if(!deleteRequest.WholeExcerise){
			if (exercise.UserDeletionRequest(userInfo, deleteRequest.ID, currentExercise, exercise.ExerciseInstance)) {
				fmt.Printf("Exercise instance deletion successful\n") 
				return c.SendString("OK");
			}
			fmt.Printf("Exercise instance deletion unsuccessful, couldn't find ID\n")
			return c.Status(400).SendString("Invalid ID");
	}
	exercise.UserDeletionRequest(userInfo, 0, currentExercise, exercise.EntireExercise)
	initaliseListOfKeys()
	return c.SendString("OK");
}

func userGetsAllExerciseNames(c *fiber.Ctx) error{
	return c.JSON(keys)
}

func getJSONOfExceriseAll(c *fiber.Ctx) error{
	return c.JSON(userInfo.Exercises)
}

func getJSONOfExcerise(c *fiber.Ctx) error{
	getRequest := &GetRequest{}
	if err := c.BodyParser(getRequest); err != nil{
		return err;
	}
	exists := findExercise(getRequest.ExceriseName)
	if(!exists){
		return c.Status(401).SendString("Invalid exercise")
	}
	return c.JSON(exercise.FetchExerciseObject(userInfo, currentExercise))
}

/*
	Just spits out the data to the user...could be useful from a logging point of view for them
	Therefore this is in a non JSON format and is very simple
*/
func userRequestToViewExerciseLog(c *fiber.Ctx) error{
	getRequest := &GetRequest{}
	if err := c.BodyParser(getRequest); err != nil{
		return err;
	}
	exists := updateCurrentExercise(getRequest.ExceriseName)
	if(!exists){
		return c.Status(401).SendString("Invalid exercise")
	}
	return c.JSON(exercise.ViewAnExercise(userInfo, currentExercise, exercise.StatsFormat(getRequest.StatsFormat)));
}

//Temp User set up....this will need to be changed when I support a user more than myself
	//With the addition of an ID to represent the user...rather than looking for the name = x x.json.
		//This will also need a rework.....we will need to store the username as a constant eventually.  
func tempLoadUser() * exercise.UsersExercise{
	
		var ue = &exercise.UsersExercise{Exercises: map[string]exercise.Exercise{}};
		var data []byte;
		for{
			username := "Adam"
			fmt.Println(username+".json")
			dat, err := os.ReadFile("userFolder/"+username+".json")
			if err != nil{
				fmt.Printf("Not a valid username, please try again\n")
				continue;
			}
			data = dat;
			break;
		}
		
		json.Unmarshal([]byte(data), &ue)
		// backupfile()
		return ue;
	}

func setupRoutes(app *fiber.App) {
	app.Get("/healthcheck", func(c *fiber.Ctx) error{
		return c.SendString("OK")
	})
	app.Get("/api/getExerciseLog", userRequestToViewExerciseLog)
	app.Get("/api/getExerciseAll", userGetsAllExerciseNames)
	app.Get("/api/getJSONOfExcerise", getJSONOfExcerise)
	app.Get("/api/getJSONOfExceriseAll", getJSONOfExceriseAll)
	app.Post("/api/addNewExerciseInstant", addNewExerciseInstant)
	app.Post("/api/deleteEntireExercise", deleteEntireExercise)
}

func main(){
	userInfo = tempLoadUser()
	initaliseListOfKeys()
	//coreFunctionLoop()

	fmt.Printf("Hello world!")

	app := fiber.New()

	todos := []ToDo{}

	app.Use(cors.New(cors.Config{
		AllowOrigins: "http://localhost:5173",
		AllowHeaders: "Origin, Content-Type, Accept",
	}))

	//fmt.Println(todos[0])

	

	setupRoutes(app)

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