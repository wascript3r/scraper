package hasher

import (
	"fmt"

	"github.com/wascript3r/gocipher/encoder"
	"github.com/wascript3r/gocipher/sha256"
	"github.com/wascript3r/scraper/api/pkg/listing"
)

type Hasher struct{}

func New() *Hasher {
	return &Hasher{}
}

func (h *Hasher) HashSoldRecord(sr *listing.SoldRecord) string {
	str := fmt.Sprintf("%s-%s-%s-%s", sr.User, sr.Price, sr.Quantity, sr.Date)
	bs := encoder.HexEncode(sha256.Compute([]byte(str)))
	return string(bs)
}
