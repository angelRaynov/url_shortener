package url

type StoreFinder interface {
	Store(short, long string) error
	FindShort(long string) (string,error)
	FindLong(short string) (string,error)
	IsShortURLUnique(short string) (bool, error)
}

