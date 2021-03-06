// Package zen -
package zen

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"strconv"
	"strings"
)

type brain struct {
	Data Data
	UI   UI
}

func isNilFixed(i interface{}) bool {
	if i == nil {
		return true

	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

// NewBrain -
// Ignore linter complaint returning a non-nexported type
// nolint:golint
func NewBrain(data Data, ui UI) (*brain, error) {
	if isNilFixed(data) {
		return nil, fmt.Errorf("nil data object supplied")
	}
	if isNilFixed(ui) {
		return nil, fmt.Errorf("nil UI upplied")
	}
	return &brain{Data: data, UI: ui}, nil
}

// Data -
type Data interface {
	Initialise()
	FindMatches(group, term, value string) ([]map[string][]string, error)
	FindRelated(string, map[string][]string) (map[string][]map[string][]string, error) // map[group][]map[fieldname][][values]
	GetGroups() ([]string, error)
	GetTerms(string) (map[string]struct{}, error)
}

// UI -
type UI interface {
	WelcomeMenu() (string, error)
	DataMenu([]string)
	GetCommand() (string, error)
	GroupQuestion() (string, error)
	TermQuestion() (string, error)
	ValueQuestion() (string, error)
	ShowResults([]map[string][]map[string][]string)
	ShowTerms(map[string][]string)
	ShowError(string)
}

func (b *brain) getGroup(maxTries int) (string, error) {
	groups, _ := b.Data.GetGroups()
	groupNum := 0
	i := 0
	for {
		if i >= maxTries {
			b.UI.ShowError("Too many attempts, sorry.")
			b.Shutdown()
		}
		group, err := b.UI.GroupQuestion()
		if err != nil {
			return "", fmt.Errorf("search group question error %w", err)
		}
		if group == "q" || group == "quit" {
			b.Shutdown()
		}
		groupNum, err = strconv.Atoi(group)
		if err != nil || groupNum < 1 || groupNum > len(groups) {
			b.UI.ShowError(fmt.Sprintf("Invalid selection, please choose between 1 and %d, or 'quit' at any time.", len(groups)))
			i++
			continue
		}
		break
	}
	return groups[groupNum-1], nil
}

func (b *brain) getTerm(maxTries int, group string) (string, error) {
	terms, err := b.Data.GetTerms(group)
	if err != nil {
		// Group doesn't exist - how did /that/ happen?
		b.Shutdown()
	}
	i := 0
	termName := ""
	for {
		if i >= maxTries {
			b.UI.ShowError("Too many attempts, sorry.")
			return "", fmt.Errorf("too many attempts")
		}
		termName, err = b.UI.TermQuestion()
		if err != nil {
			return "", fmt.Errorf("search term question error %w", err)
		}
		if termName == "q" || termName == "quit" {
			return "", fmt.Errorf("quit signal")
		}
		if _, ok := terms[termName]; !ok {
			b.UI.ShowError("Invalid selection, please choose a valid search term, or 'quit' at any time.")
			i++
			continue
		}
		break
	}
	return termName, nil
}

func (b *brain) getValue(maxTries int, group, term string) (string, error) {
	termValue, err := b.UI.ValueQuestion()
	if err != nil {
		return "", fmt.Errorf("search value question error %w", err)
	}
	return termValue, nil
}

// MaxTries - Maximum times a user can try to enter data for a single question
const MaxTries = 3

func (b *brain) Forever() {
	b.Data.Initialise()
	for {
		option, err := b.UI.WelcomeMenu()
		if err != nil {
			b.UI.ShowError(fmt.Sprintf("Data Get Groups error %v\n", err))
		}
		if option == "1" {
			// Search for tickets
			b.searchTickets()

		} else if option == "2" {
			// Show all of the terms
			b.showTerms()
		} else if strings.ToLower(option) == "q" || strings.ToLower(option) == "quit" {
			b.Shutdown()
		} else {
			b.UI.ShowError("Valid options are '1' or '2', please try again.")
		}

	}
}

// Shudown -
func (b *brain) Shutdown() {
	os.Exit(0)
}

func (b *brain) showTerms() {
	groups, err := b.Data.GetGroups()
	if err != nil {
		b.UI.ShowError(fmt.Sprintf("Data Get Groups error %v\n", err))
	}
	found := map[string][]string{}
	for _, group := range groups {
		g, err := b.Data.GetTerms(group) // map[string]struct{}
		if err != nil {
			b.UI.ShowError(fmt.Sprintf("Data Get Terms error %v\n", err))
		}
		found[group] = []string{}
		for k := range g {
			found[group] = append(found[group], k)
		}
		sort.Strings(found[group])
	}
	b.UI.ShowTerms(found)
}

func (b *brain) searchTickets() {
	groups, err := b.Data.GetGroups()
	if err != nil {
		b.UI.ShowError(fmt.Sprintf("Data Get Groups error %v\n", err))
	}
	b.UI.DataMenu(groups)
	group, err := b.getGroup(MaxTries)
	if err != nil {
		// log.Print(err)
		b.Shutdown()
	}
	term, err := b.getTerm(MaxTries, group)
	if err != nil {
		// log.Print(err)
		b.Shutdown()
	}
	value, err := b.getValue(MaxTries, group, term)
	if err != nil {
		// log.Fatal(err)
		b.Shutdown()
	}

	matches, _ := b.Data.FindMatches(group, term, value)
	found := make([]map[string][]map[string][]string, len(matches)+1)
	found[0] = map[string][]map[string][]string{}
	found[0]["matches"] = []map[string][]string{}
	found[0]["matches"] = append(found[0]["matches"], matches...)
	for idx := range matches {
		d, _ := b.Data.FindRelated(group, matches[idx])
		for k := range d {
			found[idx][k] = d[k]
		}
	}
	b.UI.ShowResults(found)
}
