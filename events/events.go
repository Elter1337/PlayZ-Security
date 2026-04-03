package events

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/summrs-dev-team/summrs-premium/database"
	"github.com/summrs-dev-team/summrs-premium/utils"
)

func antiNukeEnabled(data database.GuildData) bool {
	enabled, ok := data["antinuke-enabled"].(bool)
	if !ok {
		return false
	}
	return enabled
}

func AntiInvite(s *discordgo.Session, m *discordgo.MessageCreate) {
	data, err := database.Database.FindData(m.GuildID)
	switch {
	case err != nil:
		return
	case data["anti-invite"] == "off":
		return
	case database.Database.IsWhitelisted(m.GuildID, "whitelisted-invite-channels", m.ChannelID, nil) || utils.HasPerms(s, m.Message, m.GuildID, m.Author.ID, discordgo.PermissionManageMessages):
		return
	}

	if strings.Contains(m.Content, "discord.gg/") {
		_ = s.ChannelMessageDelete(m.ChannelID, m.ID)
	}
}

func AntiMassMention(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m == nil || m.Author == nil || m.Author.Bot || m.GuildID == "" {
		return
	}

	content := strings.ToLower(m.Content)
	mentionsEveryone := strings.Contains(content, "@everyone")
	mentionsHere := strings.Contains(content, "@here")
	if !mentionsEveryone && !mentionsHere {
		return
	}

	if !database.Database.IsSettingEnabled(m.GuildID, "antinuke-enabled") {
		return
	}

	data, err := database.Database.FindData(m.GuildID)
	if err != nil {
		return
	}

	if mentionsEveryone {
		if !toggleWithDefault(data, "anti-everyone-mention", true) {
			mentionsEveryone = false
		}
	}
	if mentionsHere {
		if !toggleWithDefault(data, "anti-here-mention", true) {
			mentionsHere = false
		}
	}
	if !mentionsEveryone && !mentionsHere {
		return
	}

	var member *discordgo.Member
	if m.Member != nil {
		member = m.Member
	} else {
		member, _ = s.GuildMember(m.GuildID, m.Author.ID)
	}
	if database.Database.IsWhitelisted(m.GuildID, "users", m.Author.ID, member) {
		return
	}
	if utils.GetGuildOwner(s, m.GuildID) == m.Author.ID {
		return
	}
	if utils.HasPerms(s, m.Message, m.GuildID, m.Author.ID, discordgo.PermissionManageMessages) {
		return
	}
	moderationType, _ := data["moderation-type"].(string)
	if moderationType == "" {
		moderationType = "ban"
	}

	started := time.Now()
	_ = s.ChannelMessageDelete(m.ChannelID, m.ID)

	reason := "Mentioning Everyone/Here"
	eventTitle := "Mass Mention"
	if mentionsEveryone && !mentionsHere {
		reason = "Mentioning Everyone"
		eventTitle = "Everyone Mention"
	} else if mentionsHere && !mentionsEveryone {
		reason = "Mentioning Here"
		eventTitle = "Here Mention"
	}

	if err := utils.HandleModerationWithType(s, m.GuildID, m.Author.ID, reason, moderationType); err != nil {
		return
	}

	logChannel, _ := data["log-channel"].(string)
	channelName := fmt.Sprintf("#%s", m.ChannelID)
	if ch, chErr := s.State.Channel(m.ChannelID); chErr == nil && ch != nil {
		channelName = fmt.Sprintf("#%s", ch.Name)
	}
	utils.SendAntiNukeLogToChannel(s, m.GuildID, logChannel, utils.AntiNukeLogData{
		EventTitle:     eventTitle,
		CriminalID:     m.Author.ID,
		Crime:          strings.ToLower(reason),
		VictimID:       m.ChannelID,
		VictimName:     channelName,
		CounterAction:  moderationType,
		ActionDuration: utils.FormatActionLatency(started, ""),
	})
}

func BanHandler(s *discordgo.Session, event *discordgo.GuildBanAdd) {
	if !database.Database.IsSettingEnabled(event.GuildID, "antinuke-enabled") {
		return
	}

	data, err := database.Database.FindData(event.GuildID)
	if err != nilabase.FindData(event.GuildID)
	if err != nil {
		return
	}

	if enabled, _ := data["anti-ban"].(bool); !enabled {
		return
	}

	utils.ReadAuditWithRecoveryData(s, event.GuildID, "banned a member", 22, data, func(entry *discordgo.AuditLogEntry) error {
		if entry.TargetID == "" {
			return fmt.Errorf("missing banned user id")
		}
		return s.GuildBanDelete(event.GuildID, entry.TargetID)
	})
}

