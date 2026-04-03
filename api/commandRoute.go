package api

import (
	"github.com/bwmarrin/discordgo"
	"github.com/summrs-dev-team/summrs-premium/commands"
)

func init() {

	commandRoute.Add("antiinvite", commandRoute.AntiInvite, &commands.Config{
		Alias:        []string{"antiinv", "noinvite"},
		Cooldown:     3,
		RequiresArgs: true,
	})

	commandRoute.Add("antinuke", commandRoute.AntiNuke, &commands.Config{
		Cooldown: 3,
	})

	commandRoute.Add("antinukestatus", commandRoute.AntiNukeStatus, &commands.Config{
		Alias:    []string{"anstatus", "nukestatus"},
		Cooldown: 3,
	})

	commandRoute.Add("avatar", commandRoute.Avatar, &commands.Config{
		Alias:    []string{"av"},
		Cooldown: 1,
	})

	commandRoute.Add("ban", commandRoute.Ban, &commands.Config{
		Cooldown:        2,
		RequiresMention: true,
		Perms:           discordgo.PermissionBanMembers,
		WhitelistedOnly: true,
	})

	commandRoute.Add("banner", commandRoute.ServerBanner, &commands.Config{
		Alias:    []string{"serverbanner", "sbanner"},
		Cooldown: 1,
	})

	commandRoute.Add("botinfo", commandRoute.BotInfo, &commands.Config{
		Alias:    []string{"bi"},
		Cooldown: 4,
	})

	commandRoute.Add("help", commandRoute.Help, &commands.Config{
		Alias:    []string{"h"},
		Cooldown: 1,
	})

	commandRoute.Add("invite", commandRoute.Invite, &commands.Config{
		Cooldown: 1,
	})

	commandRoute.Add("kick", commandRoute.Kick, &commands.Config{
		Cooldown:        2,
		RequiresMention: true,
		Perms:           discordgo.PermissionBanMembers,
		WhitelistedOnly: true,
	})

	commandRoute.Add("lockdown", commandRoute.Lockdown, &commands.Config{
		Alias:    []string{"lock"},
		Cooldown: 1,
		Perms:    discordgo.PermissionManageChannels,
	})

	commandRoute.Add("hide", commandRoute.Hide, &commands.Config{
		Alias:    []string{"chide", "channelhide"},
		Cooldown: 1,
		Perms:    discordgo.PermissionManageChannels,
	})

	commandRoute.Add("logchannel", commandRoute.LoggingChannel, &commands.Config{
		Alias:    []string{"setlogs", "log"},
		Cooldown: 5,
	})

	commandRoute.Add("clearwhitelist", commandRoute.ClearWhitelists, &commands.Config{
		Alias:    []string{"clearwhitelists", "clearwl", "resetwl", "wlclear"},
		Cooldown: 5,
	})

	commandRoute.Add("unban", commandRoute.Unban, &commands.Config{
		Cooldown:        2,
		RequiresArgs:    true,
		Perms:           discordgo.PermissionBanMembers,
		WhitelistedOnly: true,
	})

	commandRoute.Add("massunban", commandRoute.MassUnban, &commands.Config{
		Alias:    []string{"unbanall"},
		Cooldown: 30,
		Perms:    discordgo.PermissionBanMembers,
	})

	commandRoute.Add("membercount", commandRoute.MemberCount, &commands.Config{
		Alias:    []string{"mc"},
		Cooldown: 1,
	})

	commandRoute.Add("moderationtype", commandRoute.ModerationType, &commands.Config{
		Alias:    []string{"mt"},
		Cooldown: 3,
	})

	commandRoute.Add("nuke", commandRoute.Nuke, &commands.Config{
		Alias:    []string{"nk"},
		Cooldown: 30,
		Perms:    discordgo.PermissionManageChannels,
	})

	commandRoute.Add("ping", commandRoute.Ping, &commands.Config{
		Cooldown: 5,
	})

	commandRoute.Add("prefix", commandRoute.Prefix, &commands.Config{
		Alias:        []string{"setprefix"},
		Cooldown:     3,
		RequiresArgs: true,
	})

	commandRoute.Add("servericon", commandRoute.ServerIcon, &commands.Config{
		Alias:    []string{"serverpfp", "sicon", "serverpic"},
		Cooldown: 1,
	})

	commandRoute.Add("serverinfo", commandRoute.ServerInfo, &commands.Config{
		Alias:    []string{"si"},
		Cooldown: 1,
	})

	commandRoute.Add("settings", commandRoute.Settings, &commands.Config{
		Cooldown: 1,
	})

	commandRoute.Add("slowmode", commandRoute.SlowMode, &commands.Config{
		Cooldown:     1,
		RequiresArgs: true,
		Perms:        discordgo.PermissionManageChannels,
	})

	commandRoute.Add("threshold", commandRoute.Threshold, &commands.Config{
		Alias:        []string{"t"},
		Cooldown:     3,
		RequiresArgs: true,
	})

	commandRoute.Add("toggle", commandRoute.Toggle, &commands.Config{
		Cooldown:     3,
		RequiresArgs: true,
	})

	commandRoute.Add("unlockdown", commandRoute.UnLockdown, &commands.Config{
		Alias:    []string{"unlock"},
		Cooldown: 1,
		Perms:    discordgo.PermissionManageChannels,
	})

	commandRoute.Add("unhide", commandRoute.Unhide, &commands.Config{
		Alias:    []string{"cunhide", "channelunhide"},
		Cooldown: 1,
		Perms:    discordgo.PermissionManageChannels,
	})

	commandRoute.Add("unslowmode", commandRoute.UnSlowMode, &commands.Config{
		Cooldown: 1,
		Perms:    discordgo.PermissionManageChannels,
	})

	commandRoute.Add("unwhitelist", commandRoute.Unwhitelist, &commands.Config{
		Alias:           []string{"uwl", "delwl", "removewl", "dewhitelist"},
		Cooldown:        3,
		RequiresMention: true,
	})

	commandRoute.Add("unwhitelistrole", commandRoute.UnWhitelistRole, &commands.Config{
		Alias:               []string{"delrolewl", "removerolewl", "dewhitelistrole"},
		Cooldown:            3,
		RequiresRoleMention: true,
	})

	commandRoute.Add("unwhitelistinvite", commandRoute.UnWhitelistInvite, &commands.Config{
		Alias:    []string{"delwebhook", "removewebhook", "dewhitelistwebhook"},
		Cooldown: 3,
	})

	commandRoute.Add("unwhitelistwebhook", commandRoute.UnWhitelistWebhook, &commands.Config{
		Alias:    []string{"delwebhook", "removewebhook", "dewhitelistwebhook"},
		Cooldown: 3,
	})

	commandRoute.Add("userinfo", commandRoute.UserInfo, &commands.Config{
		Alias:           []string{"whois"},
		Cooldown:        1,
		RequiresMention: true,
	})

	commandRoute.Add("whitelist", commandRoute.Whitelist, &commands.Config{
		Alias:           []string{"wl"},
		Cooldown:        3,
		RequiresMention: true,
	})

	commandRoute.Add("whitelistinvite", commandRoute.WhitelistInvite, &commands.Config{
		Alias:    []string{"wlinvite", "bypassinvite"},
		Cooldown: 3,
	})

	commandRoute.Add("whitelistrole", commandRoute.WhitelistRole, &commands.Config{
		Alias:    []string{"wlrole", "addrolewhitelist", "bypassrole"},
		Cooldown: 3,
	})

	commandRoute.Add("whitelistwebhook", commandRoute.WhitelistWebhook, &commands.Config{
		Alias:    []string{"wlwebhook", "bypasswebhook"},
		Cooldown: 3,
	})

	commandRoute.Add("whitelisted", commandRoute.ViewWhitelisted, &commands.Config{
		Alias:           []string{"wled"},
		Cooldown:        3,
		WhitelistedOnly: true,
	})

}
