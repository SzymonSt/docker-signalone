package utils

import "go.mongodb.org/mongo-driver/bson"

func GenerateFilter(fields bson.M) bson.M {
	conditions := make([]bson.M, 0, len(fields))

	for field, value := range fields {
		conditions = append(conditions, bson.M{field: value})
	}

	return bson.M{"$or": conditions}
}
