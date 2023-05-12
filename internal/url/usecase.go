package url

type ShortExpander interface {
	Shorten(long string) (string, error)
	Expand(short string) (string, error)
}
