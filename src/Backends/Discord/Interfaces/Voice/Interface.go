package Voice

import (
	"errors"

	"antegr.al/chatanium-bot/v1/src/Util/Log"
	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

type Voice struct {
	conn          *discordgo.VoiceConnection
	broadcastChan chan *discordgo.Packet
	listenerChan  []chan *discordgo.Packet
	isLocked      bool
}

func (t *Voice) Init(conn *discordgo.VoiceConnection) {
	t.conn = conn
	t.isLocked = false
}

// read action for voice channel.
// it can be read by many modules, but it can be restricted by ACL. (ACL is todo)
//
// Note: even if a module is the first to open and close a Read channel,
// it will not close it if another module is using it.
func (t *Voice) Read() chan *discordgo.Packet {
	ch := make(chan *discordgo.Packet)
	t.listenerChan = append(t.listenerChan, ch)

	// if broadcast channel is not exist, create new broadcast channel
	if t.broadcastChan == nil {
		t.newBroadcastChan()
	}

	return ch
}

// write action for voice channel.
// it locked by another module, the another module can't write voice channel.
func (t *Voice) Write(pcm <-chan []int16) error {
	if t.isLocked {
		return errors.New("this voice channel is already used by another module.")
	}

	// start to write voice channel, it should be locked
	t.lock()

	// write voice channel
	dgvoice.SendPCM(t.conn, pcm)

	// if stream is closed, unlock voice channel
	t.unlock()

	return nil
}

func (t *Voice) lock() {
	t.isLocked = true
}

func (t *Voice) unlock() {
	t.isLocked = false
}

func (t *Voice) newBroadcastChan() {
	Log.Info.Printf("G:%s C:%s > starting broadcast voice channel...", t.conn.GuildID, t.conn.ChannelID)

	// start to broadcast voice channel
	go dgvoice.ReceivePCM(t.conn, t.broadcastChan)

	// manage broadcast channel
	go func() {
		for {
			data, incoming := <-t.broadcastChan

			// if no listener, close broadcast channel
			if len(t.listenerChan) == 0 {
				Log.Info.Printf("G:%s C:%s > starting broadcast voice channel...", t.conn.GuildID, t.conn.ChannelID)
				close(t.broadcastChan)
			}

			// if exist listener and incoming data, copy voice data to listeners
			if incoming {
				// copy data to listeners
				for _, ch := range t.listenerChan {
					ch <- data
				}
			} else {
				for _, ch := range t.listenerChan {
					close(ch)
				}
				return
			}
		}
	}()
}
