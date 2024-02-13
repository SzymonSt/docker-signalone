package utils

import "go.mongodb.org/mongo-driver/bson"

func CalculateNewCounter(currentScore int32, newScore int32, counter int32) int32 {
	if currentScore == 1 {
		if newScore == 1 {
			return counter
		} else if newScore == 0 {
			return counter - 1
		} else if newScore == -1 {
			return counter - 2
		}
	}

	if currentScore == -1 {
		if newScore == 1 {
			return counter + 2
		} else if newScore == 0 {
			return counter + 1
		} else if newScore == -1 {
			return counter
		}
	}

	return counter + newScore
}

func GenerateFilter(fields bson.M) bson.M {
	conditions := make([]bson.M, 0, len(fields))

	for field, value := range fields {
		conditions = append(conditions, bson.M{field: value})
	}

	return bson.M{"$and": conditions}
}
