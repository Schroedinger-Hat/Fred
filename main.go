package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type proposal struct {
	name string
	days []int
}

func (i proposal) Title() string { return i.name }
func (i proposal) Description() string {
	resultString := ""

	for index, value := range i.days {
		resultString += strconv.Itoa(value)
		if index < len(i.days)-1 {
			resultString += ","
		}
	}

	return resultString
}
func (i proposal) FilterValue() string { return i.name }

type preferredDay struct {
	day      int
	excluded []proposal
}

func (i preferredDay) Title() string { return strconv.Itoa(i.day) }
func (i preferredDay) Description() string {
	resultString := ""

	for index, value := range i.excluded {
		resultString += value.name
		if index < len(i.excluded)-1 {
			resultString += ","
		}
	}

	return resultString
}
func (i preferredDay) FilterValue() string { return strconv.Itoa(i.day) }

type model struct {
	list list.Model
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.list.View())
}

var docStyle = lipgloss.NewStyle().Margin(1, 2)

func main() {
	proposals := parseInputFile("input.txt")
	m := model{list: list.New(proposals, list.NewDefaultDelegate(), 0, 0)}
	p := tea.NewProgram(m, tea.WithAltScreen())
	m.list.Title = "List of persons with day proposal"

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}

	preferredDays := buildPreferredDays(proposals)
	m = model{list: list.New(preferredDays, list.NewDefaultDelegate(), 0, 0)}
	p = tea.NewProgram(m, tea.WithAltScreen())
	m.list.Title = "List of preferred day with excluded persons"

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func buildPreferredDays(items []list.Item) []list.Item {
	var preferredDays []list.Item
	sortedByKey := sortKeyByValue(countDays(items))
	for _, actualDay := range sortedByKey {
		preferredDay := *new(preferredDay)
		preferredDay.day = actualDay
		for _, item := range items {
			isIt := isExcluded(actualDay, item.(proposal).days)
			if isIt {
				preferredDay.excluded = append(preferredDay.excluded, item.(proposal))
			}
		}
		preferredDays = append(preferredDays, preferredDay)
	}
	return preferredDays
}

func parseInputFile(input string) []list.Item {
	file, err := os.Open(input)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return nil
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	var persons []list.Item

	for scanner.Scan() {
		persons = extractPerson(scanner, persons)
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}

	return persons
}

func extractPerson(scanner *bufio.Scanner, persons []list.Item) []list.Item {
	splitted := strings.SplitN(scanner.Text(), ":", 2)
	person := *new(proposal)
	person.name = splitted[0]
	days := strings.ReplaceAll(splitted[1], " ", "")
	for _, rawDay := range strings.Split(days, ",") {
		day, _ := strconv.Atoi(rawDay)
		person.days = append(person.days, day)
	}
	persons = append(persons, person)
	return persons
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

func countDays(items []list.Item) map[int]int {
	dayCounter := make(map[int]int)
	for _, item := range items {
		for _, day := range item.(proposal).days {
			if dayCounter[day] == 0 {
				dayCounter[day] = 1
			} else {
				dayCounter[day] += 1
			}
		}
	}
	return dayCounter
}

func isExcluded(thisDay int, days []int) bool {
	for _, day := range days {
		if day == thisDay {
			return false
		}
	}
	return true
}
