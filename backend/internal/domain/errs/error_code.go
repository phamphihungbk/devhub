package errs

const (
	StatusCodeSuccess                      = "200000" // 200000
	StatusCodeServiceCircuitBreakerTripped = "503001" // service circuit breaker tripped error code
	StatusCodeSeatBooked                   = "403001" // conflict error when trying to book a seat that is already booked
	StatusCodeSeatLocked                   = "403002" // conflict error when trying to book a seat that is already locked
)
