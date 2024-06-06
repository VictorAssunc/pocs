package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"sync"
	"time"
)

var wg sync.WaitGroup

func main() {
	for _, version := range []string{"v1", "v2", "v3"} {
		fmt.Printf("\n[INFO] Testing server script %s\n", version)

		if err := exec.Command("go", "build", "-o", version, version+".go").Run(); err != nil {
			fmt.Printf("[ERROR] failed to build server script %s: %v\n", version, err)
			return
		}

		for _, duration := range []string{"5", "10", "15"} {
			fmt.Printf("\n[INFO] Testing duration %s\n", duration)
			cmd := exec.Command("./" + version)
			if err := cmd.Start(); err != nil {
				fmt.Printf("[ERROR] failed to run %s: %v\n", version, err)
				return
			}

			wg.Add(1)
			c := make(chan struct{}, 1)
			time.Sleep(500 * time.Millisecond)
			go func() {
				defer func() {
					c <- struct{}{}
				}()

				url := "http://localhost:8888/" + duration
				// fmt.Printf("[INFO] Sending GET request to %s\n", url)

				wg.Done()
				res, err := http.Get(url)
				if err != nil {
					fmt.Printf("[ERROR] failed to get response: %v\n", err)
					return
				}
				defer res.Body.Close()

				body, err := io.ReadAll(res.Body)
				if err != nil {
					fmt.Printf("[ERROR] failed to read response body: %v\n", err)
					return
				}

				fmt.Printf("[INFO] Response: %s\n", body)
			}()

			wg.Wait()
			time.Sleep(500 * time.Millisecond)
			// fmt.Println("[INFO] Interrupting server...")
			start := time.Now()
			if err := cmd.Process.Signal(os.Interrupt); err != nil {
				fmt.Printf("[ERROR] failed to interrupt server: %v\n", err)
				return
			}

			if err := cmd.Wait(); err != nil {
				fmt.Printf("[ERROR] server exited with error: %s\n", err)
				return
			}

			// fmt.Println("[INFO] Server exited successfully")
			fmt.Printf("[INFO] Duration since interrupt signal: %s\n", time.Since(start))
			<-c
		}
	}
}
