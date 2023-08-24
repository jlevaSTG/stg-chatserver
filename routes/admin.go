package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"stg-go-websocket-server/ws"
	"time"
)

type ClientDetail struct {
	ClientID  string    `json:"client_id"`
	LoginInAt time.Time `json:"login_in_at"`
}

type ManagerDetails struct {
	NumberOfClients int64          `json:"numberOfClients"`
	Clients         []ClientDetail `json:"clients"`
}

func SetupAdminRoutes(route string, r *gin.Engine, m *ws.Manager) {
	api := r.Group(route)
	api.GET("/stats", hadnleGetConnections(m))

}

func hadnleGetConnections(m *ws.Manager) gin.HandlerFunc {
	return func(ginCtx *gin.Context) {
		m.Lock()
		defer m.Unlock()

		cDetail := make([]ClientDetail, 0)
		for _, c := range m.Clients {
			cDetail = append(cDetail, ClientDetail{ClientID: c.Id, LoginInAt: c.Joined})
		}
		ginCtx.JSON(http.StatusOK, gin.H{"serverDetail": ManagerDetails{
			NumberOfClients: m.CurrentClientCount,
			Clients:         cDetail,
		}})
	}
}
