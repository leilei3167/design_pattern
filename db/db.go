package db

/*
db包实现数据库的接口定义,并提供内存储存的实现

主要用到的设计模式
	1.简单工厂;在获取实例时,以及建造者的第一步都是简单工厂模式
	2.单例模式(懒汉);在内存存储实例中运用,初始化全局的内存存储实例
	3.建造者模式;在事务命令实例的建造中,以及数据表table的建造中使用到
	4.命令模式;在事务的实现中使用,将各种命令抽象为Command接口,并存在Transaction实例中
	5.备忘录模式;在事务的回滚功能中用到,抽象一个cmdHistory类管理执行过的Command
	5.迭代器模式

*/

/**
 * 依赖倒置原则（DIP）：
 * 1、高层模块不应该依赖低层模块，两者都应该依赖抽象
 * 2、抽象不应该依赖细节，细节应该依赖抽象
 * DIP并不是说高层模块是只能依赖抽象接口，它的本意应该是依赖稳定的接口/抽象类/具象类(数据库可能会随着开发变更,属于不稳定的)。
 * 如果一个具象类是稳定的，比如Java中的String，那么高层模块依赖它也没有问题；
 * 相反，如果一个抽象接口是不稳定的，经常变化，那么高层模块依赖该接口也是违反DIP的，这时候应该思考下接口是否抽象合理。
 * 例子：
 * 抽象出Db接口，避免高层应用依赖具体的Db实现（比如MemoryDb），符合DIP

我们可以忽略掉低层模块的细节，抽象出一个稳定的接口，然后让高层模块依赖该接口，
同时让低层模块实现该接口，从而实现了依赖关系的倒置
*/

//定义数据库接口,操作表,增删改查的基本行为,db的操作: 找到对应的数据表,执行操作
//一个数据库实例包含多个表,每个表包含多个记录

type Db interface {
	CreateTable(t *Table) error //对表的增删
	CreateTableIfNotExist(t *Table) error
	DeleteTable(tableName string) error

	Query(tableName string, primaryKey any, result any) error //key-value 存储,需创建result实例来接收结果
	Insert(tableName string, primaryKey any, record any) error
	Update(tableName string, primaryKey any, record any) error
	Delete(tableName string, primaryKey any) error

	//事务等

	CreateTransaction(name string) *Transaction
}
