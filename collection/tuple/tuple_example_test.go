package tuple

import "fmt"

func Example() {
	addr := Make2("localhost", 8080)
	fmt.Printf("%s:%d\n", addr.First, addr.Second)
	// Output:
	// localhost:8080
}
