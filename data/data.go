package data

type AccountEntity struct {
	UserId        string
	AccountNumber string
	Currency      string
}

type Repository interface {
	Save(AccountEntity) (AccountEntity, error)
	FindByUserId(string) (*AccountEntity, error)
}

type inMemoryRepository struct {
	m map[string]AccountEntity
}

func NewInMemoryRepository() *inMemoryRepository {
	return &inMemoryRepository{
		m: make(map[string]AccountEntity),
	}
}

func (r *inMemoryRepository) Save(entity AccountEntity) (AccountEntity, error) {
	r.m[entity.UserId] = entity
	return entity, nil
}

func (r *inMemoryRepository) FindByUserId(userId string) (*AccountEntity, error) {
	entity, exist := r.m[userId]

	if !exist {
		return nil, nil
	}

	return &entity, nil
}
