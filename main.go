package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"sync"
)

var streams int

func parseStream(s string) (*Bitstream, bool) {
	b := NewBitstream()
	for _, r := range s {
		if r == '+' {
			b.Append(true)
		} else if r == '-' {
			b.Append(false)
		} else {
			fmt.Printf("Invalid stream character %s\n", r)
			return nil, false
		}
	}
	return b, true
}

func main() {
	code_suggestions := []string{"++++++++", "+-+-+-+-", "++++----"}

	flag.IntVar(&streams, "streams", 3, "-streams <streamcount>")
	flag.Parse()

	s := bufio.NewScanner(os.Stdin)

	bitstreams := make([]*Bitstream, 0, streams)
	signals := make([][]byte, 0, streams)

	for i := 0; i < streams; i++ {
		if i < len(code_suggestions) {
			for {
				fmt.Printf("Code %d (%s):", i, code_suggestions[i])
				has := s.Scan()
				if !has {
					return // something bad happened, so bomb.
				}

				t := s.Text()

				if t == "" {
					t = code_suggestions[i]
				}
				b, ok := parseStream(t)
				if ok {
					bitstreams = append(bitstreams, b)
					break
				}
			}
		} else {
			for {
				fmt.Printf("Code %d:", i)
				has := s.Scan()
				if !has {
					return // something bad happened, so bomb.
				}
				t := s.Text()
				if t == "" {
					fmt.Println("No input w/ no default... please enter a bitstream.")
					continue
				}
				b, ok := parseStream(t)
				if ok {
					bitstreams = append(bitstreams, b)
					break
				}
			}
		}
	}

	for i := 0; i < streams; i++ {
		for {
			fmt.Printf("Stream %d:", i)
			has := s.Scan()
			if !has {
				return
			}

			t := s.Text()
			if t == "" {
				fmt.Println("No input w/ no default... please enter a stream.")
				continue
			}

			signals = append(signals, []byte(t))
			break
		}
	}

	decoderHandles := make([]chan float64, 0, streams)
	m := Mixer{make([]<-chan float64, 0, streams)}

	for i := 0; i < streams; i++ {
		m.In(Code(scatter(signals[i]), bitstreams[i]))
		decoderHandles = append(decoderHandles, make(chan float64))
	}

	recorder := &Recorder{make([]float64, 0, 1000)}
	out := recorder.Record(m.Mix())
	go func(out <-chan float64) {
		for v := range out {
			for _, c := range decoderHandles {
				c <- v
			}
		}

		for _, c := range decoderHandles {
			close(c)
		}
	}(out)

	var wg sync.WaitGroup

	decoded := make([][]byte, streams, streams)
	for i := 0; i < streams; i++ {
		wg.Add(1)
		go func(i int) {
			decoded[i] = gather(Decode(decoderHandles[i], bitstreams[i]))
			wg.Done()
		}(i)
	}

	wg.Wait()

	fmt.Println("CDMA Stream:")
	for _, f := range recorder.v {
		fmt.Printf("%v,", f)
	}

	fmt.Println()

	for i, d := range decoded {
		fmt.Printf("Decoded stream %d: %s\n", i, string(d))
	}
}
