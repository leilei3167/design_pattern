package db

import (
	"reflect"
	"testing"
)

type testRegion struct {
	Id   int
	Name string
}

func TestTable(t *testing.T) {
	//创建一个测试用的 testRegin存储实例
	t.Run("success", func(t *testing.T) {
		table := NewTable("testRegion").WithType(reflect.TypeOf(new(testRegion)))
		err := table.Insert(2, &testRegion{Id: 2, Name: "beijing"})
		assertNoErr(t, err)
		record := new(testRegion)

		err = table.QueryByPrimaryKey(2, record)
		assertNoErr(t, err)

		if record.Name != "beijing" {
			t.Error("query failed, want beijing, got " + record.Name)
		}

		err = table.Update(2, &testRegion{Id: 2, Name: "shanghai"})
		assertNoErr(t, err)

		err = table.QueryByPrimaryKey(2, record)
		assertNoErr(t, err)

		if record.Name != "shanghai" {
			t.Error("query failed, want shanghai, got " + record.Name)
		}

		table.Delete(2)
		err = table.QueryByPrimaryKey(2, record)
		assertErr(t, err, ErrRecordNotFound)
	})

	t.Run("fail", func(t *testing.T) {
		table := NewTable("testRegion").WithType(reflect.TypeOf(new(testRegion)))
		//插入一些数据
		err := table.Insert(1, &testRegion{Id: 2, Name: "beijing"})
		err = table.Insert(2, &testRegion{Id: 2, Name: "chengdu"})
		err = table.Insert(3, &testRegion{Id: 2, Name: "sichuan"})
		err = table.Insert(4, &testRegion{Id: 2, Name: "henan"})
		assertNoErr(t, err)

		err = table.Insert(4, &testRegion{Id: 3, Name: "henan"})
		assertErr(t, err, ErrPrimaryKeyConflict)

		err = table.Delete(1)
		assertNoErr(t, err)
		err = table.Delete(1)
		assertErr(t, err, ErrRecordNotFound)

	})

}

func assertErr(t testing.TB, got, want error) {
	t.Helper()
	if got != want {
		t.Errorf("unexpected err,got:'%v',want '%v'", got, want)
	}
}

func assertNoErr(t testing.TB, err error) {
	t.Helper()
	if err != nil {
		t.Fatal("got err: ", err)
	}
}
