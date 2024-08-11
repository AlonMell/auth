package migrator

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type Config struct {
	db              *sql.DB
	Target, Current *Migration
	Path            string
	Table           string
}

type Migration struct {
	Major, Minor int
	FileNumber   int
}

func New(db *sql.DB, path, table string, major, minor int) *Config {
	return &Config{
		db:      db,
		Table:   table,
		Path:    path,
		Target:  &Migration{major, minor, 0},
		Current: &Migration{},
	}
}

// TODO: Сделать автодобавление новой записи миграций
// TODO: Сделать парсинг файлов в горутинах (И последовательное выполнение скриптов)
// TODO: Сделать возможность отката к предыдущим версиям (добавить up/down)
// TODO: Возможно, ускорить поиск подходящего номера файла

func Migrate(cfg *Config) error {
	const op = "migrator.Migrate"

	if err := cfg.db.Ping(); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	if exists, err := cfg.tableExists(); err != nil {
		return err
	} else if !exists {
		cfg.Target.findFileNumber(cfg.Path)
		err = cfg.execBoundsFiles()
		return err
	}

	if err := cfg.setCurrentVersion(); err != nil {
		return err
	}

	if cfg.Target.Major < cfg.Current.Major ||
		cfg.Target.Major == cfg.Current.Major &&
			cfg.Target.Minor <= cfg.Current.Minor {
		return nil
	}

	cfg.Target.findFileNumber(cfg.Path)
	cfg.Current.FileNumber++

	err := cfg.execBoundsFiles()

	return err
}

func (c *Config) setCurrentVersion() error {
	const op = "migrator.setCurrentVersion"

	query := fmt.Sprintf(`
    SELECT major_version, minor_version, file_number
    FROM %s 
    ORDER BY date_applied DESC 
    LIMIT 1;`, c.Table)

	stmt, err := c.db.Prepare(query)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	err = stmt.QueryRow().Scan(&c.Current.Major, &c.Current.Minor, &c.Current.FileNumber)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

// execBoundsFiles execute all sql scripts in directory where first, second - it's bounds
func (c *Config) execBoundsFiles() error {
	const op = "migrator.exec"

	first, second := c.Current.FileNumber, c.Target.FileNumber

	files, err := os.ReadDir(c.Path)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	for first <= second {
		err = c.exec(files[first].Name())
		if err != nil {
			return fmt.Errorf("%s: %w", op, err)
		}
		first++
	}

	return nil
}

// exec execute sql script
func (c *Config) exec(name string) error {
	const op = "migrator.exec"

	path := filepath.Join(c.Path, name)
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	defer file.Close()

	// TODO: сделать чтение с горутинами
	script, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	_, err = c.db.Exec(string(script))
	return err
}

func (c *Config) tableExists() (bool, error) {
	const op = "migrator.tableExists"

	query := `
    SELECT EXISTS (
        SELECT FROM information_schema.tables 
        WHERE table_schema = 'public' 
        AND table_name = $1
    );`

	stmt, err := c.db.Prepare(query)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}
	defer stmt.Close()

	var exists bool
	err = stmt.QueryRow(c.Table).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	return exists, nil
}

// findFileNumber find last version file number
// in: 00 00
// files: 0001.00.00 0002.00.00
// out: 2
func (m *Migration) findFileNumber(path string) {
	var first int

	for i := 1; i < 10000; i++ {
		fileName := fmt.Sprintf("%04d.%02d.%02d", i, m.Major, m.Minor)
		path = filepath.Join(path, fileName)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			break
		}

		first = i
	}

	m.FileNumber = first
}
