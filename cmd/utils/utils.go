package utils

type IParse interface {
	Args(string) (string, error)
}

type Parse struct {
}

func (k Parse) Args(args string) (string, error) {

	return "", nil
}
func Parser() {
	//prs := Parse{}
	//kk, ee := prs.Args("ke")
}
