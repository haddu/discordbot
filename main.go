package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

//Stores ids of all temporary created rooms

type createdVoiseChannel struct {
	ID string
}

const token string = "NjExOTEwOTM3NDUxOTU0MTk2.XVas4g.mOOx7hUvwivq510BYfGIeKVw6xo"

func main() {

	createdVoiceChannels := make([]createdVoiseChannel, 0, 20)

	//Create new Discord session using the provided bot token.
	bot, err := discordgo.New("Bot " + token)

	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	//Register the messageCreate func as a callback for MessageCreate events.
	bot.AddHandler(messageCreate)

	//Checks if anyone joined defined voice channel
	bot.AddHandler(func(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
		voiceChannelJoinChecker(s, m, createdVoiceChannels)
	})

	bot.AddHandler(func(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
		voiceChannelLeftChecker(s, m, createdVoiceChannels)
	})

	//Open a websocket connection to Discord and begin listening
	err = bot.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	//Wait here until CTRL-C or other term signal is received
	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	//Cleanly close down the Discord session.
	bot.Close()

}

func printAllMembers(createdVoiceChannels []createdVoiseChannel) {
	fmt.Println("---------------------------------")
	fmt.Println("Created voice channels: ")
	for _, v := range createdVoiceChannels {
		fmt.Print(v.ID + ",")
	}
	fmt.Println("")
}

func voiceChannelLeftChecker(s *discordgo.Session, m *discordgo.VoiceStateUpdate, createdVoiceChannels []createdVoiseChannel) {
	if m.ChannelID == "" {
		tempChannels, er1 := s.GuildChannels(m.GuildID)

		if er1 != nil {
			fmt.Println("error getting all channels,", er1)
			return
		}

		for _, v := range tempChannels {

			if (v.ParentID == "611864915262832640") && (v.ID != "611983405680295978") {

				fmt.Println("Found empty room")

				if len(v.Recipients) == 0 {

					for i := len(createdVoiceChannels) - 1; i >= 0; i-- {
						if createdVoiceChannels[i].ID == v.ID {
							createdVoiceChannels = append(createdVoiceChannels[:i],
								createdVoiceChannels[i+1:]...)
						}
					}

					printAllMembers(createdVoiceChannels)

					s.ChannelDelete(v.ID)
				}
			}
		}

	}

}

func voiceChannelJoinChecker(s *discordgo.Session, m *discordgo.VoiceStateUpdate, createdVoiceChannels []createdVoiseChannel) {
	if m.ChannelID == "611983405680295978" {

		tempUser, er1 := s.User(m.UserID)

		if er1 != nil {
			fmt.Println("error getting userID,", er1)
			return
		}

		fmt.Println("Connected")

		tempChannel, err := s.GuildChannelCreateComplex(m.GuildID, discordgo.GuildChannelCreateData{
			Name: tempUser.Username,
			Type: 2,
			PermissionOverwrites: []*discordgo.PermissionOverwrite{
				{
					ID:    tempUser.ID,
					Type:  "memeber",
					Allow: 16,
				},
			},
			ParentID: "611864915262832640",
		})

		if err != nil {
			fmt.Println("error creating new temporal channel,", err)
			return
		}

		s.GuildMemberMove(m.GuildID, m.UserID, tempChannel.ID)

		createdVoiceChannels = append(createdVoiceChannels, createdVoiseChannel{
			ID: tempChannel.ID,
		})

		printAllMembers(createdVoiceChannels)
	}

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.HasPrefix(m.Content, ".") {

		fmt.Println(m.Content)

		//Ignore all messages created by the bot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == ".ping" {
			s.ChannelMessageSend(m.ChannelID, "pong!")
		}
	}
}
