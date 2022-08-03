package db

import "reflect"

/*
表的构建采用建造者模式
*/

// Table 定义数据表的结构,table可指代特别复杂的数据结构,可暴露一个迭代器api,供调用方便捷获取内容
// 字段全部私有,只暴露指定的方法;
type Table struct {
	name       string
	recordType reflect.Type   //这张表里存放的类型
	records    map[any]record //一张表中的所有数据
	//todo 迭代器

}

func NewTable(name string) *Table {
	return &Table{
		name:       name,
		recordType: nil,
		records:    make(map[any]record),
	}
}

func (t *Table) WithType(recordType reflect.Type) *Table {
	t.recordType = recordType
	return t
}

//todo WithTableInteratorFactory

func (t *Table) Name() string {
	return t.name
}

func (t *Table) QueryByPrimaryKey(key any, value any) error {
	record, ok := t.records[key] //根据主键查询对应的记录
	if !ok {
		return ErrRecordNotFound
	}
	return record.convertByValue(value) //然后根据具体的实例类型查询值,并写入
}

func (t *Table) Insert(key any, value any) error {
	//检查是否已经存在
	if _, ok := t.records[key]; ok {
		return ErrPrimaryKeyConflict
	}
	//不存在 根据k-v创建记录
	record, err := recordFrom(key, value)
	if err != nil {
		return err
	}
	t.records[key] = record
	return nil
}

func (t *Table) Update(key any, value any) error {
	//检查是否已经存在
	if _, ok := t.records[key]; !ok {
		return ErrRecordNotFound
	}

	//存在 则生成新纪录覆盖原记录
	record, err := recordFrom(key, value)
	if err != nil {
		return err
	}
	t.records[key] = record
	return nil
}

func (t *Table) Delete(key any) error {
	if _, ok := t.records[key]; !ok {
		return ErrRecordNotFound
	}
	delete(t.records, key) //从表中删除record
	return nil

}
