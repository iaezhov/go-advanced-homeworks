package req

import (
	"errors"
	"net/http"
	"strconv"
)

func GetIDFromPath(r *http.Request) (uint, error) {
	idStr := r.PathValue("id")
	if idStr == "" {
		return 0, errors.New("Отсутствует Id")
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id == 0 {
		return 0, errors.New("Неверный формат id")
	}

	return uint(id), nil
}
