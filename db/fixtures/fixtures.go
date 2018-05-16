package fixtures

import (
	"github.com/golang/protobuf/proto"
	"github.com/ketchuphq/ketchup/proto/ketchup/models"
)

var Pages = map[string]*models.Page{
	"1169b0f5-2006-4bda-a9e7-41b86f0cc129": &models.Page{
		Uuid:     proto.String("1169b0f5-2006-4bda-a9e7-41b86f0cc129"),
		Title:    proto.String("Jane Eyre"),
		Theme:    proto.String("none"),
		Template: proto.String("html"),
		Timestamps: &models.Timestamp{
			CreatedAt: proto.Int64(1486456483427),
			UpdatedAt: proto.Int64(1491801538232),
		},
		Contents: []*models.Content{{
			Uuid:  proto.String("affd5711-1bdd-4a35-839b-d6ba89bde435"),
			Key:   proto.String("content"),
			Value: proto.String("There was no possibility of taking a walk that day.  We had been wandering, indeed, in the leafless shrubbery an hour in the morning; but since dinner (Mrs. Reed, when there was no company, dined early) the cold winter wind had brought with it clouds so sombre, and a rain so penetrating, that further out-door exercise was now out of the question."),
			Type: &models.Content_Text{
				Text: &models.ContentText{Type: models.ContentTextType_html.Enum()},
			},
		}},
	},
	"19725193-8c5d-4b41-b270-4a8dd726ab77": &models.Page{
		Uuid:     proto.String("19725193-8c5d-4b41-b270-4a8dd726ab77"),
		Title:    proto.String("Count of Monte Cristo"),
		Theme:    proto.String("none"),
		Template: proto.String("markdown"),
		Timestamps: &models.Timestamp{
			CreatedAt: proto.Int64(1486456483427),
			UpdatedAt: proto.Int64(1491801538233),
		},
		Contents: []*models.Content{{
			Uuid:  proto.String("832a424b-7f5c-405a-98b4-1d02f236855c"),
			Key:   proto.String("content"),
			Value: proto.String("On the 24th of February, 1815, the look-out at Notre-Dame de la Garde signalled the three-master, the Pharaon from Smyrna, Trieste, and Naples.\n\nAs usual, a pilot put off immediately, and rounding the Château d’If, got on board the vessel between Cape Morgiou and Rion island.\n\nImmediately, and according to custom, the ramparts of Fort Saint-Jean were covered with spectators; it is always an event at Marseilles for a ship to come into port, especially when this ship, like the Pharaon, has been built, rigged, and laden at the old Phocee docks, and belongs to an owner of the city."),
			Type: &models.Content_Text{
				Text: &models.ContentText{Type: models.ContentTextType_markdown.Enum()},
			},
		}},
	},
	"1ae41c2f-317c-48a2-b910-f7cd231bfa13": &models.Page{
		Uuid:     proto.String("1ae41c2f-317c-48a2-b910-f7cd231bfa13"),
		Title:    proto.String("Pride and Prejudice"),
		Theme:    proto.String("ebook"),
		Template: proto.String("gutenberg.html"),
		Timestamps: &models.Timestamp{
			CreatedAt: proto.Int64(1486456483427),
			UpdatedAt: proto.Int64(1491801538233),
		},
		Contents: []*models.Content{{Uuid: proto.String("3ab6fe96-7f1f-493c-a7bc-f26b949b64cf"),
			Key:   proto.String("content"),
			Value: proto.String("It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife.\n\nHowever little known the feelings or views of such a man may be on his first entering a neighbourhood, this truth is so well fixed in the minds of the surrounding families, that he is considered the rightful property of some one or other of their daughters."),
			Type: &models.Content_Text{
				Text: &models.ContentText{Type: models.ContentTextType_markdown.Enum()},
			},
		},
		},
	},
}
var Routes = map[string]*models.Route{
	"0088ac36-463e-41c0-bf1f-d0a94750ee2c": {
		Uuid: proto.String("0088ac36-463e-41c0-bf1f-d0a94750ee2c"),
		Path: proto.String("/jane-eyre"),
		Target: &models.Route_PageUuid{
			PageUuid: "1169b0f5-2006-4bda-a9e7-41b86f0cc129",
		},
	},
	"00c83178-e489-4c9d-ac8a-a6bbdeb844c9": {
		Uuid: proto.String("00c83178-e489-4c9d-ac8a-a6bbdeb844c9"),
		Path: proto.String("/count-of-monte-cristo"),
		Target: &models.Route_PageUuid{
			PageUuid: "19725193-8c5d-4b41-b270-4a8dd726ab77",
		},
	},
	"09aac79a-6747-4670-a35c-4125e41c6f81": {
		Uuid: proto.String("09aac79a-6747-4670-a35c-4125e41c6f81"),
		Path: proto.String("/pride-and-prejudice"),
		Target: &models.Route_PageUuid{
			PageUuid: "1ae41c2f-317c-48a2-b910-f7cd231bfa13",
		},
	}}

