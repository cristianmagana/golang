package main

import (
	"fmt"
	"sync"
	"time"
)

type BankAccount struct {
	balance int
	mu      sync.RWMutex
}

func (bankAccount *BankAccount) Deposit(amount int) {
	bankAccount.mu.Lock()
	defer bankAccount.mu.Unlock()
	bankAccount.balance += amount
}

func (bankAccount *BankAccount) Withdraw(amount int) bool {
	bankAccount.mu.Lock()
	defer bankAccount.mu.Unlock()
	if bankAccount.balance >= amount {
		bankAccount.balance -= amount
		return true
	} else {
		fmt.Printf("Withdraw amount: %d is greater than current balance: %d", amount, bankAccount.balance)
		return false
	}

}

func (bankAccount *BankAccount) Balance() int {
	bankAccount.mu.RLock()
	defer bankAccount.mu.RUnlock()
	return bankAccount.balance
}

func MutexSharedBankService() {
	fmt.Println("Starting mutex bank service")

	bankAccount := &BankAccount{balance: 0}

	var wg sync.WaitGroup

	// Start balance checker goroutine
	stopBalanceChecker := make(chan bool)
	go func() {
		ticker := time.NewTicker(100 * time.Millisecond)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				fmt.Printf("Current balance: $%d\n", bankAccount.Balance())
			case <-stopBalanceChecker:
				return
			}
		}
	}()

	// Start 5 depositor goroutines (each deposits $100, 10 times)
	for i := 0; i < 5; i++ {
		wg.Add(1)
		go func(depositorID int) {
			defer wg.Done()
			for range 10 {
				bankAccount.Deposit(100)
				time.Sleep(50 * time.Millisecond) // Small delay to spread out operations
			}
		}(i)
	}

	// Start 3 withdrawer goroutines (each withdraws $50, 20 times)
	for i := 0; i < 3; i++ {
		wg.Add(1)
		go func(withdrawerID int) {
			defer wg.Done()
			for j := 0; j < 20; j++ {
				bankAccount.Withdraw(50)
				time.Sleep(30 * time.Millisecond) // Small delay to spread out operations
			}
		}(i)
	}

	// Create a timeout channel for 2 seconds
	timeout := time.After(2 * time.Second)
	done := make(chan bool)

	// Wait for all goroutines to complete or timeout
	go func() {
		wg.Wait()
		done <- true
	}()

	select {
	case <-done:
		fmt.Println("All operations completed before timeout")
	case <-timeout:
		fmt.Println("2 seconds elapsed, stopping simulation")
	}

	// Stop the balance checker
	close(stopBalanceChecker)

	// Show final balance
	fmt.Printf("Final balance: $%d\n", bankAccount.Balance())
}
