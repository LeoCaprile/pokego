package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	table     table.Model
	pokemonId string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		// Key mapping
		switch msg.String() {
		// Exit
		case "ctrl+c", "q":
			return m, tea.Quit
		// Up - Down table
		case "up", "k":
			m.table.MoveUp(1)
			m.pokemonId = m.table.SelectedRow()[0]
		case "down", "j":
			m.table.MoveDown(1)
			m.pokemonId = m.table.SelectedRow()[0]
		}

	}

	return m, nil
}

var containerStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
var containerImage = containerStyle.Height(15).Width(30)
var headerStyle = lipgloss.NewStyle().Padding(0, 1).Foreground(lipgloss.Color("252")).Bold(true)

func createPokedexTable(rows []table.Row) table.Model {

	columns := []table.Column{
		{Title: "ID", Width: 4},
		{Title: "Name", Width: 22},
	}

	t := table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
		table.WithHeight(14),
	)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("1")).
		Bold(false)
	t.SetStyles(s)

	return t
}

func (m model) View() string {
	column := lipgloss.JoinVertical(lipgloss.Top, containerImage.Render("POKE IMAGE"), containerImage.Render(m.table.View()))
	return lipgloss.JoinHorizontal(lipgloss.Top, column, containerStyle.Width(30).Height(32).Render("POKEMON INFO"+m.pokemonId))
}

type Pokemons struct {
	Count    int         `json:"count"`
	Next     string      `json:"next"`
	Previous interface{} `json:"previous"`
	Results  []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"results"`
}

func getInitialModel() model {

	res, err := http.Get("https://pokeapi.co/api/v2/pokemon?limit=150")

	if err != nil {
		log.Fatal(err)
	}

	defer res.Body.Close()

	var results Pokemons

	errCode := json.NewDecoder(res.Body).Decode(&results)

	if errCode != nil {
		log.Fatal(errCode)
	}

	tableRow := []table.Row{}

	for i, pok := range results.Results {
		tableRow = append(tableRow, []string{fmt.Sprint(i + 1), pok.Name})
	}

	return model{table: createPokedexTable(tableRow)}
}

func main() {

	if _, err := tea.NewProgram(getInitialModel()).Run(); err != nil {
		log.Fatal(err)
	}

}
