package seed

import "go.mongodb.org/mongo-driver/bson/primitive"

// AdminUnitIDHex is the hex string of the seeded admin unit's ID.
const AdminUnitIDHex = "687f569bcd2015c348afcc27"

// AdminUnitID is the ObjectID form of AdminUnitIDHex.
var AdminUnitID, _ = primitive.ObjectIDFromHex(AdminUnitIDHex)
