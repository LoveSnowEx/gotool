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
	// Load loads a table from a csv file
	Load(filename string) error
	// Save saves a table to a csv file
	Save(filename string) error
	// Fields returns the fields of the table
	Fields() Row
	// Rows returns the rows of the table
	Rows() Rows
	// SetFields sets the fields of the table
	SetFields(fields Row)
	// SetRows sets the rows of the table
	SetRows(rows Rows)
	// FieldIndex returns the index of a field, or -1 if it doesn't exist
	FieldIndex(field string) int
	// AppendRows appends rows to the table
	AppendRows(rows ...Row)
	// Sort sorts the table by a field
	Sort(field string) error
	// SortFunc sorts the table by a function
	SortFunc(func(row1, row2 Row) int)
}

type table struct {
	fields        []string
	fieldIndicies map[string]int
	rows          [][]string
}

func New() Table {
	return &table{
		fields:        make(Row, 0),
		rows:          make(Rows, 0),
		fieldIndicies: make(map[string]int),
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

func (t *table) Rows() Rows {
	return t.rows
}

func (t *table) SetFields(fields Row) {
	t.fields = fields
}

func (t *table) SetRows(rows Rows) {
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

func (t *table) AppendRows(rows ...Row) {
	t.rows = append(t.rows, rows...)
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
