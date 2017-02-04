// Package migration provides an interface for migrating a database backwards
// and forwards.
//
// You must store the path to the env.json file in the
// environment variable: JAYCONFIG
//
// Examples:
//	jay migrate:mysql make "test" # Create new migration
//	jay migrate:mysql all         # Advance all migrations
//	jay migrate:mysql reset       # Rollback all migrations
//	jay migrate:mysql refresh     # Rollback all migrations then advance all migrations
//	jay migrate:mysql status      # See last 'up' migration
//	jay migrate:mysql up          # Apply only the next 'up' migration
//	jay migrate:mysql down        # Apply only the current 'down' migration
//
//	jay migrate make "Create user table"
//	  Creates two new files in the database/migration folder using this format:
//	    * YYYYMMDD_HHMMSS.nnnnnn_create_user_table.up.sql
//	    * YYYYMMDD_HHMMSS.nnnnnn_create_user_table.down.sql
package migration

import (
	"errors"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"
)

var (
	// ErrNone is when there are no migrations in the database
	ErrNone = errors.New("No migrations yet.")
	// ErrCurrent is when the database is up-to-date
	ErrCurrent = errors.New("Database current. No changes made.")
	// ErrMissing is when the migration file cannot be found
	ErrMissing = errors.New("Migration not found.")
	// ErrTableNotCreated is when the migration cannot be created
	ErrTableNotCreated = errors.New("Could not create the migration table.")
)

// Info holds the information for the migration.
type Info struct {
	// Db is the database information
	Db Interface
	//DateFormat is the date and time format for the migration files
	DateFormat string
	// Folder is the migrations folder
	Folder string
	// List if the life of Up migrations
	List []string
	// Position is the index of the current migration in the List
	position int
	// positionName is the name of the current migration
	positionName string
	// Folder is the migrations table
	Table string
	// TemplateUp is the stub used for Up migration files when they are created
	TemplateUp string
	// TemplateDown is the stub used for Down migration files when they are created
	TemplateDown string
	// Output is the log information
	output string
}

// Position returns the index of the migration. 0 is no migration.
func (info *Info) Position() int {
	// Get the last migration from the database
	migrationName, err := info.Db.Status()

	// If migration is found
	if len(migrationName) > 0 && err == nil {
		info.position, _ = info.migrationPosition(migrationName)
		info.positionName = migrationName
	} else { // Else reset the position to 0
		info.position = 0
		info.positionName = ""
	}

	return info.position
}

// Interface defines all the functions required for a migration.
type Interface interface {
	// Extension should return an extension without a period or a blank string
	Extension() string
	// TableExist should return an error if the table does not exist
	TableExist() error
	// CreateTable should return an error if the table could not be created
	CreateTable() error
	// Status should return the name of the last migration or return an error
	// or a model.ErrNoResult error if there are no results
	Status() (string, error)
	// Migrate will run the migration and return an error if not successful
	Migrate(query string) error
	// RecordUp should record the name of the file in the database
	RecordUp(name string) error
	// RecordDown should record the name of the file in the database and make
	// any changes to the database like updating the AUTO_INCREMENT value
	RecordDown(name string) error
}

func (info *Info) log(text string) {
	info.output += text
}

// Output returns the text output of the performed operations.
func (info *Info) Output() string {
	return info.output
}

// New returns an instance of a migration after creating the migration
// table (if one doesn't exist), retrieving a list of the available
// migrations, and reading the last migration.
// You must connect to the database prior to calling this function.
func New(db Interface, table string, folder string) (*Info, error) {
	info := &Info{
		Db:         db,
		Folder:     folder,
		Table:      table,
		DateFormat: "20060102_150405.000000",
	}

	// Check for the migration table
	err := db.TableExist()
	if err != nil {
		err = db.CreateTable()
		if err != nil {
			return info, ErrTableNotCreated
		}
	}

	// Get Up migration list
	info.List, err = info.updateList()
	if err != nil {
		return info, err
	}

	return info, err
}

func (info *Info) verifyMigration(path string) bool {
	// Get the path name
	name := filepath.Base(path)
	name = strings.TrimSuffix(name, ".up"+info.Db.Extension())
	name = strings.TrimSuffix(name, ".down"+info.Db.Extension())

	// Determine if the migration matches
	return name == info.positionName
}

// Status returns the last applied migration name without the file extension.
func (info *Info) Status() string {
	// If migration is current
	if info.Position() == 0 {
		return ErrNone.Error()
	}

	// Get the name of the last applied migration (relative file path)
	path := info.List[info.Position()-1]

	// Determine if the migration is missing on disk
	if strings.HasPrefix(filepath.Base(path), info.positionName) {
		return fmt.Sprintf("Migration is missing on disk: %v", info.positionName)
	}

	// Get the name to store in the database record
	return strings.Replace(filepath.Base(path), ".up"+info.Db.Extension(), "", -1)
}

