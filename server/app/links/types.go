package links

type Scannable interface {
	Scan(...any) error
}
