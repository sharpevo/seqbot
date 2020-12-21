package lane

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

const (
	LAYOUT    = "20060102_150405"
	FILE_LANE = "lane.csv"
	SEPARATOR = ","
)

type Lane struct {
	ChipId    string
	startTime time.Time
	endTime   time.Time
}

func NewLane(chipId string) *Lane {
	return &Lane{
		ChipId: chipId,
	}
}

func (l *Lane) StartTime() string {
	return l.startTime.Format(LAYOUT)
}

func (l *Lane) EndTime() string {
	return l.endTime.Format(LAYOUT)
}

func (l *Lane) Start() error {
	l.startTime = time.Now()
	return l.addLane()
}

func (l *Lane) addLane() error {
	lines, err := readLanes(FILE_LANE)
	if err != nil {
		return err
	}
	lines = append(lines, l.String())
	return saveLanes(lines)
}

func (l *Lane) Finish() error {
	lines, err := readLanes(FILE_LANE)
	if err != nil {
		return err
	}
	index := -1
	for i := len(lines) - 1; i >= 0; i-- {
		if strings.HasPrefix(lines[i], l.ChipId) {
			index = i
			break
		}
	}
	if index < 0 {
		return fmt.Errorf("failed to get record index: %s", l.ChipId)
	}
	items := strings.Split(lines[index], SEPARATOR)
	switch len(items) {
	case 2:
		l.endTime = time.Now()
		l.startTime, err = time.ParseInLocation(
			LAYOUT, items[1], l.endTime.Location())
		if err != nil {
			return err
		}
		lines[index] = l.String()
	default:
		return fmt.Errorf("invalid lane to be done: %s", l.ChipId)
	}
	return saveLanes(lines)
}

func (l *Lane) Duration() string {
	endRound := time.Date(
		l.endTime.Year(), l.endTime.Month(), l.endTime.Day(),
		l.endTime.Hour(), l.endTime.Minute(), l.endTime.Second(),
		0, l.endTime.Location())
	return fmt.Sprintf("%v", endRound.Sub(l.startTime))
}

func (l *Lane) String() string {
	str := l.ChipId
	if !l.startTime.IsZero() {
		str = fmt.Sprintf("%s,%s", str, l.startTime.Format(LAYOUT))
	}
	if !l.endTime.IsZero() {
		str = fmt.Sprintf("%s,%s", str, l.endTime.Format(LAYOUT))
	}
	return str
}

func readLanes(lanePath string) ([]string, error) {
	var lines []string
	file, err := os.OpenFile(lanePath, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return lines, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, nil
}

func saveLanes(lines []string) error {
	file, err := os.OpenFile(FILE_LANE, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	file.Seek(0, 0)
	file.Truncate(0)
	writer := bufio.NewWriter(file)
	for _, line := range lines {
		writer.WriteString(line + "\n")
	}
	return writer.Flush()
}
