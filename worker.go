package main

import (
	"fmt"
	"net/http"
)

func Action(w http.ResponseWriter, r *http.Request) {

	c, err := r.Cookie("session_token")
	vlanTr069 := r.FormValue("vlan_tr069")
	vlanPPPoE := r.FormValue("vlan_pppoe")
	vlanIPTV := r.FormValue("vlan_iptv")
	vlanIMS := r.FormValue("vlan_ims")
	vlanvIMS := r.FormValue("vlan_vims")
	hardware := r.FormValue("hardware")
	btn := r.FormValue("btn-body")
	slot := r.FormValue("slot")
	port := r.FormValue("port")
	ont := r.FormValue("ont")

	fmt.Println("C:", c, "ERR:", err,
		"HARD:", hardware,
		"BTN:", btn,
		"slot:", slot,
		"port:", port,
		"ont:", ont,
		"TR069:", vlanTr069,
		"vlanPPPoE", vlanPPPoE,
		"vlanIPTV:", vlanIPTV,
		"vlanIMS:", vlanIMS,
		"vlanvIMS:", vlanvIMS,
	)

	// We can obtain the session token from the requests cookies, which come with every request
	if err != nil {
		if err == http.ErrNoCookie {
			// Если куки не установлены, то это unauthorized status, поэтому снова предлагаем залогиниться
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			//w.WriteHeader(http.StatusUnauthorized)
			return
		}
		// For any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := c.Value
	// На нулевой позиции находится текущее значение из списка возможных см. func InitDefData
	sessions[sessionToken].HTML.HTML_Devices[0] = hardware
	sessions[sessionToken].HTML.HTML_Slots[0] = slot
	sessions[sessionToken].HTML.HTML_Ports[0] = port
	sessions[sessionToken].HTML.HTML_Onts[0] = ont
	sessions[sessionToken].HTML.HTML_VlanTr069s[0] = vlanTr069
	sessions[sessionToken].HTML.HTML_VlanPPPoEs[0] = vlanPPPoE
	sessions[sessionToken].HTML.HTML_VlanIPTVs[0] = vlanIPTV
	sessions[sessionToken].HTML.HTML_VlanIMSs[0] = vlanIMS
	sessions[sessionToken].HTML.HTML_VlanvIMSs[0] = vlanvIMS

	UpdateHTML(sessions[sessionToken])

	// Затем мы получаем имя пользователя из нашей карты сеанса, где мы устанавливаем токен сеанса
	userSession, exists := sessions[sessionToken]

	fmt.Println("IP_OLT:", userSession.HTML.HTML_DeviceIP[hardware])
	if !exists {
		// If the session token is not present in session map, return an unauthorized error
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		//w.WriteHeader(http.StatusUnauthorized)
		return
	}
	if userSession.isExpired() {
		delete(sessions, sessionToken)
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Finally, return the welcome message to the user
	//w.Write([]byte(fmt.Sprintf("Welcome %s!", users[userSession.username].Name)))
	InfoLogger.Printf("[%s], Отрисовка формы после ввода даных в форму конфигурения", r.RemoteAddr)
	http.Redirect(w, r, "/setup", http.StatusSeeOther)

}