func BanRemoveHandler(s *discordgo.Session, event *discordgo.GuildBanRemove) {
	data, err := database.Database.FindData(event.GuildID)
	if err != nil {
		return
	}
	if !antiNukeEnabled(data) {
		return
	}
	if enabled, _ := data["anti-ban"].(bool); !enabled {
		return
	}

	utils.ReadAuditWithRecoveryData(s, event.GuildID, "unbanned a member", 23, data, func(entry *discordgo.AuditLogEntry) error {
		if entry.TargetID == "" {
			return fmt.Errorf("missing unbanned user id")
		}
		return s.GuildBanCreateWithReason(event.GuildID, entry.TargetID, utils.FormatReason("Auto Recovery"), 0)
	})
}

func ChannelCreate(s *discordgo.Session, event *discordgo.ChannelCreate) {
	if !database.Database.IsSettingEnabled(event.GuildID, "antinuke-enabled") {
		return
	}

	data, err := database.Database.FindData(event.GuildID)
	if err != nil {
		return
	}

	if enabled, _ := data["anti-channel-create"].(bool); !enabled {
		return
	}

	utils.ReadAuditWithRecoveryData(s, event.GuildID, "created a channel", 10, data, func(entry *discordgo.AuditLogEntry) error {
		if entry.TargetID == "" {
			return fmt.Errorf("missing created channel id")
		}
		_, err := s.ChannelDelete(entry.TargetID)
		return err
	})
}

func ChannelRemove(s *discordgo.Session, event *discordgo.ChannelDelete) {
	if !database.Database.IsSettingEnabled(event.GuildID, "antinuke-enabled") {
		return
	}

	data, err := database.Database.FindData(event.GuildID)
	if err != nil {
		return
	}

	if enabled, _ := data["anti-channel-delete"].(bool); !enabled {
		return
	}

	utils.ReadAuditWithRecoveryData(s, event.GuildID, "deleted a channel", 12, data, func(entry *discordgo.AuditLogEntry) error {
		if event.Channel == nil {
			return fmt.Errorf("missing deleted channel payload")
		}
		_, err := s.GuildChannelCreateComplex(event.GuildID, discordgo.GuildChannelCreateData{
			Name:                 event.Channel.Name,
			Type:                 event.Channel.Type,
			Topic:                event.Channel.Topic,
			RateLimitPerUser:     event.Channel.RateLimitPerUser,
			Position:             event.Channel.Position,
			PermissionOverwrites: event.Channel.PermissionOverwrites,
			ParentID:             event.Channel.ParentID,
			NSFW:                 event.Channel.NSFW,
		})
		return err
	})
}

func ChannelUpdate(s *discordgo.Session, event *discordgo.ChannelUpdate) {
	if !database.Database.IsSettingEnabled(event.GuildID, "antinuke-enabled") {
		return
	}

	data, err := database.Database.FindData(event.GuildID)
	if err != nil {
		return
	}
	if !toggleWithDefault(data, "anti-channel-update", true) {
		return
	}

	utils.ReadAuditWithRecoveryData(s, event.GuildID, "updated a channel", 11, data, func(entry *discordgo.AuditLogEntry) error {
		if event.BeforeUpdate == nil || event.Channel == nil {
			return fmt.Errorf("missing channel before state for recovery")
		}

		var rateLimit int
		if event.BeforeUpdate.RateLimitPerUser > 0 {
			rateLimit = event.BeforeUpdate.RateLimitPerUser
		}
		var nsfw bool
		if event.BeforeUpdate.NSFW {
			nsfw = true
		}

		_, err := s.ChannelEditComplex(event.Channel.ID, &discordgo.ChannelEdit{
			Name:                 event.BeforeUpdate.Name,
			Topic:                event.BeforeUpdate.Topic,
			NSFW:                 &nsfw,
			Bitrate:              event.BeforeUpdate.Bitrate,
			UserLimit:            event.BeforeUpdate.UserLimit,
			ParentID:             event.BeforeUpdate.ParentID,
			RateLimitPerUser:     &rateLimit,
			PermissionOverwrites: event.BeforeUpdate.PermissionOverwrites,
		})
		return err
	})
}

