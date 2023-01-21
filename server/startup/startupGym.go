package startup

import (
	"encoding/json"
	"fmt"
	"os"
	"packages/exercise"
	"packages/inputters"
	"strconv"
	"strings"
	"time"
)

var username string;

func GetDate() string{
	currentTime := time.Now()
	day := currentTime.Local().Day()
	month := currentTime.Local().Month()
	year := currentTime.Local().Year()
	return strconv.Itoa(day) +"/"+strconv.Itoa(int(month))+"/"+strconv.Itoa(year);
	
} 

func backupfile(){
	dat, _ := os.ReadFile(username+".json")
	dateFileFormat := strings.Replace(GetDate(), "/", ":", 2);
	err := os.WriteFile(username+"-BACKUP:" + dateFileFormat+ ".json", dat, 0644)
	if err != nil{
		fmt.Printf("Unsuccesful backup of: %s" ,username)
	}
}

func SaveUser(ue *exercise.UsersExercise){
	file, _ := json.MarshalIndent(ue, "", " ")
	err := os.WriteFile(username+".json", file, 0644)
    if err == nil{
		fmt.Printf("Succesful save for user %s\n", username);
	}
}

func newUser() *exercise.UsersExercise{
	username = inputters.FetchString("Please enter what you would like your username to be");
	var ue = &exercise.UsersExercise{Exercises: map[string]exercise.Exercise{}};
	return ue;
}

func loadUser() * exercise.UsersExercise{
	
	var ue = &exercise.UsersExercise{Exercises: map[string]exercise.Exercise{}};
	var data []byte;
	for{
		username = inputters.FetchString("Please enter your username:");
		fmt.Println(username+".json")
		dat, err := os.ReadFile(username+".json")
		if err != nil{
			fmt.Printf("Not a valid username, please try again\n")
			continue;
		}
		data = dat;
		break;
	}
	
	json.Unmarshal([]byte(data), &ue)
	backupfile()
	return ue;
}

func StartUp() *exercise.UsersExercise {
	fmt.Printf("Welcome!\n")
	integerInput := inputters.FetchInteger("[1] Log in [2] Create Account",2)
	fmt.Println(integerInput)
	if(integerInput ==1){
		return loadUser();
	}
	return newUser();
}