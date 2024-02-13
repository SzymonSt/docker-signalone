package utils

import "go.mongodb.org/mongo-driver/bson"

func CalculateNewCounter(currentScore int32, newScore int32, counter int32) int32 {
	return counter + (newScore - currentScore)
}

func GenerateFilter(fields bson.M, operator string) bson.M {
	conditions := make([]bson.M, 0, len(fields))

	for field, value := range fields {
		conditions = append(conditions, bson.M{field: value})
	}

	return bson.M{operator: conditions}
}
