package pagination_test

import (
	"testing"
	"time"

	"github.com/iammrsea/social-app/internal/shared/pagination"
	"github.com/stretchr/testify/assert"
)

func TestEncode_Decode_Cursor(t *testing.T) {
	t.Parallel()

	id := time.Now().UTC().Format(time.RFC3339Nano)
	encodedId := pagination.EncodeCursor(id)
	decoded, err := pagination.DecodeCursor(encodedId)
	assert.Nil(t, err)
	assert.Equal(t, id, decoded)
}
