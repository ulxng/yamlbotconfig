package messages

type Navigator interface {
	GetNextByText(text string) (Message, error)
}
