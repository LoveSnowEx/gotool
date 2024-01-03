package table

import (
	"cmp"
	"encoding/csv"
	"os"
	"slices"

	"github.com/LoveSnowEx/gotool/errors"
)

type Row = []string
type Rows = []Row

type Table interface {
	// Save saves a table to a csv file
	Save(filename string) error
	// Fields returns the fields of the table
	Fields() Row
	// FieldIndex returns the index of a field, or -1 if the field does not exist
	FieldIndex(field string) int
	// Rows returns the rows of the table
	Rows() Rows
	// Append appends rows to the table
	Append(rows ...Row)
	// Clear clears the rows of the table
	Clear()
	// Sort sorts the table by a field
	Sort(field string) error
	// SortFunc sorts the table by a function
	SortFunc(func(row1, row2 Row) int)
}

type table struct {
	fields        Row
	fieldIndicies map[string]int
	rows          Rows
}

// New creates a new table
func New(fields Row, rows ...Row) Table {
	t := &table{
		fields: fields,
		rows:   rows,
	}
	t.makeFieldIndicies()
	return t
}

// NewFromFile creates a new table from a csv file
func NewFromFile(filename string) (t Table, err error) {
	f, err := os.Open(filename)
	if err != nil {
		return
	}
	defer f.Close()
	records, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return
	}
	return New(records[0], records[1:]...), nil
}

func (t *table) makeFieldIndicies() {
	t.fieldIndicies = make(map[string]int)
	for i, field := range t.fields {
		t.fieldIndicies[field] = i
	}
}

func (t *table) Save(filename string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()
	if err = w.Write(t.fields); err != nil {
		return
	}
	if err = w.WriteAll(t.rows); err != nil {
		return
	}
	return
}

func (t *table) Fields() Row {
	return t.fields
}

func (t *table) FieldIndex(field string) int {
	if index, ok := t.fieldIndicies[field]; ok {
		return index
	}
	return -1
}

func (t *table) Rows() Rows {
	return t.rows
}

func (t *table) Append(rows ...Row) {
	t.rows = append(t.rows, rows...)
}

func (t *table) Clear() {
	t.rows = make(Rows, 0)
}

func (t *table) Sort(field string) (err error) {
	idx := t.FieldIndex(field)
	if idx == -1 {
		return errors.ErrInvalidField
	}
	t.SortFunc(func(row1, row2 Row) int {
		return cmp.Compare[string](row1[idx], row2[idx])
	})
	return
}

func (t *table) SortFunc(f func(row1, row2 Row) int) {
	slices.SortFunc[Rows](t.rows, f)
}
