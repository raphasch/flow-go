package proto_test

import (
	"testing"

	. "github.com/onsi/gomega"

	"github.com/dapperlabs/flow-go/pkg/types/proto"
	"github.com/dapperlabs/flow-go/pkg/utils/unittest"
)

func TestAccountSignature(t *testing.T) {
	RegisterTestingT(t)

	sigA := unittest.AccountSignatureFixture()

	message := proto.AccountSignatureToMessage(sigA)
	sigB := proto.MessageToAccountSignature(message)

	Expect(sigA).To(Equal(sigB))
}

func TestTransaction(t *testing.T) {
	RegisterTestingT(t)

	txA := unittest.TransactionFixture()

	message := proto.TransactionToMessage(txA)

	txB, err := proto.MessageToTransaction(message)
	Expect(err).ToNot(HaveOccurred())

	Expect(txA).To(Equal(txB))
}

func TestAccount(t *testing.T) {
	RegisterTestingT(t)

	accA := unittest.AccountFixture()

	message := proto.AccountToMessage(accA)

	accB, err := proto.MessageToAccount(message)
	Expect(err).ToNot(HaveOccurred())

	Expect(accA).To(Equal(accB))
}

func TestAccountKey(t *testing.T) {
	RegisterTestingT(t)

	keyA := unittest.AccountKeyFixture()

	message := proto.AccountKeyToMessage(keyA)

	keyB, err := proto.MessageToAccountKey(message)
	Expect(err).ToNot(HaveOccurred())

	Expect(keyA).To(Equal(keyB))
}
