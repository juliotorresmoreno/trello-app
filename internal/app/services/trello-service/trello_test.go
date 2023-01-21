package trello_service

import (
	"testing"

	"github.com/juliotorresmoreno/trello-app/configs"
	"github.com/stretchr/testify/require"
)

func TestGetBoardId(t *testing.T) {
	require := require.New(t)
	trelloConf := configs.GetConfig().Trello

	cli := NewTrelloService(trelloConf.BoardId)
	boardId, err := cli.GetBoardId()
	require.NoError(err)

	require.Equal(24, len(boardId))
}

func TestGetBoardIdE1(t *testing.T) {
	require := require.New(t)

	cli := NewTrelloService("-")
	boardId, err := cli.GetBoardId()
	require.Error(err)
	require.Empty(boardId)
}

func TestGetLists(t *testing.T) {
	require := require.New(t)

	trelloConf := configs.GetConfig().Trello

	cli := NewTrelloService(trelloConf.BoardId)
	lists, err := cli.GetLists()
	require.NoError(err)

	require.NotZero(len(lists))

	boardId, _ := cli.GetBoardId()

	for _, item := range lists {
		require.NotEmpty(item.Id)
		require.NotEmpty(item.Name)
		require.NotEmpty(item.Pos)
		require.Equal(item.IdBoard, boardId)
	}
}

func TestGetListsE1(t *testing.T) {
	require := require.New(t)

	cli := NewTrelloService("-")
	lists, err := cli.GetLists()
	require.Error(err)
	require.Zero(len(lists))
}

func TestGetLabels(t *testing.T) {
	require := require.New(t)

	trelloConf := configs.GetConfig().Trello

	cli := NewTrelloService(trelloConf.BoardId)
	labels, err := cli.GetLabels()
	require.NoError(err)

	require.NotZero(len(labels))

	for _, item := range labels {
		require.NotEmpty(item.Id)
		require.NotEmpty(item.Color)
		require.Equal(item.IdBoard, cli.boardId)
	}
}

func TestGetLabelsE1(t *testing.T) {
	require := require.New(t)

	cli := NewTrelloService("-")
	labels, err := cli.GetLabels()
	require.Error(err)
	require.Zero(len(labels))
}

func TestCreateLabel(t *testing.T) {
	require := require.New(t)

	keyLabel := "test_label"

	trelloConf := configs.GetConfig().Trello

	cli := NewTrelloService(trelloConf.BoardId)
	err := cli.CreateLabel(CreateCardLabel{
		Name:  keyLabel,
		Color: "pink",
	})
	require.NoError(err)

	labels, err := cli.GetLabels()
	require.NoError(err)
	require.Contains(labels, keyLabel)
}
func TestCreateLabelE1(t *testing.T) {
	require := require.New(t)

	keyLabel := "test_label"

	trelloConf := configs.GetConfig().Trello

	cli := NewTrelloService(trelloConf.BoardId)
	err := cli.CreateLabel(CreateCardLabel{
		Name:  keyLabel,
		Color: "test_color", // it is not allowed
	})
	require.Error(err)
}

func TestDeleteLabel(t *testing.T) {
	require := require.New(t)

	keyLabel := "test_label"
	trelloConf := configs.GetConfig().Trello

	cli := NewTrelloService(trelloConf.BoardId)

	labels, err := cli.GetLabels()
	require.NoError(err)
	require.Contains(labels, keyLabel)

	label := labels[keyLabel]

	err = cli.DeleteLabel(label.Id)
	require.NoError(err)

	labels, err = cli.GetLabels()
	require.NoError(err)
	require.NotContains(labels, keyLabel)
}

func TestDeleteLabelE1(t *testing.T) {
	require := require.New(t)

	keyLabel := "-"
	trelloConf := configs.GetConfig().Trello

	cli := NewTrelloService(trelloConf.BoardId)

	err := cli.DeleteLabel(keyLabel)
	require.Error(err)
}

func TestCreateCard(t *testing.T) {
	require := require.New(t)

	trelloConf := configs.GetConfig().Trello

	cli := NewTrelloService(trelloConf.BoardId)

	labels, err := cli.GetLabels()
	require.NoError(err)

	require.Contains(labels, "issue")
	issue := labels["issue"]
	require.Contains(labels, "test")
	test := labels["test"]

	_, err = cli.CreateCard(CreateCardScheme{
		Name:     "test_card",
		Desc:     "test card",
		IdLabels: []string{issue.Id, test.Id},
	})

	require.NoError(err)
}

func TestCreateCardE1(t *testing.T) {
	require := require.New(t)

	cli := NewTrelloService("-")

	_, err := cli.CreateCard(CreateCardScheme{
		Name:     "",
		Desc:     "",
		IdLabels: []string{},
	})

	require.Error(err)
}

func TestGetCard(t *testing.T) {
	require := require.New(t)

	trelloConf := configs.GetConfig().Trello

	cli := NewTrelloService(trelloConf.BoardId)

	boardId, _ := cli.GetBoardId()
	cards, err := cli.GetCards()
	require.NoError(err)

	for _, card := range cards {
		require.NotEmpty(card.Name)
		require.NotEmpty(card.ShortLink)
		require.NotEmpty(card.ShortURL)
		require.NotEmpty(card.URL)

		require.Equal(card.IDBoard, boardId)
	}
}

func TestGetCardE1(t *testing.T) {
	require := require.New(t)

	cli := NewTrelloService("-")

	cards, err := cli.GetCards()
	require.Error(err)
	require.Zero(len(cards))
}

func TestDeleteCard(t *testing.T) {
	require := require.New(t)

	trelloConf := configs.GetConfig().Trello

	cli := NewTrelloService(trelloConf.BoardId)

	cards, err := cli.GetCards()
	require.NoError(err)

	for _, card := range cards {
		labels := card.Labels
		for _, label := range labels {
			if label.Name == "test" {
				err = cli.DeleteCard(card.ID)
				require.NoError(err)
				break
			}
		}
	}
}
func TestDeleteCardE1(t *testing.T) {
	require := require.New(t)

	trelloConf := configs.GetConfig().Trello

	cli := NewTrelloService(trelloConf.BoardId)

	err := cli.DeleteCard("-")
	require.Error(err)
}
