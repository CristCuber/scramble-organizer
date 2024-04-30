package organizer

import "time"

type WCACompetition struct {
	CompetitionName string   `json:"name"`
	Events          []Event  `json:"events"`
	Schedule        Schedule `json:"schedule"`
}

type Event struct {
	EventID string  `json:"id"`
	Rounds  []Round `json:"rounds"`
}

type Round struct {
	RoundID              string               `json:"id"`
	Format               string               `json:"format"`
	TimeLimit            TimeLimit            `json:"timeLimit"`
	Cutoff               int                  `json:"cutoff"`
	AdvancementCondition AdvancementCondition `json:"advancementCondition"`
	ScrambleSetCount     int                  `json:"scrambleSetCount"`
}

type TimeLimit struct {
	CentiSeconds       int      `json:"centiseconds"`
	CumulativeRoundIDs []string `json:"cumulativeRoundIds"`
}

type AdvancementCondition struct {
	Type  string `json:"type"`
	Level int    `json:"level"`
}

type Schedule struct {
	StartDate    string  `json:"startDate"`
	NumberOfDays int     `json:"numberOfDays"`
	Venues       []Venue `json:"venues"`
}

type Venue struct {
	VenueID   int    `json:"id"`
	VenueName string `json:"name"`
	Rooms     []Room `json:"rooms"`
}

type Room struct {
	RoomID     int        `json:"id"`
	RoomName   string     `json:"name"`
	Activities []Activity `json:"activities"`
}

type Activity struct {
	ActivityID      int        `json:"id"`
	ActivityName    string     `json:"name"`
	ActivityCode    string     `json:"activityCode"`
	StartTime       time.Time  `json:"startTime"`
	EndTime         time.Time  `json:"endTime"`
	ChildActivities []Activity `json:"childActivities"`
}

type EventDetail struct {
	EventCode         string
	EventName         string
	EventVenue        string
	EventRoom         string
	EventRound        int
	EventAttempt      int
	EventGroupDetails []EventGroupDetail
	EventStartTime    time.Time
}

type EventGroupDetail struct {
	EventGroup       string
	EventGroupNumber int
	EventStartTime   time.Time
}
