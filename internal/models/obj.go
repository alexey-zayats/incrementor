package models

import "time"

// Obj base struct to conform sql inheritans
type Obj struct {
	GUID    string
	Created time.Time
	Updated time.Time
}
