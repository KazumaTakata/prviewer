package tui

import "charm.land/bubbles/v2/key"

const SelectFormKey = "selectHistory"

var registerNewRepositoryKey = key.NewBinding(key.WithKeys("n"), key.WithHelp("n", "新規リポジトリー登録"))
