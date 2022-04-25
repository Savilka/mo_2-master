package main

import (
	"encoding/json"
	"fmt"
	"math"
	"os"
)

type Input struct {
	X0     []float64 `json:"x0"`
	Dx     float64   `json:"dx"`
	Eps    float64   `json:"eps"`
	Lambda float64   `json:"lambda"`
	Number int       `json:"number"`
}

var numberOfFunc = 3

func searchInterval(x []float64, delta, lambda float64, x0 []float64) []float64 {
	k := make([]float64, len(x0)+1)
	var h float64
	k[0] = x0[0] + lambda*(x[0]-x0[0])
	k[1] = x0[1] + lambda*(x[1]-x0[1])
	if calculateFunction(k[0], k[1], numberOfFunc) > calculateFunction(x0[0]+(lambda+delta)*(x[0]-x0[0]), x0[1]+(lambda+delta)*(x[1]-x0[1]), numberOfFunc) {
		k[0] = lambda
		k[1] = lambda + delta
		h = delta

	} else {
		k[0] = lambda
		k[1] = lambda - delta
		h = -delta
	}
	h *= 2
	k[2] = k[0] + h
	for calculateFunction(x0[0]+k[1]*(x[0]-x0[0]), x0[1]+k[1]*(x[1]-x0[1]), numberOfFunc) > calculateFunction(x0[0]+k[2]*(x[0]-x0[0]), x0[1]+k[2]*(x[1]-x0[1]), numberOfFunc) {
		k[0] = k[1]
		k[1] = k[2]
		h *= 2
		k[2] = k[1] + h
	}

	return k[1:3]
}
func searchIntervalDescent(delta, lambda float64, x0 []float64, s []float64) ([]float64, int) {
	k := make([]float64, len(x0)+1)
	var iter = 0
	var h float64
	k[0] = x0[0] + lambda*s[0]
	k[1] = x0[1] + lambda*s[1]
	//var flag float64
	if calculateFunction(k[0], k[1], numberOfFunc) > calculateFunction(x0[0]+(lambda+delta)*s[0], x0[1]+(lambda+delta)*s[1], numberOfFunc) {
		k[0] = lambda
		k[1] = lambda + delta
		h = delta
		//flag = 1

	} else {
		k[0] = lambda
		k[1] = lambda - delta
		h = -delta
		//flag = -1
	}
	h *= 2
	k[2] = k[0] + h
	for calculateFunction(x0[0]+k[1]*s[0], x0[1]+k[1]*s[1], numberOfFunc) > calculateFunction(x0[0]+k[2]*s[0], x0[1]+k[2]*s[1], numberOfFunc) {
		k[0] = k[1]
		k[1] = k[2]
		//h = h + flag*delta
		h *= 2
		k[2] = k[1] + h
		iter += 1
	}

	return k[1:3], iter
}

func norm(x []float64) []float64 {
	n := make([]float64, 2)
	x0 := x[0]
	x1 := x[1]
	n[0] = x[0] / math.Sqrt(x0*x0+x1*x1)
	n[1] = x[1] / math.Sqrt(x0*x0+x1*x1)
	return n
}

func calculateGrad(x, y float64, n int) []float64 {
	g := make([]float64, 2)

	switch n {
	case 1:
		g[0] = 202*x - 200*y - 2
		g[1] = 200*y - 200*x
		return g
	case 2:
		g[0] = (-1)*400*x*(y-x*x) + 2*x - 2
		g[1] = 200*y - 200*x*x
		return g
	case 3:
		g[0] = (2*(x-1)*math.Exp(((-1)*(x-1)*(x-1))/9-((y-3)*(y-3))/9))/3 + (x-1)*
			math.Exp(((-1)*(x-1)*(x-1))/4-((y-1)*(y-1))/4)

		g[1] = (2*(y-3)*math.Exp((-1)*(x-1)*(x-1)/9-(y-3)*(y-3)/9))/3 + (y-1)*
			math.Exp((-1)*(x-1)*(x-1)/4-(y-1)*(y-1)/4)

		return g
	}
	return g
}
func calculateFunction(x, y float64, number int) float64 {
	switch number {
	case 1:
		return 100*(y-x)*(y-x) + (1-x)*(1-x)
	case 2:
		return 100*(y-x*x)*(y-x*x) + (1-x)*(1-x)
	case 3:
		var d, b float64
		d = 2 * math.Exp((-1)*(((x-1)/2)*((x-1)/2))-((y-1)/2)*((y-1)/2))
		b = 3 * math.Exp((-1)*(((x-1)/3)*((x-1)/3))-((y-3)/3)*((y-3)/3))
		return (-1) * (d + b)

	}
	return 0
}

