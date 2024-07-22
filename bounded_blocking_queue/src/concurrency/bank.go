package concurrency

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Account struct {
	balance int64
	blcMut  *sync.Mutex
	ID      int64
}

type Bank struct {
	users    map[int64]*Account
	usersMut *sync.Mutex
}

func InitBank() *Bank {
	return &Bank{
		users:    make(map[int64]*Account, 10),
		usersMut: &sync.Mutex{},
	}
}

func (bk *Bank) AddAccount() (id int64) {
	bk.usersMut.Lock()
	defer bk.usersMut.Unlock()
	id = GetIDAtomic()
	bk.users[id] = &Account{
		balance: 0,
		blcMut:  &sync.Mutex{},
		ID:      id,
	}
	fmt.Printf("Add account, id: %v\n", id)
	return id
}

func (bk *Bank) StoreMoney(id, money int64) {
	// bk.usersMut.Lock()
	// defer bk.usersMut.Unlock()
	v, ok := bk.users[id]
	if !ok {
		fmt.Printf("the id %v is not exist\n", id)
		return
	}
	v.blcMut.Lock()
	defer v.blcMut.Unlock()
	v.balance += money
	fmt.Printf("id: %v stored $%v\n", id, money)
}

func (bk *Bank) WithdrawMoney(id, money int64) {
	// bk.usersMut.Lock()
	// defer bk.usersMut.Unlock()
	v, ok := bk.users[id]
	if !ok {
		fmt.Printf("the id %v is not exist\n", id)
		return
	}
	v.blcMut.Lock()
	defer v.blcMut.Unlock()
	if v.balance < money {
		fmt.Printf("the id %v 's money is not enough\n", id)
		return
	}
	v.balance -= money
	fmt.Printf("id: %v withdraw $%v\n", id, money)
}

func (bk *Bank) TransferMoney(idFrom, idDst, money int64) {
	// bk.usersMut.Lock()
	// defer bk.usersMut.Unlock()
	from, ok := bk.users[idFrom]
	if !ok {
		fmt.Printf("the id %v is not exist\n", idFrom)
		return
	}
	dst, ok := bk.users[idDst]
	if !ok {
		fmt.Printf("the id %v is not exist\n", idDst)
		return
	}
	from.blcMut.Lock()
	defer from.blcMut.Unlock()
	// dst.blcMut.Lock()
	// 防止双锁加锁顺序死锁
	retry, retryTimes := 0, 100
	for retry < retryTimes {
		if dst.blcMut.TryLock() {
			break
		} else {
			retry++
			runtime.Gosched()
		}
	}
	if retry == retryTimes {
		fmt.Printf("retry too many times!\n")
		return
	}
	defer dst.blcMut.Unlock()
	if from.balance < money {
		fmt.Printf("the id %v 's money is not enough\n", idFrom)
		return
	}
	from.balance -= money
	dst.balance += money
	fmt.Printf("id: %v transfer $%v to id: %v, idFrom: %v has $%v left\n", idFrom, money, idDst, idFrom, from.balance)
}

func BankSubTest(bk *Bank, id int64, id_ int64) {
	bk.StoreMoney(id, 200)
	bk.StoreMoney(id, 200)
	bk.WithdrawMoney(id, 200)
	bk.WithdrawMoney(id, 200)
	bk.TransferMoney(id, id_, 10)
	// bk.TransferMoney(id, id_, 10)
	// bk.TransferMoney(id, id_, 10)
	// bk.TransferMoney(id, id_, 10)
	// bk.TransferMoney(id, id_, 10)
	// bk.TransferMoney(id, id_, 10)
	// bk.TransferMoney(id, id_, 10)
}

func BankTest() {
	bk := InitBank()
	id1 := bk.AddAccount()
	id2 := bk.AddAccount()
	bk.StoreMoney(id1, 1000)
	bk.StoreMoney(id2, 1000)
	bk.WithdrawMoney(id1, 200)
	bk.WithdrawMoney(id2, 200)
	bk.TransferMoney(id1, id2, 300)
	bk.TransferMoney(id2, id1, 300)
	for i := 0; i < 10; i++ {
		go BankSubTest(bk, id1, id2)
		go BankSubTest(bk, id2, id1)
		time.Sleep(500 * time.Millisecond)
		fmt.Printf("-------------------------------------\n")
	}
	time.Sleep(1 * time.Second)
}
