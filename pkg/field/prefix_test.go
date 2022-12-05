package field

import (
	"fmt"
	"testing"
)

var prefixResult AbstractFields

func BenchmarkPrefix(b *testing.B) {
	for _, fieldCount := range []int{0, 1, 10, 100, 1000} {
		b.Run(fmt.Sprintf("fieldCount%d", fieldCount), func(b *testing.B) {
			fields := dummyFields(uint(fieldCount))
			b.Run("init", func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					prefixResult = Prefix("somePrefix_", fields)
				}
			})
			b.Run("ForEachField", func(b *testing.B) {
				fields := Prefix("somePrefix_", fields)
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					fields.ForEachField(func(f *Field) bool {
						return true
					})
				}
			})
		})
	}
}
