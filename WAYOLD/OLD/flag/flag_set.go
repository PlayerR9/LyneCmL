package flag

var (
	// flag_set is the global flag set.
	flag_set *FlagSet
)

func init() {
	flag_set = NewFlagSet()
}

// FlagSet is a set of flags.
type FlagSet struct {
	// flag_list is the list of flags in the set.
	flag_list []*Flag
}

// NewFlagSet creates a new flag set.
//
// Returns:
//   - *FlagSet: A pointer to the new flag set. Never returns nil.
func NewFlagSet() *FlagSet {
	return &FlagSet{}
}

// AddFlag adds a flag to the flag set. Does nothing if the flag is nil.
//
// Parameters:
//   - flag: The flag to add.
func (fs *FlagSet) AddFlag(flag *Flag) {
	if flag == nil {
		return
	}

	fs.flag_list = append(fs.flag_list, flag)
}

// Bool creates a new bool flag. Panics if the long name is empty.
//
// Parameters:
//   - long_name: The long name of the flag.
//   - def_value: The default value of the flag.
//   - brief: A brief description of the flag.
//   - opts: A list of options for the flag.
//
// Returns:
//   - *bool: A pointer to the boolean value of the flag. Never returns nil.
func (fs *FlagSet) Bool(long_name string, def_value bool, brief string, opts ...FlagOption) *bool {
	if long_name == "" {
		panic("long name cannot be empty")
	}

	value := &bool_value{
		value: def_value,
	}

	flag := &Flag{
		long_name: long_name,
		brief:     brief,
		value:     value,
	}

	for _, opt := range opts {
		opt(flag)
	}

	fs.flag_list = append(fs.flag_list, flag)

	return &value.value
}

// Int creates a new int flag. Panics if the long name is empty.
//
// Parameters:
//   - long_name: The long name of the flag.
//   - def_val: The default value of the flag.
//   - brief: A brief description of the flag.
//   - opts: A list of options for the flag.
//
// Returns:
//   - *int: A pointer to the int value of the flag. Never returns nil.
func (fs *FlagSet) Int(long_name string, def_val int, brief string, opts ...FlagOption) *int {
	if long_name == "" {
		panic("long name cannot be empty")
	}

	value := &int_value{
		value: def_val,
	}

	flag := &Flag{
		long_name: long_name,
		brief:     brief,
		value:     value,
	}

	for _, opt := range opts {
		opt(flag)
	}

	fs.flag_list = append(fs.flag_list, flag)

	return &value.value
}

// String creates a new string flag. Panics if the long name is empty.
//
// Parameters:
//   - long_name: The long name of the flag.
//   - def_val: The default value of the flag.
//   - brief: A brief description of the flag.
//   - opts: A list of options for the flag.
//
// Returns:
//   - *string: A pointer to the string value of the flag. Never returns nil.
func (fs *FlagSet) String(long_name string, def_val string, brief string, opts ...FlagOption) *string {
	if long_name == "" {
		panic("long name cannot be empty")
	}

	value := &string_value{
		value: def_val,
	}

	flag := &Flag{
		long_name: long_name,
		brief:     brief,
		value:     value,
	}

	for _, opt := range opts {
		opt(flag)
	}

	fs.flag_list = append(fs.flag_list, flag)

	return &value.value
}

// Var creates a new custom flag given a value type. Panics if the long name is empty or the value is nil.
//
// Parameters:
//   - long_name: The long name of the flag.
//   - value: The value of the flag.
//   - brief: A brief description of the flag.
//   - opts: A list of options for the flag.
func (fs *FlagSet) Var(long_name string, value Valuer, brief string, opts ...FlagOption) {
	if long_name == "" {
		panic("long name cannot be empty")
	} else if value == nil {
		panic("value cannot be nil")
	}

	flag := &Flag{
		long_name: long_name,
		brief:     brief,
		value:     value,
	}

	for _, opt := range opts {
		opt(flag)
	}

	fs.flag_list = append(fs.flag_list, flag)
}

