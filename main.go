package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
)

type Queue struct {
	name string
	next *Queue
}

func minutesInString(min int) string {
	hours := min / 60
	hours_str := fmt.Sprint(hours)
	if hours < 10 {
		hours_str = "0" + hours_str
	}
	minutes := min % 60
	minutes_str := fmt.Sprint(minutes)
	if minutes < 10 {
		minutes_str = "0" + minutes_str
	}

	return hours_str + ":" + minutes_str
}
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
	start        int
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
		computers = append(computers, Computer{id: i + 1, time_in_work: 0, value: 0, start: 0, inWork: false})
	}
	//создадим список посетителей
	visitors := map[string]int{} //0 - нету его в очереди, -1 - ожидает, 1 - за 1 компьютером, 2 - за вторым, 3 - за третьим
	empty_comps := number_of_pc
	// queue := make([]string, number_of_pc)
	inQueue := 0
	root := &Queue{}
	//начинаем обрабатывать записи
	for i := 3; i < len(substrings); i++ {

		note := substrings[i]
		fmt.Println(note)
		notes := strings.Split(substrings[i], " ")

		//сначала выясним время прихода человека
		time_arrived := notes[0]

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
		//Если id входящей записи 1 (клиент пришел)
		if notes[1] == "1" {
			//сейчас проверим пришел ли он в часы работы клуба если нет то выводим ожибку с id 13
			if !isOpen(arrived_time, start_time, end_time) {
				fmt.Print(time_arrived, " ", 13, " NotOpenYet\n")
				continue
			}
			// если он уже не в очереди то записываем его туда если нет то выводим ошибку с id 13
			if visitors[notes[2]] != 0 {
				fmt.Print(time_arrived, " ", 13, " YouShallNotPass\n")
				continue
			}
			visitors[notes[2]] = -1
		}
		//Если id входящей записи 2 (клиент садится(меняет) стол)
		if notes[1] == "2" {
			//если не правильный ввод
			if len(notes) != 4 {
				fmt.Print("Не правильный ввод записи о том как посетитель занимает стол\n")
				return
			}
			//Получаем из запроса номер комьютеру куда хочет сесть человек
			targetComp, err := strconv.Atoi(notes[3])
			if err != nil {
				fmt.Print("Не правильный ввод записи о том как посетитель занимает стол\n")
				return
			}
			// если компьютер на который хочет сесть клиент занят(даже если им самим)
			if computers[targetComp-1].inWork {
				fmt.Print(time_arrived, " ", 13, " PlaceIsBusy\n")
				continue
			}
			//если посетитель уже за компьютером
			if visitors[notes[2]] != -1 {
				//получаем его старый компьютер
				compOfUser := visitors[notes[2]]
				//отключаем его старый компьютер
				computers[compOfUser-1].inWork = false
				empty_comps++
				//получаем время в которое посетитель сел за компьютер
				st := computers[compOfUser-1].start
				//часы за которые посетитель должен заплатить
				hours := 0
				// если время проходит через 24:00 (сел в 23:00 вышел в 2:00)
				if st > arrived_time {
					timeOnComp := 24*60 - st + arrived_time
					computers[targetComp-1].time_in_work += timeOnComp
					hours = int(math.Ceil(float64(timeOnComp) / 60))
				} else { // если не проходит через 24:00
					computers[targetComp-1].time_in_work += arrived_time - st
					hours = int(math.Ceil(float64(arrived_time-st) / 60))
				}
				// добавляем к заработу компьютера цену которую должен заплатить посетитель за старый компьютер
				computers[compOfUser-1].value += hours * price
				computers[compOfUser-1].start = 0

			}
			//получаем комьютер на который хочет сесть клиент
			visitors[notes[2]] = targetComp
			//сажаем клиента
			empty_comps--
			//запускаем компьютер и начинаем отсчет работы
			computers[targetComp-1].inWork = true
			computers[targetComp-1].start = arrived_time
		}
		if notes[1] == "3" {
			//если в очереди, но есть свободные места
			if empty_comps != 0 {
				fmt.Print(time_arrived, " ", 13, " ICanWaitNoLonger!\n")
				continue
			}
			//если в очереди, но очередь слишком большая(больше чем количество компьютеров)
			if inQueue > number_of_pc {
				fmt.Print(time_arrived, " ", 11, notes[2], "\n")
				//клиент ушел из клуба
				visitors[notes[2]] = 0
				continue
			}
			// если все подходит то добавляем клиента в конец очереди
			if inQueue == 0 {
				root = &Queue{name: notes[2], next: nil}
				inQueue++
				continue
			}
			last := root
			for last.next != nil {
				last = last.next
			}
			last.next = &Queue{name: notes[2], next: nil}
			inQueue++

		}
		if notes[1] == "4" {
			// если клиента нету в клубе
			if visitors[notes[2]] == 0 {
				fmt.Print(time_arrived, " ", 13, "ClientUnknown\n")
				continue
			}
			//если клиент сидел за компьютером то надо узнать какой компьютер он занимал
			client_pc_num := visitors[notes[2]]
			visitors[notes[2]] = 0
			//затем освободим компьютер и посчитаем его выручку
			empty_comps++
			computers[client_pc_num-1].inWork = false

			//получаем время в которое посетитель сел за компьютер
			st := computers[client_pc_num-1].start

			//часы за которые посетитель должен заплатить
			hours := 0

			// если время проходит через 24:00 (сел в 23:00 вышел в 2:00)
			if st > arrived_time {
				//считаем время до 24:00 и прибувляем время после 24:00(делим на 60, округляя в большую сторону, так как у нас все в минутах)
				timeOnComp := 24*60 - st + arrived_time
				computers[client_pc_num-1].time_in_work += timeOnComp
				hours = int(math.Ceil(float64(timeOnComp) / 60))
			} else { // если не проходит через 24:00 то просто считаем разность и делим ее на 60, округляя в большую сторону
				computers[client_pc_num-1].time_in_work += arrived_time - st
				hours = int(math.Ceil(float64(arrived_time-st) / 60))
			}

			//добавляем к заработу компьютера цену которую должен заплатить посетитель за компьютер
			computers[client_pc_num-1].value += hours * price
			//обнуляем время с которого у нас работает компьютер
			computers[client_pc_num-1].start = 0
			// затем надо посадить за компьютер 1-ого человека из очереди
			if inQueue != 0 {
				//берем имя первого в очереди человека и переводим ссылку 1 в очереди на следующего в очереди
				name := root.name
				root = root.next
				//сажаем его за освободившийся компьютер
				visitors[name] = client_pc_num
				//выводим сообщение что клиент сел за компьютер
				fmt.Print(time_arrived, " ", 12, " ", name, " ", client_pc_num, "\n")
				//занимаем компьютер
				empty_comps--
				computers[client_pc_num-1].inWork = true
				computers[client_pc_num-1].start = arrived_time

				//уменьшаем длинну очереди
				inQueue--

			}

		}

	}
	stayInClub := []string{}
	for i, v := range visitors {

		if v != 0 {
			stayInClub = append(stayInClub, i)
			if v == -1 {
				continue
			}
			empty_comps++
			computers[v-1].inWork = false
			//получаем время в которое посетитель сел за компьютер
			st := computers[v-1].start
			//часы за которые посетитель должен заплатить
			hours := 0
			// если время проходит через 24:00 (сел в 23:00 вышел в 2:00)
			if st > end_time {
				//считаем время до 24:00 и прибувляем время после 24:00(делим на 60, округляя в большую сторону, так как у нас все в минутах)
				timeOnComp := 24*60 - st + end_time
				computers[v-1].time_in_work += timeOnComp
				hours = int(math.Ceil(float64(timeOnComp) / 60))
			} else { // если не проходит через 24:00 то просто считаем разность и делим ее на 60, округляя в большую сторону
				computers[v-1].time_in_work += end_time - st
				hours = int(math.Ceil(float64(end_time-st) / 60))
			}
			//добавляем к заработу компьютера цену которую должен заплатить посетитель за компьютер
			computers[v-1].value += hours * price
			//обнуляем время с которого у нас работает компьютер
			computers[v-1].start = 0
		}
	}
	for root != nil {
		stayInClub = append(stayInClub, root.name)
		root = root.next
	}
	for _, name := range stayInClub {
		fmt.Print(end, " ", 11, " ", name, "\n")
	}

	fmt.Print(end)
	for _, comp := range computers {
		fmt.Print("\n", comp.id, " ", comp.value, " ", minutesInString(comp.time_in_work))
	}

}
