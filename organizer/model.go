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

type WCAEvent string

const (
	Event3x3x3            WCAEvent = "333"
	Event2x2x2            WCAEvent = "222"
	Event4x4x4            WCAEvent = "444"
	Event5x5x5            WCAEvent = "555"
	Event6x6x6            WCAEvent = "666"
	Event7x7x7            WCAEvent = "777"
	Event3x3x3Blindfolded WCAEvent = "333bf"
	Event3x3x3OneHanded   WCAEvent = "333oh"
	EventClock            WCAEvent = "clock"
	EventMegaminx         WCAEvent = "minx"
	EventPyraminx         WCAEvent = "pyram"
	EventSkewb            WCAEvent = "skewb"
	EventSquare1          WCAEvent = "sq1"
	Event4x4x4Blindfolded WCAEvent = "444bf"
	Event5x5x5Blindfolded WCAEvent = "555bf"
	Event3x3x3MultiBlind  WCAEvent = "333mbf"
	Event3x3x3FewestMoves WCAEvent = "333fm"
)

var eventOrder = map[WCAEvent]int{
	Event3x3x3:            1,
	Event2x2x2:            2,
	Event4x4x4:            3,
	Event5x5x5:            4,
	Event6x6x6:            5,
	Event7x7x7:            6,
	Event3x3x3Blindfolded: 7,
	Event3x3x3OneHanded:   8,
	EventClock:            9,
	EventMegaminx:         10,
	EventPyraminx:         11,
	EventSkewb:            12,
	EventSquare1:          13,
	Event4x4x4Blindfolded: 14,
	Event5x5x5Blindfolded: 15,
	Event3x3x3MultiBlind:  16,
	Event3x3x3FewestMoves: 17,
}

type WCAEvents []WCAEvent

func (events WCAEvents) Len() int      { return len(events) }
func (events WCAEvents) Swap(i, j int) { events[i], events[j] = events[j], events[i] }
func (events WCAEvents) Less(i, j int) bool {
	return eventOrder[events[i]] < eventOrder[events[j]]
}
