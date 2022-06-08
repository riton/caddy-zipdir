package fileslist

type Processor interface {
	GetFiles() ([]string, error)
}
