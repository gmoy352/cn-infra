package datasync

// KeyValIterator is an iterator for KeyVals
type KeyValIterator interface {
	// GetNext retrieves the next value from the iterator context.  The retrieved
	// value is unmarshaled into the provided argument. The allReceived flag is
	// set to true on the last KeyVal pair in the context.
	GetNext() (kv KeyVal, allReceived bool)
}

// KeyVal represents a single key-value pair
type KeyVal interface {
	WithKey
	LazyValueWithRev
}
