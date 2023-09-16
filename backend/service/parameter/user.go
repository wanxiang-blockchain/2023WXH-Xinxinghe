package parameter

//@name: register
//@url: /v1/triple_star/user/register

type RegisterReq struct {
	Name     string `json:"name"`
	Addr     string `json:"addr"`
	Password string `json:"password"`
	Gender   int    `json:"gender"`
	Memo     string `json:"memo"`
}

type GeneralOperationResp struct {
	Success bool `json:"success"`
}

//@name: login
//@url: /v1/triple_star/user/login

type LoginReq struct {
	Addr     string `json:"addr"`
	Password string `json:"password"`
}

type LoginResp struct {
	Token  string `json:"token"`
	Expiry string `json:"expiry"`
}
