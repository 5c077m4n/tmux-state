package layout

import "testing"

func BenchmarkParse(b *testing.B) {
	benchmarks := map[string]string{
		"simple":        "34b0,239x58,0,0[239x47,0,0,1,239x10,0,48,2]",
		"nested":        "b0a8,239x58,0,0[239x44,0,0,3,239x13,0,45{143x13,0,45,4,95x13,144,45,6}]",
		"no_checksum":   "239x58,0,0[239x47,0,0,1,239x10,0,48,2]",
		"deeply_nested": "6fd4,239x58,0,0[239x44,0,0,0,239x13,0,45[143x13,0,45,9,143x13,0,45[143x13,0,45,8,143x13,0,45{143x13,0,45,4,95x13,144,45,5}]]]",
	}

	for name, layoutString := range benchmarks {
		b.Run(name, func(b *testing.B) {
			for b.Loop() {
				_, err := Parse(layoutString)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}
