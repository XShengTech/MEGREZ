package middleware

type Response struct {
	Code ResCode `json:"code"`
	Msg  string  `json:"msg"`
	Data *Data   `json:"data"`
}

type Data struct {
	Result any   `json:"result"`
	Total  int64 `json:"total,omitempty"`
}
