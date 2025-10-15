package http_server

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

// HTTPServer это общая структура для HTTP сервера, совместимая с ginext
//
// Порт задаётся при запуске
type HTTPServer struct {
	router http.Handler
}

// NewHTTPServer создаёт HTTPServer с данным роутером
func NewHTTPServer(router http.Handler) *HTTPServer {
	return &HTTPServer{router: router}
}

// GracefulRun запускает HTTPServer сервер на данном порту и плавно завершает при os.Interrupt или естественной ошибке
//
// “ctx context.Context“ тоже вызывает Graceful Shutdown
//
// 1. Создать структуру сервера. При каждом запуске новая
//
// 2. Подготовить каналы с сигналами
//
// 3. Запустить фоном слушатель для shutdown - именно он и закрывает сервер в нормальных условиях
//
// 4. Запустить сам сервер и ждать, пока он не схлопнется
//
// 5. Подождать завершение слушателя и выйти
func (s *HTTPServer) GracefulRun(ctx context.Context, port int) error {
	// шаг 1.
	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: s.router,
	}

	// шаг 1.1. Каналы с сигналами о том, что 1) вышел сервер 2) вышла горутина, слушаящая os.Interrupt и ctx
	serverStopped := make(chan bool, 1)
	signalListenerExited := make(chan bool, 1)

	// шаг 2.
	go listenSignal(ctx, httpServer, serverStopped, signalListenerExited)

	// шаг 3.
	err := httpServer.ListenAndServe()
	serverStopped <- true

	// подождём пока не завершится слушатель, и потом выйдем
	<-signalListenerExited // этот сигнал всегда идёт после

	if err != nil {
		if !errors.Is(err, http.ErrServerClosed) {
			return fmt.Errorf("error while listening HTTP port '%d': %w", port, err)
		}
	}

	return nil
}

func listenSignal(ctx context.Context, httpServer *http.Server, serverStopped <-chan bool, funcExited chan<- bool) {
	// шаг 1. Graceful shutdown через сигнал - прикручиваем к основному контексту
	signalCtx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	// шаг 2. Есть 2 источника сигналов
	//
	// 1) signalCtx проверяет: настал os.Signal или "вызывающая сторона" закрыла контекст
	// 2) serverStopped проверяет, что сервер уже отрубился без нас

	select {
	// Если не <1>, а уже <2>, просто выйдем
	case <-serverStopped:
		break
	// Если всё-таки <1>, то плавно отрубаемся (где-то на фоне случится <2>)
	case <-signalCtx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		someErr := httpServer.Shutdown(shutdownCtx)
		if someErr != nil {
			log.Fatalf("error while shutting down http server: %v", someErr)
		}
	}

	funcExited <- true
}
