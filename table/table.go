package table

import (
	"bufio"
	"encoding/csv"
	"os"
	"slices"
)

type Table interface {
	// Load loads a table from a csv file
	Load(filename string) error
	// Save saves a table to a csv file
	Save(filename string) error
	// Fields returns the fields of the table
	Fields() []string
	// Rows returns the rows of the table
	Rows() [][]string
	// SetFields sets the fields of the table
	SetFields([]string)
	// SetRows sets the rows of the table
	SetRows([][]string)
	// FieldIndex returns the index of a field, or -1 if it doesn't exist
	FieldIndex(string) int
	// AppendRows appends rows to the table
	AppendRows(...[]string)
	// Sort sorts the table by a field
	Sort(string)
	// SortFunc sorts the table by a function
	SortFunc(func(row1, row2 []string) int)
}

type table struct {
	fields        []string
	fieldIndicies map[string]int
	rows          [][]string
}

func New() Table {
	return &table{
		fields: make([]string, 0),
		rows:   make([][]string, 0),
	}
}

func (t *table) Load(filename string) (err error) {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return err
	}
	t.fields = records[0]
	t.rows = records[1:]
	return
}

func (t *table) Save(filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	defer buf.Flush()
	w := csv.NewWriter(buf)
	defer w.Flush()
	err = w.WriteAll(append(
		[][]string{t.fields},
		t.rows...,
	))
	return
}

func (t *table) Fields() []string {
	return t.fields
}

func (t *table) Rows() [][]string {
	return t.rows
}

func (t *table) SetFields(fields []string) {
	t.fields = fields
}

func (t *table) SetRows(rows [][]string) {
	t.rows = rows
}

func (t *table) makeFieldIndicies() {
	t.fieldIndicies = make(map[string]int)
	for i, field := range t.fields {
		t.fieldIndicies[field] = i
	}
}

func (t *table) FieldIndex(field string) int {
	if t.fieldIndicies == nil {
		t.makeFieldIndicies()
	}
	if index, ok := t.fieldIndicies[field]; ok {
		return index
	}
	return -1
}

func (t *table) AppendRows(rows ...[]string) {
	t.rows = append(t.rows, rows...)
}

func (t *table) Sort(field string) {
	t.SortFunc(func(row1, row2 []string) int {
		switch {
		case row1[t.FieldIndex(field)] < row2[t.FieldIndex(field)]:
			return -1
		case row1[t.FieldIndex(field)] > row2[t.FieldIndex(field)]:
			return 1
		default:
			return 0
		}
	})
}

func (t *table) SortFunc(f func(row1, row2 []string) int) {
	slices.SortFunc[[][]string](t.rows, f)
}
