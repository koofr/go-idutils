package idutils

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

var _ = Describe("Generator", func() {
	Describe("NewGenerator", func() {
		It("should create new generator", func() {
			g, err := NewGenerator(0, 0)
			Expect(err).NotTo(HaveOccurred())

			id, err := g.NextId()
			Expect(err).NotTo(HaveOccurred())
			Expect(id > 0).To(BeTrue())
		})

		It("should return an error if worker id is too big", func() {
			_, err := NewGenerator(32, 0)
			Expect(err).To(HaveOccurred())
		})

		It("should return an error if worker id is negative", func() {
			_, err := NewGenerator(-1, 0)
			Expect(err).To(HaveOccurred())
		})

		It("should return an error if datacenter id is too big", func() {
			_, err := NewGenerator(0, 32)
			Expect(err).To(HaveOccurred())
		})

		It("should return an error if datacenter id is negative", func() {
			_, err := NewGenerator(0, -1)
			Expect(err).To(HaveOccurred())
		})
	})

	Describe("NextId", func() {
		It("should generate id", func() {
			g, err := NewGenerator(7, 7)
			Expect(err).NotTo(HaveOccurred())

			g.lastTimestamp = 1450772535000
			g.sequence = 2794

			g.timeGen = func() int64 {
				return 1450772535000
			}

			id, err := g.NextId()
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(Equal(int64(679215357097835243)))
			// 0000 10010110110100001110110001001100010111 | 00111 | 00111 | 101011101011
		})

		It("should return an error if clock moved backwards", func() {
			g, err := NewGenerator(7, 7)
			Expect(err).NotTo(HaveOccurred())

			g.lastTimestamp = 1450772535000

			g.timeGen = func() int64 {
				return 1450772534000
			}

			_, err = g.NextId()
			Expect(err).To(HaveOccurred())
		})

		It("should wait for new time if sequence overflows", func() {
			g, err := NewGenerator(7, 7)
			Expect(err).NotTo(HaveOccurred())

			g.lastTimestamp = 1450772534000
			g.sequence = 4095

			i := 0

			g.timeGen = func() int64 {
				if i == 0 || i == 1 {
					i += 1
					return 1450772534000
				} else {
					i += 1
					return 1450772535000
				}
			}

			id, err := g.NextId()
			Expect(err).NotTo(HaveOccurred())
			Expect(id).To(Equal(int64(679215357097832448)))
			// 0000 10010110110100001110110001001100010111 | 00111 | 00111 | 000000000000

			Expect(i).To(Equal(3))
		})
	})
})

func BenchmarkNextId(b *testing.B) {
	g, _ := NewGenerator(0, 0)

	for n := 0; n < b.N; n++ {
		g.NextId()
	}
}
