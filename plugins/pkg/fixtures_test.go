package pkg

var testGithubResponse = `[{
	"name": "v0.2.0",
	"zipball_url": "https://api.github.com/repos/ketchuphq/ketchup/zipball/v0.2.0",
	"tarball_url": "https://api.github.com/repos/ketchuphq/ketchup/tarball/v0.2.0",
	"commit": {
		"sha": "dbf737990b33a980103f9a723d254502a8686886",
		"url":
			"https://api.github.com/repos/ketchuphq/ketchup/commits/dbf737990b33a980103f9a723d254502a8686886"
	}
}, {
	"name": "v0.1.0",
	"zipball_url": "https://api.github.com/repos/ketchuphq/ketchup/zipball/v0.1.0",
	"tarball_url": "https://api.github.com/repos/ketchuphq/ketchup/tarball/v0.1.0",
	"commit": {
		"sha": "ca445d04f90ceb9a23267ece6f5c7050c2a85794",
		"url":
			"https://api.github.com/repos/ketchuphq/ketchup/commits/ca445d04f90ceb9a23267ece6f5c7050c2a85794"
	}
}]`

var testBitbucketResponse = `{
	"pagelen": 10,
	"size": 145,
	"values": [
		{
			"name": "tip",
			"links": {
				"commits": {
					"href": "https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commits/tip"
				},
				"self": {"href": "https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/refs/tags/tip"},
				"html": {"href": "https://bitbucket.org/tortoisehg/thg/commits/tag/tip"}
			},
			"tagger": null,
			"date": null,
			"message": null,
			"type": "tag",
			"target": {
				"hash": "83b7e2f14fcec35a739a99a0589e56a4cffa4f6d",
				"repository": {
					"links": {
						"self": {"href": "https://api.bitbucket.org/2.0/repositories/tortoisehg/thg"},
						"html": {"href": "https://bitbucket.org/tortoisehg/thg"},
						"avatar": {"href": "https://bitbucket.org/tortoisehg/thg/avatar/32/"}
					},
					"type": "repository",
					"name": "thg",
					"full_name": "tortoisehg/thg",
					"uuid": "{06bf263a-cea8-423d-8779-78bb87094731}"
				},
				"links": {
					"self": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commit/83b7e2f14fcec35a739a99a0589e56a4cffa4f6d"
					},
					"comments": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commit/83b7e2f14fcec35a739a99a0589e56a4cffa4f6d/comments"
					},
					"patch": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/patch/83b7e2f14fcec35a739a99a0589e56a4cffa4f6d"
					},
					"html": {
						"href":
							"https://bitbucket.org/tortoisehg/thg/commits/83b7e2f14fcec35a739a99a0589e56a4cffa4f6d"
					},
					"diff": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/diff/83b7e2f14fcec35a739a99a0589e56a4cffa4f6d"
					},
					"approve": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commit/83b7e2f14fcec35a739a99a0589e56a4cffa4f6d/approve"
					},
					"statuses": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commit/83b7e2f14fcec35a739a99a0589e56a4cffa4f6d/statuses"
					}
				},
				"author": {
					"raw": "Wagner Bruna <wbruna@softwareexpress.com.br>",
					"type": "author",
					"user": {
						"username": "wbruna",
						"display_name": "Wagner Bruna",
						"type": "user",
						"uuid": "{40716e33-3321-46a5-9f84-f57346e2d8cb}",
						"links": {
							"self": {"href": "https://api.bitbucket.org/2.0/users/wbruna"},
							"html": {"href": "https://bitbucket.org/wbruna/"},
							"avatar": {"href": "https://bitbucket.org/account/wbruna/avatar/32/"}
						}
					}
				},
				"parents": [
					{
						"hash": "c57d99f0b5707ddd6dafcd957c20b9986f7653bc",
						"type": "commit",
						"links": {
							"self": {
								"href":
									"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commit/c57d99f0b5707ddd6dafcd957c20b9986f7653bc"
							},
							"html": {
								"href":
									"https://bitbucket.org/tortoisehg/thg/commits/c57d99f0b5707ddd6dafcd957c20b9986f7653bc"
							}
						}
					}
				],
				"date": "2018-05-10T14:07:07+00:00",
				"message": "i18n: pull latest fr translations from Launchpad",
				"type": "commit"
			}
		},
		{
			"name": "4.5.3",
			"links": {
				"commits": {
					"href": "https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commits/4.5.3"
				},
				"self": {
					"href": "https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/refs/tags/4.5.3"
				},
				"html": {"href": "https://bitbucket.org/tortoisehg/thg/commits/tag/4.5.3"}
			},
			"tagger": {
				"raw": "Steve Borho <steve@borho.org>",
				"type": "author",
				"user": {
					"username": "sborho",
					"display_name": "Steve Borho",
					"type": "user",
					"uuid": "{a7b33259-0b76-44e3-9ee3-0946c696a082}",
					"links": {
						"self": {"href": "https://api.bitbucket.org/2.0/users/sborho"},
						"html": {"href": "https://bitbucket.org/sborho/"},
						"avatar": {"href": "https://bitbucket.org/account/sborho/avatar/32/"}
					}
				}
			},
			"date": "2018-04-07T20:07:58+00:00",
			"message": "Added tag 4.5.3 for changeset 07157a0943be",
			"type": "tag",
			"target": {
				"hash": "07157a0943be2f031a3cbe2d7a1e0d875ddf7731",
				"repository": {
					"links": {
						"self": {"href": "https://api.bitbucket.org/2.0/repositories/tortoisehg/thg"},
						"html": {"href": "https://bitbucket.org/tortoisehg/thg"},
						"avatar": {"href": "https://bitbucket.org/tortoisehg/thg/avatar/32/"}
					},
					"type": "repository",
					"name": "thg",
					"full_name": "tortoisehg/thg",
					"uuid": "{06bf263a-cea8-423d-8779-78bb87094731}"
				},
				"links": {
					"self": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commit/07157a0943be2f031a3cbe2d7a1e0d875ddf7731"
					},
					"comments": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commit/07157a0943be2f031a3cbe2d7a1e0d875ddf7731/comments"
					},
					"patch": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/patch/07157a0943be2f031a3cbe2d7a1e0d875ddf7731"
					},
					"html": {
						"href":
							"https://bitbucket.org/tortoisehg/thg/commits/07157a0943be2f031a3cbe2d7a1e0d875ddf7731"
					},
					"diff": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/diff/07157a0943be2f031a3cbe2d7a1e0d875ddf7731"
					},
					"approve": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commit/07157a0943be2f031a3cbe2d7a1e0d875ddf7731/approve"
					},
					"statuses": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commit/07157a0943be2f031a3cbe2d7a1e0d875ddf7731/statuses"
					}
				},
				"author": {
					"raw": "Yuya Nishihara <yuya@tcha.org>",
					"type": "author",
					"user": {
						"username": "yuja",
						"display_name": "Yuya Nishihara",
						"type": "user",
						"uuid": "{75137504-2657-4a1e-bcf3-13db1660e49e}",
						"links": {
							"self": {"href": "https://api.bitbucket.org/2.0/users/yuja"},
							"html": {"href": "https://bitbucket.org/yuja/"},
							"avatar": {"href": "https://bitbucket.org/account/yuja/avatar/32/"}
						}
					}
				},
				"parents": [
					{
						"hash": "33bf6c53ad26589b62fe9fc6242adb9cb316cae4",
						"type": "commit",
						"links": {
							"self": {
								"href":
									"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commit/33bf6c53ad26589b62fe9fc6242adb9cb316cae4"
							},
							"html": {
								"href":
									"https://bitbucket.org/tortoisehg/thg/commits/33bf6c53ad26589b62fe9fc6242adb9cb316cae4"
							}
						}
					}
				],
				"date": "2018-03-28T15:01:15+00:00",
				"message":
					"graphopt: do not build nodes just for flags() (fixes #5061)\n\nQTreeView scans items up to rowCount() at the initial layout. That was okay\non Qt4, but on Qt5, model.flags(index) is invoked for each row because Qt5\nhas an \"optimized\" path to test if a row has children, which information is\ncarried by flags().\n\nThis patch adds a light-weight replacement for graph[row], so the graphopt\ndoes not have to build a complete node object until it is requested through\nmodel.data(index).\n\nFixing graphopt.Graph.__len__() didn't go well since graphopt has to report\na pseudo length so the scrolling of the model can be \"optimized.\" That's\nthe fundamental design of the graphopt.",
				"type": "commit"
			}
		},
		{
			"name": "4.5.2",
			"links": {
				"commits": {
					"href": "https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commits/4.5.2"
				},
				"self": {
					"href": "https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/refs/tags/4.5.2"
				},
				"html": {"href": "https://bitbucket.org/tortoisehg/thg/commits/tag/4.5.2"}
			},
			"tagger": {
				"raw": "Steve Borho <steve@borho.org>",
				"type": "author",
				"user": {
					"username": "sborho",
					"display_name": "Steve Borho",
					"type": "user",
					"uuid": "{a7b33259-0b76-44e3-9ee3-0946c696a082}",
					"links": {
						"self": {"href": "https://api.bitbucket.org/2.0/users/sborho"},
						"html": {"href": "https://bitbucket.org/sborho/"},
						"avatar": {"href": "https://bitbucket.org/account/sborho/avatar/32/"}
					}
				}
			},
			"date": "2018-03-10T23:57:10+00:00",
			"message": "Added tag 4.5.2 for changeset 4bbca812fbe6",
			"type": "tag",
			"target": {
				"hash": "4bbca812fbe6fe5abab4da1201cb25d9a0be59ae",
				"repository": {
					"links": {
						"self": {"href": "https://api.bitbucket.org/2.0/repositories/tortoisehg/thg"},
						"html": {"href": "https://bitbucket.org/tortoisehg/thg"},
						"avatar": {"href": "https://bitbucket.org/tortoisehg/thg/avatar/32/"}
					},
					"type": "repository",
					"name": "thg",
					"full_name": "tortoisehg/thg",
					"uuid": "{06bf263a-cea8-423d-8779-78bb87094731}"
				},
				"links": {
					"self": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commit/4bbca812fbe6fe5abab4da1201cb25d9a0be59ae"
					},
					"comments": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commit/4bbca812fbe6fe5abab4da1201cb25d9a0be59ae/comments"
					},
					"patch": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/patch/4bbca812fbe6fe5abab4da1201cb25d9a0be59ae"
					},
					"html": {
						"href":
							"https://bitbucket.org/tortoisehg/thg/commits/4bbca812fbe6fe5abab4da1201cb25d9a0be59ae"
					},
					"diff": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/diff/4bbca812fbe6fe5abab4da1201cb25d9a0be59ae"
					},
					"approve": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commit/4bbca812fbe6fe5abab4da1201cb25d9a0be59ae/approve"
					},
					"statuses": {
						"href":
							"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commit/4bbca812fbe6fe5abab4da1201cb25d9a0be59ae/statuses"
					}
				},
				"author": {
					"raw": "Steve Borho <steve@borho.org>",
					"type": "author",
					"user": {
						"username": "sborho",
						"display_name": "Steve Borho",
						"type": "user",
						"uuid": "{a7b33259-0b76-44e3-9ee3-0946c696a082}",
						"links": {
							"self": {"href": "https://api.bitbucket.org/2.0/users/sborho"},
							"html": {"href": "https://bitbucket.org/sborho/"},
							"avatar": {"href": "https://bitbucket.org/account/sborho/avatar/32/"}
						}
					}
				},
				"parents": [
					{
						"hash": "aa367f7f486be80d60d699c5d174650533c4621b",
						"type": "commit",
						"links": {
							"self": {
								"href":
									"https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/commit/aa367f7f486be80d60d699c5d174650533c4621b"
							},
							"html": {
								"href":
									"https://bitbucket.org/tortoisehg/thg/commits/aa367f7f486be80d60d699c5d174650533c4621b"
							}
						}
					}
				],
				"date": "2018-03-10T23:11:38+00:00",
				"message": "pywin32 + Windows 10 = a whole new meaning of obnoxious",
				"type": "commit"
			}
		}
	],
	"page": 1,
	"next": "https://api.bitbucket.org/2.0/repositories/tortoisehg/thg/refs/tags?sort=-name&page=2"
}
`