/*
// flagParser is a parser for flags.
type flagParser struct {
	// flagList is the list of flags to parse.
	flagList []*Flag

	// args is the list of arguments to parse.
	args []string
}

// findFlag finds the index of the flag.
//
// Parameters:
//   - name: The name of the flag.
//   - isShort: True if the flag is a short flag. False otherwise.
//
// Returns:
//   - int: The index of the long flag. -1 if not found.
func (p *flagParser) findFlag(name string, isShort bool) int {
	var search func(flag *Flag) bool

	if isShort {
		r, _ := utf8.DecodeRuneInString(name)

		search = func(flag *Flag) bool {
			return flag.short_name == r
		}
	} else {
		search = func(flag *Flag) bool {
			return flag.long_name == name
		}
	}

	index := slices.IndexFunc(p.flagList, search)

	return index
}

func IsBooleanFlag(flag *Flag) bool {
	if flag == nil {
		return false
	}

	f, ok := flag.value.(BoolFlager)
	if !ok {
		return false
	}

	return f.IsBoolFlag()
}


// parseSingleFlag parses a single flag.
//
// Format:
//
//	--flag		: only boolean flags without arguments
//	--flag=value	: any flag that requires an argument
//	--flag value 	: any flag that requires an argument
//
//	-flag		: only boolean flags without arguments
//	-flag=value	: any flag that requires an argument
//	-flag value 	: any flag that requires an argument
//	-abc		: multiple short flags
//
// Parameters:
//   - header: The header of the flag.
//   - right: The right side of the flag.
//   - isShort: True if the flag is a short flag. False otherwise.
//
// Returns:
//   - string: The remaining right side of the flag.
//   - error: An error if the flag failed to parse.
func (fp *flagParser) parseSingleFlag(header string, right string, isShort bool) (string, error) {
	index := fp.findFlag(header, isShort)
	if index == -1 {
		return "", NewErrFlagNotFound(isShort, header)
	}

	flag := fp.flagList[index]

	ok := IsBooleanFlag(flag)
	if ok {
		if !isShort {
			return "", fmt.Errorf("flag --%s does not take an argument", flag.long_name)
		}

		err := flag.value.Set("")
		if err != nil {
			return "", NewErrInvalidFlag(isShort, flag, err)
		}

		fp.flagList = slices.Delete(fp.flagList, index, index+1)

		if isShort {
			return right, nil
		}

		var todel int

		if right == "" {
			if len(fp.args) == 1 {
				return "", NewErrFlagMissingArg(isShort, flag)
			}

			right = fp.args[1]

			todel = 2
		} else {
			todel = 1
		}

		err = flag.value.Set(right)
		if err != nil {
			return "", NewErrInvalidFlag(isShort, flag, err)
		}

		fp.args = fp.args[todel:]

		fp.flagList = slices.Delete(fp.flagList, index, index+1)

		return "", nil
	} else {
		var todel int

		if right == "" {
			if len(fp.args) == 1 {
				return "", NewErrFlagMissingArg(isShort, flag)
			}

			right = fp.args[1]

			todel = 2
		} else {
			todel = 1
		}

		err := flag.value.Set(right)
		if err != nil {
			return "", NewErrInvalidFlag(isShort, flag, err)
		}

		fp.args = fp.args[todel:]

		fp.flagList = slices.Delete(fp.flagList, index, index+1)

		return "", nil
	}
}

// parseOne parses one flag.
//
// Parameters:
//   - header: The header of the flag.
//   - isShort: True if the flag is a short flag. False otherwise.
//
// Returns:
//   - error: An error if the flag failed to parse.
func (fp *flagParser) parseOne(header string, isShort bool) error {
	fields := strings.SplitN(header, "=", 2)
	header = fields[0]

	var right string

	if len(fields) == 2 {
		right = fields[1]
	} else {
		right = ""
	}

	if !isShort {
		_, err := fp.parseSingleFlag(header, right, isShort)
		if err != nil {
			return err
		}

		return nil
	}

	runes, err := luch.StringToUtf8(header)
	if err != nil {
		return err
	}

	if len(runes) == 1 {
		right, err = fp.parseSingleFlag(header, right, true)
		if err != nil {
			ok := uc.Is[*ErrFlagNotFound](err)
			if ok {
				err = nil
			}
		} else if right != "" {
			err = fmt.Errorf("extra argument %q", right)
		}
	} else {
		for _, letter := range runes {
			right, err = fp.parseSingleFlag(string(letter), right, true)
			if err != nil {
				ok := uc.Is[*ErrFlagNotFound](err)
				if ok {
					err = fmt.Errorf("invalid merged short flags -%s", header)
				}

				return err
			}
		}

		if right != "" {
			err = fmt.Errorf("extra argument %q", right)
		}
	}

	return err
}

func (fs *FlagSet) Parse() error {

}

// isValidFlag checks if the flag is valid.
//
// Parameters:
//   - arg: The argument to check.
//
// Returns:
//   - bool: True if the flag is short. False otherwise.
//   - bool: True if the flag is valid. False otherwise.
func isValidFlag(arg string) (bool, bool) {
	ok := strings.HasPrefix(arg, LongFlagPrefix)
	if ok {
		return false, true
	}

	ok = strings.HasPrefix(arg, ShortFlagPrefix)

	return true, ok
}

type parsingResult struct {
	cutSet    []string
	data      any
	argsLeft  []string
	flagsLeft []*Flag
}

type flagResult struct {
	flagList []*Flag
	args     []string
}

func parseFlags(flagMap map[string]*Flag, args []string) (*flagResult, error) {
	var flagList []*Flag

	for _, value := range flagMap {
		flagList = append(flagList, value)
	}

	fp := &flagParser{
		flagList: flagList,
		args:     args,
	}

	for len(fp.args) > 0 {
		header := fp.args[0]

		isShort, ok := isValidFlag(header)
		if !ok {
			fr := &flagResult{
				flagList: fp.flagList,
				args:     fp.args,
			}

			return fr, nil
		}

		if isShort {
			header = strings.TrimPrefix(header, ShortFlagPrefix)
		} else {
			header = strings.TrimPrefix(header, LongFlagPrefix)
		}

		err := fp.parseOne(header, isShort)
		if err != nil {
			return nil, NewErrInvalidArg(header, isShort, err)
		}

		if len(fp.flagList) == 0 {
			fr := &flagResult{
				flagList: nil,
				args:     fp.args,
			}
			return fr, nil
		}
	}

	fr := &flagResult{
		flagList: fp.flagList,
		args:     nil,
	}

	return fr, nil
}
*/
