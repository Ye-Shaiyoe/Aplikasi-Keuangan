package main

import (
	"log"
	"os"

	"github.com/akrom/finance-backend/internal/database"
	"github.com/akrom/finance-backend/internal/handler"
	"github.com/akrom/finance-backend/internal/middleware"
	"github.com/akrom/finance-backend/internal/service"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()

	if err := database.Connect(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	if err := database.Migrate(); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	// Auto-process any due recurring transactions on startup
	if n, err := service.ProcessDueRecurring(0); err != nil {
		log.Printf("Warning: recurring processing failed: %v", err)
	} else if n > 0 {
		log.Printf("Processed %d due recurring transactions on startup", n)
	}

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"https://umkmkeuangan.vercel.app",
			"http://localhost:5173",
		},
		AllowMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders: []string{"Origin", "Content-Type", "Authorization"},
	}))

	api := r.Group("/api")
	{
		// Public auth routes
		api.POST("/auth/register", handler.Register)
		api.POST("/auth/login", handler.Login)
		api.POST("/auth/google", handler.GoogleLogin)

		// Protected routes (require auth)
		protected := api.Group("")
		protected.Use(middleware.AuthRequired())
		{
			protected.GET("/auth/me", handler.Me)

			protected.GET("/categories", handler.GetCategories)
			protected.GET("/categories/:id", handler.GetCategory)
			protected.POST("/categories", handler.CreateCategory)
			protected.PUT("/categories/:id", handler.UpdateCategory)
			protected.DELETE("/categories/:id", handler.DeleteCategory)

			protected.GET("/transactions", handler.GetTransactions)
			protected.GET("/transactions/:id", handler.GetTransaction)
			protected.POST("/transactions", handler.CreateTransaction)
			protected.PUT("/transactions/:id", handler.UpdateTransaction)
			protected.DELETE("/transactions/:id", handler.DeleteTransaction)

			protected.GET("/reports/summary", handler.GetSummary)
			protected.GET("/reports/yearly-trend", handler.GetYearlyTrend)
			protected.GET("/reports/category-trend", handler.GetCategoryTrend)
			protected.GET("/reports/analytics", handler.GetAdvancedAnalytics)

			// Savings goals
			protected.GET("/savings-goals", handler.GetSavingsGoals)
			protected.GET("/savings-goals/:id", handler.GetSavingsGoal)
			protected.POST("/savings-goals", handler.CreateSavingsGoal)
			protected.PUT("/savings-goals/:id", handler.UpdateSavingsGoal)
			protected.DELETE("/savings-goals/:id", handler.DeleteSavingsGoal)
			protected.POST("/savings-goals/:id/deposit", handler.DepositToSavingsGoal)
			protected.POST("/savings-goals/:id/withdraw", handler.WithdrawFromSavingsGoal)

			// Budgets
			protected.GET("/budgets", handler.GetBudgets)
			protected.GET("/budgets/summary", handler.GetBudgetSummary)
			protected.POST("/budgets", handler.UpsertBudget)
			protected.DELETE("/budgets/:id", handler.DeleteBudget)

			// Recurring transactions
			protected.GET("/recurring", handler.GetRecurring)
			protected.GET("/recurring/:id", handler.GetRecurringByID)
			protected.POST("/recurring", handler.CreateRecurring)
			protected.PUT("/recurring/:id", handler.UpdateRecurring)
			protected.DELETE("/recurring/:id", handler.DeleteRecurring)
			protected.POST("/recurring/process", handler.ProcessRecurringNow)

			// Machine Learning
			protected.POST("/ml/predict-category", handler.MLPredictCategory)
			protected.POST("/ml/forecast", handler.MLForecast)
			protected.POST("/ml/train", handler.MLTrainModel)
			protected.GET("/ml/health", handler.MLHealth)
		}
	}

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}



//loooooopppppp
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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
}package main

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
