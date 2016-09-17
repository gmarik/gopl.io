package main

func main() {
	println("Return without `return`:", fn(123))
}

func fn(i int) (rv int) {

	defer func() {
		if v := recover(); v != nil {
			rv = v.(int)
		}
	}()

	panic(i)
}
