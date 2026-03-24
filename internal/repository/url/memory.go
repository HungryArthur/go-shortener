package repository

type URLMemoryRepository struct {
	// сокращенный url -> куда ведет сокращенный url
	st map[string]string
}

func NewURLMemoryRepository() *URLMemoryRepository {
	return &URLMemoryRepository{
		st: make(map[string]string),
	}
}
func (r *URLMemoryRepository) Save(sourceURL, shortened string) error {
	_, ok := r.st[shortened]
	if ok {
		return ErrShortenedURLAlreadyExists
	}

	r.st[shortened] = sourceURL
	return nil
}

func (r *URLMemoryRepository) Get(shortenedURL string) (string, error) {
	srcURL, ok := r.st[shortenedURL]
	if !ok {
		return "", ErrShortenedURLDoesntExist
	}
	return srcURL, nil
}


func (r *URLMemoryRepository) Ping() error {
	return nil
}