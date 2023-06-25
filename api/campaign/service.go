package campaign

type Service interface {
	GetAll()
}
type service struct {
	repository Repository
}
func NewService(repository Repository) *service {
	return &service{}
}