func CreateGuild(s *discordgo.Session, event *discordgo.GuildCreate) {
	muteX.Lock()
	defer muteX.Unlock()

	s.State.GuildAdd(event.Guild)
	database.Database.CreateGuild(s.State.User, event.Guild)

	if _, ok := guilds[event.Guild.ID]; ok {
		return
	}

	guilds[event.Guild.ID] = event.Guild.MemberCount
	MemberCount += guilds[event.Guild.ID]
	GuildCount++
}

func DeleteGuild(s *discordgo.Session, event *discordgo.GuildDelete) {
	database.Database.DeleteGuild(event.Guild.ID)
	MemberCount -= guilds[event.Guild.ID]

	muteX.Lock()
	defer muteX.Unlock()
	delete(guilds, event.Guild.ID)
	GuildCount--
}

func GuildUpdate(s *discordgo.Session, event *discordgo.GuildUpdate) {
	if !database.Database.IsSettingEnabled(event.ID, "antinuke-enabled") {
		return
	}

	guildData, err := database.Database.FindData(event.ID)
	if err != nil {
		return
	}

	if toggleWithDefault(guildData, "anti-guild-update", true) {
		utils.ReadAuditWithData(s, event.Guild.ID, "updated guild settings", 1, guildData)
	}

	entry, _, auditErr := utils.FindAudit(s, event.Guild.ID, 1)
	if entry == nil || (auditErr != nil && auditErr.Error() != "Whitelisted") {
		return
	}

	for _, change := range entry.Changes {
		if change.Key == nil {
			continue
		}

		switch *change.Key {
		case discordgo.AuditLogChangeKeyName:
			if enabled, _ := guildData["anti-name-change"].(bool); !enabled || change.OldValue == nil {
				continue
			}
			if auditErr != nil && auditErr.Error() == "Whitelisted" {
				_, _ = database.Database.SetData("$set", event.Guild.ID, "guild-name", change.NewValue.(string))
				continue
			}
			_, _ = s.GuildEdit(event.Guild.ID, &discordgo.GuildParams{Name: guildData["guild-name"].(string)})

		case discordgo.AuditLogChangeKeyVanityURLCode:
			if enabled, _ := guildData["anti-vanity-steal"].(bool); !enabled || change.OldValue == nil {
				continue
			}
			if auditErr != nil && auditErr.Error() == "Whitelisted" {
				_, _ = database.Database.SetData("$set", event.Guild.ID, "vanity-url", change.NewValue.(string))
				continue
			}
			jsonData := []byte(fmt.Sprintf(`{"code":"%s"}`, guildData["vanity-url"]))
			_, _ = utils.MakeRequest("PATCH", fmt.Sprintf("https://discord.com/api/v10/guilds/%s/vanity-url", event.Guild.ID), s.Token, jsonData)

		case discordgo.AuditLogChangeKeyWidgetEnabled:
			if enabled, _ := guildData["anti-widget-spam"].(bool); !enabled || auditErr != nil {
				continue
			}
			if oldVal, ok := change.OldValue.(bool); ok {
				if newVal, ok2 := change.NewValue.(bool); ok2 && oldVal == newVal {
					continue
				}
			}
			_ = utils.HandleModeration(s, event.Guild.ID, entry.UserID, "tried to spam guild widget changes")
		}
	}
}

func KickHandler(s *discordgo.Session, event *discordgo.GuildMemberRemove) {
	if !database.Database.IsSettingEnabled(event.GuildID, "antinuke-enabled") {
		return
	}

	data, err := database.Database.FindData(event.GuildID)
	if err != nil {
		return
	}

	if enabled, _ := data["anti-kick"].(bool); enabled {
		utils.ReadAuditWithData(s, event.GuildID, "kicked a member", 20, data)
	}
	if enabled, _ := data["anti-prune"].(bool); enabled {
		utils.ReadAuditWithData(s, event.GuildID, "ran member prune", 21, data)
	}
}

