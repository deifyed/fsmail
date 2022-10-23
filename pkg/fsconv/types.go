package fsconv

type Message struct {
	Recipient string
	Cc        []string
	Subject   string
	Body      string
}
