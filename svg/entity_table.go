package svg

import "fmt"

// entityTable holds the names a drawing defines in its own doctype for itself:
// a name is written once there and every use of it elsewhere stands for the
// text it was given. A name whose text lives in another document is kept with
// a mark, since that document is not read and the name therefore does not
// resolve to anything.
type entityTable map[string]entity

type entity struct {
	text     string
	external bool
}

// declare adds an entity to the table. A name defined twice is an error
// rather than a silent replacement — a file saying two things with one name
// does not say either — and the five entities XML defines for every document
// are not a file's to define again.
func (t entityTable) declare(name string, e entity) error {
	switch name {
	case "lt", "gt", "amp", "apos", "quot":
		return fmt.Errorf("antui/svg: the entity &%s; is one XML defines for every document, and a file may not define it anew", name)
	}
	if _, seen := t[name]; seen {
		return fmt.Errorf("antui/svg: the entity &%s; is defined twice", name)
	}
	t[name] = e
	return nil
}

// declareEntity keeps what a `<!ENTITY …>` said, in the table its kind belongs
// to: a drawing's own names, or the doctype's parameter ones, which hold
// pieces of the doctype itself rather than of the drawing's text.
func (p *xmlParser) declareEntity(name string, e entity, parameter bool) error {
	table := p.ents
	if parameter {
		table = p.pes
	}
	if table == nil {
		table = entityTable{}
	}
	if err := table.declare(name, e); err != nil {
		return err
	}
	if parameter {
		p.pes = table
	} else {
		p.ents = table
	}
	return nil
}
