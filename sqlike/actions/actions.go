package actions

// FindOne :
func FindOne() SelectOneStatement {
	return &FindOneActions{}
}

// Find :
func Find() SelectStatement {
	return &FindActions{}
}

// ReplaceInto :
func ReplaceInto() SelectStatement {
	return &ReplaceIntoActions{}
}

// UpdateOne :
func UpdateOne() UpdateOneStatement {
	return &UpdateOneActions{}
}

// UpdateMany :
func UpdateMany() UpdateStatement {
	return &UpdateActions{}
}

// Delete :
func Delete() DeleteStatement {
	return &DeleteActions{}
}
