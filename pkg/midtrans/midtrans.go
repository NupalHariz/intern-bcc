package midtrans

import (
	"errors"
	"intern-bcc/domain"
	"os"

	"github.com/midtrans/midtrans-go"
	"github.com/midtrans/midtrans-go/coreapi"
)

type IMidTrans interface {
	ChargeTransaction(newTransaction domain.Transactions) (*coreapi.ChargeResponse, error)
	VerifyPayment(transctionIdString string) (bool, error)
}

type MidTrans struct {
	serverKey string
}

func MidTransInit() IMidTrans {
	serverKey := os.Getenv("MIDTRANS_SERVER_KEY")

	return &MidTrans{
		serverKey: serverKey,
	}
}

func (m *MidTrans) ChargeTransaction(newTransaction domain.Transactions) (*coreapi.ChargeResponse, error) {
	c := coreapi.Client{}
	c.New(m.serverKey, midtrans.Sandbox)

	chargeReq := &coreapi.ChargeReq{
		TransactionDetails: midtrans.TransactionDetails{
			OrderID:  newTransaction.Id.String(),
			GrossAmt: int64(newTransaction.Price),
		},
	}

	switch newTransaction.PaymentType {
	case "gopay":
		chargeReq.PaymentType = coreapi.PaymentTypeGopay
	case "mandiri":
		chargeReq.PaymentType = coreapi.PaymentTypeEChannel
		chargeReq.EChannel = &coreapi.EChannelDetail{
			BillInfo1: "Payment",
			BillInfo2: "Online purchase",
		}
	case "bca":
		chargeReq.PaymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: midtrans.BankBca,
		}
	case "bni":
		chargeReq.PaymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: midtrans.BankBni,
		}
	case "bri":
		chargeReq.PaymentType = coreapi.PaymentTypeBankTransfer
		chargeReq.BankTransfer = &coreapi.BankTransferDetails{
			Bank: midtrans.BankBri,
		}
	}

	coreApiRes, err := c.ChargeTransaction(chargeReq)
	if err != nil {
		return coreApiRes, err
	}
	return coreApiRes, nil
}

func (m *MidTrans) VerifyPayment(transctionIdString string) (bool, error) {
	c := coreapi.Client{}
	c.New(m.serverKey, midtrans.Sandbox)

	transactionStatusRespone, err := c.CheckTransaction(transctionIdString)
	if err != nil {
		return false, err
	}

	switch transactionStatusRespone.TransactionStatus {
	case "settlement":
		return true, nil
	case "expire":
		return false, errors.New("payment already expired")
	case "failure":
		return false, errors.New("an error occured during transaction process")
	}

	return false, nil
}
