package main

import (
	"database/sql"
	"fmt"
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

func getTableInfo(db *sql.DB, tableName string) ([]string, error) {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var columns []string
	for rows.Next() {
		var cid int
		var name, typ string
		var notnull, pk int
		var dflt_value interface{}
		if err := rows.Scan(&cid, &name, &typ, &notnull, &dflt_value, &pk); err != nil {
			return nil, err
		}
		columns = append(columns, name)
	}
	return columns, nil
}

func getTableData(db *sql.DB, tableName string) ([][]string, error) {
	rows, err := db.Query(fmt.Sprintf("SELECT * FROM %s", tableName))
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	var result [][]string
	values := make([]interface{}, len(columns))
	valuePtrs := make([]interface{}, len(columns))
	for i := range columns {
		valuePtrs[i] = &values[i]
	}

	for rows.Next() {
		err := rows.Scan(valuePtrs...)
		if err != nil {
			return nil, err
		}

		row := make([]string, len(columns))
		for i, val := range values {
			if val == nil {
				row[i] = "NULL"
			} else {
				row[i] = fmt.Sprintf("%v", val)
			}
		}
		result = append(result, row)
	}
	return result, nil
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

	// Create tables view
	tables := tview.NewTable()
	tables.
		SetBorders(true).
		SetTitle(" Tables (Enter to view, Esc to go back) ").
		SetTitleAlign(tview.AlignLeft)
	tables.SetSelectable(true, false)

	// Create data view
	data := tview.NewTable()
	data.
		SetBorders(true).
		SetTitleAlign(tview.AlignLeft)
	data.SetSelectable(true, true)

	// Set headers for tables list
	tables.SetCell(0, 0, tview.NewTableCell("Table Name").
		SetTextColor(tcell.ColorYellow).
		SetSelectable(false))

	// Get and display tables
	tableList, err := getTables(db)
	if err != nil {
		log.Fatal(err)
	}

	for i, tableName := range tableList {
		tables.SetCell(i+1, 0, tview.NewTableCell(tableName))
	}

	// Current view state
	currentView := tables
	//var selectedTable string

	// Handle table selection
	updateDataView := func(tableName string) {
		//selectedTable = tableName
		data.Clear()
		data.SetTitle(fmt.Sprintf(" Table: %s (Esc to go back) ", tableName))

		// Get columns
		columns, err := getTableInfo(db, tableName)
		if err != nil {
			log.Printf("Error getting columns: %v", err)
			return
		}

		// Set headers
		for i, col := range columns {
			data.SetCell(0, i, tview.NewTableCell(col).
				SetTextColor(tcell.ColorYellow).
				SetSelectable(false))
		}

		// Get and display data
		rows, err := getTableData(db, tableName)
		if err != nil {
			log.Printf("Error getting data: %v", err)
			return
		}

		for i, row := range rows {
			for j, cell := range row {
				data.SetCell(i+1, j, tview.NewTableCell(cell))
			}
		}
	}

	// Handle keyboard input
	app.SetInputCapture(func(event *tcell.EventKey) *tcell.EventKey {
		switch event.Key() {
		case tcell.KeyEnter:
			if currentView == tables {
				row, _ := tables.GetSelection()
				if row > 0 {
					tableName := tableList[row-1]
					updateDataView(tableName)
					currentView = data
					app.SetRoot(data, true)
				}
			}
		case tcell.KeyEscape:
			if currentView == data {
				currentView = tables
				app.SetRoot(tables, true)
			}
		}
		return event
	})

	// Select first row by default (after header)
	tables.Select(1, 0)

	if err := app.
		SetRoot(tables, true).
		EnableMouse(true).
		Run(); err != nil {
		log.Fatal(err)
	}
}
