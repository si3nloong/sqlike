package actions

// FindOne :
func FindOne() SelectOneStatement {
	return &FindOneActions{}
}

// Find :
func Find() SelectStatement {
	return &FindActions{}
}

// Paginate :
func Paginate() PaginateStatement {
	return &PaginateActions{}
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

// Copy :
func Copy() CopyStatement {
	return &CopyActions{}
}
