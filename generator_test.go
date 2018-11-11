package idutils

import (
	"testing"
	"time"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Generator", func() {
	new := func() *Generator {
		g, err := NewGenerator(7, 7)
		Expect(err).NotTo(HaveOccurred())
		return g
	}

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
			g := new()

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
			g := new()

			g.lastTimestamp = 1450772535000

			g.timeGen = func() int64 {
				return 1450772534000
			}

			_, err := g.NextId()
			Expect(err).To(HaveOccurred())
		})

		It("should wait for new time if sequence overflows", func() {
			g := new()

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

	Describe("IdToTimestamp", func() {
		It("should convert id to timestamp", func() {
			g := new()

			timestampBefore := time.Now().UnixNano() / 1000000

			id, err := g.NextId()
			Expect(err).NotTo(HaveOccurred())

			timestampAfter := time.Now().UnixNano() / 1000000

			timestamp := IdToTimestamp(id)
			Expect(timestamp).To(BeNumerically(">=", timestampBefore))
			Expect(timestamp).To(BeNumerically("<=", timestampAfter))
		})
	})

	Describe("IdToTime", func() {
		It("should convert id to time", func() {
			g := new()

			timeBefore := time.Now()

			time.Sleep(2 * time.Millisecond)

			id, err := g.NextId()
			Expect(err).NotTo(HaveOccurred())

			time.Sleep(2 * time.Millisecond)

			timeAfter := time.Now()

			idTime := IdToTime(id)

			Expect(timeBefore.Before(idTime)).To(BeTrue())
			Expect(timeAfter.After(idTime)).To(BeTrue())
		})
	})

	Describe("IdEndOfTimestamp", func() {
		It("should convert timestamp to the last possible id", func() {
			timestamp := int64(1541883369255)

			id := IdEndOfTimestamp(timestamp)
			Expect(id).To(Equal(int64(1061361893660164095)))

			newTimestamp := IdToTimestamp(id)
			Expect(newTimestamp).To(Equal(timestamp))

			afterTimestamp := IdToTimestamp(id + 1)
			Expect(afterTimestamp).To(Equal(timestamp + 1))
		})
	})

	Describe("IdEndOfTime", func() {
		It("should convert time to the last possible id", func() {
			t := time.Date(2018, 11, 10, 20, 56, 9, 255000000, time.UTC)

			id := IdEndOfTime(t)
			Expect(id).To(Equal(int64(1061361893660164095)))

			newT := IdToTime(id)
			Expect(newT).To(Equal(t))
		})
	})

	Describe("IdStartOfTimestamp", func() {
		It("should convert timestamp to the first possible id", func() {
			timestamp := int64(1541883369255)

			id := IdStartOfTimestamp(timestamp)
			Expect(id).To(Equal(int64(1061361893655969792)))

			newTimestamp := IdToTimestamp(id)
			Expect(newTimestamp).To(Equal(timestamp))

			beforeTimestamp := IdToTimestamp(id - 1)
			Expect(beforeTimestamp).To(Equal(timestamp - 1))
		})
	})

	Describe("IdStartOfTime", func() {
		It("should convert time to the first possible id", func() {
			t := time.Date(2018, 11, 10, 20, 56, 9, 255000000, time.UTC)

			id := IdStartOfTime(t)
			Expect(id).To(Equal(int64(1061361893655969792)))

			newT := IdToTime(id)
			Expect(newT).To(Equal(t))
		})
	})

	Describe("IdAddDuration", func() {
		It("should add a time Duration to the id", func() {
			id := int64(1061361893655969792)
			t := IdToTime(id)
			d := 1 * time.Hour

			newId := IdAddDuration(id, d)
			Expect(newId).To(Equal(int64(1061376993150369792)))

			newT := IdToTime(newId)
			Expect(newT).To(Equal(t.Add(d)))
		})
	})
})

func BenchmarkNextId(b *testing.B) {
	g, _ := NewGenerator(0, 0)

	for n := 0; n < b.N; n++ {
		g.NextId()
	}
}
