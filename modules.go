package app

//Modules DI container
type Modules []interface{}

//Add object to container
func (m Modules) Add(v ...interface{}) Modules {
	for _, mod := range v {
		switch mod.(type) {
		case Modules:
			m = m.Add(mod.(Modules)...)
		default:
			m = append(m, mod)
		}
	}
	return m
}
