package XX99

type Handler struct {
	usecase UseCase
}

func NHand(u UseCase) *Handler {
	return &Handler{usecase: u}
}
