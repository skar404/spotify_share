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
	ready bool

	InlineKeyboard [][]InlineKeyboardButtonReq `json:"inline_keyboard"`
}

func (s *InlineKeyboardReq) Ready() {
	s.ready = true
}

type AnswerCallbackReq struct {
	Text      string `json:"text"`
	ShowAlert bool   `json:"show_alert"`
	Url       string `json:"url"`
	CacheTime string `json:"cache_time"`
}
