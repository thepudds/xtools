// package go121 contains tests that are disabled for go 1.20
package go121

import (
	"context"
	"sync"
	"time"
)

// Cases that exercise skipping uninteresting last statements.
func _() {

	// We are able to analyze a go or defer statement closure followed by
	// a simple last statement that we understand.
	var a int
	for i := range "loop" {
		go func() {
			print(i) // want "loop variable i captured by func literal"
		}()
		i = a
	}

	// We understand a series of simple last statements.
	for i := range "loop" {
		defer func() {
			print(i) // want "loop variable i captured by func literal"
		}()
		a = 0
		a = i + 1
	}

	// We allow a trailing 'continue' as an uninteresting last statement.
	for i := range "loop" {
		go func() {
			print(i) // want "loop variable i captured by func literal"
		}()
		continue
	}

	// We do not allow a trailing 'break' as an uninteresting last statement.
	for i := range "loop" {
		go func() {
			print(i)
		}()
		break
	}

	// Compound last statements are examined recursively to check if we understand
	// each contained statement and expression.
	for i := range "loop" {
		go func() {
			print(i) // want "loop variable i captured by func literal"
		}()
		if i > 1 {
			continue
		} else if i > 2 {
			a = i
		} else {
			a = (i + 2) * 3
		}
	}

	// We understand certain builtins, like append.
	var s []int
	for i := range "loop" {
		go func() {
			print(i) // want "loop variable i captured by func literal"
		}()
		s = append(s, 0)
	}

	// We give up on functions we don't understand.
	foo := func() int { return 0 }
	for i := range "loop" {
		go func() {
			print(i)
		}()
		s = append(s, foo())
	}

	// We do not allow a trailing panic after a go or defer statememt.
	for i := range "loop" {
		defer func() {
			print(i)
		}()
		panic("last statement")
	}

	// We understand simple uses of struct literals.
	type pair struct{ a, b int }
	var pairs []pair
	for i := range "loop" {
		go func() {
			print(i) // want "loop variable i captured by func literal"
		}()
		pairs = append(pairs, pair{a: i, b: i})
	}

	// We understand pointers to struct literals.
	var ppairs []*pair
	for i := range "loop" {
		defer func() {
			print(i) // want "loop variable i captured by func literal"
		}()
		ppairs = append(ppairs, &pair{a: i, b: i})
	}

	// We allow trailing field selectors for structs.
	for i := range "loop" {
		go func() {
			print(i) // want "loop variable i captured by func literal"
		}()
		p := pair{}
		p.a = i
		pairs = append(pairs, p)
	}

	// We do not allow trailing field selectors for pointers to
	// structs, which can panic.
	for i := range "loop" {
		go func() {
			print(i)
		}()
		p := &pair{}
		p.a = i
		ppairs = append(ppairs, p)
	}

	// We understand if the underlying type is a pointer to struct,
	// and do not allow field selectors in that case.
	type ppair *pair
	var pp ppair
	for i := range "loop" {
		go func() {
			print(i)
		}()
		pp = &pair{}
		pp.a = i
		ppairs = append(ppairs, pp)
	}

	// We understand compound assignment.
	for i := range "loop" {
		go func() {
			print(i) // want "loop variable i captured by func literal"
		}()
		a, b := 1, len(s)
		s = append(s, a, b)
	}

	// We understand variable declaration statements.
	for i := range "loop" {
		go func() {
			print(i) // want "loop variable i captured by func literal"
		}()
		// Examples taken from https://go.dev/ref/spec#Variable_declarations
		var i int
		var U, V, W float64
		var k = 0
		var x, y float32 = -1, -2
		var (
			j       int
			u, v, s = 2.0, 3.0, "bar"
		)
		_, _, _, _, _, _, _, _, _, _, _ = i, U, V, W, k, x, y, j, u, v, s
	}

	// We give up on functions we don't understand when used in a variable declaration.
	for i := range "loop" {
		go func() {
			print(i)
		}()
		var j int = foo()
		_ = j
	}

	// We understand trailing const and type declarations.
	for i := range "loop" {
		go func() {
			print(i) // want "loop variable i captured by func literal"
		}()
		type myInt int
		const j myInt = 0
		_ = j
	}

	// We understand trailing go and defer statements.
	for i := range "loop" {
		go func() {
			print(i) // want "loop variable i captured by func literal"
		}()
		defer func() {
			print(i) // want "loop variable i captured by func literal"
		}()
		go func() {
			print(i) // want "loop variable i captured by func literal"
		}()
	}

	// We understand trailing go statements nested in a compound statements like for.
	for i := range "outer" {
		go func() {
			print(i) // want "loop variable i captured by func literal"
		}()
		for j := range "inner" {
			go func() {
				print(i) // want "loop variable i captured by func literal"
				print(j) // want "loop variable j captured by func literal"
			}()
		}
	}

	// We analyze a go statement with complex arguments (e.g., calling foo here),
	// but we do not analyze a go statement preceding a go statement with complex arguments.
	for i := range "loop" {
		go func() {
			print(i)
		}()
		go func(int) {
			print(i) // want "loop variable i captured by func literal"
		}(foo())
		go func() {
			print(i) // want "loop variable i captured by func literal"
		}()
	}

	// We understand sync.WaitGroup.Add.
	var wg sync.WaitGroup
	for i := range "loop" {
		wg.Add(1)
		go func() {
			defer wg.Done()
			print(i) // want "loop variable i captured by func literal"
		}()

		wg.Add(1)
		go func() {
			defer wg.Done()
			print(i) // want "loop variable i captured by func literal"
		}()
	}

	// We understand nested loops and various compound statements prior to
	// a go or defer statement.
	for i := range "outer" {
		for j := range "inner" {
			if j < 1 {
				a++
				go func() {
					print(i) // want "loop variable i captured by func literal"
				}()
				a++
			} else {
				go func() {
					print(i)
				}()
				a++
				print("we don't catch the error above because of this statement")
				a++
			}
			a++
		}
		a++
	}

	// We understand various compound statements trailing a go or defer statement.
	for i := range "outer" {
		for j := range "inner" {
			switch j {
			case 1:
				go func() {
					print(i) // want "loop variable i captured by func literal"
				}()
				if a := i; a > 0 {
					a++
				}
			case 2:
				go func() {
					print(i) // want "loop variable i captured by func literal"
				}()
				for k := 0; k < 10; k++ {
					a++
				}
			case 3:
				defer func() {
					print(i) // want "loop variable i captured by func literal"
				}()
				for k := range "loop" {
					a = k
				}
			case 4:
				go func() {
					print(i) // want "loop variable i captured by func literal"
				}()
				switch i {
				case 1:
					i++
				}
			case 5:
				var a interface{}
				go func() {
					print(i) // want "loop variable i captured by func literal"
				}()
				switch a.(type) {
				case int:
					i++
				}
			case 6:
				var a interface{}
				go func() {
					print(i) // want "loop variable i captured by func literal"
				}()
				switch a := a.(type) {
				case int:
					i += a
				}
			}
		}
	}

	// We give up on statements and expressions we don't understand within various
	// compound statements that trail a go or defer statement. In these examples,
	// the call to foo is not understood.
	for i := range "outer" {
		for j := range "inner" {
			switch j {
			case 1: // call foo in an if init statement
				go func() {
					print(i)
				}()
				if a = foo(); a > 0 {
					a++
				}
			case 2: // call foo in a for loop init statement
				go func() {
					print(i)
				}()
				for k := foo(); k < 10; k++ {
					a++
				}
			case 3: // call foo in a range expression
				go func() {
					print(i)
				}()
				for k := range append(s, foo()) {
					a += k
				}
			case 4: // call foo in a switch statement tag
				go func() {
					print(i)
				}()
				switch foo() {
				case 0:
					a++
				}
			case 5: // call foo in a switch statement case clause expression
				go func() {
					print(i)
				}()
				switch i {
				case foo():
					a++
				}
			case 6: // call foo in a type switch init statement
				var a interface{}
				go func() {
					print(i)
				}()
				switch a = foo(); a.(type) {
				case int:
					i++
				}
			case 7: // call foo in a type switch case clause body
				var a interface{}
				go func() {
					print(i)
				}()
				switch a.(type) {
				case int:
					i = foo()
				}
			case 8: // call bar in the type assert expression of a type switch
				bar := func() interface{} { return nil }
				go func() {
					print(i)
				}()
				switch bar().(type) {
				case int:
					i++
				}
			}
		}
	}

	// Cases that exercise examining the same nested compound statement twice
	// to determine if it is understood as a simple trailing statement.
	for i := range "outer" {
		for j := range "inner" {
			go func() {
				print(i) // want "loop variable i captured by func literal"
			}()
			if i > 1 {
				j++
				if i > 2 {
					j++
				}
			}
		}
	}

	for i := range "outer" {
		for j := range "inner" {
			go func() {
				print(i)
			}()
			if i > 1 {
				j++
				if i > 2 {
					print("last statement")
				}
			}
		}
	}

	for i := range "outer" {
		for j := range "inner" {
			go func() {
				print(i)
			}()
			if i > 1 {
				j++
			} else {
				if i > 2 {
					print("last statement")
				}
			}
		}
	}

	// A case that relies on resetting our tracked loop variables when
	// we encounter a complex statement in a range expression or similar.
	// In this case, foo could in theory (via a wait or similar)
	// make the use of the outer loop variable i safe, but foo
	// cannot make the use of the inner loop variable j safe.
	for i := range "outer" {
		go func() {
			print(i)
		}()
		for j := range append(s, foo()) {
			defer func() { print(i) }()
			defer func() { print(j) }() // want "loop variable j captured by func literal"
		}
	}

	// Cases that nest range and statements within func literals.
	for i := range "outer" {
		i++
		_ = func() {
			for j := range "inner" {
				go func() { print(j) }() // want "loop variable j captured by func literal"
			}
		}
		go func() {
			for k := 0; k < 10; k++ {
				go func() { print(k) }() // want "loop variable k captured by func literal"
			}
		}()
	}

	// Some additional cases we currently purposefully disallow.

	// We disallow a trailing division, which can panic.
	for i := range "loop" {
		a, b := 1, 0
		if i > 1 {
			go func() {
				print(i)
			}()
			i = a / b
		} else {
			go func() {
				print(i)
			}()
			i /= b
		}
	}

	// We disallow a trailing use of a function we don't understand within a compound literal.
	for i := range "loop" {
		go func() {
			print(i)
		}()
		pairs = append(pairs, pair{a: foo()})
	}

	// We disallow trailing methods.
	for i := range "loop" {
		now := time.Now()
		go func() {
			print(i)
		}()
		_ = time.Since(now)
	}

	// Some additional cases we do not yet handle.

	// We do not yet understand trailing conversions.
	for i := range "loop" {
		go func() {
			print(i)
		}()
		_ = float64(i)
	}

	// We do not yet understand trailing slice literals or map literals.
	for i := range "loop" {
		if i > 1 {
			go func() {
				print(i)
			}()
			_ = []int{1, 2, 3}
		} else {
			go func() {
				print(i)
			}()
			_ = map[int]int{0: 0}
		}
	}

	// We do not yet understand nested structs.
	type nested struct {
		inner struct{ pair }
	}
	var n nested
	for i := range "loop" {
		go func() {
			print(i)
		}()
		_ = n.inner.pair
	}
}

