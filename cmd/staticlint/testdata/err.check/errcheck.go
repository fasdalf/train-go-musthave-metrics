package errcheck

import (
	"fmt"
	"os"
)

func mulfunc(i int) (int, error) {
	return i * 2, nil
}

func errCheckFunc() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	mulfunc(5)           // want "expression returns unchecked error"
	res, _ := mulfunc(5) // want "assignment with unchecked error"
	fmt.Println(res)     // want "expression returns unchecked error"
	os.Exit(0)
}

func TestFunc() {
	var i int
	myfunc := func() error {
		return nil
	}
	myfunc() // want "expression returns unchecked error"
	if true {
		i := 7
		i, _ = mulfunc(i) // want "assignment with unchecked error"
	}

	(func() {})()
	i = *(func() *int { j := 5; return &j })()
	i, _ = i+1, myfunc() // want "assignment with unchecked error"
}
