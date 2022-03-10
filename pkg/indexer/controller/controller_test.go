package controller

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"testing"
	"time"

	warcindexer "github.com/crossedbot/go-warc-indexer"
	"github.com/crossedbot/simplecdxj"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	"github.com/crossedbot/hermes-archiver/pkg/indexer/models"
	"github.com/crossedbot/hermes-archiver/pkg/mocks"
)

func TestFindRecords(t *testing.T) {
	surt := "https://(com,example,:8443)/hello/world"
	types := []string{"response"}
	before := int64(1622551389)
	after := int64(1622550000)
	match := models.TextMatchExact
	limit := 1
	expected := models.Records{
		Count: 1,
		Results: []models.Record{{
			Id:        "abc123",
			Surt:      surt,
			Timestamp: time.Unix(before-1000, 0),
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
		}},
	}
	mockCtrl := gomock.NewController(t)
	defer mockCtrl.Finish()
	mockDb := mocks.NewMockCdxjRecords(mockCtrl)
	mockDb.EXPECT().
		Find(surt, types, before, after, match.String(), limit).
		Return(expected, nil)
	ctlr := &controller{db: mockDb}
	actual, err := ctlr.FindRecords(
		surt,
		[]simplecdxj.RecordType{simplecdxj.ResponseRecordType},
		strconv.FormatInt(before, 10),
		strconv.FormatInt(after, 10),
		match,
		limit,
	)
	require.Nil(t, err)
	require.Equal(t, expected, actual)
}

func TestGetRecord(t *testing.T) {
	expected := models.Record{
		Id:        "abc123",
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
		Get(expected.Id).
		Return(expected, nil)
	ctlr := &controller{db: mockDb}
	actual, err := ctlr.GetRecord(expected.Id)
	require.Nil(t, err)
	require.Equal(t, expected, actual)
}
