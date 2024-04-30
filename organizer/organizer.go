package organizer

import (
	"archive/zip"
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

type Org struct{}

type IOrganizer interface {
	OrganizScramble()
}

func NewOrganizer() *Org {
	return &Org{}
}

func (o *Org) OrganizScramble() error {
	inputScanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Enter competition id:")
	inputScanner.Scan()
	input := inputScanner.Text()

	competition, err := getWCIF(input)
	if err != nil {
		panic(err)
	}

	var eventDetail []EventDetail
	passcodeMap := make(map[string]string)

	for _, venue := range competition.Schedule.Venues {
		for _, room := range venue.Rooms {
			for _, activity := range room.Activities {
				if !strings.HasPrefix(activity.ActivityCode, "other-") {
					splitActivityCode := strings.Split(activity.ActivityCode, "-")

					activityCode := splitActivityCode[0]
					activityRound := splitActivityCode[1]

					roundNumber, err := strconv.Atoi(string(activityRound[1]))
					if err != nil {
						fmt.Println("Error:", err)
						return err
					}

					ed := EventDetail{
						EventCode:      activityCode,
						EventName:      EventCodeToFullMap[activityCode],
						EventVenue:     venue.VenueName,
						EventRoom:      room.RoomName,
						EventRound:     roundNumber,
						EventStartTime: activity.StartTime,
					}

					if activityCode == "333fm" || activityCode == "333mbf" {
						activityAttempt := splitActivityCode[2]
						attemptNumber, err := strconv.Atoi(string(activityAttempt[1]))
						if err != nil {
							fmt.Println("Error:", err)
							return err
						}

						ed.EventAttempt = attemptNumber

					} else {
						for _, child := range activity.ChildActivities {
							splitChildActivityCode := strings.Split(child.ActivityCode, "-")

							childActivityGroup := splitChildActivityCode[2]
							groupNumber, err := strconv.Atoi(string(childActivityGroup[1]))
							if err != nil {
								fmt.Println("Error:", err)
								return err
							}

							egd := EventGroupDetail{
								EventGroup:       getAlphabetFromOrder(groupNumber),
								EventGroupNumber: groupNumber,
								EventStartTime:   child.StartTime,
							}

							ed.EventGroupDetails = append(ed.EventGroupDetails, egd)
						}
					}

					eventDetail = append(eventDetail, ed)

				}

			}
		}
	}

	sort.Slice(eventDetail, func(i, j int) bool {
		return eventDetail[i].EventStartTime.Before(eventDetail[j].EventStartTime)
	})

	scrambleZipName := competition.CompetitionName + ".zip"
	if _, err := os.Stat(scrambleZipName); os.IsNotExist(err) {
		fmt.Println("Zip file does not exist")
		return err
	}

	computerDisplayZipName := competition.CompetitionName + " - Computer Display PDFs.zip"
	passcodeFileName := competition.CompetitionName + " - Computer Display PDF Passcodes - SECRET.txt"

	extractionFolder := "extracted_files"
	err = unzipScrambleFile(scrambleZipName, extractionFolder, computerDisplayZipName, passcodeFileName)
	if err != nil {
		panic(err)
	}

	extractedPasscode := extractionFolder + "/" + passcodeFileName

	passcodeFile, err := os.Open(extractedPasscode)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return err
	}
	defer passcodeFile.Close()

	scanner := bufio.NewScanner(passcodeFile)

	for i := 0; i < 9; i++ {
		if !scanner.Scan() {
			fmt.Println("File have less than 10 lines.")
			return err
		}
	}

	for scanner.Scan() {
		line := scanner.Text()

		splitStr := strings.Split(line, ":")

		splitEvent := strings.Split(splitStr[0], " Round ")

		splitGroup := strings.Split(splitEvent[1], " Scramble Set ")

		eventCode := EventFullToCodeMap[splitEvent[0]]
		round := splitEvent[1][0:1]

		passcodeKey := eventCode + "-r" + round

		if eventCode == "333fm" || eventCode == "333mbf" {
			splitAtempt := strings.Split(splitEvent[1], " Attempt ")
			passcodeKey = passcodeKey + "-a" + splitAtempt[1][0:1]
		} else {
			groupChar := splitGroup[1][0:1]
			groupNumb := strconv.Itoa(getOrderFromAlphabet(groupChar))
			passcodeKey = passcodeKey + "-g" + groupNumb
		}

		passcodeMap[passcodeKey] = splitStr[1]
	}

	scrambleDir := "scramble"
	scheduleOrder := 1
	displayScrambleDir := extractionFolder + "/scramble_display"
	var currentDate time.Time

	err = os.MkdirAll(scrambleDir, os.ModePerm)
	if err != nil {
		return err
	}

	sortedPasscodeFile, err := os.OpenFile(scrambleDir+"/sortedPasscode.txt", os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer sortedPasscodeFile.Close()

	passcodeWriter := bufio.NewWriter(sortedPasscodeFile)

	for _, sortedEvent := range eventDetail {
		currentDateString := currentDate.Format("2006-01-02")

		if scheduleOrder == 1 || currentDateString != sortedEvent.EventStartTime.Format("2006-01-02") {
			currentDate = sortedEvent.EventStartTime
			currentDateString = currentDate.Format("2006-01-02")

			_, err = passcodeWriter.WriteString(currentDateString + "\n")
			if err != nil {
				return err
			}
		}

		eventDir := scrambleDir + "/" + sortedEvent.EventVenue + "/" + sortedEvent.EventRoom + "/" + currentDateString + "/" + strconv.Itoa(scheduleOrder) + "_" + sortedEvent.EventName + "_r" + strconv.Itoa(sortedEvent.EventRound)
		if sortedEvent.EventCode == "333fm" || sortedEvent.EventCode == "333mbf" {
			eventDir = eventDir + "_a" + strconv.Itoa(sortedEvent.EventAttempt)
		}

		err = os.MkdirAll(eventDir, os.ModePerm)
		if err != nil {
			return err
		}

		if sortedEvent.EventCode == "333fm" || sortedEvent.EventCode == "333mbf" {
			scrambleFileName := sortedEvent.EventName + " Round " + strconv.Itoa(sortedEvent.EventRound) + " Scramble Set A Attempt " + strconv.Itoa(sortedEvent.EventAttempt) + ".pdf"

			scrambleSource := displayScrambleDir + "/" + scrambleFileName
			scrambleDest := eventDir + "/" + scrambleFileName

			err = moveFile(scrambleSource, scrambleDest)
			if err != nil {
				return err
			}

			passcodeKey := sortedEvent.EventCode + "-r" + strconv.Itoa(sortedEvent.EventRound) + "-a" + strconv.Itoa(sortedEvent.EventAttempt)

			passcode := passcodeMap[passcodeKey]
			_, err = passcodeWriter.WriteString(sortedEvent.EventName + " Round" + strconv.Itoa(sortedEvent.EventRound) + " Attempt" + strconv.Itoa(sortedEvent.EventAttempt) + ": " + passcode + "\n")
			if err != nil {
				return err
			}
		} else {
			for _, child := range sortedEvent.EventGroupDetails {
				scrambleFileName := sortedEvent.EventName + " Round " + strconv.Itoa(sortedEvent.EventRound) + " Scramble Set " + child.EventGroup + ".pdf"

				scrambleSource := displayScrambleDir + "/" + scrambleFileName
				scrambleDest := eventDir + "/" + scrambleFileName

				err = moveFile(scrambleSource, scrambleDest)
				if err != nil {
					return err
				}

				passcodeKey := sortedEvent.EventCode + "-r" + strconv.Itoa(sortedEvent.EventRound) + "-g" + strconv.Itoa(child.EventGroupNumber)

				passcode := passcodeMap[passcodeKey]
				_, err = passcodeWriter.WriteString(sortedEvent.EventName + " Round" + strconv.Itoa(sortedEvent.EventRound) + " Set(" + child.EventGroup + "): " + passcode + "\n")
				if err != nil {
					return err
				}
			}
		}

		scheduleOrder = scheduleOrder + 1

		_, err = passcodeWriter.WriteString("\n")
		if err != nil {
			return err
		}
	}

	if err := passcodeWriter.Flush(); err != nil {
		return err
	}

	defer func() {
		err := os.RemoveAll(extractionFolder)
		if err != nil {
			fmt.Println("Error removing extraction folder:", err)
		} else {
			fmt.Println("Extraction folder removed.")
		}
	}()

	return nil
}

func unzipScrambleFile(zipFileName string, extractionFolder string, displayFileName string, passcodeFileName string) error {
	err := os.MkdirAll(extractionFolder, os.ModePerm)
	if err != nil {
		return err
	}

	scrambleZipReader, err := zip.OpenReader(zipFileName)
	if err != nil {
		return err
	}
	defer scrambleZipReader.Close()

	for _, scrambleFile := range scrambleZipReader.File {
		if filepath.Base(scrambleFile.Name) == displayFileName || filepath.Base(scrambleFile.Name) == passcodeFileName {
			zippedFile, err := scrambleFile.Open()
			if err != nil {
				return err
			}
			defer zippedFile.Close()

			extractPath := filepath.Join(extractionFolder, filepath.Base(scrambleFile.Name))

			extractedFile, err := os.OpenFile(extractPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, scrambleFile.Mode())
			if err != nil {
				return err
			}
			defer extractedFile.Close()

			_, err = io.Copy(extractedFile, zippedFile)
			if err != nil {
				return err
			}

			if filepath.Base(scrambleFile.Name) == displayFileName {
				displayScrambleDir := extractionFolder + "/scramble_display"

				err := os.MkdirAll(displayScrambleDir, os.ModePerm)
				if err != nil {
					return err
				}

				displayScrambleZipReader, err := zip.OpenReader(extractionFolder + "/" + displayFileName)
				if err != nil {
					return err
				}
				defer scrambleZipReader.Close()

				for _, displayScrambleFile := range displayScrambleZipReader.File {
					displayZippedFile, err := scrambleFile.Open()
					if err != nil {
						return err
					}
					defer displayZippedFile.Close()

					extracDisplayPath := filepath.Join(displayScrambleDir, filepath.Base(displayScrambleFile.Name))

					extractedDisplayFile, err := os.OpenFile(extracDisplayPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, displayScrambleFile.Mode())
					if err != nil {
						return err
					}
					defer extractedDisplayFile.Close()

					_, err = io.Copy(extractedDisplayFile, displayZippedFile)
					if err != nil {
						return err
					}
				}
			}
		}
	}

	return nil
}

func getWCIF(competitionId string) (WCACompetition, error) {
	wcifURL := "https://www.worldcubeassociation.org/api/v0/competitions/{competitionID}/wcif/public"

	thisCompURL := strings.Replace(wcifURL, "{competitionID}", competitionId, -1)

	resp, err := http.Get(thisCompURL)
	if err != nil {
		fmt.Printf("error when call wcif: %v\n", err)
		return WCACompetition{}, err
	}
	defer resp.Body.Close()

	wcifByte, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("error when read body: %v\n", err)
		return WCACompetition{}, err
	}

	var Competition WCACompetition

	json.Unmarshal(wcifByte, &Competition)

	return Competition, nil
}

func getAlphabetFromOrder(num int) string {
	if num < 1 || num > 26 {
		return ""
	}
	return string(rune('A' + num - 1))
}

func moveFile(src, dst string) error {
	err := os.Rename(src, dst)
	if err != nil {
		return fmt.Errorf("error moving file: %s", err)
	}

	return nil
}

func getOrderFromAlphabet(alphabet string) int {
	if len(alphabet) != 1 || alphabet[0] < 'A' || alphabet[0] > 'Z' {
		return -1 // Return -1 for invalid input
	}
	return int(alphabet[0] - 'A' + 1)
}
