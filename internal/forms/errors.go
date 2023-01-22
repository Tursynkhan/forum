package forms

type Err map[string][]string

func (e Err) Add(field, message string) {
	e[field] = append(e[field], message)
}

func (e Err) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}
