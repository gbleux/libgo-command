package command

const (
	flagValueName = "VAL"
)

type Flag struct {
	desc string
	env  string

	valueName string
	valueReq  bool
	value     Value
}

func (f *Flag) Bool(defaultValue bool) *bool {
	out := defaultValue
	value := &boolValue{&out}

	f.value = value

	return &out
}

func (f *Flag) BoolVar(out *bool) {
	value := &boolValue{out}

	f.value = value
}

func (f *Flag) Int(defaultValue int) *int {
	out := defaultValue
	value := &intValue{&out}

	f.value = value

	return &out
}

func (f *Flag) IntVar(out *int) {
	value := &intValue{out}

	f.value = value
}

func (f *Flag) File(defaultValue string) *string {
	out := defaultValue
	value := &fileValue{&out, false}

	f.value = value

	return &out
}

func (f *Flag) FileVar(out *string) {
	value := &fileValue{out, false}

	f.value = value
}

func (f *Flag) Dir(defaultValue string) *string {
	out := defaultValue
	value := &fileValue{&out, true}

	f.value = value

	return &out
}

func (f *Flag) DirVar(out *string) {
	value := &fileValue{out, true}

	f.value = value
}

func (f *Flag) Var(out Value) {
	f.value = out
}

func (f *Flag) Void(emulateBool bool) {
	f.value = &voidValue{emulateBool}
}

func (f *Flag) EnvironmentValue(name string) *Flag {
	f.env = name

	return f
}

func (f *Flag) Value(name string, required bool) *Flag {
	f.valueName = name
	f.valueReq = required

	return f
}

func (f *Flag) String() string {
	return f.desc
}

func newFlag(description string) *Flag {
	value := &voidValue{}

	return &Flag{description, "", flagValueName, false, value}
}
