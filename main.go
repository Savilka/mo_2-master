package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type Input struct {
	X0     []float32 `json:"x0"`
	Dx     float32   `json:"dx"`
	Eps    float32   `json:"eps"`
	Lambda float32   `json:"lambda"`
}

func searchInterval(x []float32, delta, lambda float32, x0 []float32) []float32 {
	k := make([]float32, len(x0)+1)
	var h float32
	k[0] = x0[0] + lambda*(x[0]-x0[0])
	k[1] = x0[1] + lambda*(x[1]-x0[1])
	if calculateFunction(k[0], k[1], 1) > calculateFunction(x0[0]+(lambda+delta)*(x[0]-x0[0]), x0[1]+(lambda+delta)*(x[1]-x0[1]), 1) {
		k[0] = lambda
		k[1] = lambda + delta
		h = delta

	} else {
		k[0] = lambda
		k[1] = lambda - delta
		h = -2 * delta
	}
	h *= 2
	k[2] = k[0] + h
	for calculateFunction(x0[0]+k[1]*(x[0]-x0[0]), x0[1]+k[1]*(x[1]-x0[1]), 1) > calculateFunction(x0[0]+k[2]*(x[0]-x0[0]), x0[1]+k[2]*(x[1]-x0[1]), 1) {
		k[0] = k[1]
		k[1] = k[2]
		h *= 2
		k[2] = k[1] + h
	}

	return k[1:3]
}

func calculateFunction(x, y float32, number int) float32 {
	switch number {
	case 1:
		return 100*(y-x)*(y-x) + (1-y)*(1-y)
	case 2:

	case 3:

	}
	return 0
}

func search(dx float32, x0 Input) ([]float32, int) {
	var flag int
	x := make([]float32, len(x0.X0))

	copy(x, x0.X0)
	for i := 0; i < len(x0.X0); i++ {
		cx := make([]float32, len(x))
		copy(cx, x)
		cx[i] += dx
		if calculateFunction(cx[0], cx[1], 1) <= calculateFunction(x[0], x[1], 1) {
			x[i] += dx
			flag = 1
		} else {
			cx[i] -= 2 * dx
			if calculateFunction(cx[0], cx[1], 1) <= calculateFunction(x[0], x[1], 1) {
				x[i] += dx
				flag = 1
			}
		}
	}
	if flag == 0 {
		return x0.X0, flag
	}
	return x, flag
}
func GoldenSearch(x []float32) float32 {
	var lamda float32

	return lamda
}
func main() {
	inputRaw, err := os.ReadFile("input.json")
	if err != nil {
		return
	}
	var input Input
	err = json.Unmarshal(inputRaw, &input)
	if err != nil {
		return
	}

	for input.Dx >= input.Eps {
		x, flag := search(input.Dx, input)
		if flag == 0 {
			input.Dx /= 10

		}
		lam := searchInterval(x, input.Dx, input.Lambda, input.X0)
		fmt.Println(lam)
		copy(input.X0, x)
	}

}
