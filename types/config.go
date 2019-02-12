package types

type Http struct {
	ListenAddress string
}

type Mail struct {
	Host string
	From string
	To   string
}

type App struct {
	Name string
	Dir  string
}

type Config struct {
	Dev    bool
	Secret string
	Http   Http
	Mail   Mail
	Apps   []App
}
