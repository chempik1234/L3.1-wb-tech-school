package transport

import "github.com/wb-go/wbf/ginext"

// AssembleRouter is the function you'd call in `main.go` to get THE app router
func AssembleRouter(notifyHandler *NotifyHandler) *ginext.Engine {
	router := ginext.New("release")

	router.POST("/notify", notifyHandler.CreateNotification)
	router.GET("/notify/:id", notifyHandler.GetNotification)
	router.DELETE("/notify/:id", notifyHandler.DeleteNotification)

	return router
}
