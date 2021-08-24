package domain

type RabbitMQMeta struct {
	ChannelName  string
	ExchangeType string
	ExchangeName string
	QueueName    string
	RoutingKey   string
	BindingKey   string
}

type MessageUserLogin struct {
	ID     string `json:"id"`
	Source string `json:"source"`
}
