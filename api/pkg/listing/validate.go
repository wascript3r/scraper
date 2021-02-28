package listing

type Validate interface {
	RawRequest(s interface{}) error
}