// Illustration of why we do not report captured loop iteration variables
// with a go statement trailed by something that waits, might wait, or we can't prove won't wait.
// In this example, the func literal always completes prior to the loop variable changing.
func _() {
	for i := range "loop" {
		wg := &sync.WaitGroup{}
		wg.Add(1)
		go func() {
			defer wg.Done()
			print(i)
		}()
		wg.Wait()
	}
}

// Example from #57173.
func _() {
	process := func(s string) {}

	var processMore []string
	for _, s := range []string{"one", "two", "three"} {
		go func() { process(s) }() // want "loop variable s captured by func literal"
		switch s {
		case "two":
			continue
		default:
			processMore = append(processMore, s)
		}
	}
}

// Real-world example from #21412, slightly simplified.
// This relies on handling certain statements following a go statement,
// including this use of append.
type Batch struct{ Requests []interface{} }

func _(ctx context.Context, batch *Batch) error {
	handle := func(ctx context.Context, request interface{}) error { return nil }
	combineErrors := func(err1 error, err2 error) error { return nil }

	var chans []chan error
	for _, request := range batch.Requests {
		c := make(chan error)
		go func(c chan error) {
			c <- handle(ctx, request) // want "loop variable request captured by func literal"
		}(c)
		chans = append(chans, c)
	}

	var err error
	for _, c := range chans {
		err = combineErrors(err, <-c)
	}
	return err
}
