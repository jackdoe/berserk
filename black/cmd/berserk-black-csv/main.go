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

	if f.Replace {
		models.MustDeleteAll(db, f.Dataset.DatasetKey)
	}

	err := models.UpsertDataset(db, f.Dataset)
	if err != nil {
		log.Fatal(err)
	}

	defer models.MustRecount(db, f.Dataset)

	docs := make(chan models.Document, 0)
	done := make(chan bool)
	go func() {
		models.InsertMany(db, f.BatchSize, docs)
		done <- true
	}()

	var header []string
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("%v, err: %v", record, err.Error())
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

		d, err := models.ToDocument(f.DocKeyField, f.DocTagsField, dict)
		d.DatasetKey = f.Dataset.DatasetKey
		docs <- d
	}
	close(docs)
	<-done
}
