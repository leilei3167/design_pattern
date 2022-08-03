package db

import (
	"math/rand"
	"time"
)

/*迭代器模式*/

type TableIterator interface {
	HasNext() bool
	Next(next any) error
}

// TableIteratorFactory 迭代器工厂,具有Create方法能返回迭代器接口实例的实例就能作为工厂
type TableIteratorFactory interface {
	Create(table *Table) TableIterator
}

//随机迭代器
type tableIteratorImpl struct {
	records []record
	cursor  int
}

func (t *tableIteratorImpl) HasNext() bool {
	return t.cursor < len(t.records)
}

func (t *tableIteratorImpl) Next(next any) error {
	r := t.records[t.cursor]
	t.cursor++
	if err := r.convertByValue(next); err != nil {
		return err
	}
	return nil
}

//创建随机迭代器的工厂
type randomTableIteratorFactory struct{}

func (r *randomTableIteratorFactory) Create(table *Table) TableIterator {
	var records []record
	for _, r := range table.records {
		records = append(records, r)
	}
	rrand := rand.New(rand.NewSource(time.Now().UnixNano()))
	rrand.Shuffle(len(records), func(i, j int) {
		records[i], records[j] = records[j], records[i]
	})
	return &tableIteratorImpl{
		records: records,
		cursor:  0,
	}

}
