package models

import "time"

type SessionData struct {
	UserID      int       // ID de l'utilisateur associé à la session
	Username    string    // Nom d'utilisateur de l'utilisateur
	IsLoggedIn  bool      // Indique si l'utilisateur est connecté
	CreatedTime time.Time // Heure de création de la session
}
