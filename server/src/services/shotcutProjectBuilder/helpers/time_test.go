package helpers

import (
	"time"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("ShotcutProjectBuilder: Utils / ShotcutFormatDuration", func() {
	It("should format duration correctly", func() {
		Expect(ShotcutFormatDuration(0)).To(Equal("00:00:00.000"))
		Expect(ShotcutFormatDuration(500 * time.Millisecond)).To(Equal("00:00:00.500"))
		Expect(ShotcutFormatDuration(1 * time.Second)).To(Equal("00:00:01.000"))
		Expect(ShotcutFormatDuration(1 * time.Minute)).To(Equal("00:01:00.000"))
		Expect(ShotcutFormatDuration(599960 * time.Millisecond)).To(Equal("00:09:59.960"))
		Expect(ShotcutFormatDuration(1 * time.Hour)).To(Equal("01:00:00.000"))
		Expect(ShotcutFormatDuration(2 * time.Hour)).To(Equal("02:00:00.000"))
		Expect(ShotcutFormatDuration(3661000 * time.Millisecond)).To(Equal("01:01:01.000"))
	})
})
