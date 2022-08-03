package db

import "sync"

type memoryDb struct {
	//一个数据库中有多个表
	tables sync.Map
}

var memoryDbInstance = &memoryDb{tables: sync.Map{}}

func MemoryDbInstance() *memoryDb { return memoryDbInstance }

func (m *memoryDb) CreateTable(t *Table) error {
	//先检查是否有同名表
	if _, ok := m.tables.Load(t.Name()); ok {
		return ErrTableAlreadyExist
	}
	m.tables.Store(t.Name(), t) //不存在则存入
	return nil
}

func (m *memoryDb) CreateTableIfNotExist(t *Table) error {
	if _, ok := m.tables.Load(t.Name()); ok {
		return nil
	}
	m.tables.Store(t.Name(), t) //不存在则存入
	return nil
}

func (m *memoryDb) DeleteTable(tableName string) error {
	//检查是否有这个表
	if _, ok := m.tables.Load(tableName); !ok {
		return ErrTableNotExist
	}
	m.tables.Delete(tableName)
	return nil
}

func (m *memoryDb) Query(tableName string, primaryKey any, result any) error {
	//首先要找到对应的表
	table, ok := m.tables.Load(tableName)
	if !ok {
		return ErrTableNotExist
	}
	//用找到的table进行查询
	return table.(*Table).QueryByPrimaryKey(primaryKey, result)
}

func (m *memoryDb) Insert(tableName string, primaryKey any, record any) error {
	table, ok := m.tables.Load(tableName)
	if !ok {
		return ErrTableNotExist
	}
	return table.(*Table).Insert(primaryKey, record)
}

func (m *memoryDb) Update(tableName string, primaryKey any, record any) error {
	table, ok := m.tables.Load(tableName)
	if !ok {
		return ErrTableNotExist
	}
	return table.(*Table).Update(primaryKey, record)
}

func (m *memoryDb) Delete(tableName string, primaryKey any) error {
	table, ok := m.tables.Load(tableName)
	if !ok {
		return ErrTableNotExist
	}
	return table.(*Table).Delete(primaryKey)
}

func (m *memoryDb) CreateTransaction(name string) *Transaction {
	return NewTransaction(name, m)
}

func (m *memoryDb) Clear() {
	m.tables = sync.Map{}
}
