package template

import "time"

// defaultLoadingTimeout is how long a program may take to load before it
// closes itself, which is the number [SetLoadingTimeout] starts at and the
// one a program that says nothing gets: long enough for a stylesheet, a
// font and a database of icons, short enough that a program which never
// says Ready does not sit on a loading picture for ever.
const defaultLoadingTimeout = 5 * time.Second
