package handler

import (
	"net/http"

	"github.com/waynecraig/wechat-token-hub/internal/tokens"
)

// Ticket handles requests to the /ticket path
func Ticket(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	ticketType := query.Get("type")
	rotateTicket := query.Get("rotate_ticket")
	ticket, err := tokens.GetTicket(ticketType, rotateTicket)
	if err != nil {
		// return error if get ticket fail
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}
	// return the access token
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(ticket))
}
