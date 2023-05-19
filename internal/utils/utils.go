package utils

import "go.mongodb.org/mongo-driver/bson/primitive"

func GetObjectId(hexObjectIdSting string) (primitive.ObjectID, error) {
    id := primitive.NilObjectID
    if hexObjectIdSting != "" {
        return primitive.ObjectIDFromHex(hexObjectIdSting)
    }
    return id, nil
}
