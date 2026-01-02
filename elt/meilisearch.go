package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	"github.com/meilisearch/meilisearch-go"
)

type meiliUpload struct {
	table string
	index string
}

func uploadToMeili(src *sqlx.DB, meili meilisearch.ServiceManager, operations []meiliUpload) error {
	errs := make([]error, 0)
	for _, op := range operations {
		data, err := selectSearchable(src, op.table)
		if err != nil {
			log.Printf("❌ meilisearch: select %s: %v", op.table, err)
			errs = append(errs, err)
			continue
		}
		if err = uploadToMeilisearch(meili, op.index, data); err != nil {
			log.Printf("❌ meilisearch: loading %s: %v", op.index, err)
			errs = append(errs, err)
			continue
		}
		log.Printf("✅ meilisearch: loading %s successfull", op.index)
	}
	if len(errs) > 0 {
		return listOfErrors(errs)
	}
	return nil
}

type record map[string]any

func (r *record) Scan(val any) error {
	switch v := val.(type) {
	case []byte:
		json.Unmarshal(v, &r)
		return nil
	case string:
		json.Unmarshal([]byte(v), &r)
		return nil
	default:
		return fmt.Errorf("Unsupported type: %T", v)
	}
}

func selectSearchable(db *sqlx.DB, table string) ([]map[string]any, error) {
	var rows []record
	query := `
		SELECT TO_JSONB(%s)
		FROM %s
	`
	query = fmt.Sprintf(query, table, table)
	if err := db.Select(&rows, query); err != nil {
		return nil, fmt.Errorf("failed to select searchable courses: %w", err)
	}
	result := make([]map[string]any, len(rows))
	for i, row := range rows {
		result[i] = map[string]any(row)
	}
	return result, nil
}

func uploadToMeilisearch(client meilisearch.ServiceManager, indexName string, data []map[string]any) error {
	index := client.Index(indexName)
	index.DeleteAllDocuments()
	task, err := index.AddDocuments(data)
	if err != nil {
		return err
	}
	_ = task.TaskUID
	return nil
}
