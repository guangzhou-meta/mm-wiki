package controllers

import (
	"encoding/json"
	"github.com/phachon/mm-wiki/app/models"
	"strings"
)

type OpenApiController struct {
	BaseController
}

func (this *OpenApiController) Index() {
	result := make(map[string]string)
	this.jsonSuccess("", result)
}

func (this *OpenApiController) Search() {
	keyword := strings.TrimSpace(this.GetString("keyword", ""))
	type Doc struct {
		Title string
		Url   string
	}
	result := make(map[string]*Doc)
	if keyword == "" {
		this.jsonSuccess("", result)
		return
	}
	add := func(doc map[string]string) {
		id := doc["document_id"]
		result[id] = &Doc{
			Title: doc["name"],
			Url:   "/document/index?document_id=" + id,
		}
	}
	documents, err := models.DocumentModel.GetDocumentsByContent(keyword)
	if err == nil {
		for _, v := range documents {
			add(v)
		}
	}
	documents, err = models.DocumentModel.GetDocumentsByTags(keyword)
	if err == nil {
		for _, v := range documents {
			add(v)
		}
	}
	documents, err = models.DocumentModel.GetDocumentsByLikeName(keyword)
	if err == nil {
		for _, v := range documents {
			add(v)
		}
	}
	resultList := make([]*Doc, 0, len(result))
	for _, v := range result {
		resultList = append(resultList, v)
	}
	writer := this.Ctx.ResponseWriter
	writer.WriteHeader(200)
	data, _ := json.Marshal(resultList)
	_, _ = writer.Write(data)
	writer.Flush()
	//this.jsonSuccess("", resultList)
}
