package main

import "time"

func main() {
	for true {
		A()
		time.Sleep(time.Second * 3)
		B()
		time.Sleep(time.Second * 3)
		C()
		time.Sleep(time.Second * 3)
		D()
		time.Sleep(time.Second * 3)
	}
}

//go:noinline
func A() {

}

//go:noinline
func B() {

}

//go:noinline
func C() {

}

//go:noinline
func D() {

}
