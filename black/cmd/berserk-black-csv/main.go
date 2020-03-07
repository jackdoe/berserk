package main

import (
	"encoding/csv"
	"io"
	"log"
	"os"

	"github.com/jackdoe/berserk/black/pkg/common"
	"github.com/jackdoe/berserk/black/pkg/models"
)

func main() {
	f := common.ParseFlags()

	db := common.MustOpenDB()
	defer db.Close()

	r := csv.NewReader(os.Stdin)
	r.Comma = rune(f.CSVDelim)
	r.LazyQuotes = true
	r.TrimLeadingSpace = true

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	err = models.UpsertDataset(tx, f.Dataset)
	if err != nil {
		tx.Rollback()
		log.Fatal(err)
	}

	if f.Replace {
		_, err = tx.Exec("DELETE FROM documents WHERE dataset_key = ?", f.Dataset.DatasetKey)
		if err != nil {
			tx.Rollback()
			log.Fatal(err)
		}
	}
	rollback := false
	var header []string

	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("%v, err: %v", record, err.Error())
			rollback = true
			break
		}

		if header == nil {
			header = record
			continue
		}

		dict := map[string]interface{}{}
		for i := range header {
			dict[header[i]] = record[i]
		}

		d, err := models.ToDocument(f.DocKeyField, dict)

		d.DatasetKey = f.Dataset.DatasetKey

		if err != nil {
			log.Printf("%v, err: %v", record, err.Error())
			rollback = true
			break
		}

		err = models.UpsertDocument(tx, d)
		if err != nil {
			log.Printf("%v, err: %v", record, err.Error())
			rollback = true
			break
		}
	}

	if rollback {
		err = tx.Rollback()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(1)
	} else {
		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}
}
