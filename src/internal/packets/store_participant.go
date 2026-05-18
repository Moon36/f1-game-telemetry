package packets

// StoreParticipant is the struct used to store enriched participant data in Redis
type StoreParticipant struct {
	M_aiControlled uint8  `json:"m_aiControlled"` // Whether the vehicle is AI (1) or Human (0) controlled
	M_driverName   string `json:"m_driverName"`   // Driver id - see appendix, 255 if network human
	M_networkId    uint8  `json:"m_networkId"`    // Network id – unique identifier for network players
	M_teamName     string `json:"m_teamName"`     // Team id - see appendix
	M_myTeam       uint8  `json:"m_myTeam"`       // My team flag – 1 = My Team, 0 = otherwise
	M_raceNumber   uint8  `json:"m_raceNumber"`   // Race number of the car
	M_nationality  string `json:"m_nationality"`  // Nationality of the driver
	M_name         string `json:"m_name"`         /* Name of participant in UTF-8 format – null terminated. Will be
	truncated with "..." (U+2026) if too long */
	M_yourTelemetry   uint8  `json:"m_yourTelemetry"`   // The player's UDP setting, 0 = restricted, 1 = public
	M_showOnlineNames uint8  `json:"m_showOnlineNames"` // The player's show online names setting, 0 = off, 1 = on
	M_platformName    string `json:"m_platformName"`    /* 1 = Steam, 3 = PlayStation, 4 = Xbox, 6 = Origin
	255 = unknown */

	// Add additional fields as needed, use pointers for new/optional fields and `omitempty` JSON tag
}
