package common

///////////////////////////////////////////////////////

type Fixer interface {
	// Fix applies a thorough check on the component and fixes any issues.
	Fix()
}
