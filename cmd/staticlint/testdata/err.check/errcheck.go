package errcheck

import (
	"fmt"
	"os"
)

func mulfunc(i int) (int, error) {
	return i * 2, nil
}
func mulfuncP(i int) (*int, error) {
	ii := i * 2
	return &ii, nil
}
func funcP(i int) *int {
	ii := i * 2
	return &ii
}

func errCheckFunc() {
	// формулируем ожидания: анализатор должен находить ошибку,
	// описанную в комментарии want
	mulfunc(5)           // want "expression returns unchecked error"
	res, _ := mulfunc(5) // want "assignment with unchecked error"
	go mulfunc(6)        // want "go statement with unchecked error"
	defer mulfunc(6)     // want "defer with unchecked error"
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
	ii, _ := mulfuncP(i) // want "assignment with unchecked error"
	ii = funcP(*ii)
	i = *ii
	i, _ = i+1, myfunc() // want "assignment with unchecked error"
}
