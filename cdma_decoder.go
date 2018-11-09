package main

func Decode(in <-chan float64, b *Bitstream) <-chan bool {
	out := make(chan bool)
	go func() {
		var sum float64
		for {
			sum = 0
			bits := b.Stream()
			for bit := range bits {
				observation, ok := <-in

				if !ok {
					close(out)
					return
				}

				var bitExtrapolate float64
				if bit {
					bitExtrapolate = 1.0
				} else {
					bitExtrapolate = -1.0
				}

				product := bitExtrapolate * observation
				sum += product
			}

			if sum > 0 {
				out <- true
			} else {
				out <- false
			}
		}
	}()
	return out
}