func MemberJoin(s *discordgo.Session, event *discordgo.GuildMemberAdd) {
	started := time.Now()

	MemberCount++
	s.State.MemberAdd(event.Member)

	if !event.Member.User.Bot {
		return
	}

	if !database.Database.IsSettingEnabled(event.GuildID, "antinuke-enabled") {
		return
	}

	data, err := database.Database.FindData(event.GuildID)
	if err != nil {
		return
	}
	if enabled, _ := data["anti-bot"].(bool); !enabled {
		return
	}
	moderationType, _ := data["moderation-type"].(string)
	if moderationType == "" {
		moderationType = "ban"
	}

	entry, _, err := utils.FindAudit(s, event.GuildID, 28)
	if entry == nil || err != nil {
		return
	}

	_ = utils.HandleModerationWithType(s, event.GuildID, entry.UserID, "Anti Bot Invite", moderationType)
	_ = utils.HandleModerationWithType(s, event.GuildID, event.User.ID, "Anti Bot Invite", moderationType)
	logChannel, _ := data["log-channel"].(string)
	utils.SendAntiNukeLogToChannel(s, event.GuildID, logChannel, utils.AntiNukeLogData{
		EventTitle:     "Bot Invite",
		CriminalID:     entry.UserID,
		Crime:          fmt.Sprintf("invited bot (%s)", event.User.Username),
		VictimID:       event.User.ID,
		VictimName:     fmt.Sprintf("%s#%s", event.User.Username, event.User.Discriminator),
		CounterAction:  moderationType,
		ActionDuration: utils.FormatActionLatency(started, entry.ID),
	})
}

func MemberLeave(s *discordgo.Session, event *discordgo.GuildMemberRemove) {
	MemberCount--
	s.State.MemberRemove(event.Member)
}

func MemberRoleUpdate(s *discordgo.Session, event *discordgo.GuildMemberUpdate) {
	started := time.Now()

	if !database.Database.IsSettingEnabled(event.GuildID, "antinuke-enabled") {
		return
	}

	data, err := database.Database.FindData(event.GuildID)
	if err != nil {
		return
	}
	if enabled, _ := data["anti-member-role"].(bool); !enabled {
		return
	}
	moderationType, _ := data["moderation-type"].(string)
	if moderationType == "" {
		moderationType = "ban"
	}

	entry, change, err := utils.FindAudit(s, event.GuildID, 25)
	if err != nil || change == nil {
		return
	}

	changes, ok := change.([]interface{})
	if !ok || len(changes) == 0 {
		return
	}

	roleID, _ := changes[0].(map[string]interface{})["id"].(string)
	if roleID == "" {
		return
	}

	guildRole, err := s.State.Role(event.GuildID, roleID)
	if err != nil || guildRole.Permissions&discordgo.PermissionAdministrator == 0 {
		return
	}

	if err = s.GuildMemberRoleRemove(event.GuildID, entry.TargetID, roleID); err != nil {
		return
	}

	if err = utils.HandleModerationWithType(s, event.GuildID, entry.UserID, "gave admin role to a member", moderationType); err != nil {
		return
	}

	logChannel, _ := data["log-channel"].(string)
	victimName := fmt.Sprintf("User `%s`", entry.TargetID)
	if victim, victimErr := s.User(entry.TargetID); victimErr == nil && victim != nil {
		victimName = fmt.Sprintf("%s#%s", victim.Username, victim.Discriminator)
	}
	utils.SendAntiNukeLogToChannel(s, event.GuildID, logChannel, utils.AntiNukeLogData{
		EventTitle:     "Admin Role Update",
		CriminalID:     entry.UserID,
		Crime:          "gave administrator role to a member",
		VictimID:       entry.TargetID,
		VictimName:     victimName,
		CounterAction:  moderationType,
		ActionDuration: utils.FormatActionLatency(started, entry.ID),
	})
}

func Ready(s *discordgo.Session, event *discordgo.Ready) {
	_ = s.UpdateStreamingStatus(2, fmt.Sprintf(">help"), "https://twitch.tv/discord")
	fmt.Printf("Connected to shard #%d\n", s.ShardID)
}

func RoleCreate(s *discordgo.Session, event *discordgo.GuildRoleCreate) {
	if !database.Database.IsSettingEnabled(event.GuildID, "antinuke-enabled") {
		return
	}

	data, err := database.Database.FindData(event.GuildID)
	if err != nil {
		return
	}
	if enabled, _ := data["anti-role-create"].(bool); !enabled {
		return
	}
	utils.ReadAuditWithRecoveryData(s, event.GuildID, "created a role", 30, data, func(entry *discordgo.AuditLogEntry) error {
		if entry.TargetID == "" {
			return fmt.Errorf("missing created role id")
		}
		return s.GuildRoleDelete(event.GuildID, entry.TargetID)
	})
}

