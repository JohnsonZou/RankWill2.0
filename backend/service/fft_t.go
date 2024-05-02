package service

import (
	"fmt"
	"github.com/mjibson/go-dsp/fft"
	"golang.org/x/exp/rand"
	"math"
)

const (
	EXPAND_SIZE      = 1
	MAX_RATING       = 4000 * EXPAND_SIZE
	RATING_PRECISION = 400 * EXPAND_SIZE
)

func preCalcConvolution(oldRating []float64) []float64 {
	f := make([]float64, 2*MAX_RATING+1)
	for i := -MAX_RATING; i <= MAX_RATING; i++ {
		f[i+MAX_RATING] = 1 / (1 + math.Pow10((i)/(RATING_PRECISION)))
	}

	g := make([]float64, 2*MAX_RATING+1)
	for _, r := range oldRating {
		idx := int(math.Round(r * EXPAND_SIZE))
		g[idx+MAX_RATING]++
	}

	convolution := fftConvolve(f, g)

	for i := 5500; i < 5800; i += 1 {
		fmt.Println(i, convolution[i])
	}
	return convolution
}

// getExpectedRank 基于预先计算的卷积值获取期望排名。
func getExpectedRank(convolution []float64, x int) float64 {
	return convolution[x+MAX_RATING] + 0.5
}

// getEquationLeft 获取基于预先计算的卷积值的方程左侧。
func getEquationLeft(convolution []float64, x int) float64 {
	return convolution[x+MAX_RATING] + 1
}

// binarySearchExpectedRating 执行二分查找以找到给定平均排名的期望评分。
func binarySearchExpectedRating(convolution []float64, meanRank float64) int {
	lo, hi := 0, MAX_RATING
	for lo < hi {
		mid := (lo + hi) / 2
		if getEquationLeft(convolution, mid) < meanRank {
			hi = mid
		} else {
			lo = mid + 1
		}
	}
	return lo
}

// getExpectedRating 基于当前排名、评分和预先计算的卷积计算期望评分。
func getExpectedRating(rank int, rating float64, convolution []float64) float64 {
	expectedRank := getExpectedRank(convolution, int(math.Round(rating*float64(EXPAND_SIZE))))
	meanRank := math.Sqrt(expectedRank * float64(rank))
	return float64(binarySearchExpectedRating(convolution, meanRank)) / float64(EXPAND_SIZE)
}

// fftDelta 使用快速傅里叶变换（FFT）计算Elo评分变化。
func fftDelta(ranks []int, ratings []float64) []float64 {
	convolution := preCalcConvolution(ratings)
	expectedRatings := make([]float64, len(ranks))
	for i, rank := range ranks {
		rating := ratings[i]
		expectedRatings[i] = getExpectedRating(rank, rating, convolution)
	}
	return expectedRatings
}

// fftConvolve performs convolution of two signals using FFT.
func fftConvolve(f, g []float64) []float64 {
	// Determine the size of the result.
	n := len(f) + len(g) - 1
	// Pad the input signals to the next power of two for FFT.
	m := nextPowerOfTwo(n)
	// Create the complex arrays for FFT input and output.
	fFFT := make([]complex128, m)
	gFFT := make([]complex128, m)
	for i, val := range f {
		fFFT[i] = complex(val, 0)
	}
	for i, val := range g {
		gFFT[i] = complex(val, 0)
	}
	// Perform the FFTs.
	fft.FFT(fFFT)
	fft.FFT(gFFT)
	// Multiply the FFTs.
	for i := range fFFT {
		fFFT[i] *= gFFT[i]
	}
	// Perform the inverse FFT.
	fft.IFFT(fFFT)
	// Convert the result back to a float64 slice.
	result := make([]float64, n)
	for i := 0; i < n; i++ {
		result[i] = real(fFFT[i])
	}
	return result
}

// nextPowerOfTwo returns the next power of two that is greater than or equal to n.
func nextPowerOfTwo(n int) int {
	if n&(n-1) == 0 {
		return n
	}
	return 1 << (uint(math.Log2(float64(n))) + 1)
}

func TestFft() {
	// 示例数据

	ranks := []int{}
	ratings := []float64{}

	for i := 1; i <= 10000; i++ {
		ranks = append(ranks, i)
		ratings = append(ratings, float64(rand.Intn(3000)))
	}

	// 计算delta ratings
	deltaRatings := fftDelta(ranks, ratings)
	fmt.Println("exp Ratings:", len(deltaRatings))
}
