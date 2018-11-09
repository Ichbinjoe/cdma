package main

func scatter(b []byte) <-chan bool {
	out := make(chan bool)
	go func() {
		for _, v := range b {
			for i := 7; i >= 0; i-- {
				if v>>uint(i)&0x1 == 0x1 {
					out <- true
				} else {
					out <- false
				}
			}
		}
		close(out)
	}()
	return out
}

func gather(in <-chan bool) (ret []byte) {
	var v byte
	for {
		v = 0
		for i := 0; i < 8; i++ {
			bit, ok := <-in
			if !ok {
				return
			}
			v <<= 1
			if bit {
				v |= 0x1
			}
		}
		ret = append(ret, v)
	}
}
