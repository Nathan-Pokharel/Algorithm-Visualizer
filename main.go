package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
)

var arrayStates [][]int

func shellSort(arr []int) {
	arrayStates = [][]int{}
	n := len(arr)
	gap := n / 2
	for gap > 0 {
		for i := gap; i < n; i++ {
			temp := arr[i]
			j := i
			for j >= gap && arr[j-gap] > temp {
				arr[j] = arr[j-gap]
				j -= gap
			}
			arr[j] = temp
		}
		arrayStates = append(arrayStates, append([]int{}, arr...))
		gap /= 2
	}
}

func selectionSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		minIndex := i
		for j := i + 1; j < n; j++ {
			if arr[j] < arr[minIndex] {
				minIndex = j
			}
		}
		if minIndex != i {
			arr[i], arr[minIndex] = arr[minIndex], arr[i]
		}
		arrayStates = append(arrayStates, append([]int{}, arr...))
	}
}

func MergeSort(arr []int, temp []int, left, right int) {
	if left < right {
		mid := (left + right) / 2
		MergeSort(arr, temp, left, mid)
		MergeSort(arr, temp, mid+1, right)
		merge(arr, temp, left, mid, right)
	}
}

func merge(arr []int, temp []int, left, mid, right int) {
	i := left
	j := mid + 1
	k := left
	for i <= mid && j <= right {
		if arr[i] <= arr[j] {
			temp[k] = arr[i]
			i++
		} else {
			temp[k] = arr[j]
			j++
		}
		k++
	}
	for i <= mid {
		temp[k] = arr[i]
		k++
		i++
	}
	for j <= right {
		temp[k] = arr[j]
		k++
		j++
	}
	for k := left; k <= right; k++ {
		arr[k] = temp[k]
	}
	arrayStates = append(arrayStates, append([]int{}, arr...))
}

func quickSort(arr []int, low, high int) {
	if low < high {
		pivotPos := partition(arr, low, high)
		quickSort(arr, low, pivotPos-1)
		quickSort(arr, pivotPos+1, high)
	}
}

func partition(arr []int, low, high int) int {
	pivot := arr[high]
	i := low - 1
	for j := low; j < high; j++ {
		if arr[j] <= pivot {
			i++
			arr[i], arr[j] = arr[j], arr[i]
		}
	}
	arr[i+1], arr[high] = arr[high], arr[i+1]
	arrayStates = append(arrayStates, append([]int{}, arr...))
	return i + 1
}

func heapSort(arr []int) {
	n := len(arr)
	for i := n/2 - 1; i >= 0; i-- {
		heapify(arr, n, i)
	}
	for i := n - 1; i > 0; i-- {
		arr[0], arr[i] = arr[i], arr[0]
		arrayStates = append(arrayStates, append([]int{}, arr...))
		heapify(arr, i, 0)
	}
}

func heapify(arr []int, n, i int) {
	largest := i
	left := 2*i + 1
	right := 2*i + 2
	if left < n && arr[left] > arr[largest] {
		largest = left
	}

	if right < n && arr[right] > arr[largest] {
		largest = right
	}
    
	if largest != i {
		arr[i], arr[largest] = arr[largest], arr[i]
		heapify(arr, n, largest)
	}
}

func InsertionSort(arr []int) {
	n := len(arr)
	for i := 1; i < n; i++ {
		key := arr[i]
		j := i - 1

		for j >= 0 && arr[j] > key {
			arr[j+1] = arr[j]
			j--
		}
		arr[j+1] = key
		stateCopy := make([]int, len(arr))
		copy(stateCopy, arr)
		arrayStates = append(arrayStates, stateCopy)
	}
}

func BubbleSort(arr []int) {
	n := len(arr)
	for i := 0; i < n-1; i++ {
		swapped := false
		for j := 0; j < n-i-1; j++ {
			if arr[j] > arr[j+1] {
				arr[j], arr[j+1] = arr[j+1], arr[j]
				swapped = true
			}
		}
		stateCopy := make([]int, len(arr))
		copy(stateCopy, arr)
		arrayStates = append(arrayStates, stateCopy)
		if !swapped {
			break
		}
	}
}

func SortHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		var data struct {
			Algorithm string `json:"algorithm"`
			Heights   []int  `json:"heights"`
		}

		err := json.NewDecoder(r.Body).Decode(&data)
		if err != nil {
			http.Error(w, "Error decoding JSON data", http.StatusBadRequest)
			return
		}

		switch data.Algorithm {
		case "bubble":
			arrayStates = [][]int{}
			fmt.Println("Bubble Sort Called")
			BubbleSort(data.Heights)
		case "insertion":
			arrayStates = [][]int{}
			fmt.Println("Insertion Sort Called")
			InsertionSort(data.Heights)
		case "heap":
			arrayStates = [][]int{}
			fmt.Println("Heap Sort Called")
			heapSort(data.Heights)
		case "quick":
			arrayStates = [][]int{}
			fmt.Println("Quick Sort Called")
			quickSort(data.Heights, 0, len(data.Heights)-1)
		case "merge":
			arrayStates = [][]int{}
			fmt.Println("Merge Sort Called")
			MergeSort(data.Heights, make([]int, len(data.Heights)), 0, len(data.Heights)-1)
		case "selection":
			arrayStates = [][]int{}
			fmt.Println("Selection Sort Called")
			selectionSort(data.Heights)
		case "shell":
			arrayStates = [][]int{}
			fmt.Println("Shell Sort Called")
			shellSort(data.Heights)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(arrayStates)
		return
	}
}

func Index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/visualizer.html")
	if err != nil {
		http.Error(w, fmt.Sprintf("Error parsing template: %v", err), http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error executing template: %v", err), http.StatusInternalServerError)
		return
	}

}

func main() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/sort", SortHandler)

	fmt.Printf("Server is running on port %s...\n", "8080")

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		fmt.Printf("Error starting the server: %v\n", err)
	}
}
