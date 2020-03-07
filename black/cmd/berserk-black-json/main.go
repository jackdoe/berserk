package main

import (
	"bufio"
	"encoding/json"
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

	if f.Replace {
		_, err := db.Exec("DELETE FROM documents WHERE dataset_key = $1", f.Dataset.DatasetKey)
		if err != nil {
			log.Fatal(err)
		}
	}

	err := models.UpsertDataset(db, f.Dataset)
	if err != nil {
		log.Fatal(err)
	}

	defer models.MustRecount(db, f.Dataset)

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	r := bufio.NewReader(os.Stdin)
	rollback := false
	i := 0
	docs := []models.Document{}
	for {
		data, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("%v, err: %v", string(data), err.Error())
			rollback = true
			break
		}

		var dict map[string]interface{}
		err = json.Unmarshal(data, &dict)
		if err != nil {
			log.Printf("%v, err: %v", string(data), err.Error())
			rollback = true
			break
		}

		d, err := models.ToDocument(f.DocKeyField, f.DocTagsField, dict)
		if err != nil {
			log.Printf("%v, err: %v", dict, err.Error())
			rollback = true
			break
		}

		d.DatasetKey = f.Dataset.DatasetKey
		docs = append(docs, d)

		i++
		if i%f.BatchSize == 0 {
			err = models.UpsertDocument(tx, docs)
			if err != nil {
				log.Printf("%v, err: %v", dict, err.Error())
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

			docs = []models.Document{}
			log.Printf("%d ...", i)
		}
	}

	err = models.UpsertDocument(tx, docs)
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
		os.Exit(1)
	} else {
		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}
	}
}
