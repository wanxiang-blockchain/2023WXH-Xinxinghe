package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/sirupsen/logrus"
	"log"
	"math/big"
	"os"
	"time"
	"triple_star/dao"
	"triple_star/service/contract"
	"triple_star/service/parameter"
	merror "triple_star/util/util_error"
)

// -------------------- contractor ---------------------------
type storage struct {
	conn    *ethclient.Client
	Storage *contract.Storage
}

const (
	prvKey = "5917e2ac1cd4678d82cd070ac289e3328c1df53f0f6c2b9ac85a2356e9ab9b08"
	addr   = "0x9d66d2d4d87aafddb5bc2efdc51e6c86abd78662"
	url    = "https://mainnet.infura.io/v3/5c01ca17a44d4f479c520e74ef29dd68"
)

var storageObj = &storage{}

func (s *storage) New() {
	conn, err := ethclient.Dial(url)
	if err != nil {
		logrus.WithField("err-msg", err).
			Errorln("Failed to connect to the Ethereum client")
	}
	s.conn = conn
	// Instantiate the contract and display its name
	// NOTE update the deployment address!
	store, err := contract.NewStorage(common.HexToAddress(addr), conn)
	if err != nil {
		logrus.WithField("err-msg", err).
			Errorln("failed to instantiate storageObj contract")
	}
	s.Storage = store
}

func (s *storage) Store(hash, updateAt string) error {
	PrivateKey, _ := crypto.HexToECDSA(prvKey)

	// Create an authorized transactor and call the store function
	nonce, _ := s.conn.NonceAt(context.Background(), common.HexToAddress("你私钥对应的账户地址"), nil)
	gasPrice, _ := s.conn.SuggestGasPrice(context.Background())
	// id
	auth, err := bind.NewKeyedTransactorWithChainID(PrivateKey, big.NewInt(5))
	auth.GasLimit = uint64(300000)
	auth.Nonce = new(big.Int).SetUint64(nonce)
	auth.GasPrice = gasPrice
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	// Call the store() function
	tx, err := s.Storage.Store(auth, hash, updateAt)
	if err != nil {
		logrus.WithField("err-msg", err).
			Errorln("Failed to update value: %v", err)
		return err
	}
	logrus.Infoln("Update pending: 0x%x\n", tx.Hash())
	return nil
}

type payment struct {
	conn    *ethclient.Client
	Payment *contract.Payment
}

var paymentObj = &payment{}

func (p *payment) New() {
	conn, err := ethclient.Dial(url)
	if err != nil {
		logrus.WithField("err-msg", err).
			Errorln("Failed to connect to the Ethereum client")
	}
	p.conn = conn
	// Instantiate the contract and display its name
	// NOTE update the deployment address!
	pm, err := contract.NewPayment(common.HexToAddress(addr), conn)
	if err != nil {
		logrus.WithField("err-msg", err).
			Errorln("failed to instantiate storageObj contract")
	}
	p.Payment = pm
}

func (p *payment) Pay(payer string, payee string, amount uint64) error {
	PrivateKey, _ := crypto.HexToECDSA(prvKey)

	// Create an authorized transactor and call the store function
	nonce, _ := p.conn.NonceAt(context.Background(), common.HexToAddress("你私钥对应的账户地址"), nil)
	gasPrice, _ := p.conn.SuggestGasPrice(context.Background())
	// id
	auth, err := bind.NewKeyedTransactorWithChainID(PrivateKey, big.NewInt(5))
	auth.GasLimit = uint64(300000)
	auth.Nonce = new(big.Int).SetUint64(nonce)
	auth.GasPrice = gasPrice
	if err != nil {
		log.Fatalf("Failed to create authorized transactor: %v", err)
	}
	// Call the store() function
	payerAddr := common.Address([]byte(payer))
	payeeAddr := common.Address([]byte(payee))
	amountInt := new(big.Int).SetUint64(amount)
	tx, err := p.Payment.Pay(auth, payerAddr, payeeAddr, amountInt)
	if err != nil {
		logrus.WithField("err-msg", err).
			Errorln("Failed to update value: %v", err)
		return err
	}
	logrus.Infoln("Update pending: 0x%x\n", tx.Hash())
	return nil
}

// -------------------- data process --------------------------

type data struct{}

var Data = &data{}

const dir = "uploads/"

