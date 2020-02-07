package excel

type Attributes map[string]string

type Files map[string]*File

type Elements []Attributes

func (e Elements) Get(key, value string) (element Attributes) {
	for _, element := range e {
		if element[key] == value {
			return element
		}
	}
	return
}
