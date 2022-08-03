package db

import (
	"reflect"
	"sync"
	"testing"
)

func TestTransaction(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		db := &memoryDb{tables: sync.Map{}}
		err := db.CreateTable(NewTable("region1").
			WithType(reflect.TypeOf(new(testRegion))))
		err = db.CreateTable(NewTable("region2").
			WithType(reflect.TypeOf(new(testRegion))))
		assertNoErr(t, err)

		tx := db.CreateTransaction("region_tx")
		tx.Begin()
		err = tx.Exec(NewInsertCmd("region1").WithPrimaryKey(1).
			WithRecord(&testRegion{Id: 1, Name: "beijing"}))
		err = tx.Exec(NewInsertCmd("region2").WithPrimaryKey(2).
			WithRecord(&testRegion{Id: 2, Name: "shanghai"}))
		assertNoErr(t, err)

		err = tx.Commit()
		assertNoErr(t, err)

		result := new(testRegion)
		db.Query("region1", 1, result)
		if result.Name != "beijing" {
			t.Error(result.Name)
		}

	})
	t.Run("fail and rollback", func(t *testing.T) {
		db := &memoryDb{tables: sync.Map{}}
		err := db.CreateTable(NewTable("region1").
			WithType(reflect.TypeOf(new(testRegion))))
		err = db.CreateTable(NewTable("region2").
			WithType(reflect.TypeOf(new(testRegion))))
		assertNoErr(t, err)

		tx := db.CreateTransaction("region_tx")
		tx.Begin()
		err = tx.Exec(NewInsertCmd("region1").WithPrimaryKey(1).
			WithRecord(&testRegion{Id: 1, Name: "beijing"}))
		err = tx.Exec(NewInsertCmd("region2").WithPrimaryKey(2).
			WithRecord(&testRegion{Id: 2, Name: "shanghai"}))
		err = tx.Exec(NewInsertCmd("region1").WithPrimaryKey(3).
			WithRecord(&testRegion{Id: 3, Name: "chengdu"}))
		err = tx.Exec(NewInsertCmd("region2").WithPrimaryKey(4).
			WithRecord(&testRegion{Id: 4, Name: "chongqing"}))
		err = tx.Exec(NewDeleteCmd("region2").WithPrimaryKey(5)) //fail here

		err = tx.Commit()
		assertErr(t, err, ErrRecordNotFound)

		result := new(testRegion)
		err = db.Query("region1", 1, result)
		assertErr(t, err, ErrRecordNotFound)
	})
}
