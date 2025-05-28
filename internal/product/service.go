package product

type ProductService interface {
	
}

type productService struct  {
	repository ProductRepository
}

func NewProductService(repository ProductRepository) ProductService {
	return &productService{
		repository: repository,
	}
}