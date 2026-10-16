package main

import (
	"fmt"
	"net"
	"os"
	"os/exec"
	"strings"
)

func main() {
	CreateLogTable()
	var choice string

	for {
		clearScreen()
		fmt.Println("\033[1;35m🛡️< СИСТЕМНЫЙ СТРАЖ IVORY >🛡️\033[0m")
		fmt.Println("Life     - Пульс системы ♥︎")
		fmt.Println("HW       - Информация о процессоре 🔲")
		fmt.Println("Clean    - Чистка логов 🧹")
		fmt.Println("TP       - Менеджер процессов 🗃️")
		fmt.Println("InfoD    - Информация о дисках 💽")
		fmt.Println("NetworkR - Проверить пинг 📡")
		fmt.Println("NetInfo  - Сетевые интерфейсы 🌐")
		fmt.Println("ShowL    - Показать историю логов 📋")
		fmt.Println("ConnBD   - Проверить подключение к PostgreSQL 🗄️")
		fmt.Println("Exit     - Выход ⭕")
		fmt.Println("\033[1;35m<==============================>\033[0m")

		fmt.Print(">")
		fmt.Scanln(&choice)

		switch choice {
		case "Exit":
			fmt.Println("🏮 Завершение работы программы...")
			return

		case "Life":
			fmt.Println("Запуск модуля \033[34m>Пульс системы<\033[0m")
			showPulse()

		case "NetworkR":
			fmt.Println("Запуск модуля \033[34m>Сетевой радар<\033[0m")
			checkPing()

		case "NetInfo":
			fmt.Println("Запуск модуля \033[34m>Сетевые интерфейсы<\033[0m")
			showNetworkInfo()

		case "ShowL":
			fmt.Println("Запуск модуля \033[34m>История логов из базы данных<\033[0m")
			PrintLogsFromDB()

		case "ConnBD":
			fmt.Println("Запуск модуля \033[34m>Подключение к базе данных<\033[0m")
			ConnectDB()

		case "Clean":
			fmt.Println("Запуск модуля \033[34m>Чистка логов<\033[0m")
			cleanLogs()

		case "InfoD":
			fmt.Println("Запуск модуля \033[34m>Информация о дисках<\033[0m")
			showDisks()

		case "TP":
			fmt.Println("Запуск модуля \033[34m>Менеджер процессов<\033[0m")
			manageProcesses()

		case "HW":
			fmt.Println("Запуск модуля \033[34m>Информация о процессоре<\033[0m")
			showCPU()

		default:
			fmt.Println("\033[31mНеизвестная команда ❓\033[0m")
		}
	}
}

func showPulse() {
	data, err := os.ReadFile("/proc/uptime")
	if err != nil {
		fmt.Println("\033[31mОшибка чтения системных данных!\033[0m")
		return
	}

	var totalSeconds float64

	fmt.Sscanf(string(data), "%f", &totalSeconds)

	hours := int(totalSeconds) / 3600
	minutes := (int(totalSeconds) % 3600) / 60

	memData, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		fmt.Println("\033[31mОшибка чтения данных памяти!\033[0m")
		return
	}

	var memTotal, memAvailable int64

	fmt.Sscanf(string(memData), "MemTotal: %d kB\nMemFree: %d kB\nMemAvailable: %d kB", &memTotal, &memAvailable, &memAvailable)

	totalGB := float64(memTotal) / 1024 / 1024
	availableGB := float64(memAvailable) / 1024 / 1024
	usedGB := totalGB - availableGB

	fmt.Println("\n\033[1;36m--- РЕЗУЛЬТАТЫ МОНИТОРИНГА ---\033[0m")
	fmt.Printf("⏱️  Время работы:  %d ч. %d мин.\n", hours, minutes)
	fmt.Printf("Память (RAM): %.2f ГБ свободно / %.2f ГБ всего (Занято: %.2f ГБ)\n\n", availableGB, totalGB, usedGB)
	waitForKey()
}

