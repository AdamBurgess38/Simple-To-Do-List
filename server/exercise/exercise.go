package exercise

import (
	"fmt"
	"strconv"
	"strings"
	"time"
)

type UsersExercise struct{
	Exercises map[string]Exercise
}

type Exercise struct {
	ID int
    Iterations []Iteration
}

type UserInput struct {
	ExceriseName string `json:"name"`;
	Reps []float64 `json:"reps"`
	Weights []float64 `json:"weights"`
	Sets int `json:"sets"`
	Weight float64 `json:"weight"`
	DaysAgo int `json:"daysAgo"`
	Note string `json:"note"`
}

type Iteration struct {
	Reps []float64 
	Weights []float64
	Variances []float64
	ID int
	Sets int
	Weight float64
	Date string
	Note string
	TotalWeight float64
	AverageWeight float64
	AverageRep float64
	AverageWeightRepTotal float64
}

type StatsFormat int64

//0...n
const (
	StandardStats StatsFormat = iota
	AverageOverall     
	SimpleStats     
	MostRecent       
)

type DeletionType int64

const (
	EntireExercise DeletionType = iota
	ExerciseInstance        
)

func ViewAnExercise(ue *UsersExercise, requestedExercise string, t StatsFormat) string{
	entry, ok := ue.Exercises[requestedExercise];	
	if !ok {
		return "Invalid excerise";
	}
    switch t {
		case StandardStats:
			return fetchStandardStats(entry)
		case AverageOverall:
			return fetchAverageOverall(entry)
		case SimpleStats:
			return fetchSimpleStats(entry)
		case MostRecent:
			return fetchMostRecent(entry)
    }
	return "Invalid stats type"
}

func BoltOnSeperator() string{
	return "-----------------------------------------------------------\n";
} 

func FetchExerciseObject(ue *UsersExercise, requestedExercise string) []Iteration{
	entry, ok := ue.Exercises[requestedExercise];	
	if !ok {
		return []Iteration{};
	}
	return generateExerciseObject(entry)
}

func generateExerciseObject(ex Exercise) []Iteration{
	return ex.Iterations;

}

func fetchStandardStats(ex Exercise) string{
	returnString := ""
	for _ , iter := range ex.Iterations {
		returnString += instanceConverter(iter)
		returnString += BoltOnSeperator();
	}
	return returnString;
}

func instanceConverter(iter Iteration) string {
	return "ID: " + strconv.Itoa(iter.ID) + "\n" + 
	"Date: " + iter.Date + "\n" +
	"Planned Weight: " + strconv.FormatFloat(iter.Weight, 'f', 2, 64) + "\n" +
	"Number of Sets: " + strconv.Itoa(iter.Sets) + "\n" +
	"Weights per set: " + arrayToString(iter.Weights) + "\n" +
	"Reps per set: " + arrayToString(iter.Reps) + "\n" +
	"Average Weight: " + strconv.FormatFloat(iter.AverageWeight, 'f', 2, 64)+ "\n" +
	"Average reps: " + strconv.FormatFloat(iter.AverageRep, 'f', 2, 64) + "\n" +
	"Average total: " + strconv.FormatFloat(iter.AverageWeightRepTotal, 'f', 2, 64)+ "\n" +
	"Total Weight pushed: " + strconv.FormatFloat(iter.TotalWeight, 'f',2,64) + "\n"+
	"Note: " + iter.Note +"\n";
}

func fetchAverageOverall(ex Exercise) string{
	returnString := "";
	for _ , iter := range ex.Iterations {
		returnString += "Date: " + iter.Date + " Average Weight: " + strconv.FormatFloat(iter.AverageWeight, 'f', 2, 64) +" Average Rep: " + strconv.FormatFloat(iter.AverageRep, 'f', 2, 64) +"\n"
	}
	return returnString;
}

func fetchSimpleStats(ex Exercise) string{
	returnString := "";
	for _ , iter := range ex.Iterations {
		returnString += "ID: " + strconv.Itoa(iter.ID) + " Date: " + iter.Date + " Weights: " + arrayToString(iter.Weights) +" Reps: " + arrayToString(iter.Reps) +"\n"
	}
	return returnString;
}

