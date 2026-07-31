package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// Karakter yang akan ditampilkan di layar
const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789@#$%^&*()_+-=[]{}|;:,.<>?"

// Kode warna ANSI untuk terminal
const (
	ColorReset     = "\033[0m"
	ColorBrightGreen = "\033[1;32m"
	ColorDimGreen    = "\033[0;32m"
	ColorWhite       = "\033[1;37m"
	ClearScreen      = "\033[2J\033[H"
)

func main() {
	// Membersihkan layar terminal awal
	fmt.Print(ClearScreen)
	
	// Menangani tombol Ctrl+C agar keluar dengan rapi
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		fmt.Print(ColorReset + ClearScreen)
		fmt.Println("System Offline. Safe travels, Hacker!")
		os.Exit(0)
	}()

	// Menentukan jumlah kolom (sesuai lebar terminal standar)
	columns := 80
	rand.Seed(time.Now().UnixNano())

	// Channel untuk mengirim karakter dari goroutine
	ch := make(chan struct {
		col  int
		char byte
		bold bool
	})

	// Jalankan goroutine untuk setiap kolom
	for i := 0; i < columns; i++ {
		go func(col int) {
			for {
				// Durasi acak agar kecepatan setiap kolom berbeda
				time.Sleep(time.Duration(rand.Intn(100)+30) * time.Millisecond)
				
				char := chars[rand.Intn(len(chars))]
				bold := rand.Float32() > 0.8 // 20% peluang karakter bersinar terang

				ch <- struct {
					col  int
					char byte
					bold bool
				}{col: col, char: char, bold: bold}
			}
		}(i)
	}

	// Loop utama untuk merender karakter ke terminal
	for {
		data := <-ch
		
		// Set posisi kursor acak untuk simulasi hujan
		row := rand.Intn(30) + 1
		
		// Pindah kursor ANSI ke \033[Row;ColumnH
		fmt.Printf("\033[%d;%dH", row, data.col+1)

		// Efek warna
		if data.bold {
			fmt.Printf("%s%c%s", ColorWhite, data.char, ColorReset)
		} else {
			fmt.Printf("%s%c%s", ColorBrightGreen, data.char, ColorReset)
		}
	}
}
