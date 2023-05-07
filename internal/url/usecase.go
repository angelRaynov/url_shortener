package url

type ShortExpander interface {
	Shorten(long string) string
	Expand(short string) (string,error)
}
