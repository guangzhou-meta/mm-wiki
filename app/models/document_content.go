package models

import (
	"database/sql"
	"fmt"
	"github.com/phachon/mm-wiki/app/utils"
	"strings"
	"time"
)

const (
	Table_Document_Content_Name = "document_content"

	Document_Content_Slice_Size = 3072
)

type DocumentContent struct {
}

var DocumentContentModel = DocumentContent{}

func (d *DocumentContent) UpdateDocumentContent(documentId string, spaceId string, parentId string, documentContent string) (err error) {
	err = d.DeleteDB(documentId)
	if err != nil {
		return
	}
	db := G.DB()
	if len(documentContent) == 0 {
		return nil
	}
	contents := utils.MarkdownToPlainText([]byte(documentContent), Document_Content_Slice_Size)
	if len(contents) == 0 {
		return nil
	}
	var seg int
	var ctn string
	now := time.Now().Unix()
	for segment, content := range contents {
		seg = segment
		ctn = content
		_, _ = db.Exec(db.AR().Insert(Table_Document_Content_Name, map[string]interface{}{
			"document_id": documentId,
			"space_id":    spaceId,
			"parent_id":   parentId,
			"plain_text":  ctn,
			"segment":     seg,
			"is_delete":   0,
			"create_time": now,
			"update_time": now,
		}))
	}
	return nil
}

func (d *DocumentContent) getFullLikeSql(keyword string) string {
	conditions := []string{
		fmt.Sprintf("(b.plain_text like '%%%s%%' and a.plain_text like '%%%s%%')", keyword, keyword),
	}

	var i = 1
	var l = len(keyword)
	var sqlTemplate = "(b.plain_text like '%%%s' and a.segment = (b.segment + 1) and a.plain_text like '%s%%')"
	for i < l {
		conditions = append(conditions, fmt.Sprintf(sqlTemplate, keyword[0:i], keyword[i:]))
		i++
	}

	return fmt.Sprintf("select a.document_id document_id from mw_%s a left join mw_%s b on %s where b.reference_id is not null group by a.document_id", Table_Document_Content_Name, Table_Document_Content_Name, strings.Join(conditions, " or "))
}

func (d *DocumentContent) GetDocumentIdsByContent(keyword string) (list []string, err error) {
	db := G.DB()
	likeSql := d.getFullLikeSql(keyword)
	rs, err := db.Query(db.AR().Raw(likeSql))
	if err != nil {
		return
	}
	for _, r := range rs.Row() {
		list = append(list, r)
	}
	return
}

func (d *DocumentContent) DeleteDB(documentId string) (err error) {
	db := G.DB()
	var tx *sql.Tx
	tx, err = db.Begin(db.Config)
	if err != nil {
		return
	}
	_, err = db.ExecTx(db.AR().Delete(Table_Document_Content_Name, map[string]interface{}{
		"document_id": documentId,
	}), tx)
	if err != nil {
		return tx.Rollback()
	}
	err = tx.Commit()
	return
}
