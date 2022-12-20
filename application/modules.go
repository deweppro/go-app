package application

// Modules DI container
type Modules []interface{}

// Add object to container
func (m Modules) Add(v ...interface{}) Modules {
	for _, mod := range v {
		switch v := mod.(type) {
		case Modules:
			m = m.Add(v...)
		default:
			m = append(m, mod)
		}
	}
	return m
}
