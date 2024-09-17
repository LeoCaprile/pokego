package main

import (
	"fmt"
	"log"
	"pokego/client"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	table           table.Model
	pokemons        map[string]client.Pokemon
	selectedPokemon client.Pokemon
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
			pokemonId := m.table.SelectedRow()[0]

			if pokemon, ok := m.pokemons[pokemonId]; !ok {
				pokemon := client.GetPokemon(pokemonId)
				m.pokemons[pokemonId] = pokemon
				m.selectedPokemon = pokemon
			} else {
				m.selectedPokemon = pokemon
			}

		case "down", "j":
			m.table.MoveDown(1)
			pokemonId := m.table.SelectedRow()[0]

			if pokemon, ok := m.pokemons[pokemonId]; !ok {
				pokemon := client.GetPokemon(pokemonId)
				m.pokemons[pokemonId] = pokemon
				m.selectedPokemon = pokemon
			} else {
				m.selectedPokemon = pokemon
			}
		}
	}

	return m, nil
}

var containerStyle = lipgloss.NewStyle().Border(lipgloss.NormalBorder())
var containerBorder = containerStyle.Height(13).Width(30)
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
	column := lipgloss.JoinVertical(lipgloss.Top, containerBorder.Render(m.table.View()), containerBorder.Render(m.selectedPokemon.GetPokemonDescriptionView()))
	return lipgloss.JoinHorizontal(lipgloss.Top, m.selectedPokemon.GetImageView(80), column)
}

func getInitialModel() model {

	results := client.GetPokemonList()

	tableRow := []table.Row{}

	for i, pok := range results.Results {
		tableRow = append(tableRow, []string{fmt.Sprint(i + 1), strings.Title(pok.Name)})
	}

	initialPokemon := client.GetPokemon("1")

	return model{table: createPokedexTable(tableRow), selectedPokemon: initialPokemon, pokemons: map[string]client.Pokemon{
		"1": initialPokemon,
	}}
}

func main() {

	if _, err := tea.NewProgram(getInitialModel(), tea.WithAltScreen()).Run(); err != nil {
		log.Fatal(err)
	}

}
