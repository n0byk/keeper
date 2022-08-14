package internal

var (
	Publish        = "publish"
	Subscribe      = "subscribe"
	Unsubscribe    = "unsubscribe"
	Register       = "register"
	Auth           = "auth"
	LobbyCreate    = "create_lobby"
	GetPublicToken = "get_public_token"
	GetLobbyId     = "lobby_id"
	LobbyInvite    = "invite_lobby"
	SetData        = "set_data"
)

type TRegistrationRequest struct {
	Action   string `json:"action"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TAuthRequest struct {
	Action   string `json:"action"`
	Login    string `json:"login"`
	Password string `json:"password"`
}

type TCreateLobbyRequest struct {
	Action    string `json:"action"`
	Token     string `json:"token"`
	LobbyName string `json:"lobby_name"`
}

type TGetPublicTokenRequest struct {
	Action string `json:"action"`
	Login  string `json:"login"`
}

type TLobbyIdRequest struct {
	Action    string `json:"action"`
	LobbyName string `json:"lobby_name"`
}

type TAddToLobbyRequest struct {
	Action      string `json:"action"`
	LobbyName   string `json:"lobby_name"`
	Token       string `json:"token"`
	PublicToken string `json:"public_token"`
}

type TSetDataRequest struct {
	Action    string `json:"action"`
	Token     string `json:"token"`
	LobbyName string `json:"lobby_name"`
	LobbyId   string `json:"lobby_id"`
	RowData   string `json:"row_data"`
}
