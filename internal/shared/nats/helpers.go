package xnats

const (
	SubjectVersionHeaderKey = "subjectversion"
)

// SubjectVersionFromMap returns the value of the AggregateVersionHeaderKey from the given map.
func SubjectVersionFromMap(m map[string]interface{}) string {
	if v, ok := m[SubjectVersionHeaderKey]; ok {
		return v.(string)
	}
	return ""
}