func checkPing() {
	cmd := exec.Command("ping", "-c", "10", "8.8.8.8")

	fmt.Println("\n📡 Сканирование сети... Пожалуйста, подождите...")

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("\033[31m❌ Ошибка: Сеть полностью недоступна или сервер не отвечает\033[0m\n")
		InsertLogIntoDB("КРИТИЧЕСКИЙ СБОЙ: Сеть полностью недоступна")
		waitForKey()
		return
	}

	fmt.Print(string(output))

	var min, avg, max, mdev float64
	outputStr := string(output)

	for _, line := range strings.Split(outputStr, "\n") {
		if strings.HasPrefix(line, "rtt") {
			fmt.Sscanf(line, "rtt min/avg/max/mdev = %f/%f/%f/%f", &min, &avg, &max, &mdev)
		}
	}
	fmt.Println("\n\033[1;36m--- АНАЛИЗ СЕТЕВОГО РАДАРА ---\033[0m")
	fmt.Printf("📊 Средняя задержка (Ping): %.1f мс\n", avg)

	if avg > 300 {
		fmt.Printf("\033[31m🔴 Всё плохо! Сеть сильно лагает (Пинг: %.1f мс). Возможны обрывы!\033[0m\n\n", avg)
		InsertLogIntoDB(fmt.Sprintf("Внимание: Сеть сильно лагает! Высокий пинг: %.1f мс", avg))
	} else if avg > 120 {
		fmt.Printf("\033[33m🟡 Есть проблемы! Задержка повышена (Пинг: %.1f мс).\033[0m\n\n", avg)
	} else if avg == 0 {
		fmt.Println("\033[31m❌ Не удалось рассчитать статистику задержки.\033[0m\n")
	} else {
		fmt.Printf("\033[32m✅ Интернет работает стабильно! Задержка отличная (Пинг: %.1f мс).\033[0m\n\n", avg)
	}
	waitForKey()
}

func cleanLogs() {
	var confirm string
	fmt.Print("\033[33mВы уверены, что хотите очистить лог (y/n): \033[0m")
	fmt.Scanln(&confirm)

	if confirm == "y" || confirm == "Y" {
		ClearLogTable()
	} else {
		fmt.Println("Очистка отменена. Логи в безопасности.\n")
		waitForKey()
	}
}

func showDisks() {
	cmd := exec.Command("df", "-h", "--total")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("\n\033[1;36m--- СОСТОЯНИЕ НАКОПИТЕЛЕЙ ---\033[0m")

	err := cmd.Run()
	if err != nil {
		fmt.Println("\033[31m❌ Ошибка при получении информации о дисках!\033[0m\n")
	}
	fmt.Println()
	waitForKey()
}

func manageProcesses() {
	cmd := exec.Command("bash", "-c", "ps -eo pid,cmd,%mem --sort=-%mem | head -n 15")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	fmt.Println("\n\033[1;36m--- ТОП-15 НАГРУЖАЮЩИХ ПРОЦЕССОВ (ПО RAM) ---\033[0m")

	err := cmd.Run()
	if err != nil {
		fmt.Println("\033[31m❌ Ошибка при получении списка процессов!\033[0m\n")
	}
	fmt.Println()

	var pid string
	fmt.Print("\033[33mВведите PID процесса для завершения (или 0 для отмены): \033[0m")
	fmt.Scanln(&pid)

	if pid == "0" || pid == "" {
		fmt.Println("Возврат в главное меню...\n")
		return
	}
	killCmd := exec.Command("kill", pid)
	err = killCmd.Run()

	if err != nil {
		fmt.Println("\033[31m❌ Не удалось завершить процесс. Возможно, не хватает прав или PID не существует!\033[0m\n")
	} else {
		fmt.Printf("\033[32m✅ Процесс с PID %s успешно завершен!\033[0m\n\n", pid)
	}
	waitForKey()
}

func clearScreen() {
	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func waitForKey() {
	var input string
	fmt.Print("\033[33mEnter - вернуться в меню\033[0m")
	fmt.Scanln(&input)
}

func showCPU() {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		fmt.Println("\033[31m❌ Ошибка чтения данных процессора\033[0m")
		waitForKey()
		return
	}

	fmt.Println("\n\033[1;36m--- ИНФОРМАЦИЯ О ПРОЦЕССОРЕ ---\033[0m")

	cpuText := string(data)

	coresCount := 0
	modelName := ""
	cpuMHz := ""

	for _, line := range strings.Split(cpuText, "\n") {

		if strings.HasPrefix(line, "processor") {
			coresCount++
		}

		if strings.HasPrefix(line, "model name") && modelName == "" {
			modelName = line
		}

		if strings.HasPrefix(line, "cpu MHz") && cpuMHz == "" {
			cpuMHz = line
		}
	}

	if modelName != "" {
		fmt.Println("🔲", modelName)
	}

	fmt.Printf("🔳 Количество ядер (потоков): %d\n", coresCount)

	if cpuMHz != "" {
		fmt.Println("⚡", cpuMHz)
	}

	fmt.Println()
	waitForKey()
}

func showNetworkInfo() {
	interfaces, err := net.Interfaces()
	if err != nil {
		fmt.Println("\033[31m❌ Ошибка получения сетевых интерфейсов\033[0m")
		waitForKey()
		return
	}

	fmt.Println("\n\033[1;36m--- СЕТЕВЫЕ ИНТЕРФЕЙСЫ ---\033[0m")

	for _, iface := range interfaces {
		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {

			parts := strings.Split(addr.String(), "/")
			ipAddress := parts[0]

			if strings.Contains(ipAddress, ":") {
				continue
			}

			fmt.Printf("🌐 %s -> %s\n", iface.Name, ipAddress)

		}
	}
	fmt.Println()
	waitForKey()
}
