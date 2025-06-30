package messages

type Message struct {
	Text    string   `yaml:"text"`
	Buttons []Button `yaml:"buttons"` // todo добавить верстку строк
	Answers []Answer `yaml:"answers"`
	Image   string   `yaml:"image"`
	File    string   `yaml:"file"`
}

type Button struct {
	Text string `yaml:"text"`
	Code string `yaml:"code"`
	Link string `yaml:"link"`
}

type Answer struct {
	Text    string `yaml:"text"`
	Contact bool   `yaml:"request_contact"`
	Link    string `yaml:"link"`
}

func (a Answer) CallbackUnique() string {
	return a.Text
}
