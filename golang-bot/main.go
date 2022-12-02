package main

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"io/ioutil"
	"strconv"
)

var (
	Token     string
	BotPrefix string

	config *configStruct
)

// Ban reason -> duration
//
//	const CRIME_AND_PUNISHMENT struct = []time.Time{
//		ADVERTISMENT
//		FLOOD
//	 Behavior?
//
// }
const (
	Error1        = "Error while connecting starting bot."
	Error2        = "Error while opening bot."
	Error3        = "Invalid data (should be in range(7, 365)."
	Error4        = "Error while banning user."
	Error5        = "You have no permissions to call this function."
	Error6        = "Error while filtering message."
	Error7        = "No channel with this ID"
	Error8        = "Error while checking guild. Probably wrong ID"
	Error9        = "Converting error?"
	DefaultReport = "Something went wrong. Working on it."
)

type configStruct struct {
	Token     string `json: "Token"`
	BotPrefix string `json: "BotPrefix"`
}

// discordgo.Session -> will handle all the interactions with the Discord servers

func ReadConfig() error {
	file, err := ioutil.ReadFile("./config.json")

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	fmt.Println(string(file))

	err = json.Unmarshal(file, &config)

	if err != nil {
		fmt.Println(err.Error())
		return err
	}

	Token = config.Token
	BotPrefix = config.BotPrefix

	return nil
}

var BotId string
var goBot *discordgo.Session

func Start() {
	// discord, err := discordgo.New("Bot " + "authentication token")
	goBot, err := discordgo.New("Bot " + config.Token)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	u, err := goBot.User("@me")

	if err != nil {
		fmt.Println("Something went wrong. Check start function.")
		return
	}

	BotId = u.ID

	// FINISHED HERE
	// Mapping For Administration

	goBot.AddHandler(info)
	goBot.AddHandler(banHandler)
	goBot.AddHandler(messageFilterHandler)
	goBot.AddHandler(pingPongHandler)
	goBot.AddHandler(banHandler)

	err = goBot.Open()

	if err != nil {
		fmt.Print(err)
	}

	fmt.Println("Bot is running fine!")

}

func pingPongHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == BotId {
		return
	}
	// m -> message
	if m.Content == BotPrefix+"ping" {
		_, _ = s.ChannelMessageSend(m.ChannelID, "pong")
	}
}

func main() {
	err := ReadConfig()
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	Start()

	<-make(chan struct{})
	return
}

// Starting go discord bot
func banHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	//discordgo.GuildMemberParams{}
	// GuildRoleCreate
	// guild -> member -> permissions
	if m.Content == BotPrefix+"ban" {
		Author := m.Author.ID
		parsedAuthorID, err := strconv.Atoi(Author)
		if err != nil {
			sendReportToOwner(s, m, Error9)
		}
		GuildID := m.ChannelID
		//channel, err := s.Channel(GuildID)

		//if err != nil {
		//	sendReportToOwner(s, m, Error7)
		//}

		//userId, err := strconv.Atoi(Author)
		//
		//if err != nil {
		//	sendReportToOwner(s, m, DefaultReport)
		//}

		//UserPermission := channel.Members[userId].UserID
		GuildStruct, err := s.Guild(GuildID)
		if err != nil {
			sendReportToOwner(s, m, Error8)
		}
		UserPermission := GuildStruct.Members[parsedAuthorID].Permissions
		//s := discordgo.GuildApplicationCommandPermissions[GuildID].Permissions
		if UserPermission == discordgo.PermissionAdministrator {
			if m.Type == 19 {
				mes, err := s.ChannelMessageSend(m.ChannelID, "How long this user need to be banned?")
				parsedInt, _ := strconv.Atoi(mes.Content)
				if parsedInt < 7 || parsedInt > 365 {
					sendReportToOwner(s, m, Error3)
				}

				if err != nil {
					sendReportToOwner(s, m, DefaultReport)
				}

				errr := s.GuildBanCreate(m.ChannelID, m.Author.ID, parsedInt)
				if errr != nil {
					sendReportToOwner(s, m, Error4)
				}
				// m.ReferencedMessage.Member

				if err != nil {
					// Reporting error
					sendReportToOwner(s, m, DefaultReport)
				}
			} else {
				s.ChannelMessageSend(m.ChannelID, "You should reply to message of nearly banned user. ")
			}
		} else {
			sendReportToOwner(s, m, Error5)
		}
	}
}

// Drink vodka is reference to some bad words
func messageFilterHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Content == "let's go drink vodka" || m.Content == "let's drink some vodka" {
		_, err := s.ChannelMessageSendReply(m.ChannelID, "lezz goo", nil)
		if err != nil {
			sendReportToOwner(s, m, Error6)
		}
	}
}

func info(s *discordgo.Session, m *discordgo.MessageCreate) *discordgo.Message {
	if m.Content == BotPrefix+"info" {
		mes, err := s.ChannelMessageSend(m.ChannelID, "Info about this bot: I love Slipknot!!!")
		if err != nil {
			s.ChannelMessageSend(m.ChannelID, DefaultReport)
		}
		return mes
	}
	return
}

func sendReportToOwner(s *discordgo.Session, m *discordgo.MessageCreate, trouble string) string {
	// Errors sending to me directly
	_, err := s.ChannelMessageSend("1047801998210244719", trouble)
	if err != nil {
		fmt.Println("Some troubles with bot...")
	}
	return "Sended."
}
