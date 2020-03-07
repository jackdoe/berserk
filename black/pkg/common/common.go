package common

import (
	"database/sql"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/jackdoe/berserk/black/pkg/models"
)

type FlagArgs struct {
	Dataset      models.Dataset
	DocKeyField  string
	DocTagsField string
	Replace      bool
	BatchSize    int
	CSVDelim     rune
}

func split(s string) []string {
	out := []string{}
	for _, item := range strings.Split(s, ",") {
		if len(item) > 0 {
			out = append(out, item)
		}
	}
	return out

}

func MustOpenDB() *sql.DB {
	uri := os.Getenv("DATABASE_URL")
	if uri == "" {
		uri = "postgres://postgres:postgres@localhost/black"
	}

	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	return db
}
func ParseFlags() FlagArgs {
	replace := flag.Bool("ds-replace", false, "delete and insert vs upsert")
	dsLicense := flag.String("ds-license", "", "path to the license file")
	dsName := flag.String("ds-name", "", "friendly dataset name")
	delim := flag.String("csv-delim", ",", "csv delimiter")
	batchSize := flag.Int("batch-size", 1000, "insert batch size")
	dsKey := flag.String("ds-key", "", "dataset key")
	dsTags := flag.String("ds-tags", "", "dataset tags")

	docKeyField := flag.String("doc-key-field", "", "document key field")
	docTagsField := flag.String("doc-tags-field", "", "document tags field")

	flag.Parse()

	license, err := ioutil.ReadFile(*dsLicense)
	if err != nil {
		log.Fatal(err)
	}

	if *dsName == "" {
		log.Fatal("-ds-name is required")
	}

	if *dsKey == "" {
		log.Fatal("-ds-key is required")
	}

	if *docKeyField == "" {
		log.Fatal("-doc-key-field is required")
	}

	if *docTagsField == "" {
		log.Fatal("-doc-tags-field is required")
	}

	sdelim := *delim

	return FlagArgs{
		Dataset: models.Dataset{
			Name:       *dsName,
			License:    string(license),
			DatasetKey: *dsKey,
			Tags:       split(*dsTags),
		},
		BatchSize:    *batchSize,
		Replace:      *replace,
		DocKeyField:  *docKeyField,
		DocTagsField: *docTagsField,
		CSVDelim:     rune(sdelim[0]),
	}
}
