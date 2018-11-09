package main

func Amplifier(signal <-chan float64, amp float64) <-chan float64 {
	out := make(chan float64)
	go func() {
		s := <-signal
		s *= amp
		out <- s
	}()
	return out
}

type Mixer struct {
	mixers []<-chan float64
}

func (m *Mixer) In(signal <-chan float64) {
	m.mixers = append(m.mixers, signal)
}

func (m *Mixer) Mix() <-chan float64 {
	out := make(chan float64)
	go func() {
		for {
			v := 0.0
			h := false
			for _, m := range m.mixers {
				x, o := <-m
				if o {
					v += x
					h = true
				}
			}
			if !h {
				close(out)
				return
			}
			out <- v
		}
	}()
	return out
}
