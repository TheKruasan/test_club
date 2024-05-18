package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
)

func isOpen(arive, start, end int) bool {
	//Создадим переменную которая покажет работает ли наш клуб в ночь (например 23:00 - 03:00)
	is_night := false
	if start > end {
		is_night = true
	}
	if is_night {
		if arive >= start || arive < end {
			return true
		}
		return false
	}
	if arive >= start && arive < end {
		return true
	}
	return false
}

type Computer struct {
	id           int
	time_in_work int
	value        int
	inWork       bool
}

func main() {
	//считываем название файла с консли
	args := os.Args
	if len(args) < 2 {
		fmt.Println("Usage: go run main.go <filename>\n")
		os.Exit(1)
	}
	fileName := args[1]
	//читаем файл
	file_data, err := ioutil.ReadFile(fileName)
	if err != nil {
		log.Fatal(err)
	}
	//разделяем файл по строкам
	str := string(file_data)
	substrings := strings.Split(str, "\r\n")
	//получаем количество компьютеров в клубе
	number_of_pc, err := strconv.Atoi(substrings[0])
	if err != nil {
		fmt.Print("Enter the correct number of computers\n")
		log.Fatal(err)
	}
	fmt.Println(number_of_pc)

	// получаем время начала работы и конца работы
	start, end := substrings[1][:5], substrings[1][6:]

	//переводим время начала работы в минуты для удобного сравнения
	start_hour, err := strconv.Atoi(start[:2])
	start_minute, e := strconv.Atoi(start[3:])
	if err != nil {
		fmt.Println("Ошибка парсинга часов:\n", err)
		return
	}
	if e != nil {
		fmt.Println("Ошибка парсинга минут:\n", err)
		return
	}
	start_time := start_hour*60 + start_minute

	//переводим время конца работы в минуты для удобного сравнения
	end_hour, err := strconv.Atoi(end[:2])
	end_minute, e := strconv.Atoi(end[3:])
	if err != nil {
		fmt.Println("Ошибка парсинга часов:\n", err)
		return
	}
	if e != nil {
		fmt.Println("Ошибка парсинга минут:\n", err)
		return
	}
	end_time := end_hour*60 + end_minute
	fmt.Println(start_time, "-", end_time)

	//считываем почасовую оплату
	price_str := substrings[2]
	price, err := strconv.Atoi(price_str)
	if err != nil {
		fmt.Println("Ошибка при чтении цены\n", err)
		return
	}
	fmt.Println("Цена ", price)

	fmt.Println(start)
	// Создадим пусты столы в количестве number_of_pc
	computers := []Computer{}
	for i := 0; i < number_of_pc; i++ {
		computers = append(computers, Computer{id: i + 1, time_in_work: 0, value: 0, inWork: false})
	}

	//начинаем обрабатывать записи
	for i := 3; i < len(substrings); i++ {
		queue := map[string]int{}
		note := substrings[i]
		print("\n")
		fmt.Println(note)
		notes := strings.Split(substrings[i], " ")

		//сначала выясним время прихода человека
		time_arrived := notes[0]
		fmt.Print(time_arrived, " ")

		//переводим время прихода в минуты для удобного сравнения
		arrived_hour, err := strconv.Atoi(time_arrived[:2])
		arrived_minute, e := strconv.Atoi(time_arrived[3:])
		if err != nil {
			fmt.Println("Ошибка парсинга часов:\n", err)
			return
		}
		if e != nil {
			fmt.Println("Ошибка парсинга минут:\n", err)
			return
		}
		arrived_time := arrived_hour*60 + arrived_minute

		if notes[1] == "1" {
			//сейчас проверим пришел ли он в часы работы клуба если нет то выводим ожибку с id 13
			if !isOpen(arrived_time, start_time, end_time) {
				fmt.Print(13, " NotOpenYet")
				continue
			}
			// если он уже не в очереди то записываем его туда если нет то выводим ошибку с id 13
			if queue[notes[2]] != 0 {
				fmt.Print(13, " YouShallNotPass")
				continue
			}
			queue[notes[2]] = 1
		}
		if notes[1] == "2" {
			if len(notes) != 4 {
				fmt.Print("Не правильный ввод записи о том как посетитель занимает стол")
				return
			}

		}

	}
	fmt.Print("\n", end)

}
