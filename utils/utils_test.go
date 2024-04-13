package utils

import (
	"testing"

	. "github.com/onsi/gomega"
)

func TestFormatFloat_Basic(t *testing.T) {
	we := NewGomegaWithT(t)

	num := 0.251
	resultStr := PrettyFormatFloat(num, -1)
	we.Expect(resultStr).To(Equal("0.251"))
	resultStr = PrettyFormatFloat(num, 2)
	we.Expect(resultStr).To(Equal("0.25"))

	num = 3.00
	resultStr = PrettyFormatFloat(num, -1)
	we.Expect(resultStr).To(Equal("3"))

	num = 3.10
	resultStr = PrettyFormatFloat(num, -1)
	we.Expect(resultStr).To(Equal("3.1"))

	num = 5.98765000
	resultStr = PrettyFormatFloat(num, -1)
	we.Expect(resultStr).To(Equal("5.98765"))
	resultStr = PrettyFormatFloat(num, 2)
	we.Expect(resultStr).To(Equal("5.99"))
	resultStr = PrettyFormatFloat(num, 3)
	we.Expect(resultStr).To(Equal("5.988"))
}
