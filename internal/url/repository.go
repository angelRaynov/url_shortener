package url

type StoreFinder interface {
	Store(short, long string) error
	FindShort(long string) (string, error)
	FindLong(short string) (string, error)
	IsShortURLUnique(short string) (bool, error)
}

type GetCacher interface {
	GetShort(long string) (string, error)
	GetLong(short string) (string, error)
	Cache(short, long string) error
}
