package zk

import (
	"github.com/consensys/gnark/frontend"
	"math/big"
)

func Sigma(api frontend.API, in frontend.Variable) frontend.Variable {
	in2 := api.Mul(in, in)
	in4 := api.Mul(in2, in2)
	return api.Mul(in4, in)
}

func Ark(api frontend.API, in []frontend.Variable, c []*big.Int, r int) []frontend.Variable {
	out := make([]frontend.Variable, len(in))
	for i, v := range in {
		out[i] = api.Add(v, c[i+r])
	}
	return out
}

func Mix(api frontend.API, in []frontend.Variable, m [][]*big.Int) []frontend.Variable {
	t := len(in)
	out := make([]frontend.Variable, t)
	for i := 0; i < t; i++ {
		lc := frontend.Variable(0)
		for j := 0; j < t; j++ {
			lc = api.Add(lc, api.Mul(m[j][i], in[j]))
		}
		out[i] = lc
	}
	return out
}

func MixLast(api frontend.API, in []frontend.Variable, m [][]*big.Int, s int) frontend.Variable {
	t := len(in)
	out := frontend.Variable(0)
	for j := 0; j < t; j++ {
		out = api.Add(out, api.Mul(m[j][s], in[j]))
	}
	return out
}

func MixS(api frontend.API, in []frontend.Variable, s []*big.Int, r int) []frontend.Variable {
	t := len(in)
	out := make([]frontend.Variable, t)
	lc := frontend.Variable(0)
	for i := 0; i < t; i++ {
		lc = api.Add(lc, api.Mul(s[(t*2-1)*r+i], in[i]))
	}
	out[0] = lc
	for i := 1; i < t; i++ {
		out[i] = api.Add(in[i], api.Mul(in[0], s[(t*2-1)*r+t+i-1]))
	}
	return out
}

func PoseidonEx(api frontend.API, inputs []frontend.Variable, initialState frontend.Variable, nOuts int) []frontend.Variable {
	nInputs := len(inputs)
	out := make([]frontend.Variable, nOuts)
	nRoundsPC := [16]int{56, 57, 56, 60, 60, 63, 64, 63, 60, 66, 60, 65, 70, 60, 64, 68}
	t := nInputs + 1
	nRoundsF := 8
	nRoundsP := nRoundsPC[t-2]
	c := POSEIDON_C(t)
	s := POSEIDON_S(t)
	m := POSEIDON_M(t)
	p := POSEIDON_P(t)

	state := make([]frontend.Variable, t)
	for j := 0; j < t; j++ {
		if j == 0 {
			state[0] = initialState
		} else {
			state[j] = inputs[j-1]
		}
	}
	state = Ark(api, state, c, 0)

	for r := 0; r < nRoundsF/2-1; r++ {
		for j := 0; j < t; j++ {
			state[j] = Sigma(api, state[j])
		}
		state = Ark(api, state, c, (r+1)*t)
		state = Mix(api, state, m)
	}

	for j := 0; j < t; j++ {
		state[j] = Sigma(api, state[j])
	}
	state = Ark(api, state, c, nRoundsF/2*t)
	state = Mix(api, state, p)

	for r := 0; r < nRoundsP; r++ {

		state[0] = Sigma(api, state[0])

		state[0] = api.Add(state[0], c[(nRoundsF/2+1)*t+r])
		newState0 := frontend.Variable(0)
		for j := 0; j < len(state); j++ {
			mul := api.Mul(s[(t*2-1)*r+j], state[j])
			newState0 = api.Add(newState0, mul)
		}

		for k := 1; k < t; k++ {
			state[k] = api.Add(state[k], api.Mul(state[0], s[(t*2-1)*r+t+k-1]))
		}
		state[0] = newState0
	}

	for r := 0; r < nRoundsF/2-1; r++ {
		for j := 0; j < t; j++ {
			state[j] = Sigma(api, state[j])
		}
		state = Ark(api, state, c, (nRoundsF/2+1)*t+nRoundsP+r*t)
		state = Mix(api, state, m)
	}

	for j := 0; j < t; j++ {
		state[j] = Sigma(api, state[j])
	}

	for i := 0; i < nOuts; i++ {
		out[i] = MixLast(api, state, m, i)
	}
	return out
}

func Poseidon(api frontend.API, inputs []frontend.Variable) frontend.Variable {
	out := PoseidonEx(api, inputs, 0, 1)
	return out[0]
}
