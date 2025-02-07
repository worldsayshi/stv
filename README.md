Visualize an sqlite database using tview in go.

# Considerations

- This app will display a lot of tview tables with various content. It should work quite similarly to k9s, but aimed at displaying sqlite data instead of k8s data.

# TODO

- [x] init.sql script that creates a users table
- [x] Justfile script that has an init-sql recipe that creates an sqlite db and runs the init.sql script
- [x] main.go script that takes the path to an sqlite db file as argument and list all tables as a tview Table.
- [x] Allow moving highlighting of a row with arrow keys. The table should be focused at startup.
- [x] Pressing enter on a row in the tables table should bring you to a tview table displaying that table
- [ ] Only show either the table listing all the sqlite tables or the table showing the data. Navigate to the data table by selecting with enter, also make sure to move the focus. Navigate back by pressing esc, which closes the data table and moves focus back to the table of tables.
- [ ] Pressing ':' at any time should make a popup appear. Let's call this pop-up 'EntitySelector'. The behavior will be similar to the ':' selector in k9s.
- [ ] The EntitySelector the user can type the name of an entity. Initially entities will be '.tables' for selecting the view of all tables or the name of a table, to display that table.