func RoleRemove(s *discordgo.Session, event *discordgo.GuildRoleDelete) {
	if !database.Database.IsSettingEnabled(event.GuildID, "antinuke-enabled") {
		return
	}

	data, err := database.Database.FindData(event.GuildID)
	if err != nil {
		return
	}
	if enabled, _ := data["anti-role-delete"].(bool); !enabled {
		return
	}
	utils.ReadAuditWithRecoveryData(s, event.GuildID, "deleted a role", 32, data, func(entry *discordgo.AuditLogEntry) error {
		recoveredName := "Auto-Recovered Role"
		for _, change := range entry.Changes {
			if change.Key == nil {
				continue
			}
			if *change.Key != discordgo.AuditLogChangeKeyName || change.OldValue == nil {
				continue
			}
			if name, ok := change.OldValue.(string); ok && strings.TrimSpace(name) != "" {
				recoveredName = name
			}
			break
		}

		_, err := s.GuildRoleCreate(event.GuildID, &discordgo.RoleParams{
			Name: recoveredName,
		})
		return err
	})
}

func RoleUpdate(s *discordgo.Session, event *discordgo.GuildRoleUpdate) {
	if !database.Database.IsSettingEnabled(event.GuildID, "antinuke-enabled") {
		return
	}

	data, err := database.Database.FindData(event.GuildID)
	if err != nil {
		return
	}
	if !toggleWithDefault(data, "anti-role-update", true) {
		return
	}

	utils.ReadAuditWithRecoveryData(s, event.GuildID, "updated a role", 31, data, func(entry *discordgo.AuditLogEntry) error {
		if entry.TargetID == "" {
			return fmt.Errorf("missing updated role id")
		}

		params := &discordgo.RoleParams{}
		edited := false
		for _, change := range entry.Changes {
			if change.Key == nil || change.OldValue == nil {
				continue
			}

			switch *change.Key {
			case discordgo.AuditLogChangeKeyName:
				if oldName, ok := change.OldValue.(string); ok {
					params.Name = oldName
					edited = true
				}
			case discordgo.AuditLogChangeKeyPermissions:
				if perms, ok := toInt64(change.OldValue); ok {
					params.Permissions = &perms
					edited = true
				}
			case discordgo.AuditLogChangeKeyMentionable:
				if mentionable, ok := change.OldValue.(bool); ok {
					params.Mentionable = &mentionable
					edited = true
				}
			case discordgo.AuditLogChangeKeyHoist:
				if hoist, ok := change.OldValue.(bool); ok {
					params.Hoist = &hoist
					edited = true
				}
			case discordgo.AuditLogChangeKeyColor:
				if color, ok := toInt(change.OldValue); ok {
					params.Color = &color
					edited = true
				}
			}
		}

		if !edited {
			return nil
		}

		_, err = s.GuildRoleEdit(event.GuildID, entry.TargetID, params)
		return err
	})
}

func WebhookCreate(s *discordgo.Session, event *discordgo.WebhooksUpdate) {
	if !database.Database.IsSettingEnabled(event.GuildID, "antinuke-enabled") {
		return
	}

	data, err := database.Database.FindData(event.GuildID)
	if err != nil {
		return
	}
	if enabled, _ := data["anti-webhook-create"].(bool); !enabled {
		return
	}
	moderationType, _ := data["moderation-type"].(string)
	if moderationType == "" {
		moderationType = "ban"
	}

	webhooks, err := s.ChannelWebhooks(event.ChannelID)
	if err != nil {
		return
	}

	selfMember, err := s.State.Member(event.GuildID, s.State.User.ID)
	if err != nil {
		selfMember, err = s.GuildMember(event.GuildID, s.State.User.ID)
		if err != nil {
			return
		}
	}

	for _, webhook := range webhooks {
		started := time.Now()

		targetMember, err := s.State.Member(event.GuildID, webhook.User.ID)
		if err != nil {
			targetMember, err = s.GuildMember(event.GuildID, webhook.User.ID)
			if err != nil {
				continue
			}
		}

		if database.Database.IsWhitelisted(event.GuildID, "users", webhook.User.ID, targetMember) ||
			database.Database.IsWhitelisted(event.GuildID, "whitelisted-webhook-channels", event.ChannelID, nil) {
			continue
		}
		if utils.GetGuildOwner(s, event.GuildID) == webhook.User.ID {
			continue
		}

		_ = s.WebhookDelete(webhook.ID)

		targetHighest := utils.HighestRole(s, event.GuildID, targetMember)
		selfHighest := utils.HighestRole(s, event.GuildID, selfMember)
		if targetHighest == nil || selfHighest == nil {
			continue
		}

		if !utils.IsAbove(selfHighest, targetHighest) || !utils.HasPerms(s, nil, event.GuildID, selfMember.User.ID, discordgo.PermissionBanMembers) {
			continue
		}

		if err := utils.HandleModerationWithType(s, event.GuildID, webhook.User.ID, "created a webhook", moderationType); err != nil {
			continue
		}

		logChannel, _ := data["log-channel"].(string)
		victimName := fmt.Sprintf("Channel `%s`", event.ChannelID)
		if ch, chErr := s.State.Channel(event.ChannelID); chErr == nil && ch != nil {
			victimName = fmt.Sprintf("#%s", ch.Name)
		}
		utils.SendAntiNukeLogToChannel(s, event.GuildID, logChannel, utils.AntiNukeLogData{
			EventTitle:     "Webhook Create",
			CriminalID:     webhook.User.ID,
			Crime:          "created a webhook",
			VictimID:       event.ChannelID,
			VictimName:     victimName,
			CounterAction:  moderationType,
			ActionDuration: utils.FormatActionLatency(started, ""),
		})
	}
}

