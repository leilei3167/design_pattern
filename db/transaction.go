package db

/*
db的事务实现用到了命令模式 和 备忘录模式

命令模式:
	- 命令模式（Command Pattern）是一种数据驱动的设计模式，它属于行为型模式。
请求以命令的形式包裹在对象中，并传给调用对象。调用对象寻找可以处理该命令的合适的对象，
并把该命令传给相应的对象，该对象执行命令。
	关键点就是将命令抽象为接口,某个对象持有这些命令接口
	- 主要解决：在软件系统中，行为请求者与行为实现者通常是一种紧耦合的关系，但某些场合，
比如需要对行为进行记录、撤销或重做、事务等处理时，这种无法抵御变化的紧耦合的设计就不太合适。

注意:在go中,Command的形式不局限于是一个接口,是自定义的函数类型也是完全可以的,只要符合某个函数签名的 就可以
作为命令来发送执行


*/

/*
备忘录模式:
	备忘录模式（Memento Pattern）在不破坏封装性的前提下保存一个对象的某个状态，以便在适当的时候恢复对象。
备忘录模式属于行为型模式。


*/

/*
数据库的事务是命令模式和备忘录模式的经典使用场景

*/

// Command 命令接口
type Command interface {
	// Exec 执行增改删命令
	Exec() error
	// Undo 回滚
	Undo()
	//设置关联的数据库
	setDb(db Db)
}

// Transaction 事务的实现
type Transaction struct {
	name string
	db   Db
	cmds []Command //把具体的命令,封装到事务中,Transaction就包装了命令可已被传递,在某处可以执行
}

func NewTransaction(name string, db Db) *Transaction {
	return &Transaction{
		name: name,
		db:   db,
		cmds: nil,
	}
}

//调用顺序为 begin->exec->...->commit

func (t *Transaction) Begin() {
	t.cmds = make([]Command, 0)
}

func (t *Transaction) Exec(cmd Command) error {
	if t.cmds == nil {
		return ErrTransactionNotBegin //确保提前执行了Begin
	}

	cmd.setDb(t.db)              //每条命令的关联数据库和此次事务的一致
	t.cmds = append(t.cmds, cmd) //将命令添加到执行队列中
	return nil
}

func (t *Transaction) Commit() error {
	//创建cmdHistory对象来存储历史记录
	history := &cmdHistory{history: make([]Command, 0, len(t.cmds))}
	//一次执行任务队列中的任务
	for _, cmd := range t.cmds {
		history.add(cmd)
		if err := cmd.Exec(); err != nil { //执行第一条就出现错误? 将无法实现第一条命令的回滚
			history.rollback()
			return err
		}
	}
	return nil
}

/*备忘录模式 核心在于记录执行的历史记录
此处就是将执行的命令存在切片中,要回滚时倒叙遍历
*/

type cmdHistory struct { //cmdHistory就是管理快照的类,具有添加和回滚的功能
	history []Command
}

func (c *cmdHistory) add(cmd Command) {
	c.history = append(c.history, cmd)
}

func (c *cmdHistory) rollback() {
	for i := len(c.history) - 1; i >= 0; i-- {
		c.history[i].Undo() //从最末尾向前执行每条命令的回滚操作,实现整个事务的回滚
	}
}

//一些命令的实例,可以很方便的添加各种命令

// InsertCmd 插入命令
type InsertCmd struct {
	db         Db
	tableName  string
	primaryKey any
	newRecord  any
}

func (i *InsertCmd) Exec() error { //复用插入的代码
	return i.db.Insert(i.tableName, i.primaryKey, i.newRecord)
}

func (i *InsertCmd) Undo() {
	//回滚就是将插入的给删掉
	i.db.Delete(i.tableName, i.primaryKey)
}

func (i *InsertCmd) setDb(db Db) {
	i.db = db
}

/*命令建造*/

func NewInsertCmd(tableName string) *InsertCmd {
	return &InsertCmd{tableName: tableName}
}

func (i *InsertCmd) WithPrimaryKey(primaryKey any) *InsertCmd {
	i.primaryKey = primaryKey
	return i
}

func (i *InsertCmd) WithRecord(record interface{}) *InsertCmd {
	i.newRecord = record
	return i
}

// UpdateCmd 更新命令
type UpdateCmd struct {
	db         Db
	tableName  string
	primaryKey interface{}
	newRecord  interface{}
	oldRecord  interface{} //记录更新前的数据
}

func NewUpdateCmd(tableName string) *UpdateCmd {
	return &UpdateCmd{tableName: tableName}
}

func (u *UpdateCmd) WithPrimaryKey(primaryKey interface{}) *UpdateCmd {
	u.primaryKey = primaryKey
	return u
}

func (u *UpdateCmd) WithRecord(record interface{}) *UpdateCmd {
	u.newRecord = record
	return u
}

func (u *UpdateCmd) Exec() error {
	//更新前应该先查询(Update的实现已验证了是否存在),此处目的是将修改前的数据拿出来
	if err := u.db.Query(u.tableName, u.primaryKey, u.oldRecord); err != nil {
		return err
	}
	return u.db.Update(u.tableName, u.primaryKey, u.newRecord)
}

func (u *UpdateCmd) Undo() { //因为保存了更新前的数据,恢复即可
	u.db.Update(u.tableName, u.primaryKey, u.oldRecord)
}

func (u *UpdateCmd) setDb(db Db) {
	u.db = db
}

// DeleteCmd 删除命令
type DeleteCmd struct {
	db         Db
	tableName  string
	primaryKey interface{}
	oldRecord  interface{} //保存删除前的数据
}

func NewDeleteCmd(tableName string) *DeleteCmd {
	return &DeleteCmd{tableName: tableName}
}

func (d *DeleteCmd) WithPrimaryKey(primaryKey interface{}) *DeleteCmd {
	d.primaryKey = primaryKey
	return d
}

func (d *DeleteCmd) Exec() error {
	if err := d.db.Query(d.tableName, d.primaryKey, d.oldRecord); err != nil {
		return err
	}
	return d.db.Delete(d.tableName, d.primaryKey)
}

func (d *DeleteCmd) Undo() { //撤销就是插入被删除的数据
	d.db.Insert(d.tableName, d.primaryKey, d.oldRecord)
}

func (d *DeleteCmd) setDb(db Db) {
	d.db = db
}
