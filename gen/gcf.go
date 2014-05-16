package main

func gcf(xs ...int) int {

	allfactors := [][]int{}

	for _, x := range xs {
		factors := []int{}
		p, s := primes()

		for x > 1 {
			prime := <-p

			for x%prime == 0 {
				factors = append(factors, prime)
				x /= prime
			}
		}

		allfactors = append(allfactors, factors)
		close(s)
	}

	// common := intersect(allfactors...)
	common := allfactors[0]
	for _, set := range allfactors[1:] {
		common = intersect(common, set)
	}

	product := 1
	for _, x := range common {
		product *= x
	}
	return product
}

func intersect(a, b []int) []int {
	apos, bpos := 0, 0
	var res []int

	for apos < len(a) && bpos < len(b) {
		x := a[apos]
		y := b[bpos]
		if x == y {
			res = append(res, x)
			apos++
			bpos++
		} else {
			if x < y {
				for apos < len(a) && a[apos] < y {
					apos++
				}
			} else {
				for bpos < len(b) && b[bpos] < x {
					bpos++
				}
			}
		}
	}
	return res
}

func primes() (<-chan int, chan<- bool) {
	var out = make(chan int)
	var stop = make(chan bool)

	go func() {

		var primes = []int{}
		x := 2
		for {
			select {
			case out <- x:
				primes = append(primes, x)
			L:
				for {
					x++
					for _, prime := range primes {
						if x%prime == 0 {
							continue L
						}
					}
					break L
				}
			case <-stop:
				return
			}
		}
	}()

	return out, stop
}
