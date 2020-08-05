package moderation

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jonas747/discordgo"
	"github.com/jonas747/yagpdb/common"
)

type ModlogAction struct {
	Prefix string
	Emoji  string
	Color  int

	Footer string
}

func (m ModlogAction) String() string {
	str := m.Emoji + m.Prefix
	if m.Footer != "" {
		str += " (" + m.Footer + ")"
	}

	return str
}

var (
	MAMute       = ModlogAction{Prefix: "Wyciszono", Emoji: "🔇", Color: 0x57728e}
	MAUnmute     = ModlogAction{Prefix: "Usunięto wyciszenie", Emoji: "🔊", Color: 0x62c65f}
	MAKick       = ModlogAction{Prefix: "Wyrzucono", Emoji: "👢", Color: 0xf2a013}
	MABanned     = ModlogAction{Prefix: "Zbanowano", Emoji: "🔨", Color: 0xd64848}
	MAUnbanned   = ModlogAction{Prefix: "Odbanowano", Emoji: "🔓", Color: 0x62c65f}
	MAWarned     = ModlogAction{Prefix: "Ostrzeżono", Emoji: "⚠", Color: 0xfca253}
	MAGiveRole   = ModlogAction{Prefix: "", Emoji: "➕", Color: 0x53fcf9}
	MARemoveRole = ModlogAction{Prefix: "", Emoji: "➖", Color: 0x53fcf9}
)

func CreateModlogEmbed(config *Config, author *discordgo.User, action ModlogAction, target *discordgo.User, reason, logLink string) error {
	channelID := config.IntActionChannel()
	config.GetGuildID()
	if channelID == 0 {
		return nil
	}

	emptyAuthor := false
	if author == nil {
		emptyAuthor = true
		author = &discordgo.User{
			ID:            0,
			Username:      "Unknown",
			Discriminator: "????",
		}
	}

	if reason == "" {
		reason = "(powód nie został podany)"
	}

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%s#%s (ID %d)", author.Username, author.Discriminator, author.ID),
			IconURL: discordgo.EndpointUserAvatar(author.ID, author.Avatar),
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: discordgo.EndpointUserAvatar(target.ID, target.Avatar),
		},
		Color: action.Color,
		Description: fmt.Sprintf("**%s%s %s**#%s *(ID %d)*\n📄**Powód:** %s",
			action.Emoji, action.Prefix, target.Username, target.Discriminator, target.ID, reason),
	}

	if logLink != "" {
		embed.Description += " ([Logi](" + logLink + "))"
	}

	if action.Footer != "" {
		embed.Footer = &discordgo.MessageEmbedFooter{
			Text: action.Footer,
		}
	}

	m, err := common.BotSession.ChannelMessageSendEmbed(channelID, embed)
	if err != nil {
		if common.IsDiscordErr(err, discordgo.ErrCodeMissingAccess, discordgo.ErrCodeMissingPermissions, discordgo.ErrCodeUnknownChannel) {
			// disable the modlog
			config.ActionChannel = ""
			config.Save(config.GetGuildID())
			return nil
		}
		return err
	}

	if emptyAuthor {
		placeholder := fmt.Sprintf("Przypisz autora i powód używając komendy **'reason %d twoj-powod`**", m.ID)
		updateEmbedReason(nil, placeholder, embed)
		_, err = common.BotSession.ChannelMessageEditEmbed(channelID, m.ID, embed)
	}
	return err
}

var (
	logsRegex = regexp.MustCompile(`\(\[Logi\]\(.*\)\)`)
)

func updateEmbedReason(author *discordgo.User, reason string, embed *discordgo.MessageEmbed) {
	const checkStr = "📄**Powód:**"

	index := strings.Index(embed.Description, checkStr)
	withoutReason := embed.Description[:index+len(checkStr)]

	logsLink := logsRegex.FindString(embed.Description)
	if logsLink != "" {
		logsLink = " " + logsLink
	}

	embed.Description = withoutReason + " " + reason + logsLink

	if author != nil {
		embed.Author = &discordgo.MessageEmbedAuthor{
			Name:    fmt.Sprintf("%s#%s (ID %d)", author.Username, author.Discriminator, author.ID),
			IconURL: discordgo.EndpointUserAvatar(author.ID, author.Avatar),
		}
	}
}