func (d *data) Add(para *parameter.DataAddReq) (*parameter.GeneralOperationResp, error) {
	bs, err := os.ReadFile(dir + para.File)
	if err != nil {
		return nil, &merror.Error{Code: merror.ServerInternalError, Desc: err.Error()}
	}

	hash := sha256.New()
	hash.Write(bs)
	sum := hash.Sum(nil)
	hexStr := hex.EncodeToString(sum)
	dt := time.Now().Format(time.DateTime)
	_ = storageObj.Store(hexStr, dt)

	// insert into database---------------
	// base info
	bs, _ = json.Marshal(&para.Mark)
	base := &dao.DataInfo{
		Hash:     hexStr,
		Name:     para.Name,
		Addr:     para.Addr,
		Mark:     string(bs),
		File:     para.File,
		Memo:     para.Memo,
		Category: para.Category,
		CreateAt: dt,
	}
	_ = dao.Data.Insert(base)
	// label
	for _, mark := range para.Mark {
		label := &dao.LabelInfo{
			Name: mark,
			Memo: mark,
		}
		_ = dao.Label.Insert(label)
	}
	// data label
	di := dao.Data.GetByHash(hexStr)
	labels := dao.Label.GetByNameList(para.Mark)
	for _, label := range labels {
		info := &dao.DataLabelInfo{
			LabelId: label.ID,
			DataId:  di.ID,
		}
		_ = dao.DataLabel.Insert(info)
	}
	return &parameter.GeneralOperationResp{Success: true}, nil
}

func (d *data) Query(para *parameter.DataQueryReq) (*parameter.DataQueryResp, error) {
	labels := dao.Label.GetByNameList(para.Mark)
	ids := make([]uint, 0, len(labels))
	for _, label := range labels {
		ids = append(ids, label.ID)
	}

	dls := dao.DataLabel.GetByIdList(ids, para.Type)
	ids = make([]uint, 0, len(dls))
	for _, dl := range dls {
		ids = append(ids, dl.DataId)
	}

	infos := dao.Data.GetByIdList(ids)
	items := make([]*parameter.DataInfo, 0, len(infos))
	for _, info := range infos {
		items = append(items, d.getDataInfo(info))
	}
	return &parameter.DataQueryResp{DataList: items}, nil
}

func (d *data) Buy(para *parameter.DataBuyReq) (*parameter.GeneralOperationResp, error) {
	info := dao.Data.GetById(para.Id)
	err := paymentObj.Pay(para.Addr, info.Addr, uint64(info.Price))
	if err != nil {
		logrus.WithField("err-msg", err).Errorln("pay failed")
	}

	buyer := dao.User.GetByAddr(para.Addr)
	dt := time.Now().Format(time.DateTime)
	record := &dao.DataUserBuyInfo{
		UserId: 0,
		Addr:   buyer.Addr,
		DataId: info.ID,
		BuyAt:  dt,
	}
	err = dao.DataUserBuy.Insert(record)
	if err != nil {
		return nil, &merror.Error{
			Code: dbError,
			Desc: "insert into database failed",
		}
	}

	return &parameter.GeneralOperationResp{Success: true}, nil
}

func (d *data) getDataInfo(info *dao.DataInfo) *parameter.DataInfo {
	mark := make([]string, 0)
	_ = json.Unmarshal([]byte(info.Mark), &mark)
	ret := &parameter.DataInfo{
		Id:       info.ID,
		Hash:     info.Hash,
		Name:     info.Name,
		Mark:     mark,
		Addr:     info.Addr,
		File:     info.File,
		CreateAt: info.CreateAt,
		Memo:     info.Memo,
		Category: info.Category,
		Price:    info.Price,
	}
	return ret
}
func (d *data) QueryOneSelf(para *parameter.UserDataQueryReq) (*parameter.UserDataQueryResp, error) {
	infos := dao.Data.GetByAddr(para.Addr)
	uploads := make([]*parameter.DataInfo, 0)
	for _, info := range infos {
		uploads = append(uploads, d.getDataInfo(info))
	}

	dubs := dao.DataUserBuy.GetByAddr(addr)
	ids := make([]uint, 0)
	for _, dub := range dubs {
		ids = append(ids, dub.DataId)
	}
	infos = dao.Data.GetByIdList(ids)
	buys := make([]*parameter.DataInfo, 0)
	for _, info := range infos {
		buys = append(buys, d.getDataInfo(info))
	}

	return &parameter.UserDataQueryResp{
		Uploads: uploads,
		Buys:    buys,
	}, nil
}
