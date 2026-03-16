package url

type URLService interface {
	Save(url string) (string, error)
	Get(string) (string, error)
}

type URLHandler struct {
	service URLService
}

func NewURLHandler(service URLService) *URLHandler {
	return &URLHandler{
		service: service,
	}
}
