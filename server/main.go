package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"packages/exercise"
	"packages/logger"
	"packages/sort"
	"strconv"

	"golang.org/x/exp/slices"
)

type GetRequest struct {
	StatsFormat  int    `json:statsFormat`
	ExceriseName string `json:exceriseName`
}

type CompareRequest struct {
	ExceriseAttribute int64  `json:exceriseAttribute`
	ExceriseName      string `json:exceriseName`
}

type DeleteRequest struct {
	WholeExcerise bool   `json:wholeExcerise`
	ID            int    `json:ID`
	ExceriseName  string `json:exceriseName`
}

func initaliseListOfKeys() {
	var kys []string
	for key, _ := range userInfo.Exercises {
		kys = append(kys, key)
	}
	kys = sort.Sort(kys)
	keys = kys
	initialiseExceriseNamesAndIDs()
}

func printAllExceriseNames() {
	for v, d := range keys {
		fmt.Printf("%d. %s\n", v, d)
	}
}

func initialiseExceriseNamesAndIDs() {
	for v, d := range keys {
		if entry, ok := userInfo.Exercises[d]; ok {
			entry.ID = v
			userInfo.Exercises[d] = entry
		}
	}
}

var userInfo *exercise.UsersExercise
var currentExercise string
var keys []string

/*
Returns false is the current excerise doesn't exist
*/
func updateCurrentExercise(newName string) bool {
	return findExercise(newName)
}

func findExercise(exerciseRequested string) bool {
	for v, d := range keys {
		//fmt.Println(v, exerciseRequested, d)
		if strconv.Itoa(v) == exerciseRequested || d == exerciseRequested {
			currentExercise = d
			return true
		}
	}
	currentExercise = exerciseRequested
	return false
}

func addNewExerciseInstantMainHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data exercise.UserInput

	err = json.Unmarshal(body, &data)

	defer r.Body.Close()
	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	exists := updateCurrentExercise(data.ExceriseName)

	if exercise.UserRequestNewIteration(userInfo, currentExercise, *exercise.UserTempIteration(data.Reps, data.Weights, data.Sets, data.Weight,
		data.DaysAgo, data.Note)) {
		logger.LogInfo("Instance of ", currentExercise, " has been successfully added")
		if !exists {
			initaliseListOfKeys()
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte("Fail to add exercise iteration."))
	logger.LogInfo(http.StatusBadRequest, "Fail to add exercise iteration.")

}

func deleteEntireExerciseMainHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data DeleteRequest

	err = json.Unmarshal(body, &data)

	defer r.Body.Close()
	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	exists := findExercise(data.ExceriseName)
	if !exists {
		logger.LogInfo("Invalid exercise", data.ExceriseName)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !data.WholeExcerise {
		if exercise.UserDeletionRequest(userInfo, data.ID, currentExercise, exercise.ExerciseInstance) {
			logger.LogInfo("Exercise", data.ExceriseName, "instance deletion successful (ID) ", data.ID)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Exercise instance odeletion successful"))
			return
		}
		logger.LogInfo("Exercise instance deletion unsuccessful (", data.ExceriseName, "), couldn't find ID: ", data.ID, " \n")
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid ID " + strconv.Itoa(data.ID) + " For " + data.ExceriseName))
		return
	}

	exercise.UserDeletionRequest(userInfo, 0, currentExercise, exercise.EntireExercise)
	initaliseListOfKeys()
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Exercise successfully deleted"))
}

func userGetsAllExerciseNamesMainHTTP(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(keys)

	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func getJSONOfExceriseAllMainHTTP(w http.ResponseWriter, r *http.Request) {
	jsonData, err := json.Marshal(userInfo.Exercises)

	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonData)
}

func getJSONOfExceriseMainHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data GetRequest

	err = json.Unmarshal(body, &data)

	defer r.Body.Close()
	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	exists := updateCurrentExercise(data.ExceriseName)

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("This exercise name does not exist"))
		logger.LogInfo(http.StatusBadRequest, "Excerise name does not exist for user")
	}

	jsonData, err := json.Marshal(exercise.FetchExerciseObject(userInfo, currentExercise))

	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsonData)
}

func userRequestToViewExerciseLogMainHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data GetRequest

	err = json.Unmarshal(body, &data)

	defer r.Body.Close()
	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	exists := updateCurrentExercise(data.ExceriseName)

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("This exercise name does not exist"))
		logger.LogInfo(http.StatusBadRequest, "Excerise name does not exist for user")
	}

	jsonData, err := json.Marshal(exercise.ViewAnExercise(userInfo, currentExercise, exercise.StatsFormat(data.StatsFormat)))

	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsonData)

}

/*
	NOT YET IMPLEMENTED
*/
func getJSONComparisionMainHTTP(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}

	var data CompareRequest

	err = json.Unmarshal(body, &data)

	defer r.Body.Close()
	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	exists := updateCurrentExercise(data.ExceriseName)

	if !exists {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid exercise"))
		logger.LogInfo(http.StatusBadRequest, "Excerise name does not exist for user")
		return
	}
	attributeExists := slices.Contains(exercise.AttributesList, data.ExceriseAttribute)
	if !attributeExists {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid attribute"))
		logger.LogInfo(http.StatusBadRequest, "Invalid Attribute")
		return
	}
	//Need to link into GenerateComparisionObject function
	returnArray := exercise.GenerateComparisionObject(userInfo, data.ExceriseName, exercise.ExerciseAttribute(data.ExceriseAttribute))
	fmt.Println(returnArray)

	jsonData, err := json.Marshal(returnArray)

	if err != nil {
		logger.LogError(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	w.Write(jsonData)

}

//Temp User set up....this will need to be changed when I support a user more than myself
//With the addition of an ID to represent the user...rather than looking for the name = x x.json.
//This will also need a rework.....we will need to store the username as a constant eventually.
func tempLoadUser() *exercise.UsersExercise {

	var ue = &exercise.UsersExercise{Exercises: map[string]exercise.Exercise{}}
	var data []byte
	for {
		username := "Adam"
		fmt.Println(username + ".json")
		dat, err := os.ReadFile("userFolder/" + username + ".json")
		if err != nil {
			fmt.Printf("Not a valid username, please try again\n")
			continue
		}
		data = dat
		break
	}

	json.Unmarshal([]byte(data), &ue)
	// backupfile()
	return ue
}

func setUpRoutesMainLine() {
	http.HandleFunc("/healthcheck", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Add("Access-Control-Allow-Methods", "GET,HEAD,OPTIONS,POST,PUT")
		w.Header().Add("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
		w.Write([]byte{})
	})

	http.HandleFunc("/api/getExerciseLog", userRequestToViewExerciseLogMainHTTP)
	http.HandleFunc("/api/getExerciseAll", userGetsAllExerciseNamesMainHTTP)
	http.HandleFunc("/api/getJSONOfExcerise", getJSONOfExceriseMainHTTP)
	http.HandleFunc("/api/getJSONOfExceriseAll", getJSONOfExceriseAllMainHTTP)

	http.HandleFunc("/api/getJSONComparision", getJSONComparisionMainHTTP)

	http.HandleFunc("/api/addNewExerciseInstant", addNewExerciseInstantMainHTTP)

	http.HandleFunc("/api/deleteEntireExercise", deleteEntireExerciseMainHTTP)
}

func main() {
	logger.InitLogger()
	userInfo = tempLoadUser()
	initaliseListOfKeys()

	setUpRoutesMainLine()
	http.ListenAndServe(":8080", nil)

}