func EmojiUpdate(s *discordgo.Session, event *discordgo.GuildEmojisUpdate) {
	data, err := database.Database.FindData(event.GuildID)
	if err != nil || !antiNukeEnabled(data) {
		return
	}

	auditTypes := []struct {
		Type   int
		Key    string
		Reason string
	}{
		{Type: 60, Key: "anti-emoji-create", Reason: "created an emoji"},
		{Type: 61, Key: "anti-emoji-update", Reason: "updated an emoji"},
		{Type: 62, Key: "anti-emoji-delete", Reason: "deleted an emoji"},
	}

	for _, candidate := range auditTypes {
		if !toggleWithDefault(data, candidate.Key, true) {
			continue
		}
		entry, _, err := utils.FindAudit(s, event.GuildID, candidate.Type)
		if entry == nil {
			continue
		}
		if err != nil {
			return
		}
		utils.ReadAuditWithData(s, event.GuildID, candidate.Reason, candidate.Type, data)
		return
	}
}

func StickerUpdate(s *discordgo.Session, event *discordgo.GuildStickersUpdate) {
	data, err := database.Database.FindData(event.GuildID)
	if err != nil || !antiNukeEnabled(data) {
		return
	}

	auditTypes := []struct {
		Type   int
		Key    string
		Reason string
	}{
		{Type: 90, Key: "anti-sticker-create", Reason: "created a sticker"},
		{Type: 91, Key: "anti-sticker-update", Reason: "updated a sticker"},
		{Type: 92, Key: "anti-sticker-delete", Reason: "deleted a sticker"},
	}

	for _, candidate := range auditTypes {
		if !toggleWithDefault(data, candidate.Key, true) {
			continue
		}
		entry, _, err := utils.FindAudit(s, event.GuildID, candidate.Type)
		if entry == nil {
			continue
		}
		if err != nil {
			return
		}
		utils.ReadAuditWithData(s, event.GuildID, candidate.Reason, candidate.Type, data)
		return
	}
}

var (
	guilds      = make(map[string]int)
	GuildCount  int
	MemberCount int
	muteX       = &sync.RWMutex{}
)

func toInt64(v interface{}) (int64, bool) {
	switch typed := v.(type) {
	case float64:
		return int64(typed), true
	case int64:
		return typed, true
	case int:
		return int64(typed), true
	case string:
		parsed, err := strconv.ParseInt(typed, 10, 64)
		return parsed, err == nil
	default:
		return 0, false
	}
}

func toInt(v interface{}) (int, bool) {
	switch typed := v.(type) {
	case float64:
		return int(typed), true
	case int64:
		return int(typed), true
	case int:
		return typed, true
	case string:
		parsed, err := strconv.Atoi(typed)
		return parsed, err == nil
	default:
		return 0, false
	}
}

func toggleWithDefault(data database.GuildData, key string, fallback bool) bool {
	raw, ok := data[key]
	if !ok {
		return fallback
	}
	value, ok := raw.(bool)
	if !ok {
		return fallback
	}
	return value
}
