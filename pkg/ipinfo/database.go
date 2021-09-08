package ipinfo

type Database interface {
	Area(string) string
}
