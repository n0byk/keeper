package distribution

func ActionExchange(msg Message) (string, error) {

	switch msg.Action {
	case "registration":
		return "Asdasdasdasgfg", nil
	default:
		return "", nil
	}

}
