package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"packages/exercise"
	"packages/inputters"
	"packages/sort"
	"packages/startup"
	"strconv"
	"strings"
	"time"

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
func updateCurrentExercise(statement string) bool{
	return findExercise(inputters.FetchString(statement));
}

func findExercise(exerciseRequested string) bool {
	for v, d := range keys{
		// fmt.Println(v, exerciseRequested, d)
		if(strconv.Itoa(v) == exerciseRequested || d == exerciseRequested){
			currentExercise = d;
			return true;
		}
	}
	return false;
}

func userInputDeleteExercise(){
	printAllExceriseNames()
	if(!updateCurrentExercise("State the name of the exercise you would like to delete")){
		fmt.Printf("Exercise does not exist\n")
		return; 
	}
	if(inputters.FetchBoolean(("Would you like the entire exercise" +  currentExercise + "or a single iteration of it?"))){
		exercise.ViewAnExercise(userInfo, currentExercise, exercise.SimpleStats)
		ID := inputters.FetchInteger("Please enter the ID you would like to delete", 10000);
		if (exercise.UserDeletionRequest(userInfo, ID, currentExercise, exercise.ExerciseInstance)) {
			fmt.Printf("Exercise instance deletion successful\n") 
			return;
		}
		fmt.Printf("Exercise instance deletion unsuccessful, couldn't find ID\n")
		return;
	}
	exercise.UserDeletionRequest(userInfo, 0, currentExercise, exercise.EntireExercise)
}


func userInputNewExercise(){
	updateCurrentExercise("State the name of the exercise")
	fmt.Println("Adding iteration of " , currentExercise);
	weight := inputters.FetchDouble("Please state the weight you worked at");
	constantWeight := inputters.FetchBoolean("Was the weight constant throughout?");
	constantReps := inputters.FetchBoolean("Were the reps constant throughout?");
	var reps []float64
	var weights []float64
	var sets int
	if(constantWeight && constantReps){
		rep := inputters.FetchDouble("And what was this value?");
		sets := inputters.FetchInteger("For how many sets?",1000);
		for sets != 0{
			reps = append(reps, rep);
			weights = append(weights, weight)
			sets --;
		}
		
	}
	if(!constantWeight && constantReps){
		rep := inputters.FetchDouble("And what was this value?");
		weights = inputters.FetchArray("Please state the weight throughout the sets")
		for x := 0; x < len(weights); x++{
			reps = append(reps, rep)
		}
	}
	if(constantWeight && !constantReps){
		reps = inputters.FetchArray("Please state the reps throughout the sets")
		for x := 0; x < len(weights); x++{
			weights = append(weights, weight)
		}
	}
	if(!constantWeight && !constantReps){
		reps = inputters.FetchArray("Please state the reps throughout the sets")
		weights = inputters.FetchArray("Please state the weight throughout the sets")
	}

	var note string = ""
	if(inputters.FetchBoolean("Would you like to leave a note?")){
		note += inputters.FetchString("Please enter what you would like the note to be?")
	}

	var daysAgo = 0
	if(!inputters.FetchBoolean("Did you perform this exercise today?")){
		daysAgo = inputters.FetchInteger("How many days ago did you perform this exercise?", 365);
	}
	if(exercise.UserRequestNewIteration(userInfo, currentExercise, *exercise.UserTempIteration(reps, weights, sets, weight, 
		strings.Replace((time.Now().Local().AddDate(0, 0, -daysAgo)).Format("01-02-2006"),"/", ":", 2),note))){
			fmt.Println("Instance of " , currentExercise , " has been successfully added")
			return;
	}

	fmt.Println("Instance of " , currentExercise , " has been unsuccessfully added due to the number of reps and weights not alligning")
	
}

/*
Backend only version.
*/
func userRequestToViewAnExercise(){
	exists := updateCurrentExercise("What exercise would you like to view?")
	if(!exists){
		fmt.Printf("This exercise does not exist in your records, sorry");
	}
	choice := exercise.StatsFormat(inputters.FetchInteger("What format would you like?\n[1] Standard stats [2] Average overall [3] Simple stats  [4] Most recent", 4)-1);
	fmt.Println(exercise.ViewAnExercise(userInfo, currentExercise, choice))
}


func userGetsAllExerciseNames(c *fiber.Ctx) error{
	return c.JSON(keys)
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
	exists := findExercise(getRequest.ExceriseName)
	if(!exists){
		fmt.Printf("This exercise does not exist in your records, sorry");
	}
	fmt.Printf("we get to here")
	choice := exercise.StatsFormat(getRequest.StatsFormat)
	repsone := exercise.ViewAnExercise(userInfo, currentExercise, choice)
	fmt.Println(repsone)

	return c.JSON(repsone)
}

func userRequestToViewAllExercises(){
	choice := inputters.FetchInteger("[1] All names [2] All iterations of every exercise",2)
	if(choice == 1){
		printAllExceriseNames()
		return; 
	}
	for _, d := range keys{
		fmt.Println("Exercise " , d, "Instances:\n" ,exercise.ViewAnExercise(userInfo, d, exercise.StandardStats))
	}

}

func coreFunctionLoop(){
	var choice int
	for{
		choice = inputters.FetchInteger("[1] Add exercise instance [2] View an exercise [3] Delete an exercise [4] View all exercises [5] Save", 5)
		switch choice {
		case 1:
			userInputNewExercise()
		case 2:
			userRequestToViewAnExercise()
		case 3:
			userInputDeleteExercise()
		case 4:
			userRequestToViewAllExercises()
		case 5:
			startup.SaveUser(userInfo)
		}
	}
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
	app.Get("/api/getExerciseLog", userRequestToViewExerciseLog)
	app.Get("/api/getExerciseAll", userGetsAllExerciseNames)
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

	app.Get("/healthcheck", func(c *fiber.Ctx) error{
		return c.SendString("OK")
	})

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