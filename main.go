package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"
)

type PeopleAvailability struct {
	Name string
	Days []int
}

func main() {
	peopleAvailability := parseInputFile("input.txt")
	dayCounter := countDays(peopleAvailability)
	preferredDays := sortKeyByValue(dayCounter)
	excludedInPreferredDay := whosExcludedIn(preferredDays[0], peopleAvailability)

	fmt.Println("---")
	fmt.Println(peopleAvailability)
	fmt.Println(dayCounter)
	fmt.Println(preferredDays)
	fmt.Println(excludedInPreferredDay)
	fmt.Println("---")
}

func whosExcludedIn(thisDay int, peopleAvailability []PeopleAvailability) []string {
	var excluded []string

	for _, person := range peopleAvailability {
		isIt := isExcluded(thisDay, person.Days)
		if isIt {
			excluded = append(excluded, person.Name)
		}
	}

	return excluded
}

func isExcluded(thisDay int, days []int) bool {
	for _, day := range days {
		if day == thisDay {
			return false
		}
	}
	return true
}

func sortKeyByValue(input map[int]int) []int {
	type KeyValue struct {
		Key   int
		Value int
	}
	var keyValues []KeyValue

	for key, value := range input {
		keyValues = append(keyValues, KeyValue{Key: key, Value: value})
	}
	sort.Slice(keyValues, func(i, j int) bool {
		return keyValues[i].Value > keyValues[j].Value
	})

	var output []int
	for _, keyValue := range keyValues {
		output = append(output, keyValue.Key)
	}

	return output
}

func countDays(peopleAvailability []PeopleAvailability) map[int]int {
	dayCounter := make(map[int]int)
	for _, personAvailability := range peopleAvailability {
		for _, day := range personAvailability.Days {
			if dayCounter[day] == 0 {
				dayCounter[day] = 1
			} else {
				dayCounter[day] += 1
			}
		}
	}
	return dayCounter
}

func parseInputFile(input string) []PeopleAvailability {
	file, err := os.Open(input)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var persons []PeopleAvailability

	for scanner.Scan() {
		persons = extractPerson(scanner, persons)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return persons
}

func extractPerson(scanner *bufio.Scanner, persons []PeopleAvailability) []PeopleAvailability {
	splitted := strings.SplitN(scanner.Text(), ":", 2)
	person := new(PeopleAvailability)
	person.Name = splitted[0]
	days := strings.ReplaceAll(splitted[1], " ", "")
	for _, rawDay := range strings.Split(days, ",") {
		day, _ := strconv.Atoi(rawDay)
		person.Days = append(person.Days, day)
	}
	persons = append(persons, *person)
	return persons
}
