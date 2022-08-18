package main

import (
	"fmt"
	"math"
	"strconv"
)

func main () {
	var n int
	var index [3]string
	for i := 0; i < 2^n; i++ {
		// n= 3 then lables are set {000, 001, 010, 011, 100, 101, 110, 111}
		//  depth = 3
		//  step : generate label - n lenth string {0,1}
		//index = math.Logb(n) - 0  -> 000
		//                     - 1  -> 001
		//					    2    -> 00
		//                     - 2  -> 10  -> 010
		//					   - 4  -> 11  --> 011
		/* calculate  {
			check strings from right until find '1'
			change '1' to '0' and delete the edge right bit
			
		*/ 
		//  
		// caculate index patent by E' formula 
		// if index is in tree then calculate  node[i]
		// 		{

		//       }
		// esle add parent to tree
		index = math.Abs(n)[:2]
		index = strconv.FormatComplex(x byte[3])

		
 
	}

}