package actions

// Find :
func Find() SelectStatement {
	return &FindActions{}
}

// Update :
func Update() UpdateStatement {
	return &UpdateActions{}
}

// Delete :
func Delete() DeleteStatement {
	return &DeleteActions{}
}
