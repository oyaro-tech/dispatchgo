package flags

import "flag"

type Flags struct {
	debug    *bool
	routines *int
	playbook *string
}

func New() *Flags {
	return &Flags{
		debug:    flag.Bool("debug", false, "enable debug output"),
		routines: flag.Int("routines", 10, "number of go routines"),
		playbook: flag.String("playbook", "./playbook.yaml", "a playbook yaml file"),
	}
}

func (f *Flags) Parse() {
	flag.Parse()
}

func (f *Flags) Playbook() string {
	return *f.playbook
}

func (f *Flags) Debug() bool {
	return *f.debug
}

func (f *Flags) Threads() int {
	return *f.routines
}