func search(dx float64, x0 Input) ([]float64, int) {
	var flag int
	x := make([]float64, len(x0.X0))
	copy(x, x0.X0)
	for i := 0; i < len(x0.X0); i++ {
		cx := make([]float64, len(x))
		copy(cx, x)
		cx[i] += dx
		if calculateFunction(cx[0], cx[1], numberOfFunc) <= calculateFunction(x[0], x[1], numberOfFunc) {
			x[i] += dx
			flag = 1
		} else {
			cx[i] -= 2 * dx
			if calculateFunction(cx[0], cx[1], numberOfFunc) <= calculateFunction(x[0], x[1], numberOfFunc) {
				x[i] -= dx
				flag = 1
			}
		}
	}
	if flag == 0 {
		return x0.X0, flag
	}
	return x, flag
}
func GoldenSearchDescent(x []float64, eps float64, x0 []float64, s []float64) float64 {
	var a, b, x1, x2, phi1, phi2, f1, f2 float64
	var iter int
	a = x[0]
	b = x[1]
	phi1 = (3 - math.Sqrt(5.0)) / 2
	phi2 = (math.Sqrt(5.0) - 1) / 2
	x1 = a + (b-a)*phi1
	x2 = a + (b-a)*phi2
	f1 = calculateFunction(x0[0]+(x1)*s[0], x0[1]+(x1)*s[1], numberOfFunc)
	f2 = calculateFunction(x0[0]+(x2)*s[0], x0[1]+(x2)*s[1], numberOfFunc)
	iter += 2
	for (math.Abs(a - b)) > eps {
		if f1 > f2 {
			a = x1
			x1 = x2
			f1 = f2
			x2 = a + phi2*(b-a)
			f2 = calculateFunction(x0[0]+(x2)*s[0], x0[1]+(x2)*s[1], numberOfFunc)
		} else {
			b = x2
			x2 = x1
			f2 = f1
			x1 = a + phi1*(b-a)
			f1 = calculateFunction(x0[0]+(x1)*s[0], x0[1]+(x1)*s[1], numberOfFunc)
		}

		iter++
	}
	return (a + b) / 2
}
func GoldenSearch(x, lam []float64, eps float64, x0 []float64) float64 {
	var a, b, x1, x2, phi1, phi2, f1, f2 float64
	var iter int
	a = lam[0]
	b = lam[1]
	phi1 = (3 - math.Sqrt(5.0)) / 2
	phi2 = (math.Sqrt(5.0) - 1) / 2
	x1 = a + (b-a)*phi1
	x2 = a + (b-a)*phi2
	f1 = calculateFunction(x0[0]+(x1)*(x[0]-x0[0]), x0[1]+(x1)*(x[1]-x0[1]), numberOfFunc)
	f2 = calculateFunction(x0[0]+(x2)*(x[0]-x0[0]), x0[1]+(x2)*(x[1]-x0[1]), numberOfFunc)
	iter += 2
	for (math.Abs(a - b)) > eps {
		if f1 > f2 {
			a = x1
			x1 = x2
			f1 = f2
			x2 = a + phi2*(b-a)
			f2 = calculateFunction(x0[0]+(x2)*(x[0]-x0[0]), x0[1]+(x2)*(x[1]-x0[1]), numberOfFunc)
		} else {
			b = x2
			x2 = x1
			f2 = f1
			x1 = a + phi1*(b-a)
			f1 = calculateFunction(x0[0]+(x1)*(x[0]-x0[0]), x0[1]+(x1)*(x[1]-x0[1]), numberOfFunc)
		}

		iter++
	}
	return (a + b) / 2
}
func main() {
	var iter = 0
	inputRaw, err := os.ReadFile("input.json")
	if err != nil {
		return
	}
	var input Input
	err = json.Unmarshal(inputRaw, &input)
	if err != nil {
		return
	}

	switch input.Number {
	case 1:
		for input.Dx >= input.Eps {
			x, flag := search(input.Dx, input)
			if flag == 0 {
				input.Dx /= 10
				continue
			}
			lam := searchInterval(x, input.Dx, input.Lambda, input.X0)
			input.Lambda = GoldenSearch(x, lam, input.Eps, input.X0)
			input.X0[0] = input.X0[0] + input.Lambda*(x[0]-input.X0[0])
			input.X0[1] = input.X0[1] + input.Lambda*(x[1]-input.X0[1])
			iter++
			fmt.Println(input.X0, iter)
		}
	case 2:
		g := calculateGrad(input.X0[0], input.X0[1], numberOfFunc)
		x0prev := 1.
		x1prev := 1.
		for math.Sqrt(x0prev*x0prev+x1prev*x1prev) > input.Eps {
			s := norm(g)
			lam, _ := searchIntervalDescent(input.Dx, input.Lambda, input.X0, s)
			input.Lambda = GoldenSearchDescent(lam, input.Eps, input.X0, s)
			x0prev = input.X0[0]
			x1prev = input.X0[1]
			input.X0[0] = input.X0[0] + input.Lambda*s[0]
			input.X0[1] = input.X0[1] + input.Lambda*s[1]
			x0prev -= input.X0[0]
			x1prev -= input.X0[1]
			g = calculateGrad(input.X0[0], input.X0[1], numberOfFunc)

			fmt.Println(input.X0, iter)
			iter++
		}
	}
}
