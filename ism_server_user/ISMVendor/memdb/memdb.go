/**
 * @ Author: ISM Web组态软件
 * @ Create Time: 2023-01-09 08:53:22
 * @ Modified by: ISM Web组态软件
 * @ Modified time: 2023-04-03 08:57:51
 * @ Description: 此源码版权归 www.ismctl.com 所有,个人私自不得二次销售。
 */

package memdb

import (
	"sync"

	"github.com/iddxc/memodb/storage"
)

var Memdb *storage.Store

func MemdbServer() {
	Memdb = storage.InitStore("data/db/Redis.db", 20)

	var wait sync.WaitGroup
	wait.Add(1)
	go Memdb.Run()
	defer wait.Wait()
}
