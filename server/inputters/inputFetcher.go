package inputters

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

//Will be used for later reformatting....have this printed off the singular input function we will have eventually.
func reprintStatement(statement, typeNeededStatement string){
	fmt.Printf("%s \n%s", typeNeededStatement, statement);
}
//This will need enums etc etc later on 

func FetchDouble(statement string) float64{
	valueFound := false;
	var returnValue float64;
	for(!valueFound) {
		fmt.Printf("%s\n",statement);
		userInput  := bufio.NewReader(os.Stdin)
		userVal, err := userInput.ReadString('\n')
		
		if err != nil {
			fmt.Printf("Enter a valid double\n");
			continue;
		}
		input := strings.TrimSpace(userVal)
		returnValue, err = strconv.ParseFloat(input, 64)
		if err != nil {
			fmt.Printf("Enter a valid double\n");
			continue;
		}
		valueFound = true;
	}
	return returnValue;
}

func FetchBoolean(statement string) bool{
	for{
		fmt.Printf("%s",statement + " (Y\\N)\n");
		userInput  := bufio.NewReader(os.Stdin)
		userVal, err := userInput.ReadString('\n')
		if err != nil {
			fmt.Printf("Please enter Yes (Y) or No (no)\n");
		}
		input := strings.ToUpper(strings.TrimSpace(userVal))

		if (input == "YES" || input == "Y"){
			return true;
		}
		if (input == "NO" || input == "N"){
			return false;
		}
		fmt.Printf("Please enter Yes (Y) or No (no)\n");
	}
}

func FetchInteger(statement string, intRange int) int{
	valueFound := false;
	var returnValue int;
	for(!valueFound) {
		fmt.Printf("%s\n",statement);
		userInput  := bufio.NewReader(os.Stdin)
		userVal, err := userInput.ReadString('\n')
		if err != nil {
			fmt.Printf("Enter a value between 0 and %d\n" , intRange);
			continue;
		}

		input := strings.TrimSpace(userVal)
		returnValue, err = strconv.Atoi(input)
		if err != nil {
			fmt.Printf("Enter a value between 0 and %d\n" , intRange);
			continue;
		}

		if(returnValue <= 0 || returnValue > intRange){
			fmt.Printf("Enter a value between 0 and %d\n" , intRange);
			continue;
		}
		// fmt.Println(returnValue)
		valueFound = true;
	}
	return returnValue;
}

func FetchString(statement string) string{
	var returnValue string;
	var err error;
	for {
		fmt.Printf("%s\n",statement);
		userInput  := bufio.NewReader(os.Stdin)
		returnValue, err = userInput.ReadString('\n')
		if err != nil {
			fmt.Printf("Please enter a valid string\n");
			continue;
		}
		break;

	}
	
	return returnValue[0:len(returnValue)-1];
}

func FetchArray(statement string) []float64{
	var inputString string;
	var err error;
	for  {
		fmt.Printf("%s\n",statement);
		userInput  := bufio.NewReader(os.Stdin)
		inputString, err = userInput.ReadString('\n')
		if err != nil {
			fmt.Printf("Please enter a valid array\n");
			continue;
		}
		spaceResult := ArraySplitter(inputString, " ");
		commaResult := ArraySplitter(inputString, ",");
		if(commaResult[0] == -1 && spaceResult[0] == -1){
			fmt.Printf("Please enter a valid array, and therefore numbers only\n");
			continue;
		}
		if(spaceResult[0] != -1){
			return spaceResult;
		}
		return commaResult;
	}
}

func ArraySplitter(temp string, spiltter string) []float64{
	var finalArray []float64
	var integersArray []float64
	inputString := strings.Split(temp[0:len(temp)-1], spiltter)

	for i := 0; i < len(inputString); i++ {
		xstr,_ := strconv.ParseFloat(inputString[i], 64)
		if(xstr == 0){
			return []float64{-1};
		}
		integersArray = append(integersArray , xstr )
	}

	finalArray = integersArray
	return finalArray;
}
