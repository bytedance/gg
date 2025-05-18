// Package goption provides a generic representation of optional value.
//
// Every [O] contains a value or nothing.
//
// # Simplifying the "if v, ok := ...; ok {...}" pattern
//
// Use [os.LookupEnv] as example:
//
// The trivial way:
//
//	sh, ok := os.LookupEnv("SHELL")
//	if !ok {
//	    // Do something.
//	}
//	return sh
//
// Use optional value:
//
//	// The env must be present, otherwise panic.
//	return Of(os.LookupEnv("SHELL")).Value()
//
//	// Return zero value when the env is not present.
//	return Of(os.LookupEnv("SHELL")).Value()
//
//	// Return a custom value when the env is not present.
//	return Of(os.LookupEnv("SHELL")).ValueOr("/bin/sh")
//
// # JSON
//
// [O] implements [encoding/json.Marshaler] and [encoding/json.Ummarshaler], so
// you can use it in JSON marshaling/unmarshaling.
// See [goption.O.MarshalJSON] and [goption.O.UnmarshalJSON].
package goption
