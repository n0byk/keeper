package distribution

type RegistrationResponse struct {
	Token string
}

func ActionExchange(sbj string, msg interface{}) (interface{}, error) {

	switch sbj {
	case "keeper.registration":
		return RegistrationResponse{Token: "asdasdas"}, nil
	default:
		return "", nil
	}

}
