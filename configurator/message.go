package configurator

type Message struct {
	Text    string   `yaml:"text"`
	Buttons []Button `yaml:"buttons"`
	Answers []string `yaml:"answers"`
	Image   string   `yaml:"image"`
	File    string   `yaml:"file"`
}

type Button struct {
	Text string `yaml:"text"`
	Code string `yaml:"code"`
	Link string `yaml:"link"`
}
