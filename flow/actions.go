package flow

// набор стандартных экшнов на шаги флоу. Кастомные тоже можно использовать, но эти зафиксированы
type Action = string

const (
	SendMessage    Action = "send_message"
	CollectText    Action = "collect_text"
	CollectContact Action = "collect_contact"
)
