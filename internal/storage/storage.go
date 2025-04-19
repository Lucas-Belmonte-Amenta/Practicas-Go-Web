package storage

type Storage interface {
	Read(emptyListEntity any) error
	Write(emptyListEntity any) error
}
