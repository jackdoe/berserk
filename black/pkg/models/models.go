package models

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/lib/pq"
)

type Document struct {
	DocumentKey string          `json:"document_key"`
	DatasetKey  string          `json:"dataset_key"`
	Data        json.RawMessage `json:"data"`
	Tags        pq.StringArray  `json:"tags"`
	CreatedAt   time.Time       `json:"created_at"`
	UpdatedAt   time.Time       `json:"updated_at"`
}

type Dataset struct {
	DocumentCount int64          `json:"document_count"`
	DatasetKey    string         `json:"dataset_key"`
	Tags          pq.StringArray `json:"tags"`
	License       string         `json:"license"`
	Name          string         `json:"name"`
	CreatedAt     time.Time      `json:"created_at"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

var debug = os.Getenv("DEBUG") == "true"

func MustRecount(db *sql.DB, ds Dataset) {
	_, err := db.Exec(`update datasets set document_count = (select count(*) from documents where dataset_key=$1)`, ds.DatasetKey)
	if err != nil {
		log.Fatal(err)
	}
}

func UpsertDataset(tx *sql.DB, ds Dataset) error {
	q := `INSERT INTO datasets(name, license, key, tags, created_at, updated_at)
                     VALUES($1,$2,$3,$4,NOW(),NOW())
              ON CONFLICT (key) DO
                  UPDATE SET
                     tags = excluded.tags,
                     license = excluded.license,
                     created_at = excluded.created_at,
                     updated_at = excluded.updated_at,
                     name = excluded.name`
	r, err := tx.Exec(q, ds.Name, ds.License, ds.DatasetKey, pq.StringArray(ds.Tags))
	if err == nil && debug {
		ra, _ := r.RowsAffected()
		log.Printf("[DS] inserting: %v, rows affected: %v", ds, ra)
	}

	if err != nil {
		log.Printf("[DS] %v, %v", ds, err)
	}

	return err
}

func UpsertDocument(tx *sql.Tx, d []Document) error {
	// blah
	if len(d) == 0 {
		return nil
	}
	place := ""
	vals := []interface{}{}

	for i, doc := range d {
		x := i * 4
		place += fmt.Sprintf("($%d,$%d,$%d,$%d,NOW(),NOW()),", x+1, x+2, x+3, x+4)
		vals = append(vals, doc.DatasetKey, doc.DocumentKey, pq.StringArray(doc.Tags), doc.Data)
	}

	place = strings.TrimSuffix(place, ",")
	query := `
                   INSERT INTO documents(dataset_key, document_key, tags, data, created_at, updated_at)
                          VALUES ` + place + `
                   ON CONFLICT (dataset_key, document_key) DO
                       UPDATE SET
                          tags = excluded.tags,
                          created_at = excluded.created_at,
                          updated_at = excluded.updated_at,
                          data = excluded.data`

	r, err := tx.Exec(query, vals...)

	if err == nil && debug {
		ra, _ := r.RowsAffected()
		log.Printf("[DOC] inserting: %v, rows affected: %v", len(d), ra)
	}
	if err != nil {
		log.Printf("[DOC] %v %v, %v", query, d, err)
	}

	return err
}

func ToDocument(keyField string, tagsKeyField string, in map[string]interface{}) (Document, error) {
	_, ok := in[keyField]
	if !ok {
		return Document{}, fmt.Errorf("bad document, no key field(%s) in %v", keyField, in)
	}

	k, err := extractKey(in[keyField])
	if err != nil {
		return Document{}, err
	}

	if len(k) == 0 {
		return Document{}, fmt.Errorf("bad document, bad key value field(%s) %v", keyField, in)
	}

	tags := extractTags(in[tagsKeyField])

	b, err := json.Marshal(in)
	if err != nil {
		return Document{}, err
	}
	return Document{
		Tags:        pq.StringArray(tags),
		DocumentKey: k,
		Data:        b,
	}, nil
}

func MustDeleteAll(db *sql.DB, dsKey string) {
	_, err := db.Exec("DELETE FROM documents WHERE dataset_key = $1", dsKey)
	if err != nil {
		log.Fatal(err)
	}
}

func FindDS(db *sql.DB, dsKey string) (Dataset, error) {
	var ds Dataset
	row := db.QueryRow("SELECT document_count, license, key, name, tags, created_at, updated_at from datasets WHERE key = $1", dsKey)
	err := row.Scan(&ds.DocumentCount, &ds.License, &ds.DatasetKey, &ds.Name, &ds.Tags, &ds.CreatedAt, &ds.UpdatedAt)
	return ds, err
}

func FindAllDS(db *sql.DB) ([]Dataset, error) {
	rows, err := db.Query("SELECT document_count, license, key, name, tags, created_at, updated_at from datasets ORDER BY name")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := []Dataset{}
	for rows.Next() {
		var ds Dataset

		err := rows.Scan(&ds.DocumentCount, &ds.License, &ds.DatasetKey, &ds.Name, &ds.Tags, &ds.CreatedAt, &ds.UpdatedAt)
		if err != nil {
			return nil, err
		}
		out = append(out, ds)
	}
	return out, nil
}

func InsertMany(db *sql.DB, batchSize int, in chan Document) {
	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	rollback := false
	i := 0

	docs := []Document{}
	for d := range in {
		docs = append(docs, d)

		i++
		if i%batchSize == 0 {
			err = UpsertDocument(tx, docs)
			if err != nil {
				log.Printf("%v, err: %v", d, err.Error())
				rollback = true
				break
			}

			err = tx.Commit()
			if err != nil {
				log.Fatal(err)
			}

			tx, err = db.Begin()
			if err != nil {
				log.Fatal(err)
			}

			docs = []Document{}
			log.Printf("%d ...", i)
		}
	}

	err = UpsertDocument(tx, docs)
	if err != nil {
		log.Printf("%v, err: %v", docs, err.Error())
		rollback = true
	}

	if rollback {
		err = tx.Rollback()
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("rolled back")
	} else {
		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}
	}
}

func ReplaceSQL(old, searchPattern string) string {
	tmpCount := strings.Count(old, searchPattern)
	for m := 1; m <= tmpCount; m++ {
		old = strings.Replace(old, searchPattern, "$"+strconv.Itoa(m), 1)
	}
	return old
}

func extractKey(in interface{}) (string, error) {
	switch v := in.(type) {
	case string:
		return v, nil
	case int:
		return fmt.Sprintf("%d", v), nil
	case int32:
		return fmt.Sprintf("%d", v), nil
	case int64:
		return fmt.Sprintf("%d", v), nil
	case int16:
		return fmt.Sprintf("%d", v), nil
	case uint32:
		return fmt.Sprintf("%d", v), nil
	case uint64:
		return fmt.Sprintf("%d", v), nil
	case uint16:
		return fmt.Sprintf("%d", v), nil
	case float32:
		return strconv.FormatFloat(float64(v), 'f', 0, 64), nil
	case float64:
		return strconv.FormatFloat(float64(v), 'f', 0, 64), nil
	default:
		return "", fmt.Errorf("unknown type: %v", in)
	}
}

func extractTags(in interface{}) []string {
	switch v := in.(type) {
	case map[string]interface{}:
		out := make([]string, 0, len(v))
		for k, _ := range v {
			out = append(out, k)
		}
		return out
	case string:
		return []string{v}
	case []interface{}:
		out := make([]string, len(v))
		for i, val := range v {
			out[i] = fmt.Sprintf("%v", val)
		}
		return out
	case []string:
		return v
	case nil:
		return []string{"_"}
	default:
		return []string{fmt.Sprintf("%v", v)}
	}

}
