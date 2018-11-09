package main

type Bitstream struct {
	stream []uint8
	l      uint
}

func NewBitstream() *Bitstream {
	return &Bitstream{
		stream: nil,
		l:      0,
	}
}

func (b *Bitstream) Append(v bool) {
	b.l++
	c := b.l >> 3
	if int(c) >= len(b.stream) {
		b.stream = append(b.stream, 0)
	}
	if v {
		b.stream[c] |= 0x1 << (b.l & 0x7)
	}
}

func (b *Bitstream) Stream() <-chan bool {
	c := make(chan bool)
	go func() {
		for i := uint(0); i < b.l; i++ {
			if (b.stream[i>>3]>>(i&0x7))&0x1 == 0x1 {
				c <- true
			} else {
				c <- false
			}
		}
		close(c)
	}()
	return c
}

func Code(in <-chan bool, s *Bitstream) <-chan float64 {
	out := make(chan float64)
	go func() {
		for v := range in {
			// take each bit in our bitstream and multiply it with our bit
			stream := s.Stream()
			for b := range stream {
				if v == b {
					out <- 1.0
				} else {
					out <- -1.0
				}
			}
		}
		close(out)
	}()
	return out
}
