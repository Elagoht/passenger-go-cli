package schemas

type ResponseStatus struct {
	Status bool `json:"initialized"`
}

type ResponseRegister struct {
	Recovery string `json:"recovery"`
}

type ResponseLogin struct {
	Token string `json:"token"`
}