// updateList returns the list of Up migrations.
func (info *Info) updateList() ([]string, error) {
	return filepath.Glob(filepath.Join(info.Folder, "*.up"+info.Db.Extension()))
}

// migrtionPosition returns the position of the migration or an error.
func (info *Info) migrationPosition(current string) (int, error) {
	for i := 0; i < len(info.List); i++ {
		if strings.Contains(info.List[i], current) {
			return i + 1, nil
		}
	}
	return 0, ErrMissing
}

// Create writes two new migration files to the folder with timestamps and descriptions.
func (info *Info) Create(description string) error {
	// Remove spaces and convert to lowercase
	desc := strings.ToLower(strings.Replace(description, " ", "_", -1))

	// Set the timestamp
	now := time.Now().Format(info.DateFormat)
	prefix := fmt.Sprintf("%v_%v", now, desc)

	// Create full paths
	up := filepath.Join(info.Folder, prefix+".up"+info.Db.Extension())
	down := filepath.Join(info.Folder, prefix+".down"+info.Db.Extension())

	// Create up file
	err := ioutil.WriteFile(up, []byte(info.TemplateUp), 0644)
	if err != nil {
		return err
	}

	info.output += fmt.Sprintf("Migration created: %v\n", up)

	// Create down file
	err = ioutil.WriteFile(down, []byte(info.TemplateDown), 0644)
	if err != nil {
		return err
	}

	info.output += fmt.Sprintf("Migration created: %v\n", down)

	// Update migration list
	info.List, err = info.updateList()
	if err != nil {
		return err
	}

	return nil
}

// UpOne applies only the next migration.
func (info *Info) UpOne() error {
	// If migration is current
	if info.Position() >= len(info.List) {
		return ErrCurrent
	}

	// Start at next position
	err := info.up()
	if err != nil {
		return err
	}

	info.output += "  | Migration up complete\n"

	return nil
}

// UpAll applies all migrations that have not been applied.
func (info *Info) UpAll() error {
	// If migration is current
	if info.Position() >= len(info.List) {
		return ErrCurrent
	}

	// Start at next position
	for i := info.Position(); i < len(info.List); i++ {
		err := info.up()
		if err != nil {
			return err
		}
	}

	info.output += "  | Migration up complete\n"

	return nil
}

// up reads the query and passes it to the database.
func (info *Info) up() error {
	// Get the name of the next Up file to migrate
	file := info.List[info.Position()]

	// Read the file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// Get the name to store in the database record
	name := strings.Replace(filepath.Base(file), ".up"+info.Db.Extension(), "", -1)

	// Run the migration
	err = info.Db.Migrate(string(data))
	if err != nil {
		return err
	}

	// Record a successful result
	err = info.Db.RecordUp(name)
	if err != nil {
		return err
	}

	info.output += fmt.Sprintf("+ | Applied: %v\n", name)

	return nil
}

// DownOne removes only the last migration.
func (info *Info) DownOne() error {
	// If migration is current
	if info.Position() == 0 {
		return ErrNone
	}

	// Start at current position
	err := info.down()
	if err != nil {
		return err
	}

	info.output += "  | Migration down complete\n"
	return nil
}

// DownAll removes all migrations.
func (info *Info) DownAll() error {
	// If migration is current
	if info.Position() == 0 {
		return ErrNone
	}

	// Start at current position
	for i := info.Position() - 1; i >= 0; i-- {
		err := info.down()
		if err != nil {
			return err
		}
	}

	info.output += "  | Migration down complete\n"

	return nil
}

// down reads the query and passes it to the database.
func (info *Info) down() error {
	// Get the name of the current Down file to migrate
	fileUp := info.List[info.Position()-1]

	// Change the extension to reference the down file
	file := strings.Replace(fileUp, ".up"+info.Db.Extension(), ".down"+info.Db.Extension(), -1)

	// Verify the migration
	if !info.verifyMigration(file) {
		return errors.New("Migration is missing: " + file)
	}

	// Read the file
	data, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// Get the name to store in the database record
	name := strings.Replace(filepath.Base(file), ".down"+info.Db.Extension(), "", -1)

	// Run the migration
	err = info.Db.Migrate(string(data))
	if err != nil {
		return err
	}

	// Record a successful result
	err = info.Db.RecordDown(name)
	if err != nil {
		return err
	}

	info.output += fmt.Sprintf("- | Removed: %v\n", name)

	return nil
}
