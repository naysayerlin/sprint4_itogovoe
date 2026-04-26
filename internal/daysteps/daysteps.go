package daysteps

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Yandex-Practicum/tracker/internal/spentcalories"
)

const (
	// Длина одного шага в метрах
	stepLength = 0.65
	// Количество метров в одном километре
	mInKm = 1000
)

func parsePackage(data string) (int, time.Duration, error) {
	sliceStrok := strings.Split(data, ",")
	if len(sliceStrok) < 2 {
		return 0, 0, errors.New("Нехватка параметров")
	}
	stepCount, err := strconv.Atoi(sliceStrok[0])
	if err != nil {
		return 0, 0, errors.New("Ошибка конвертации шагов")
	}
	if stepCount <= 0 {
		return 0, 0, errors.New("Вы бездельник-седун, пройдитесь!")
	}
	timeOfWalk, err := time.ParseDuration(sliceStrok[1])
	if err != nil {
		return 0, 0, fmt.Errorf("Ошибка конвертации времени: %w\n", err)
	}
	if timeOfWalk <= 0 {
		return 0, 0, errors.New("Вы сегодня не ходили, пройдитесь")
	}
	return stepCount, timeOfWalk, nil

}

func DayActionInfo(data string, weight, height float64) string {
	stepCount, timeOfWalk, err := parsePackage(data)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	if stepCount <= 0 {
		return ""
	}
	distance := stepLength * float64(stepCount)
	distInKm := distance / float64(mInKm)
	calories, err := spentcalories.WalkingSpentCalories(stepCount, weight, height, timeOfWalk)
	return fmt.Sprintf("Количество шагов: %d.\nДистанция составила %.2f км.\nВы сожгли %.2f ккал.\n", stepCount, distInKm, calories)

}
