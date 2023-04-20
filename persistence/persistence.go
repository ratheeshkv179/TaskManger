package persistence

type Persistence interface {
	Get()
	GetOne()
	UpSert()
	Delete()
	Create()
}
