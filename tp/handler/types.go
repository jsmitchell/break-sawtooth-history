package handler

import "github.com/google/uuid"

type PayloadAction int
const WRITE_BLOCK PayloadAction = 1

type BreakSawtoothPayload struct {
	Action PayloadAction
	Id uuid.UUID
	Block []byte
}
