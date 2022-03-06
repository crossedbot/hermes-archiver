package indexer

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"testing"
	"time"

	warcindexer "github.com/crossedbot/go-warc-indexer"
	"github.com/crossedbot/simplecdxj"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/crossedbot/hermes-archiver/pkg/indexer/models"
	"github.com/crossedbot/hermes-archiver/pkg/mocks"
)

func TestStore(t *testing.T) {
	expectedId := "123e4567-e89b-12d3-a456-426614174000"
	record := models.Record{
		Surt:      "https://(com,example,:8443)/hello/world",
		Timestamp: time.Unix(int64(1622551389), 0),
		Type:      simplecdxj.ResponseRecordType,
		Content: warcindexer.JsonBlock{
			Uri:              "https://example.com",
			Ref:              fmt.Sprintf("warcfile:%s#%d", "world.warc", 123),
			Sha:              "xOsha256Ox",
			Hsc:              200,
			Mct:              "text/html",
			Rid:              "<urn:uuid:B0B3862C-B271-4670-A4B5-B127576C6118>",
			Locator:          "ipfs/abc-123/def-456",
			Title:            "Hello World",
			EncryptionKeyID:  base64.URLEncoding.EncodeToString([]byte("QqqxKStb")),
			EncryptionMethod: "aes-gcm",
			EncryptionNonce:  base64.URLEncoding.EncodeToString([]byte("qYiUZUNB")),
		},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDb := mocks.NewMockCdxjRecords(mockCtrl)
	mockDb.EXPECT().
		Set(record).
		Return(expectedId, nil)
	b, err := json.Marshal(record.Content)
	require.Nil(t, err)
	ids, err := store(mockDb, simplecdxj.CDXJ{
		Records: []*simplecdxj.Record{{
			SURT:      record.Surt,
			Timestamp: record.Timestamp,
			Type:      record.Type,
			Content:   b,
		}},
	})
	require.Nil(t, err)
	require.Len(t, ids, 1)
	require.Equal(t, expectedId, ids[0])
}
