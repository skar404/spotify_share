package telegram

type InlineKeyboardButtonReq struct {
	Text string `json:"text"`
	Url  string `json:"url"`
	//LoginUrl LoginUser `json:"login_url"`
	//CallbackData                 string `json:"callback_data"`
	//SwitchInlineQuery            string `json:"switch_inline_query"`
	//SwitchInlineQueryCurrentChat string `json:"switch_inline_query_current_chat"`
	//CallbackGame CallbackGame `json:"callback_game"`
	//Pay bool `json:"pay"`
}

type InlineKeyboardReq struct {
	InlineKeyboard []InlineKeyboardButtonReq `json:"inline_keyboard"`
}
