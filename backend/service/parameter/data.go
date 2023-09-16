package parameter

//@name: upload file
//@url: /v1/triple_star/user/upload

type UploadResp struct {
	Files []string `json:"files"`
}

//@name: add one data
//@url: /v1/triple_star/data/add

type DataAddReq struct {
	Name     string   `json:"name"`     // 名称
	Addr     string   `json:"addr"`     // 钱包地址
	Mark     []string `json:"mark"`     // 标签，可以多个
	File     string   `json:"file"`     // 文件名，zip压缩文件
	Memo     string   `json:"memo"`     // 数据简介，255个字以内
	Category string   `json:"category"` // 大类：人类基因/土壤微生物
	Price    float64  `json:"price"`    // 价格
}

//@name: query data
//@url: /v1/triple_star/data/query

type DataQueryReq struct {
	Mark []string `json:"mark"` // 标签
	Type int      `json:"type"` // 类型：人类基因/土壤微生物
}

type DataInfo struct {
	Id       uint     `json:"id"`                 // id
	Hash     string   `json:"hash,omitempty"`     // 哈希值
	Name     string   `json:"name,omitempty"`     // 名称
	Mark     []string `json:"mark,omitempty"`     // 标签
	Addr     string   `json:"addr"`               // 钱包地址
	File     string   `json:"file,omitempty"`     // 文件名
	CreateAt string   `json:"createAt,omitempty"` // 创建时间
	Memo     string   `json:"memo,omitempty"`     // 说明，简介
	Category string   `json:"category"`           // 类型：人类基因/土壤微生物
	Price    float64  `json:"price"`              // 价格
}

type DataQueryResp struct {
	DataList []*DataInfo `json:"dataList"` // 数据列表
}

//@name: buy one data
//@url: /v1/triple_star/data/buy

type DataBuyReq struct {
	Id   uint   `json:"id"`   // 数据ID
	Addr string `json:"addr"` // 购买者地址
}

//@name: query the user's data
//@url: /v1/triple_star/data/queryOneself

type UserDataQueryReq struct {
	Addr string
}

type UserDataQueryResp struct {
	Uploads []*DataInfo
	Buys    []*DataInfo
}
