// Package gresult provides a generic union type for value or error.
//
// Every [R] contains a value, representing success and containing a value,
// or an error, representing failure and containing an error.
//
// # Simplifying the "if v, err := ...; err != nil {...}" pattern
//
// Use [os.Open] as example:
//
// The trivial way:
//
//	fd, err := os.Open("/tmp/error.log")
//	if err != nil {
//	    // Do something.
//	}
//	return fd
//
// Use result:
//
//	// The file must be present, otherwise panic.
//	return Of(os.Open("/tmp/error.log")).Value()
//
//	// Return zero value when the file is not present.
//	return Of(os.Open("/tmp/error.log")).Value()
//
//	// Return a custom file object when the file is not present.
//	return Of(os.Open("/tmp/error.log")).ValueOr(os.Stderr)
//
// # JSON
//
// [R] implements [encoding/json.Marshaler] and [encoding/json.Ummarshaler], so
// you can use it in JSON marshaling/unmarshaling.
// See [gresult.R.MarshalJSON] and [gresult.R.UnmarshalJSON].
package gresult
