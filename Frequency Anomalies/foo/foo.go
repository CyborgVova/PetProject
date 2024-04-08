package foo

import (
	"fmt"
	"log"
	"math"
	"os"
)

func SaveData(num int32, expVal, stdDiv float64) {
	file, err := os.OpenFile("report.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("Unable to open file:", err)
		os.Exit(1)
	}
	defer file.Close()
	file.WriteString(fmt.Sprintf("Обработано: %d Mean: %f Sd: %f\n", num, expVal, stdDiv))
}

func Job(i int32, k float64, arrf []float64) (mean, sd float64) {
	i++
	mean = GetMean(arrf)
	sd = GetSD(arrf)
	if i%30 == 0 {
		SaveData(i, mean, sd)
	}
	return
}

func Anomalies(freq, mean, rng float64) bool {
	if freq < mean-rng || freq > mean+rng {
		log.Println("Аномалия: ", freq)
		return true
	}
	return false
}

func GetMean(arr []float64) (res float64) {
	var out float64 = 0
	for _, v := range arr {
		out += float64(v)
	}
	return out / float64(len(arr))
}

func GetSD(arr []float64) float64 {
	if len(arr) == 1 {
		return 0
	}
	var out float64
	avg := GetMean(arr)
	for _, v := range arr {
		out += math.Pow((float64(v) - avg), 2)
	}
	out = out / float64(len(arr)-1)
	return math.Sqrt(out)
}