func fetchMostRecent(ex Exercise) string{
	returnString := ""
	latestIndex := len(ex.Iterations);
	iter := ex.Iterations[latestIndex-1]
	
	returnString += instanceConverter(iter)
	returnString += BoltOnSeperator();
	
	return returnString;
}

func arrayToString(array []float64) string{
	returnString := ""
	for _ , x := range array{
		returnString += strconv.FormatFloat(x, 'f', 2, 64) +","
	}
	return returnString[0:len(returnString)-1];

}

func UserTempIteration(reps, weights []float64, sets int, weight float64, dateDifference int, note string) *UserInput{
	return &(
		UserInput{
			Reps: reps, 
			Weights: weights, 
			Sets : sets, 
			Weight: weight,
			DaysAgo: dateDifference,
			Note: note,
			});
}

func NewIteration(reps []float64, weights []float64, variances []float64, ID int, sets int, weight float64, date string, note string, totalWeight float64, averageRep float64, averageWeight float64, averageWeightRepTotal float64) *Iteration{
	return &(
		Iteration{
			Reps: reps, 
			Weights: weights, 
			Variances: variances,
			ID : ID, 
			Sets : sets, 
			Weight: weight,
			Date: date,
			Note: note, 
			TotalWeight: totalWeight,
			AverageRep: averageRep,
			AverageWeight: averageWeight,
			AverageWeightRepTotal: averageWeightRepTotal,
			});
}

func initialiseExcerise(ue *UsersExercise,name string){
	ue.Exercises[name] = Exercise{};
}

func addIteration(ue *UsersExercise, name string, x Iteration){
	fmt.Println(x)
	entry, ok := ue.Exercises[name];
	fmt.Println(entry)
	if ok {
		entry.Iterations = append(entry.Iterations,x)
	}	
	fmt.Println(entry)
	ue.Exercises[name] = entry;
}

func Map[T, V any](ts []T, fn func(T) V) []V {
    result := make([]V, len(ts))
    for i, t := range ts {
        result[i] = fn(t)
    }
    return result
}

func UserRequestNewIteration(ue *UsersExercise, name string, x UserInput) bool{
	if(len(x.Reps) != len(x.Weights)){
		return false;
	}
	entry, ok := ue.Exercises[name];

	if !ok {
		initialiseExcerise(ue, name)
	}	

	var newID int = len(entry.Iterations);
	var foundID bool = false

	for !foundID {
		foundID = true;
		for _, x := range entry.Iterations{
			if(x.ID == newID){
				newID++
				foundID = false;
				break;
			}
		}
	}
	var totalWeight float64 = 0
	var totalWeightRep float64
	var totalReps float64 = 0
	for i , w := range x.Weights{
		totalWeightRep += w * x.Reps[i]
		totalReps += x.Reps[i]
		totalWeight += w
	}


	addIteration(ue, name, *NewIteration(
		x.Reps, x.Weights, Map(x.Weights, func(item float64) float64 { return item - x.Weight }), newID, 
		len(x.Reps), 
		x.Weight, 
		strings.Replace((time.Now().Local().AddDate(0, 0, -x.DaysAgo)).Format("01-02-2006"),"/", ":", 2), 
		x.Note, totalWeightRep, totalReps/float64(len(x.Reps)), totalWeight/float64(len(x.Weights)), totalWeightRep/float64(len(x.Weights))))
	return true;
}

func UserDeletionRequest(ue *UsersExercise, requestedID int, name string, deleteType DeletionType) bool{
	switch deleteType {
		case EntireExercise:
			deleteEntireExercise(ue, name);
			return true;
		case ExerciseInstance:
			return deleteExerciseInstance(ue, requestedID, name)
    }
	return true;
}

func deleteExerciseInstance(ue *UsersExercise, requestedID int, name string) bool{
	for i , w := range ue.Exercises[name].Iterations{
		if(w.ID == requestedID){
			if entry, ok := ue.Exercises[name]; ok{
				result := remove(ue.Exercises[name].Iterations, i)
				entry.Iterations = result;
				ue.Exercises[name] = entry;
				return true;
			}
		}
	}
	return false;
}

func remove[T any](slice []T, s int) []T {
	return append(slice[:s], slice[s+1:]...)
}

func deleteEntireExercise(ue *UsersExercise, name string){
	delete(ue.Exercises, name)
}