package query

type Validate interface {
	RawRequest(s interface{}) error
}
