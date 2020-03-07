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

	r := bufio.NewReader(os.Stdin)
	for {
		data, err := r.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Printf("%v, err: %v", string(data), err.Error())
			break
		}

		var dict map[string]interface{}
		err = json.Unmarshal(data, &dict)
		if err != nil {
			log.Printf("%v, err: %v", string(data), err.Error())
			continue
		}

		d, err := models.ToDocument(f.DocKeyField, f.DocTagsField, dict)
		if err != nil {
			log.Printf("%v, err: %v", dict, err.Error())
			continue
		}

		d.DatasetKey = f.Dataset.DatasetKey
		docs <- d
	}
	close(docs)
	<-done
}
