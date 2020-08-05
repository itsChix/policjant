package cah

import (
	"strings"

	"github.com/jonas747/cardsagainstdiscord"
	"github.com/jonas747/dcmd"
	"github.com/jonas747/yagpdb/bot"
	"github.com/jonas747/yagpdb/commands"
	"github.com/sirupsen/logrus"
)

func (p *Plugin) AddCommands() {

	cmdCreate := &commands.YAGCommand{
		Name:        "Create",
		CmdCategory: commands.CategoryFun,
		Aliases:     []string{"c"},
		Description: "Tworzy grę Karty Przeciwko Ludzkości (Cards Against Humanity) na obecnym kanale, aby dodać paczki napisz nazwę paczki np. r/cah c nazwapaczki, lub * aby grać z wszystkimi. (dodaj -v aby grać z głosowaniem.).",
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "packs", Type: dcmd.String, Default: "main", Help: "Paczki są odzielane przez spację, aby użyć wszystkich dopisz *."},
		},
		ArgSwitches: []*dcmd.ArgDef{
			{Switch: "v", Name: "Tryb głosowania - gracze głosują na najlepszą kombinację."},
		},
		RunFunc: func(data *dcmd.Data) (interface{}, error) {
			voteMode := data.Switch("v").Bool()
			pStr := data.Args[0].Str()
			packs := strings.Fields(pStr)

			_, err := p.Manager.CreateGame(data.GS.ID, data.CS.ID, data.Msg.Author.ID, data.Msg.Author.Username, voteMode, packs...)
			if err == nil {
				logrus.Info("[cah] Stworzono nową grę: ", data.CS.ID, ":", data.GS.ID)
				return "", nil
			}

			if cahErr := cardsagainstdiscord.HumanizeError(err); cahErr != "" {
				return cahErr, nil
			}

			return "", err
		},
	}

	cmdEnd := &commands.YAGCommand{
		Name:        "End",
		CmdCategory: commands.CategoryFun,
		Description: "Kończy grę CAH na obecnym kanale.",
		RunFunc: func(data *dcmd.Data) (interface{}, error) {
			isAdmin, err := bot.AdminOrPermMS(data.CS.ID, data.MS, 0)
			if err == nil && isAdmin {
				err = p.Manager.RemoveGame(data.CS.ID)
			} else {
				err = p.Manager.TryAdminRemoveGame(data.Msg.Author.ID)
			}

			if err != nil {
				if cahErr := cardsagainstdiscord.HumanizeError(err); cahErr != "" {
					return cahErr, nil
				}

				return "", err
			}

			return "Zatrzymano grę", nil
		},
	}

	cmdKick := &commands.YAGCommand{
		Name:         "Kick",
		CmdCategory:  commands.CategoryFun,
		RequiredArgs: 1,
		Arguments: []*dcmd.ArgDef{
			&dcmd.ArgDef{Name: "user", Type: dcmd.UserID},
		},
		Description: "Wyrzuca gracza z gry na obecnym kanale.",
		RunFunc: func(data *dcmd.Data) (interface{}, error) {
			userID := data.Args[0].Int64()
			err := p.Manager.AdminKickUser(data.Msg.Author.ID, userID)
			if err != nil {
				if cahErr := cardsagainstdiscord.HumanizeError(err); cahErr != "" {
					return cahErr, nil
				}

				return "", err
			}

			return "Użytkownik usunięty", nil
		},
	}

	cmdPacks := &commands.YAGCommand{
		Name:         "Packs",
		CmdCategory:  commands.CategoryFun,
		RequiredArgs: 0,
		Description:  "Lista paczek",
		RunFunc: func(data *dcmd.Data) (interface{}, error) {
			resp := "Dostępne paczki: \n\n"
			for _, v := range cardsagainstdiscord.Packs {
				resp += "`" + v.Name + "` - " + v.Description + "\n"
			}

			return resp, nil
		},
	}

	container := commands.CommandSystem.Root.Sub("cah")
	container.NotFound = commands.CommonContainerNotFoundHandler(container, "")

	container.AddCommand(cmdCreate, cmdCreate.GetTrigger())
	container.AddCommand(cmdEnd, cmdEnd.GetTrigger())
	container.AddCommand(cmdKick, cmdKick.GetTrigger())
	container.AddCommand(cmdPacks, cmdPacks.GetTrigger())
}
