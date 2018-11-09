package main

type Recorder struct {
	v []float64
}

func (r *Recorder) Record(in <-chan float64) <-chan float64 {
	out := make(chan float64)
	go func() {
		for v := range in {
			r.v = append(r.v, v)
			out <- v
		}
		close(out)
	}()
	return out
}
