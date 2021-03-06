package tri

import (
	"time"
	"reflect"
	"errors"
	"fmt"
	"unicode"
)

// Validate checks to ensure the contents of this node type satisfy constraints.
// Brief only contains one thing, so we make sure it has it - one string. This string may not contain any type of control characters, and is limited to 80 characters in length.
func (r *Brief) Validate() error {

	R := (*r)
	if len(R) != 1 {
		return errors.New("Brief field must have (only) one item")
	}
	s, ok := R[0].(string)
	if !ok {
		return errors.New("Brief's mandatory field is not a string")
	}
	if len(s) > 80 {
		return errors.New("Brief's text may not be over 80 characters in length")
	}
	for i, x := range s {
		if unicode.IsControl(x) {
			return fmt.Errorf("Brief text may not contain any control characters, one found at position %d", i)
		}
	}
	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// This validator only has to check the elements of the slice are zero or more Command items, and a valid name at index 0.
func (r *Command) Validate() error {

	R := *r
	if len(R) < 1 {
		return errors.New("empty Command")
	}
	s, ok := R[0].(string)
	if !ok {
		return fmt.Errorf("first element of Command must be a string")
	}
	if e := ValidName(s); e != nil {
		return fmt.Errorf("error in name of Command: %v", e)
	}
	// validSet is an array of 4 elements that represent the presence of the 4 mandatory parts.
	var validSet [2]bool
	brief, handler := 0, 1
	var singleSet [4]bool
	usage, short, help, examples := 0, 1, 2, 3
	for i, x := range R[1:] {
		switch c := x.(type) {
		case Short:
			if singleSet[short] {
				return fmt.Errorf("only one Short field allowed in Command")
			}
			singleSet[short] = true

			e := c.Validate()
			if e != nil {
				return fmt.Errorf("error in Command at index %d: %v", i, e)
			}
		case Brief:
			if validSet[brief] {
				return fmt.Errorf("only one Brief permitted in a Command, second found at index %d", i)
			}
			validSet[brief] = true
			e := c.Validate()
			if e != nil {
				return fmt.Errorf("error in Command at index %d: %v", i, e)
			}
		case Usage:
			if singleSet[usage] {
				return fmt.Errorf("only one Usage field allowed in Command")
			}
			singleSet[usage] = true
			e := c.Validate()
			if e != nil {
				return fmt.Errorf("error in Command at index %d: %v", i, e)
			}
		case Help:
			if singleSet[help] {
				return fmt.Errorf("only one Help field allowed in Command")
			}
			singleSet[help] = true

			e := c.Validate()
			if e != nil {
				return fmt.Errorf("error in Command at index %d: %v", i, e)
			}
		case Examples:
			if singleSet[examples] {
				return fmt.Errorf("only one Examples field allowed in Command")
			}
			singleSet[examples] = true

			e := c.Validate()
			if e != nil {
				return fmt.Errorf("error in Command at index %d: %v", i, e)
			}
		case Var:
			e := c.Validate()
			if e != nil {
				return fmt.Errorf("error in Command at index %d: %v", i, e)
			}
		case Trigger:
			e := c.Validate()
			if e != nil {
				return fmt.Errorf("error in Command at index %d: %v", i, e)
			}
		case func(*Tri) int:
			if validSet[handler] {
				return fmt.Errorf("only one Handler permitted in a Command, second found at index %d", i)
			}
			validSet[handler] = true
			if c == nil {
				return fmt.Errorf("nil handler in Command found at index %d", i)
			}
		default:
			return fmt.Errorf("invalid type present in Command: %v", reflect.TypeOf(c))
		}
	}
	if !validSet[brief] {
		return errors.New("Brief field must be present")
	}
	if !validSet[handler] {
		return errors.New("Command must have a handler")
	}
	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// This validator only triggers the validator on its elements.
func (r *Commands) Validate() error {

	R := (*r)
	for i, x := range R {
		e := x.Validate()
		if e != nil {
			return fmt.Errorf("error in element %d of Commands list: %v", i, e)
		}
	}
	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// The only constraint on the Default subtype is that it contains at only one element, the value is checked for correct typing by the Commands validator.
func (r *Default) Validate() error {

	R := (*r)
	if len(R) != 1 {
		return errors.New("the Default container must only contain one element")
	}
	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// The constraint of DefaultCommand is that it has at least one element, and that the 0 element is a string. The check for the command name's presence in the Commands set is in the Tri validator.
func (r *DefaultCommand) Validate() error {

	R := (*r)
	if len(R) != 1 {
		return errors.New(
			"the DefaultCommand element must contain only one element")
	}
	s, ok := R[0].(string)
	if !ok {
		return errors.New("element 0 of DefaultCommand must be a string")
	}
	if e := ValidName(s); e != nil {
		return fmt.Errorf("error in DefaultCommand: %v", e)
	}
	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// RunAfter is a simple flag that indicates by existence of an empty value, so it is an error if it has anything inside it.
func (r *DefaultOn) Validate() error {

	R := *r
	if len(R) > 0 {
		return errors.New(
			"DefaultOn may not contain anything, empty declaration only")
	}
	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// The constraints of examples are minimum two elements and all elements are strings. The intent is the even numbered items are snippets showing invocation and a description string of the same format as Brief{}.
func (r *Examples) Validate() error {

	R := *r
	if len(R) < 2 {
		return errors.New("Examples field may not be empty")
	}
	if len(R)%2 != 0 {
		return fmt.Errorf(
			"Examples must be in pairs, odd number of elements found")
	}
	for i, x := range R {
		_, ok := x.(string)
		if !ok {
			return fmt.Errorf(
				"Examples elements may only be strings, element %d is not a string", i)
		}
	}

	for i := 1; i <= len(R)-1; i += 2 {
		if len(R[i-1].(string)) > 40 {
			return errors.New(
				"Examples example text may not be over 4 characters in length")
		}
		if len(R[i].(string)) > 80 {
			return errors.New(
				"Examples explainer text may not be over 80 characters in length")
		}
		for i, x := range R[i].(string) {
			if unicode.IsControl(x) {
				return fmt.Errorf(
					"Examples even numbered field string may not contain control characters, one found at index %d", i)
			}
		}
	}
	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// A group must contain one string, anything else is invalid. It also has the same limitation as a name - only letters.
func (r *Group) Validate() error {

	R := *r
	if len(R) != 1 {
		return errors.New("Group must (only) contain one element")
	}
	s, ok := R[0].(string)
	if !ok {
		return errors.New("Group element must be a string")
	}
	if e := ValidName(s); e != nil {
		return fmt.Errorf("error in name of Group: %v", e)
	}
	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// Help may only contain one string. It will be parsed as markdown format and possibly can be set to style it with ANSI codes.
func (r *Help) Validate() error {

	R := *r
	if len(R) != 1 {
		return errors.New("Help field must contain (only) one item")
	}
	_, ok := R[0].(string)
	if !ok {
		return errors.New("Help field element is not a string")
	}
	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// RunAfter is a simple flag that indicates by existence of an empty value, so it is an error if it has anything inside it.
func (r *RunAfter) Validate() error {

	R := *r
	if len(R) > 0 {
		return errors.New(
			"RunAfter may not contain anything, empty declaration only")
	}
	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// Short names contain only a single Rune variable.
func (r *Short) Validate() error {

	R := *r
	if len(R) != 1 {
		return errors.New("Short name item must contain (only) one item")
	}
	s, ok := R[0].(rune)
	if !ok {
		return errors.New("Short's element must be a rune (enclose in '')")
	}
	if !(unicode.IsDigit(s) || unicode.IsLetter(s)) {
		return errors.New("Short element is not a letter or number")
	}
	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// Slot may only contain one type of element. The type check is in the Var, here we only ensure the slots contain pointers to the same type, the parser will put the final parsed value in all of them. Multiple variables are permitted here to enable the configuration of more than one application.
func (r *Slot) Validate() error {

	R := *r
	var slotTypes []reflect.Type
	for _, x := range R {
		slotTypes = append(slotTypes, reflect.TypeOf(x))
	}
	for i, x := range slotTypes {
		if i > 0 {
			if slotTypes[i] != slotTypes[i-1] {
				return fmt.Errorf("slot contains more than one type of variable, found %v at index %d", x, i)
			}
		}
	}
	for _, x := range R {
		if reflect.ValueOf(x).Kind() != reflect.Ptr {
			return fmt.Errorf("slot contains non-pointer type")
		}
	}

	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// Terminates is a flag value, and may not contain anything.
func (r *Terminates) Validate() error {

	R := *r
	if len(R) > 0 {
		return errors.New("Terminates type may not contain anything")
	}
	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// A Tri, the base type, in a declaration must contain a name as first element, a Brief, Version and a Commands item, and only one of each. Also, this and several other subtypes of Tri.
func (r *Tri) Validate() error {
	R := *r
	if len(R) < 3 {
		return errors.New("a Tri must contain at least 3 elements: name, Brief and Version")
	}
	// validSet is an array of 4 elements that represent the presence of the 4 mandatory parts.
	var validSet [2]bool
	brief, version := 0, 1
	var singleSet [3]bool
	defcom, commands := 0, 1
	n, ok := R[0].(string)
	if !ok {
		return errors.New("first element of a Tri must be a string")
	}
	if e := ValidName(n); e != nil {
		return fmt.Errorf("error in name of Tri: %v", e)
	}

	// The mandatory elements also may not be repeated:
	for i, x := range R {
		if i == 0 {
			continue
		}
		switch y := x.(type) {
		case Brief:
			if validSet[brief] {
				return fmt.Errorf(
					"Tri contains more than one Brief, second found at index %d", i)
			}
			validSet[brief] = true
			if e := y.Validate(); e != nil {
				return fmt.Errorf("Tri field %d: %s", i, e)
			}
		case Version:
			if validSet[version] {
				return fmt.Errorf(
					"Tri contains more than one Version, second found at index %d", i)
			}
			validSet[version] = true
			if e := y.Validate(); e != nil {
				return fmt.Errorf("Tri field %d: %s", i, e)
			}
			validSet[version] = true
		case Commands:
			if singleSet[commands] {
				return fmt.Errorf(
					"Tri contains more than one Commands, second found at index %d", i)
			}
			singleSet[commands] = true
			e := y.Validate()
			if e != nil {
				return fmt.Errorf("error in Tri field %d: %s", i, e)
			}
		case Var:
			e := y.Validate()
			if e != nil {
				return fmt.Errorf("error in Tri at index %d: %v", i, e)
			}
		case Trigger:
			e := y.Validate()
			if e != nil {
				return fmt.Errorf("error in Tri at index %d: %v", i, e)
			}
		case DefaultCommand:
			if singleSet[defcom] {
				return fmt.Errorf(
					"Tri contains more than one DefaultCommand, second found at index %d", i)
			}
			singleSet[defcom] = true
			e := y.Validate()
			if e != nil {
				return fmt.Errorf("Tri field %d: %s", i, e)
			}
			commname := y[0].(string)
			// DefaultCommand must match in its name one of the Command items in also present Commands array
			foundComm := false
			foundDefComm := false
			for _, a := range R {
				switch c := a.(type) {
				case Commands:
					foundComm = true
					for _, b := range c {
						if b[0].(string) == commname {
							foundDefComm = true
						}
					}
				default:
				}
			}
			if !foundComm {
				return errors.New("DefaultCommand with no Commands array present")
			} else if !foundDefComm {
				return errors.New("DefaultCommand found with no matching Command")
			}

		default:
			return fmt.Errorf(
				"Tri contains an element type it may not contain at index %d", i)
		}
	}
	switch {
	case !validSet[brief]:
		return errors.New("Tri is missing its Brief field")
	case !validSet[version]:
		return errors.New("Tri is missing its Version field")
	}
	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// Trigger must contain (one) name, Brief and Handler, and nothing other than these and Short, Usage, Help, Default, Terminates, RunAfter.
func (r *Trigger) Validate() error {

	R := *r
	if len(R) < 3 {
		return errors.New(
			"Trigger must contain a name, Brief and Handler at minimum")
	}
	name, ok := R[0].(string)
	if !ok {
		return errors.New("first element of Trigger must be the name")
	} else if e := ValidName(name); e != nil {
		return fmt.Errorf("Invalid Name in Trigger at index 0: %v", e)
	}
	// validSet is an array that represent the presence of the mandatory parts.
	var validSet [2]bool
	brief, handler := 0, 1
	var singleSet [7]bool
	short, usage, help, defon, terminates, runafter, group := 0, 1, 2, 3, 4, 5, 6
	for i, x := range R[1:] {

		switch y := x.(type) {

		case Brief:
			if validSet[brief] {
				return fmt.Errorf("Trigger may (only) contain one Brief, second found at index %d", i)
			} else {
				validSet[brief] = true
			}
			if e := y.Validate(); e != nil {
				return fmt.Errorf(
					"Trigger contains invalid element at %d - %s", i, e)
			}

		case func(*Tri) int:
			if validSet[handler] {
				return fmt.Errorf(
					"Trigger may (only) contain one Handler, second found at index %d", i)
			} else {
				validSet[handler] = true
			}
			if y == nil {
				return fmt.Errorf("Handler at index %d may not be nil", i)
			}

		case Short:
			if singleSet[short] {
				return fmt.Errorf("Trigger may only contain one Short, extra found at index %d", i)
			}
			singleSet[short] = true
			if e := y.Validate(); e != nil {
				return fmt.Errorf(
					"Trigger contains invalid element at %d - %s", i, e)
			}

		case Usage:
			if singleSet[usage] {
				return fmt.Errorf("Trigger may only contain one Usage, extra found at index %d", i)
			}
			singleSet[usage] = true
			if e := y.Validate(); e != nil {
				return fmt.Errorf(
					"Trigger contains invalid element at %d - %s", i, e)
			}

		case Help:
			if singleSet[help] {
				return fmt.Errorf("Trigger may only contain one Help, extra found at index %d", i)
			}
			singleSet[help] = true
			if e := y.Validate(); e != nil {
				return fmt.Errorf(
					"Trigger contains invalid element at %d - %s", i, e)
			}

		case DefaultOn:
			if singleSet[defon] {
				return fmt.Errorf("Trigger may only contain one DefaultOn, extra found at index %d", i)
			}
			singleSet[defon] = true
			if e := y.Validate(); e != nil {
				return fmt.Errorf(
					"Trigger contains invalid element at %d - %s", i, e)
			}

		case Terminates:
			if singleSet[terminates] {
				return fmt.Errorf("Trigger may only contain one Terminates, extra found at index %d", i)
			}
			singleSet[terminates] = true
			if e := y.Validate(); e != nil {
				return fmt.Errorf(
					"Trigger contains invalid element at %d - %s", i, e)
			}

		case RunAfter:
			if singleSet[runafter] {
				return fmt.Errorf("Trigger may only contain one RunAfter, extra found at index %d", i)
			}
			singleSet[runafter] = true
			if e := y.Validate(); e != nil {
				return fmt.Errorf(
					"Trigger contains invalid element at %d - %s", i, e)
			}

		case Group:
			if singleSet[group] {
				return fmt.Errorf(
					"Trigger may only contain one Group, extra found at index %d", i)
			}
			singleSet[group] = true
			if e := y.Validate(); e != nil {
				return fmt.Errorf(
					"Trigger contains invalid element at %d - %s", i, e)
			}

		default:
			return fmt.Errorf(
				"found invalid item type at element %d in a Trigger", i)
		}
	}
	if !(validSet[brief] && validSet[handler]) {
		return errors.New("Trigger must contain one each of Brief and Handler")
	}
	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// Usage fields contain only one string of no more than 80 characters and no control characters.
func (r *Usage) Validate() error {
	R := *r
	if len(R) != 1 {
		return errors.New("Usage field must contain (only) one element")
	}
	s, ok := R[0].(string)
	if !ok {
		return errors.New("Usage field element is not a string")
	}
	if ll := len(s); ll > 80 {
		return fmt.Errorf("Usage string is %d chars long, may not be longer than 80", ll)
	}
	for i, x := range s {
		if unicode.IsControl(x) {
			return fmt.Errorf(
				"Usage field string may not contain control characters, one found at index %d", i)
		}
	}
	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// Var must contain name, Brief and Slot, and optionally, Short, Usage, Help and Default. The type in the Slot and the Default must be the same.
func (r *Var) Validate() error {

	R := *r
	if len(R) < 3 {
		return errors.New(
			"Var must contain a name, Brief and Slot at minimum")
	}
	name, ok := R[0].(string)
	if !ok {
		return errors.New("first element of Var must be the name")
	} else if e := ValidName(name); e != nil {
		return fmt.Errorf("Invalid Name in Var at index 0: %v", e)
	}
	// validSet is an array that represent the presence of the mandatory parts.
	var validSet [2]bool
	brief, slot := 0, 1
	// singleSet is an array representing the optional elements that may not be more than one inside a Var
	var singleSet [5]bool
	short, usage, help, def, group := 0, 1, 2, 3, 4
	for i, x := range R[1:] {

		switch y := x.(type) {

		case Brief:
			if validSet[brief] {
				return fmt.Errorf("Var may must (only) contain one Brief, second found at index %d", i)
			} else {
				validSet[brief] = true
			}
			if e := y.Validate(); e != nil {
				return fmt.Errorf(
					"Var contains invalid element at %d - %s", i, e)
			}

		case Short:
			if singleSet[short] {
				return fmt.Errorf("Var may only contain one Short, extra found at index %d", i)
			}
			singleSet[short] = true
			if e := y.Validate(); e != nil {
				return fmt.Errorf(
					"Var contains invalid element at %d - %s", i, e)
			}

		case Usage:
			if singleSet[usage] {
				return fmt.Errorf("Var may only contain one Usage, extra found at index %d", i)
			}
			singleSet[usage] = true
			if e := y.Validate(); e != nil {
				return fmt.Errorf(
					"Var contains invalid element at %d - %s", i, e)
			}

		case Help:
			if singleSet[help] {
				return fmt.Errorf("Var may only contain one Help, extra found at index %d", i)
			}
			singleSet[help] = true
			if e := y.Validate(); e != nil {
				return fmt.Errorf(
					"Var contains invalid element at %d - %s", i, e)
			}

		case Default:
			if singleSet[def] {
				return fmt.Errorf("Var may only contain one Default, extra found at index %d", i)
			}
			singleSet[def] = true
			if e := y.Validate(); e != nil {
				return fmt.Errorf(
					"Var contains invalid element at %d - %s", i, e)
			}
			for _, z := range R {
				s, ok := z.(Slot)
				if ok {
					switch s[0].(type) {
					case *string:
						_, ok := y[0].(string)
						if !ok {
							return errors.New("slot is not same type as default")
						}
					case *int:
						_, ok := y[0].(int)
						if !ok {
							return errors.New("slot is not same type as default")
						}
					case *uint32:
						_, ok := y[0].(uint32)
						if !ok {
							return errors.New("slot is not same type as default")
						}
					case *float64:
						_, ok := y[0].(float64)
						if !ok {
							return errors.New("slot is not same type as default")
						}
					case *[]string:
						_, ok := y[0].([]string)
						if !ok {
							return errors.New("slot is not same type as default")
						}
					case *time.Duration:
						_, ok := y[0].(time.Duration)
						if !ok {
							return errors.New("slot is not same type as default")
						}
					}
					// *s[0] = *y[0]
				}
			}

		case Slot:
			if validSet[slot] {
				return fmt.Errorf("Var may only contain one Slot, extra found at index %d", i)
			}
			validSet[slot] = true
			if e := y.Validate(); e != nil {
				return fmt.Errorf(
					"Var contains invalid element at %d - %s", i, e)
			}

		case Group:
			if singleSet[group] {
				return fmt.Errorf(
					"Var may only contain one Group, extra found at index %d", i)
			}
			singleSet[group] = true
			if e := y.Validate(); e != nil {
				return fmt.Errorf(
					"Var contains invalid element at %d - %s", i, e)
			}
		default:
			return fmt.Errorf(
				"found invalid item type at element %d in a Var", i)
		}
	}
	if !(validSet[brief] && validSet[slot]) {
		return errors.New("Var must contain one each of Brief and Slot")
	}
	// TODO: check that Default value can be assigned to dereferenced Slot variable

	return nil
}

// Validate checks to ensure the contents of this node type satisfy constraints.
// A version item contains three integers and an optional (less than 16 character) string, and the numbers may not be more than 99.
func (r *Version) Validate() error {

	R := *r
	if len(R) > 4 {
		return errors.New("Version field may not contain more than 4 fields")
	}
	if len(R) < 3 {
		return errors.New("Version field must contain at least 3 fields")
	}
	for i, x := range R[:3] {
		n, ok := x.(int)
		if !ok {
			return fmt.Errorf("Version field %d is not an integer: %d", i, n)
		}
		if n > 99 {
			return fmt.Errorf("Version field %d value is over 99: %d", i, n)
		}
	}
	if len(R) > 3 {
		s, ok := R[3].(string)
		if !ok {
			return fmt.Errorf("optional field 4 of Version is not a string")
		} else {
			for i, x := range s {
				if !(unicode.IsLetter(x) || unicode.IsDigit(x)) {
					return fmt.Errorf(
						"optional field 4 of Version contains other than letters and numbers at position %d: '%v,", i, x)
				}
			}
		}
	}
	return nil
}

// ValidName checks that a Tri name element that should be a name only contains letters.
func ValidName(s string) error {

	if len(s) < 3 {
		return errors.New("name is less than 3 characters long")
	}
	for i, x := range s {
		if !unicode.IsLetter(x) {
			return fmt.Errorf(
				"element %d, '%v' of name is not a letter", i, x)
		}
	}
	return nil
}
