package controller

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"

	mocksrpl "github.com/crossedbot/hermes-archiver/pkg/replayer/mocks"
	"github.com/crossedbot/hermes-archiver/pkg/replayer/models"
)

func TestReplay(t *testing.T) {
	id := "abc123"
	key := []byte("supersecret")
	expected := models.Replay{}
	mockCtlr := gomock.NewController(t)
	defer mockCtlr.Finish()
	mockRpl := mocksrpl.NewMockReplayer(mockCtlr)
	mockRpl.EXPECT().
		Replay(id, key).
		Return(expected, nil)
	ctlr := &controller{key: key, rpl: mockRpl}
	actual, err := ctlr.Replay(id)
	require.Nil(t, err)
	require.Equal(t, expected, actual)
}
