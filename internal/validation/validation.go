package validation

func Start(pErr chan<- error) {
	// Not implemented yet

	// Start validation workers

	// Worker does:
	// ==== move example ====
	// - pull message from queue
	// - locks user
	// - checks usersLastSeq - if less, then done
	// - updates gameTick state for char position
	// - unlock
	// ==== attack example (hit) ====
	// - pull message from queue
	// - locks user
	// - checks usersLastSeq - if less, then done + unlock
	// -
	// - updates gameTick state for char
	// - unlock

	//

	select {}
}
