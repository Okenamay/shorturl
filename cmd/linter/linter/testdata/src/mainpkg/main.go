package main

import (
	"log"
	"os"
)

// otherFunc - функция не в main, здесь вызовы запрещены
//
//lint:ignore U1000 "This function is used for analysis testing only"
func otherFunc() {
	panic("I am a panic")     // want "do not use panic"
	log.Fatal("I am fatal")   // want "os.Exit/log.Fatal call outside main.main"
	log.Fatalf("I am fatalf") // want "os.Exit/log.Fatal call outside main.main"
	os.Exit(1)                // want "os.Exit/log.Fatal call outside main.main"
}

// main - пакет main, функция main.
// panic все равно запрещен, а log.Fatal/os.Exit - разрешены.
func main() {
	panic("I am a panic in main") // want "do not use panic"

	log.Fatal("I am fatal in main")   // OK
	log.Fatalf("I am fatalf in main") // OK
	os.Exit(1)                        // OK

	// Проверка на кастомные функции
	var myLog logger
	myLog.Fatal("это не log.Fatal") // OK
}

// Кастомный тип, чтобы проверить, что анализатор не срабатывает
// на методы с таким же именем
type logger struct{}

func (l *logger) Fatal(s string) {}
