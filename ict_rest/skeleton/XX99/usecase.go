package XX99

type useCase struct {
	repo Repository
}

func NCase(r Repository) UseCase {
	return &useCase{repo: r}
}
