db_file := "db/database.sqlite"

# Initialize the SQLite database with schema
init-sql:
    @mkdir -p db
    @rm -f {{db_file}}
    @sqlite3 {{db_file}} < db/init.sql
    @echo "Database initialized at {{db_file}}"
