package extension

import (
	"fmt"
	"regexp"
	"strconv"
	"time"

	"github.com/beevik/ntp"
	"github.com/bwmarrin/discordgo"
	"github.com/mattn/go-colorable"
	logHandler "github.com/sirupsen/logrus"
)

// 날짜에 해당하는 한국어 요일 반환
func GetKoreanWeekday(date time.Time) string {
	var koreanWeekday string

	switch int(date.Weekday()) {
	case 0:
		koreanWeekday = "일"
	case 1:
		koreanWeekday = "월"
	case 2:
		koreanWeekday = "화"
	case 3:
		koreanWeekday = "수"
	case 4:
		koreanWeekday = "목"
	case 5:
		koreanWeekday = "금"
	case 6:
		koreanWeekday = "토"
	}

	return "(" + koreanWeekday + ")"
}

// 급식 시간대 코드에 해당하는 이름 반환
func GetMealName(mealCode int) string {
	var mealName string

	switch mealCode {
	case 1:
		mealName = "조식"

	case 2:
		mealName = "중식"

	case 3:
		mealName = "석식"

	default:
		mealName = "급식"
	}

	return mealName
}

// 한국어 단어에 을/를 추가
func AddKoreanObjectParticle(str string) string {
	start := rune('가')
	intervalEnd := rune('갛')
	end := rune('힣')

	mod := intervalEnd - start + 1
	hasFinalConsonant := false

	for _, r := range str {
		if !(r >= start && r <= end) {
			hasFinalConsonant = false
		}

		hasFinalConsonant = (r-start)%mod != 0
	}

	if hasFinalConsonant {
		return str + "을"
	} else {
		return str + "를"
	}
}

// 1~6의 정수인지 확인
func IsGradeNumber(str string) bool {
	var digitCheck = regexp.MustCompile(`^[0-9]+$`)
	num, _ := strconv.Atoi(str)

	return digitCheck.MatchString(str) && num >= 1 && num <= 6
}

// 1~16의 정수인지 확인
func IsClassNumber(str string) bool {
	var digitCheck = regexp.MustCompile(`^[0-9]+$`)
	num, _ := strconv.Atoi(str)

	return digitCheck.MatchString(str) && num >= 1 && num <= 16
}

// 메시지 전송 및 예외 처리
func ChannelMessageSend(s *discordgo.Session, channelId string, format string, args ...interface{}) {
	_, err := s.ChannelMessageSend(channelId, fmt.Sprintf(format, args...))

	if err != nil {
		Log().Warningln(err)
	}
}

// 임베드 전송 및 예외 처리
func ChannelMessageSendEmbed(s *discordgo.Session, channelId string, embed *discordgo.MessageEmbed) {
	_, err := s.ChannelMessageSendEmbed(channelId, embed)

	if err != nil {
		Log().Warningln(err)
	}
}

// ntp 서버에서 현재 시각 가져와 한국 시간대로 변경
func NtpTimeKorea() (time.Time, error) {
	date, err := ntp.Time("time.windows.com")

	if err != nil {
		return date, err
	}

	loc, err := time.LoadLocation("")

	if err != nil {
		return date, err
	}

	return date.In(loc).Add(time.Hour * 9), err
}

func Log() *logHandler.Entry {
	logHandler.SetFormatter(&logHandler.TextFormatter{
		ForceColors: true,
	})
	logHandler.SetOutput(colorable.NewColorableStdout())
	logHandler.SetLevel(logHandler.DebugLevel)
	var lo = logHandler.WithFields(logHandler.Fields{})
	return lo
}
