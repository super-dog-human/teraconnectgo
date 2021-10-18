package domain

import (
	"time"

	"cloud.google.com/go/datastore"
)

// DeleteOrderはエンティティの削除予約を作成します。LessonまたはUserの関連エンティティの定時削除に使用されることを想定しています。
type DeleteOrder struct {
	ID         int64 `datastore:"-"`
	EntityName string
	TargetID   int64 `datastore:",noindex"`
	Created    time.Time
}

func CreateDeleteLessonOrderInTransaction(tx *datastore.Transaction, lessonID int64) error {
	order := new(DeleteOrder)
	order.EntityName = "Lesson"
	order.TargetID = lessonID
	order.Created = time.Now()

	key := datastore.IncompleteKey("DeleteOrder", nil)
	if _, err := tx.Put(key, order); err != nil {
		return err
	}

	return nil
}