var JSON = `{
	"pages": [
		{
			"uuid": "19725193-8c5d-4b41-b270-4a8dd726ab77",
			"title": "Count of Monte Cristo",
			"theme": "none",
			"template": "markdown",
			"timestamps": {
				"createdAt": "1486456483427",
				"updatedAt": "1491801538233"
			},
			"contents": [
				{
					"uuid": "832a424b-7f5c-405a-98b4-1d02f236855c",
					"key": "content",
					"value": "On the 24th of February, 1815, the look-out at Notre-Dame de la Garde signalled the three-master, the Pharaon from Smyrna, Trieste, and Naples.\n\nAs usual, a pilot put off immediately, and rounding the Château d’If, got on board the vessel between Cape Morgiou and Rion island.\n\nImmediately, and according to custom, the ramparts of Fort Saint-Jean were covered with spectators; it is always an event at Marseilles for a ship to come into port, especially when this ship, like the Pharaon, has been built, rigged, and laden at the old Phocee docks, and belongs to an owner of the city.",
					"text": {
						"type": "markdown"
					}
				}
			]
		}, {
			"uuid": "1169b0f5-2006-4bda-a9e7-41b86f0cc129",
			"title": "Jane Eyre",
			"theme": "none",
			"template": "html",
			"timestamps": {
				"createdAt": "1486456483427",
				"updatedAt": "1491801538232"
			},
			"contents": [
				{
					"uuid": "affd5711-1bdd-4a35-839b-d6ba89bde435",
					"key": "content",
					"value": "There was no possibility of taking a walk that day.  We had been wandering, indeed, in the leafless shrubbery an hour in the morning; but since dinner (Mrs. Reed, when there was no company, dined early) the cold winter wind had brought with it clouds so sombre, and a rain so penetrating, that further out-door exercise was now out of the question.",
					"text": {
						"type": "html"
					}
				}
			]
		},
		{
			"uuid": "1ae41c2f-317c-48a2-b910-f7cd231bfa13",
			"title": "Pride and Prejudice",
			"theme": "ebook",
			"template": "gutenberg.html",
			"timestamps": {
				"createdAt": "1486456483427",
				"updatedAt": "1491801538233"
			},
			"contents": [
				{
					"uuid": "3ab6fe96-7f1f-493c-a7bc-f26b949b64cf",
					"key": "content",
					"value": "It is a truth universally acknowledged, that a single man in possession of a good fortune, must be in want of a wife.\n\nHowever little known the feelings or views of such a man may be on his first entering a neighbourhood, this truth is so well fixed in the minds of the surrounding families, that he is considered the rightful property of some one or other of their daughters.",
					"text": {
						"type": "markdown"
					}
				}
			]
		}
	],
	"routes": [
		{
			"uuid": "0088ac36-463e-41c0-bf1f-d0a94750ee2c",
			"path": "/jane-eyre",
			"pageUuid": "1169b0f5-2006-4bda-a9e7-41b86f0cc129"
		},
		{
			"uuid": "00c83178-e489-4c9d-ac8a-a6bbdeb844c9",
			"path": "/count-of-monte-cristo",
			"pageUuid": "19725193-8c5d-4b41-b270-4a8dd726ab77"
		},
		{
			"uuid": "09aac79a-6747-4670-a35c-4125e41c6f81",
			"path": "/pride-and-prejudice",
			"pageUuid": "1ae41c2f-317c-48a2-b910-f7cd231bfa13"
		}
	]
}`
