package iic

import (
	"fmt"
	"math/rand"
	"os"

	"github.com/irifrance/gini/z"
	"github.com/irifrance/reach"
	"github.com/irifrance/reach/iic/internal/cnf"
)

type sifter struct {
	cnf      *cnf.T
	sat      *satmon
	pri      *reach.Primer
	initVals []int8
	ms, ns   []z.Lit
}

func newSifter(f *cnf.T, sat *satmon, pri *reach.Primer, iv []int8) *sifter {
	return &sifter{cnf: f, sat: sat, pri: pri, initVals: iv}
}

func (s *sifter) sift(dst []z.Lit, c cnf.Id) (toAdd []z.Lit, timeOk bool) {
	s.ms = s.ms[:0]
	s.ms = append(s.ms, s.cnf.Lits(c)...)
	if debugSift {
		fmt.Printf("sifting %v: ", s.ms)
	}
	orgLen := len(s.ms)
	if orgLen == 1 {
		return dst, true
	}
	yMap := make(map[z.Lit]struct{}, len(s.ms))

	for {
		for _, m := range s.ms {
			mp := s.pri.Prime(m)
			s.sat.Assume(mp.Not())
		}
		switch s.sat.Try() {
		case 0:
			return dst, false
		case 1:
			fmt.Printf("orgLen %d now %v\n", orgLen, s.ms)
			for _, m := range s.ms {
				fmt.Printf("prime: %s\n", s.pri.Prime(m))
			}
			s.sat.sat.Write(os.Stdout)
			os.Stdout.Sync()
			panic("wilma!")
			return dst, false
		case -1:
			s.extractNs(yMap)

			if len(s.ns) == orgLen {
				if debugSift {
					fmt.Printf("=> no effect.\n")
				}
				return dst, true
			}
			if len(s.ns) == len(s.ms) {
				if debugSift {
					fmt.Printf("=> %v\n", s.ms)
				}
				dst = append(dst, s.ms...)
				dst = append(dst, 0)
				return dst, true
			}
		}
		s.ms, s.ns = s.ns, s.ms
		rand.Shuffle(len(s.ms), func(i, j int) {
			s.ms[i], s.ms[j] = s.ms[j], s.ms[i]
		})
	}
}

func (s *sifter) extractNs(ym map[z.Lit]struct{}) {
	s.whyMap(ym)
	s.ns = s.ns[:0]
	for _, m := range s.ms {
		mp := s.pri.Prime(m).Not()
		if _, ok := ym[mp]; ok {
			s.ns = append(s.ns, m)
		}
	}
	s.clearMap(ym)
	s.ensureInit()
}

func (s *sifter) whyMap(m map[z.Lit]struct{}) {
	s.ns = s.sat.Why(s.ns[:0])
	for _, n := range s.ns {
		m[n] = struct{}{}
	}
}

func (s *sifter) clearMap(m map[z.Lit]struct{}) {
	for k := range m {
		delete(m, k)
	}
}

func (s *sifter) ensureInit() {
	for _, n := range s.ns {
		if s.initVals[n.Var()] == n.Sign() {
			return
		}
	}
	for _, m := range s.ms {
		if s.initVals[m.Var()] == m.Sign() {
			s.ns = append(s.ns, m)
			return
		}
	}
	fmt.Printf("ns %v\n", s.ns)
	panic("no init")
}
