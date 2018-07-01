package yonder

//empty struct (size 0, for signalling)
type empty struct{}

//Empty is an export of the empty type
type Empty empty

//Signal is a size 0 object used as a signal in channels
type Signal chan Empty
