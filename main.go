package main

import (
	"encoding/json"
	"fmt"

	"math"
	"os"
)

// needs to be redone
type Input struct {
	X0     []float64 `json:"x0"`
	Dx     float64   `json:"dx"`
	Eps    float64   `json:"eps"`
	Lambda float64   `json:"lambda"`
	Number int       `json:"number"`
	R      float64   `json:"penalty"`
}

// checking the condition of penalt functions
func penalty(x, y, eps float64, number int, r float64) bool {
	switch number {
	case 1:
		if 2*math.Abs(Functiong(x, y)) <= 1e-6 {
			return true
		} else {
			return false
		}
	case 2:
		if 2*math.Abs(Functionh(x, y)) <= 1e-6 {
			return true
		} else {
			return false
		}
	default:
		return true
	}

}

var numberOfFunc = 1 // mb dont need

func Functiong(x, y float64) float64 { //first function
	return -x - y + 5
}
func FunctionG(x, y float64) float64 {
	return math.Pow((Functiong(x, y)+math.Abs(Functiong(x, y)))/2, 2)
}
func Functionh(x, y float64) float64 { // second function
	return y - x - 2
}
func FunctionH(x, y float64) float64 {
	return math.Pow(math.Abs(Functionh(x, y)), 1)
}

//given function
func calculateFunction(x, y float64) float64 {

	return math.Pow(x+y, 2) + 4*math.Pow(y, 2)
}

// calculate Q(x,r)
func penaltyFunction(x, y, r float64, number int) float64 {
	switch number {
	case 1:
		return calculateFunction(x, y) + r*FunctionG(x, y)

	case 2:
		return calculateFunction(x, y) + r*FunctionH(x, y)

	default:
		return 0
	}
}

func GoldenSearch(x, lam []float64, eps, r float64, x0 []float64, iter *int64) float64 {
	var a, b, x1, x2, phi1, phi2, f1, f2 float64

	a = lam[0]
	b = lam[1]
	phi1 = (3 - math.Sqrt(5.0)) / 2
	phi2 = (math.Sqrt(5.0) - 1) / 2
	x1 = a + (b-a)*phi1
	x2 = a + (b-a)*phi2
	f1 = penaltyFunction(x0[0]+(x1)*(x[0]-x0[0]), x0[1]+(x1)*(x[1]-x0[1]), r, numberOfFunc)
	f2 = penaltyFunction(x0[0]+(x2)*(x[0]-x0[0]), x0[1]+(x2)*(x[1]-x0[1]), r, numberOfFunc)
	*iter += 2
	for (math.Abs(a - b)) > eps {
		if f1 > f2 {
			a = x1
			x1 = x2
			f1 = f2
			x2 = a + phi2*(b-a)
			f2 = penaltyFunction(x0[0]+(x2)*(x[0]-x0[0]), x0[1]+(x2)*(x[1]-x0[1]), r, numberOfFunc)
			*iter++
		} else {
			b = x2
			x2 = x1
			f2 = f1
			x1 = a + phi1*(b-a)
			f1 = penaltyFunction(x0[0]+(x1)*(x[0]-x0[0]), x0[1]+(x1)*(x[1]-x0[1]), r, numberOfFunc)
			*iter++
		}

	}
	return (a + b) / 2
}

func searchInterval(x []float64, delta, lambda, r float64, x0 []float64, iter *int64) []float64 {
	k := make([]float64, len(x0)+1)
	var h float64
	k[0] = x0[0] + lambda*(x[0]-x0[0])
	k[1] = x0[1] + lambda*(x[1]-x0[1])
	if penaltyFunction(k[0], k[1], r, numberOfFunc) > penaltyFunction(x0[0]+(lambda+delta)*(x[0]-x0[0]), x0[1]+(lambda+delta)*(x[1]-x0[1]), r, numberOfFunc) {
		k[0] = lambda
		k[1] = lambda + delta
		h = delta
		*iter += 2

	} else {
		k[0] = lambda
		k[1] = lambda - delta
		h = -delta
		*iter += 2
	}
	h *= 2
	k[2] = k[0] + h
	for penaltyFunction(x0[0]+k[1]*(x[0]-x0[0]), x0[1]+k[1]*(x[1]-x0[1]), r, numberOfFunc) > penaltyFunction(x0[0]+k[2]*(x[0]-x0[0]), x0[1]+k[2]*(x[1]-x0[1]), r, numberOfFunc) {
		k[0] = k[1]
		k[1] = k[2]
		h *= 2
		k[2] = k[1] + h
		*iter += 2
	}
	//k[1] = k[2]
	return k[1:3]
}

func search(dx float64, x0 Input, iter *int64) ([]float64, int) {
	var flag int
	x := make([]float64, len(x0.X0))
	copy(x, x0.X0)
	for i := 0; i < len(x0.X0); i++ {
		cx := make([]float64, len(x))
		copy(cx, x)
		cx[i] += dx
		if penaltyFunction(cx[0], cx[1], x0.R, numberOfFunc) <= penaltyFunction(x[0], x[1], x0.R, numberOfFunc) {
			x[i] += dx
			flag = 1
			*iter += 2
		} else {
			cx[i] -= 2 * dx
			if penaltyFunction(cx[0], cx[1], x0.R, numberOfFunc) <= penaltyFunction(x[0], x[1], x0.R, numberOfFunc) {
				x[i] -= dx
				flag = 1
				*iter += 2
			}
		}
	}
	if flag == 0 {
		return x0.X0, flag
	}
	return x, flag
}
func main() {
	var iter = 0
	var iterFunction int64

	iterFunction = 0
	inputRaw, err := os.ReadFile("input.json")
	if err != nil {
		return
	}
	var input Input
	err = json.Unmarshal(inputRaw, &input)
	if err != nil {
		return
	}
	var step = input.Dx
	for !penalty(input.X0[0], input.X0[1], input.Eps, numberOfFunc, input.R) && input.R < 10000000000 {
		input.Dx = step
		input.R *= 2
		for input.Dx >= input.Eps {

			x, flag := search(input.Dx, input, &iterFunction)
			if flag == 0 {
				input.Dx /= 10
				continue
			}
			lam := searchInterval(x, input.Dx, input.Lambda, input.R, input.X0, &iterFunction)
			input.Lambda = GoldenSearch(x, lam, input.Eps, input.R, input.X0, &iterFunction)
			input.X0[0] = input.X0[0] + input.Lambda*(x[0]-input.X0[0])
			input.X0[1] = input.X0[1] + input.Lambda*(x[1]-input.X0[1])

			iter++
			fmt.Println(input.X0, penaltyFunction(input.X0[0], input.X0[1], input.R, numberOfFunc), input.R)
		}
	}
	fmt.Println(input.X0[0], input.X0[1])
}
