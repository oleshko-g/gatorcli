package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/oleshko-g/gatorcli/internal/config"
	"github.com/oleshko-g/gatorcli/internal/database"
)

type state struct {
	cfg *config.Config
	db  *database.Queries
}

func setState() state {
	var s state
	cfg, err := config.Read()
	if err != nil {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("error reading config"))
		os.Exit(1)
	}
	s.cfg = &cfg
	return s
}

func parseArgs() command {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, fmt.Sprintf("error not enough arguments were provided."))
		os.Exit(1)
	}
	return command{
		name: os.Args[1],
		args: os.Args[2:],
	}
}

func openPostgresDB(dbURL string) (*sql.DB, error) {
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

func main() {
	state := setState()
	db, errDB := openPostgresDB(state.cfg.DataBaseURL)
	if errDB != nil {
		fmt.Fprintln(os.Stderr, errDB)
		os.Exit(1)
	}

	state.db = database.New(db)
	cmds := NewCommands()
	cmds.register("login", loginHandler)
	cmd := parseArgs()
	err := cmds.run(&state, cmd)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

}
