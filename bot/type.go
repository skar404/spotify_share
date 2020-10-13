package bot

type Command struct {
	Name string
	Args []string
}

const (
	TemplateMessageStart  = "Привет, для работы данного бота необходимо авторизоваться черезе Spotify"
	TemplateMessageHelp   = "// FIXME - команда help"
	TemplateMessageLogout = "Вы успешно отключили Spotify"
	TemplateMessageLogin  = "Вы успешно авторизовались через Spotify\n" +
		"для работы бота введите имя бота %s и выбирите трек из списка."
)
