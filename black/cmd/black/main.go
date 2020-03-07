package main

import (
	"log"
	"net/http"

	"github.com/jackdoe/berserk/black/pkg/common"
	"github.com/jackdoe/berserk/black/pkg/models"

	jsoniter "github.com/json-iterator/go"

	"github.com/lib/pq"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ScanRequest struct {
	DatasetKey   string   `json:"dataset_key" binding:"required" uri:"dataset_key"`
	Tags         []string `json:"tags"`
	DocumentKeys []string `json:"document_keys"`
	Offset       uint64   `json:"offset" uri:"offset"`
	Limit        uint64   `json:"limit" uri:"limit"`
}

type LookupRequest struct {
	DatasetKey  string `json:"dataset_key" binding:"required"`
	DocumentKey string `json:"document_key" binding:"required"`
}

func main() {
	r := gin.Default()
	r.Use(cors.Default())

	db := common.MustOpenDB()
	defer db.Close()

	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"alive": true,
		})
	})

	r.GET("/", func(c *gin.Context) {
		ds, err := models.FindAllDS(db)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.IndentedJSON(200, gin.H{
			"datasets": ds,
		})
	})

	stream := func(c *gin.Context, j ScanRequest) {
		ds, err := models.FindDS(db, j.DatasetKey)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		if j.Limit == 0 {
			j.Limit = 100
		}

		query := "SELECT dataset_key, document_key, data, tags, created_at, updated_at FROM documents WHERE dataset_key = ? "
		args := []interface{}{ds.DatasetKey}

		if len(j.Tags) != 0 {
			query += " AND tags @> ? "
			args = append(args, pq.StringArray(j.Tags))
		}

		if len(j.DocumentKeys) != 0 {
			query += " AND document_key = ANY(?) "
			args = append(args, pq.Array(j.DocumentKeys))
		}

		query += " LIMIT ? OFFSET ?"
		args = append(args, j.Limit, j.Offset)

		query = models.ReplaceSQL(query, "?")
		rows, err := db.Query(query, args...)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return

		}

		log.Printf("%v %v", query, j)
		defer rows.Close()

		s := jsoniter.NewStream(jsoniter.ConfigFastest, c.Writer, 102400)
		s.WriteObjectStart()
		s.WriteObjectField("dataset")
		s.WriteVal(ds)
		s.WriteMore()

		s.WriteObjectField("documents")
		s.WriteArrayStart()

		i := 0

		for rows.Next() {
			var d models.Document

			if i > 0 {
				s.WriteMore() // fuck
			}

			err := rows.Scan(&d.DatasetKey, &d.DocumentKey, &d.Data, &d.Tags, &d.CreatedAt, &d.UpdatedAt)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			s.WriteVal(d)

			i++
		}

		err = rows.Err()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		s.WriteArrayEnd()
		s.WriteObjectEnd()
		s.Flush()
	}

	r.GET("/s/:dataset_key/:limit/:offset", func(c *gin.Context) {
		var j ScanRequest
		if err := c.ShouldBindUri(&j); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		stream(c, j)
	})

	r.POST("/s", func(c *gin.Context) {
		var j ScanRequest
		if err := c.ShouldBindJSON(&j); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		stream(c, j)
	})

	log.Fatal(r.Run())
}
