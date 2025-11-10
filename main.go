package main

import (
	"fmt"
	"log"
	"os"
	"pokego/client"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type PokemonList = map[string]client.Pokemon

type model struct {
	table           table.Model
	pokemons        PokemonList
	selectedPokemon client.Pokemon
}

func (m model) Init() tea.Cmd {
	return nil
}

func GetNextPokemon(currentPokemonId int, pokemonList *PokemonList) {

	log.Print("GETTING NEXT POKEMON")
	log.Print("Current ID: ", currentPokemonId)
	for i := currentPokemonId; i < currentPokemonId+10; i++ {
		pokemonId := fmt.Sprint(i + 1)

		log.Print("Current foorloop ID: ", pokemonId)

		if _, ok := (*pokemonList)[pokemonId]; !ok {

			pokemon, err := client.GetPokemon(pokemonId)
			if err != nil {
				log.Print(err)
			}
			(*pokemonList)[pokemonId] = *pokemon
		}
	}

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
			m.selectedPokemon = m.pokemons[pokemonId]

		case "down", "j":
			m.table.MoveDown(1)
			pokemonId := m.table.SelectedRow()[0]
			log.Print(pokemonId, len(m.pokemons))
			log.Print("cursor", m.table.Cursor())

			if len(m.pokemons)-1 == m.table.Cursor() {
				go GetNextPokemon(m.table.Cursor(), &m.pokemons)
			}

			m.selectedPokemon = m.pokemons[pokemonId]
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

	initialPokemon, err := client.GetPokemon("1")
	if err != nil {
		log.Print(err)
	}

	model := model{table: createPokedexTable(tableRow), selectedPokemon: *initialPokemon, pokemons: PokemonList{
		"1": *initialPokemon,
	}}

	go GetNextPokemon(1, &model.pokemons)

	return model
}

func main() {
	f, err := os.OpenFile("logs", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalf("error opening file: %v", err)
	}
	defer f.Close()

	log.SetOutput(f)

	if _, err := tea.NewProgram(getInitialModel(), tea.WithAltScreen()).Run(); err != nil {
		log.Fatal(err)
	}

}
