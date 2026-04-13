package service

func Health() map[string]string {
	return map[string]string{
		"status":  "ok",
		"service": "{{SERVICE_NAME}}",
	}
}
