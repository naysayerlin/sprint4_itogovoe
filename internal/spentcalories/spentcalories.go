package spentcalories

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"
)

// Основные константы, необходимые для расчетов.
const (
	lenStep                    = 0.65 // средняя длина шага.
	mInKm                      = 1000 // количество метров в километре.
	minInH                     = 60   // количество минут в часе.
	stepLengthCoefficient      = 0.45 // коэффициент для расчета длины шага на основе роста.
	walkingCaloriesCoefficient = 0.5  // коэффициент для расчета калорий при ходьбе

)

func parseTraining(data string) (int, string, time.Duration, error) {
	sliceStrok := strings.Split(data, ",")
	if len(sliceStrok) < 3 || len(sliceStrok) > 3 {
		return 0, "", 0, errors.New("Неверное кол-во параметров")
	}
	stepCount, err := strconv.Atoi(sliceStrok[0])
	if err != nil {
		return 0, "", 0, errors.New("Ошибка конвертации шагов")
	}
	if stepCount <= 0 {
		return 0, "", 0, errors.New("Вы бездельник-седун, пройдитесь!")
	}
	timeOfWalk, err := time.ParseDuration(sliceStrok[2])
	if err != nil {
		return 0, "", 0, fmt.Errorf("Ошибка: %w\n", err)
	}
	if timeOfWalk <= 0 {
		return 0, "", 0, errors.New("Вы сегодня не ходили, пройдитесь")
	}
	return stepCount, sliceStrok[1], timeOfWalk, nil
}

func distance(steps int, height float64) float64 {
	lenOfstep := height * stepLengthCoefficient
	dist := float64(steps) * lenOfstep
	return dist / float64(mInKm)
}

func meanSpeed(steps int, height float64, duration time.Duration) float64 {
	if duration <= 0 {
		return 0
	}
	dist := distance(steps, height)
	durat := duration.Hours()
	return dist / durat
}

func TrainingInfo(data string, weight, height float64) (string, error) {
	stepCount, activity, timeOfWalk, err := parseTraining(data)
	if err != nil {
		log.Println(err)
	}
	dist := distance(stepCount, height)
	avgSpeed := meanSpeed(stepCount, height, timeOfWalk)
	switch activity {
	case "Бег":
		calories, err := RunningSpentCalories(stepCount, weight, height, timeOfWalk)
		if err != nil {
			log.Println(err)
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, timeOfWalk.Hours(), dist, avgSpeed, calories), nil
	case "Ходьба":
		calories, err := WalkingSpentCalories(stepCount, weight, height, timeOfWalk)
		if err != nil {
			log.Println(err)
		}
		return fmt.Sprintf("Тип тренировки: %s\nДлительность: %.2f ч.\nДистанция: %.2f км.\nСкорость: %.2f км/ч\nСожгли калорий: %.2f\n", activity, timeOfWalk.Hours(), dist, avgSpeed, calories), nil

	default:
		return "", errors.New("неизвестный тип тренировки")
	}

}

func RunningSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 {
		return 0, errors.New("Количество шагов не указано либо некорректно")
	}
	if weight <= 0 {
		return 0, errors.New("Вес либо равен нулю либо отрицательный")
	}
	if height <= 0 {
		return 0, errors.New("Рост либо равен нулю либо отрицательный")
	}
	if duration <= 0 {
		return 0, errors.New("Продолжительность тренировки равна нулю либо отрицательна")
	}
	avgSpeed := meanSpeed(steps, height, duration)
	durationInMins := duration.Minutes()
	return (weight * avgSpeed * durationInMins) / minInH, nil
}

func WalkingSpentCalories(steps int, weight, height float64, duration time.Duration) (float64, error) {
	if steps <= 0 || weight <= 0 || height <= 0 || duration <= 0 {
		return 0, errors.New("Ошибка данных")
	}
	avgSpeed := meanSpeed(steps, height, duration)
	durationInMins := duration.Minutes()
	calories := (weight * avgSpeed * durationInMins) / minInH
	return calories * walkingCaloriesCoefficient, nil

}
