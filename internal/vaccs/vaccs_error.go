package vaccs

type Vaccs_error struct {
	message string
}

func (v *Vaccs_error) Error() string {
	return v.message
}
