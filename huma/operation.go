package huma

import (
	"net/http"

	basehuma "github.com/danielgtaylor/huma/v2"
)

type AffectedRowsOutput struct {
	Body struct {
		AffectedRows int `json:"affectedRows" example:"1" doc:"Number of affected rows"`
	}
}

func WithTag(tag string) func(*basehuma.Operation) {
	return func(op *basehuma.Operation) {
		op.Tags = []string{tag}
	}
}

func StatusCreated(op *basehuma.Operation) {
	op.DefaultStatus = http.StatusCreated
}

func AffectedRows(rows int64) *AffectedRowsOutput {
	out := &AffectedRowsOutput{}
	out.Body.AffectedRows = int(rows)
	return out
}
