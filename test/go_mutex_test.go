package test

import (
	"github.com/sirupsen/logrus"
	"sync"
	"testing"
)

func TestGolangMutex(t *testing.T) {
	var (
		wg        = &sync.WaitGroup{}
		mu        = &sync.Mutex{}
		count     = 10000
		totalData int
	)

	t.Run("test tanpa mutex", func(t *testing.T) {
		for i := 0; i < count; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup) {
				defer wg.Done()

				totalData++
				logrus.WithField("i", i+1).Infof("menambah kan 1 data")
			}(wg)
		}

		// wait all goroutines done
		wg.Wait()
		logrus.Infof("jumlah total data %d", totalData)
	})

	t.Run("test menggunakan mutex", func(t *testing.T) {
		for i := 0; i < count; i++ {
			wg.Add(1)
			go func(wg *sync.WaitGroup, mu *sync.Mutex) {
				mu.Lock()
				defer wg.Done()
				defer mu.Unlock()

				totalData++
				logrus.WithField("i", i+1).Infof("menambah kan 1 data")
			}(wg, mu)
		}

		wg.Wait()

		// wait all goroutines done
		logrus.Infof("jumlah total data %d", totalData)
	})
}
