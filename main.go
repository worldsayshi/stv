package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/gdamore/tcell/v2"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rivo/tview"
)

func getTables(db *sql.DB) ([]string, error) {
	rows, err := db.Query(`
		SELECT name FROM sqlite_master
		WHERE type='table'
		ORDER BY name;
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tables []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		tables = append(tables, name)
	}
	return tables, nil
}

func main() {
	if len(os.Args) != 2 {
		log.Fatal("Usage: go run main.go <path-to-sqlite-db>")
	}

	db, err := sql.Open("sqlite3", os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	app := tview.NewApplication()
	table := tview.NewTable()
	table.
		SetBorders(true).
		SetTitle(" Tables ").
		SetTitleAlign(tview.AlignLeft)
	table.
		SetSelectable(true, false) // Enable row selection, disable column selection

	// Set headers
	table.SetCell(0, 0, tview.NewTableCell("Table Name").
		SetTextColor(tcell.ColorYellow).
		SetSelectable(false))

	// Get and display tables
	tables, err := getTables(db)
	if err != nil {
		log.Fatal(err)
	}

	for i, tableName := range tables {
		table.SetCell(i+1, 0, tview.NewTableCell(tableName))
	}

	// Select first row by default (after header)
	table.Select(1, 0)

	if err := app.
		SetRoot(table, true).
		SetFocus(table). // Ensure table has focus at startup
		EnableMouse(true).
		Run(); err != nil {
		log.Fatal(err)
	}
}
