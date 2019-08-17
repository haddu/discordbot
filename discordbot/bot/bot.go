package bot

import (
	"fmt"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const(
	delimiterIndent = "---------------------------------"
)

// Stores ids of all temporary created rooms
type createdVoiseChannel struct {
	ID string
}

func (ch *createdVoiseChannel) String() string {
	return ch.ID
}

type client struct {
	session *discordgo.Session
	channels []createdVoiseChannel
}

func New(token string, rooms int) (*client, error) {
	// Create new Discord session using the provided discordbot token.
	s, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, fmt.Errorf("creating Discord session error: %v", err)
	}

	return &client{
		session: s,
		channels: make([]createdVoiseChannel, 0, rooms),
	}, nil
}

func (c *client) Start() {
	// Register the messageCreate func as a callback for MessageCreate events.
	c.session.AddHandler(c.messageCreate)

	// Checks if anyone joined defined voice channel
	c.session.AddHandler(func(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
		c.voiceChannelJoinChecker(s, m, c.channels)
	})

	c.session.AddHandler(func(s *discordgo.Session, m *discordgo.VoiceStateUpdate) {
		c.voiceChannelLeftChecker(s, m, c.channels)
	})

	// Open a websocket connection to Discord and begin listening
	if err := c.session.Open(); err != nil {
		log.Fatalf("open connection error: %v", err)
	}
}

func (c *client) Close() {
	if err := c.session.Close(); err != nil {
		log.Printf("close session error: %v", err)
	}
}

func (c *client) printAllMembers(createdVoiceChannels []createdVoiseChannel) {
	channels := ""
	for _, v := range createdVoiceChannels {
		channels += fmt.Sprintf("%v,\n", v)
	}
	log.Printf("%s\nCreated voice channels: %s", delimiterIndent, channels)
}

func (c *client) voiceChannelLeftChecker(s *discordgo.Session, m *discordgo.VoiceStateUpdate, createdVoiceChannels []createdVoiseChannel) {
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

					c.printAllMembers(createdVoiceChannels)

					if _, err := s.ChannelDelete(v.ID); err != nil {
						log.Printf("delete channel error: %v", err)
					}
				}
			}
		}

	}

}

func (c *client) voiceChannelJoinChecker(s *discordgo.Session, m *discordgo.VoiceStateUpdate, createdVoiceChannels []createdVoiseChannel) {
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

		if err := s.GuildMemberMove(m.GuildID, m.UserID, tempChannel.ID); err != nil {
			log.Printf("move member error: %v", err)
		}

		createdVoiceChannels = append(createdVoiceChannels, createdVoiseChannel{
			ID: tempChannel.ID,
		})

		c.printAllMembers(createdVoiceChannels)
	}

}

func (c *client) messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if strings.HasPrefix(m.Content, ".") {

		fmt.Println(m.Content)

		// Ignore all messages created by the discordbot itself
		if m.Author.ID == s.State.User.ID {
			return
		}

		if m.Content == ".ping" {
			s.ChannelMessageSend(m.ChannelID, "pong!")
		}
	}
}
