package telegram

type User struct {
	Id                      int64  `json:"id"`
	IsBot                   bool   `json:"is_bot"`
	FirstName               string `json:"first_name"`
	LastName                string `json:"last_name"`
	Username                string `json:"username"`
	LanguageCode            string `json:"language_code"`
	CanJoinGroups           bool   `json:"can_join_groups"`
	CanReadAllGroupMessages bool   `json:"can_read_all_group_messages"`
	SupportsInlineQueries   bool   `json:"supports_inline_queries"`
}

type ChatPhoto struct {
	SmallFileId       string `json:"small_file_id"`
	SmallFileUniqueId string `json:"small_file_unique_id"`
	BigFileId         string `json:"big_file_id"`
	BigFileUniqueId   string `json:"big_file_unique_id"`
}

type ChatPermissions struct {
	CanSendMessages       bool `json:"can_send_messages"`
	CanSendMediaMessages  bool `json:"can_send_media_messages"`
	CanSendPolls          bool `json:"can_send_polls"`
	CanSendOtherMessages  bool `json:"can_send_other_messages"`
	CanAddWebPagePreviews bool `json:"can_add_web_page_previews"`
	CanChangeInfo         bool `json:"can_change_info"`
	CanInviteUsers        bool `json:"can_invite_users"`
	CanPinMessages        bool `json:"can_pin_messages"`
}

type Chat struct {
	Id          int64     `json:"id"`
	Type        string    `json:"type"`
	Title       string    `json:"title"`
	Username    string    `json:"username"`
	FirstName   string    `json:"first_name"`
	LastName    string    `json:"last_name"`
	Photo       ChatPhoto `json:"photo"`
	Description string    `json:"description"`
	InviteLink  string    `json:"invite_link"`

	PinnedMessage    *Message        `json:"pinned_message"`
	Permissions      ChatPermissions `json:"permissions"`
	SlowModeDelay    int             `json:"slow_mode_delay"`
	StickerSetName   string          `json:"sticker_set_name"`
	CanSetStickerSet bool            `json:"can_set_sticker_set"`
}

type Message struct {
	MessageId            int      `json:"message_id"`
	From                 User     `json:"from"`
	Date                 int      `json:"date"`
	Chat                 Chat     `json:"chat"`
	ForwardFrom          User     `json:"forward_from"`
	ForwardFromChat      Chat     `json:"forward_from_chat"`
	ForwardFromMessageId int      `json:"forward_from_message_id"`
	ForwardSignature     string   `json:"forward_signature"`
	ForwardSenderName    string   `json:"forward_sender_name"`
	ForwardDate          int      `json:"forward_date"`
	ReplyToMessage       *Message `json:"reply_to_message"`
	EditDate             int      `json:"edit_date"`
	MediaGroupId         string   `json:"media_group_id"`
	AuthorSignature      string   `json:"author_signature"`
	Text                 string   `json:"text"`
}

type Location struct {
	Longitude            float64 `json:"longitude"`
	Latitude             float64 `json:"latitude"`
	HorizontalAccuracy   float64 `json:"horizontal_accuracy"`
	LivePeriod           int64   `json:"live_period"`
	Heading              int64   `json:"heading"`
	ProximityAlertRadius int64   `json:"proximity_alert_radius"`
}

type InlineQuery struct {
	Id       string   `json:"id"`
	From     User     `json:"from"`
	Location Location `json:"location"`
	Query    string   `json:"query"`
	Offset   string   `json:"offset"`
}

type Update struct {
	UpdateId          int         `json:"update_id"`
	Message           Message     `json:"message"`
	EditedMessage     Message     `json:"edited_message"`
	ChannelPost       Message     `json:"channel_post"`
	EditedChannelPost Message     `json:"edited_channel_post"`
	InlineQuery       InlineQuery `json:"inline_query"`
}

type GetUpdate struct {
	Ok     bool     `json:"ok"`
	Result []Update `json:"result"`
}

type GetMe struct {
	Ok     bool `json:"ok"`
	Result User `json:"result"`
}
