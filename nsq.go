package dobby

type NsqMessage struct {
	Action   string      `json:"action" msgpack:"action"`
	Message  interface{} `json:"message" msgpack:"message"`
	CreateAt int         `json:"ct,omitempty" msgpack:"ct,omitempty"`
	APP      string      `json:"app" msgpack:"app"`
}
