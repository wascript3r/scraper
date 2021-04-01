package listing

type Hasher interface {
	HashSoldRecord(sr *SoldRecord) string
}
