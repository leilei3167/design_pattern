package db

import (
	"math/rand"
	"sort"
	"time"
)

/*迭代器模式
有时会遇到这样的需求，开发一个模块，用于保存对象；不能用简单的数组、列表，得是红黑树、跳表等较为复杂的数据结构；
有时为了提升存储效率或持久化，还得将对象序列化；
但必须给客户端提供一个易用的 API，允许方便地、多种方式地遍历对象，丝毫不察觉背后的数据结构有多复杂。

Db接口已经实现了查询接口,但我们还想提供全表查询接口,有随机和有序两种方式,并且支持客户端拓展遍历方式!

*/

// TableIterator 关键点 定义迭代器抽象接口(产品接口),允许后续客户端拓展遍历方式
type TableIterator interface {
	HasNext() bool
	Next(next any) error
}

//随机迭代器 关键点2:定义迭代器接口的实现(类似于java中的基类)
type tableIteratorImpl struct {
	//关键点3: 定义一个集合存储待遍历的记录，这里的记录已经排序好或者随机打散
	records []record
	//关键点4: 定义一个游标记录当前遍历的位置
	cursor int
}

// HasNext 关键点5: 在HasNext函数中的判断是否已经遍历完所有记录
func (t *tableIteratorImpl) HasNext() bool {
	return t.cursor < len(t.records)
}

// Next 关键点6: 在Next函数中取出下一个记录，并转换成客户端期望的对象类型，记得增加cursor
func (t *tableIteratorImpl) Next(next any) error {
	r := t.records[t.cursor]
	t.cursor++
	if err := r.convertByValue(next); err != nil {
		return err
	}
	return nil
}

// TableIteratorFactory 抽象工厂,具有Create方法能返回迭代器接口实例的实例就能作为工厂
type TableIteratorFactory interface {
	Create(table *Table) TableIterator
}

//创建随机迭代器的工厂
type randomTableIteratorFactory struct{}

func (r *randomTableIteratorFactory) Create(table *Table) TableIterator {
	var records []record
	for _, r := range table.records { //为了避免并发问题 原始数据拷贝
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
func NewRandomTableIteratorFactory() *randomTableIteratorFactory {
	return &randomTableIteratorFactory{}
}

//---------------有序迭代器-------------

// Less 如果i<j返回ture，否则返回false
type Less func(i, j interface{}) bool

// records 辅助record记录根据主键排序
type records struct {
	less Less
	rs   []record
}

func newRecords(rs []record, less Less) *records {
	return &records{
		less: less,
		rs:   rs,
	}
}
func (r *records) Len() int {
	return len(r.rs)
}

func (r *records) Less(i, j int) bool {
	return r.less(r.rs[i].primaryKey, r.rs[j].primaryKey)
}

func (r *records) Swap(i, j int) {
	tmp := r.rs[i]
	r.rs[i] = r.rs[j]
	r.rs[j] = tmp
}

// sortedTableIteratorFactory 根据主键进行排序，排序逻辑由Comparable定义
type sortedTableIteratorFactory struct {
	less Less
}

func (s *sortedTableIteratorFactory) Create(table *Table) TableIterator {
	var records []record
	for _, r := range table.records {
		records = append(records, r)
	}
	sort.Sort(newRecords(records, s.less))
	return &tableIteratorImpl{
		records: records,
		cursor:  0,
	}
}

func NewSortedTableIteratorFactory(less Less) *sortedTableIteratorFactory {
	return &sortedTableIteratorFactory{less: less}
}
