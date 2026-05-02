package packets

type StoreParticipant struct {
	M_aiControlled uint8    // Whether the vehicle is AI (1) or Human (0) controlled
	M_driverName   [48]byte // Driver id - see appendix, 255 if network human
	M_networkId    uint8    // Network id – unique identifier for network players
	M_teamName     [48]byte // Team id - see appendix
	M_myTeam       uint8    // My team flag – 1 = My Team, 0 = otherwise
	M_raceNumber   uint8    // Race number of the car
	M_nationality  [48]byte // Nationality of the driver
	M_name         [48]byte // Name of participant in UTF-8 format – null terminated. Will be truncated with "..."
	// (U+2026) if too long
	M_yourTelemetry   uint8    // The player's UDP setting, 0 = restricted, 1 = public
	M_showOnlineNames uint8    // The player's show online names setting, 0 = off, 1 = on
	M_platformName    [48]byte // 1 = Steam, 3 = PlayStation, 4 = Xbox, 6 = Origin, 255 = unknown
}
