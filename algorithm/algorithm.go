package algorithm

import (
	"TSPs_Hopfield/usages"
	"fmt"
	"math"
)

type City struct {
	Name string
	X    float64
	Y    float64
}

type TSP struct {
	Cities []City
	N      int // Number of cities
}

func NewTSP(cities []City) *TSP {
	return &TSP{
		Cities: cities,
		N:      len(cities),
	}
}

func (t *TSP) Distance(city1, city2 City) float64 {
	dx := city1.X - city2.X
	dy := city1.Y - city2.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func (t *TSP) GenerateSymmetricWeightMatrix(A, B, C, D float64) [][]float64 {
	// Create a square matrix of size N x N
	weights := make([][]float64, t.N*t.N)
	for i := range weights {
		weights[i] = make([]float64, t.N*t.N)
	}

	// Calculate the weights based on the given equations
	for X := 0; X < t.N; X++ {
		for i := 0; i < t.N; i++ {
			for Y := 0; Y < t.N; Y++ {
				for j := 0; j < t.N; j++ {
					// Kronecker delta function
					δ_ij := 0.0
					if i == j {
						δ_ij = 1.0
					}

					δ_XY := 0.0
					if X == Y {
						δ_XY = 1.0
					}

					δj_ip1 := 0.0
					if j == i+1 {
						δj_ip1 = 1.0
					}

					δj_im1 := 0.0
					if j == i-1 {
						δj_im1 = 1.0
					}

					// Calculate the weight value based on the equations
					weights[X*t.N+i][Y*t.N+j] = -A*δ_XY*(1-δ_ij) - B*δ_ij*(1-δ_XY) - C - 
						D*t.Distance(t.Cities[X], t.Cities[Y])*(δj_ip1+δj_im1)

					// Symmetric property of weight matrix
					weights[Y*t.N+j][X*t.N+i] = weights[X*t.N+i][Y*t.N+j]
				}
			}
		}
	}
	for i := 0; i < t.N*t.N; i++ {
		weights[i][i] = 0.0
	}
	return weights
}



func (t *TSP) HopfieldEnergy(states [][]int, A, B, C, D float64) float64 {
	energyObj := 0.0
	energyConst := 0.0

	// Calculate the objective function part of the energy (E_obj)
	for X := 0; X < t.N; X++ {
		for i := 0; i < t.N; i++ {
			for Y := 0; Y < t.N; Y++ {
				if Y != X {
					energyObj += 0.5 * t.Distance(t.Cities[X], t.Cities[i]) * float64(states[X][Y]) * (float64(states[Y][(i+1)%t.N]) + float64(states[Y][(i-1+t.N)%t.N]))
				}
			}
		}
	}
	// Calculate the constraint part of the energy (E_const)
	// Constraint 1: Each city is visited only once (No more than two neurons in each row outputting 1)
	for X := 0; X < t.N; X++ {
		for i := 0; i < t.N; i++ {
			for j := 0; j < t.N; j++ {
				if j != i {
					energyConst += A * 0.5 * float64(states[X][i]*states[X][j])
				}
			}
		}
	}

	// Constraint 2: Cannot visit two cities at the same time (No more than two neurons in each column outputting 1)
	for i := 0; i < t.N; i++ {
		for X := 0; X < t.N; X++ {
			for Y := 0; Y < t.N; Y++ {
				if Y != X {
					energyConst += B * 0.5 * float64(states[X][i]*states[Y][i])
				}
			}
		}
	}

	// Constraint 3: Visit all cities (Exactly N neurons outputting 1)
	countCities := 0
	for X := 0; X < t.N; X++ {
		for i := 0; i < t.N; i++ {
			countCities += states[X][i]
		}
	}
	energyConst += C * 0.5 * math.Pow(float64(countCities-t.N), 2)

	// Calculate the total energy (E_net)
	energyNet := 0.5 * energyObj + energyConst

	return energyNet
}

func (t *TSP) HopfieldEnergyGeneral(states [][]int, weights [][]float64, A, B, C, D float64) float64 {
	second_order := 0.0
	first_oder := 0.0
	for X := 0; X < t.N; X++ {
		for i := 0; i < t.N; i++ {
			first_oder += float64(states[X][i])
			for Y := 0; Y < t.N; Y++ {
				for j := 0; j < t.N; j++ {
					second_order += weights[X*t.N+i][Y*t.N+j] * float64(states[X][i]) * float64(states[Y][j])
				}
			}
		}
	}
	energy := -0.5*second_order + -C*float64(t.N) * first_oder
	return energy
}

// Dynamic function to update states until convergence
func (t *TSP) HopfieldDynamic(states [][]int, weights [][]float64, A, B, C, D, convergenceThreshold float64) {
    maxIterations := 1000 // Maximum number of iterations to prevent infinite loops

    for iteration := 0; iteration < maxIterations; iteration++ {
        energyBefore := t.HopfieldEnergyGeneral(states, weights, A, B, C, D)

        updatedStates := make([][]int, t.N)
        for i := 0; i < t.N; i++ {
            updatedStates[i] = make([]int, t.N)
        }

		for X := 0; X < t.N; X++ {
			for i := 0; i < t.N; i++ {
				sum := 0.0
				for Y := 0; Y < t.N; Y++ {
					for j := 0; j < t.N; j++ {
						sum += weights[X*t.N+i][Y*t.N+j] * float64(states[Y][j])
					}
				}
				// fmt.Println("The Sum", X, "and", i," is: ", sum, "and the Theta is: ", -C*float64(t.N))
				if sum >= -C*float64(t.N) {
					updatedStates[X][i] = 1
				} else {
					updatedStates[X][i] = 0
				}
			}
		}
		
        energyAfter := t.HopfieldEnergyGeneral(updatedStates, weights, A, B, C, D)
		fmt.Println("EnergyAfter is: ", energyAfter)

        // Check for convergence based on the change in energy
        energyChange := math.Abs(energyAfter - energyBefore)
        if energyChange < convergenceThreshold {
            break
        }

        // Update states for the next iteration
		fmt.Println("States: ",states)
        states = updatedStates
    }
}

func (t* TSP) HopfieldDynamicGeneral (states [][]int, A, B, C, D float64) float64 {
	theta := 0.0
	// Iterate through rows and columns to update neuron states
	for X := 0; X < t.N; X++ {
		for i := 0; i < t.N; i++ {
			// Step 1: Calculate the input to the neuron
			input := 0.0
			for j := 0; j < t.N; j++ {
				if i != j {
					input += A * float64(states[X][j])
				}
			}
			for Y := 0; Y < t.N; Y++ {
				if X != Y {
					input += B * float64(states[Y][i])
					input += t.Distance(t.Cities[X], t.Cities[Y]) * float64(states[Y][(i-1+t.N)%t.N] + states[Y][(i+1)%t.N])
				}
			}
			input += C * (usages.SumAllStates(states) - float64(t.N))

			// Step 2: Apply the threshold function
			if input > theta {
				states[X][i] = 1
			} else {
				states[X][i] = 0
			}
		}
	}
	return 0.0
}


// Helper function to calculate the sum of distances between adjacent cities
func sumAdjacentDistances(t *TSP, states [][]int, X, i int) float64 {
	sum := 0.0
	for Y := 0; Y < t.N; Y++ {
		if Y != X {
			sum += t.Distance(t.Cities[X], t.Cities[Y]) * float64(states[X][i]*(states[Y][(i+1)%t.N]+states[Y][(i-1+t.N)%t.N]))
		}
	}
	return sum
}



func (t *TSP) DecodeSolution(states [][]int) []City {
	// Given the final states, decode the solution to get the tour order
	// The tour order will be the order in which cities are visited
	tourOrder := make([]City, t.N)
	for i := 0; i < t.N; i++ {
		for j := 0; j < t.N; j++ {
			if states[i][j] == 1 {
				tourOrder[j] = t.Cities[i]
				break
			}
		}
	}

	return tourOrder
}

func (t *TSP) CalculateTotalTourLength(tourOrder []City) float64 {
	// Given the tour order, calculate the total tour length
	totalLength := 0.0
	for i := 0; i < t.N; i++ {
		// Calculate the distance between the current city and the next city in the tour order
		fromCity := tourOrder[i]
		toCity := tourOrder[(i+1)%t.N]
		totalLength += t.Distance(fromCity, toCity)
	}

	return totalLength
}

