package domain

import (
	"encoding/json"
	"fmt"
)

type LessonStatus int8

const (
	LessonStatusDraft   LessonStatus = 0
	LessonStatusLimited LessonStatus = 1
	LessonStatusPublic  LessonStatus = 2
)

func (r LessonStatus) String() string {
	switch r {
	case LessonStatusDraft:
		return "draft"
	case LessonStatusLimited:
		return "limited"
	case LessonStatusPublic:
		return "public"
	default:
		return "unknown"
	}
}

func (r LessonStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (s *LessonStatus) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("data should be a string, got %s", data)
	}

	var status LessonStatus
	switch str {
	case "draft":
		status = LessonStatusDraft
	case "limited":
		status = LessonStatusLimited
	case "public":
		status = LessonStatusPublic
	default:
		return fmt.Errorf("invalid LessonStatus %s", str)
	}
	*s = status
	return nil
}

type DrawingAction int8

const (
	DrawingActionDraw  DrawingAction = 0
	DrawingActionClear DrawingAction = 1
	DrawingActionShow  DrawingAction = 2
	DrawingActionHide  DrawingAction = 3
)

func (r DrawingAction) String() string {
	switch r {
	case DrawingActionDraw:
		return "draw"
	case DrawingActionClear:
		return "clear"
	case DrawingActionShow:
		return "show"
	case DrawingActionHide:
		return "hide"
	default:
		return "unknown"
	}
}

func (r DrawingAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (a *DrawingAction) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("data should be a string, got %s", data)
	}

	var action DrawingAction
	switch str {
	case "draw":
		action = DrawingActionDraw
	case "clear":
		action = DrawingActionClear
	case "show":
		action = DrawingActionShow
	case "hide":
		action = DrawingActionHide
	default:
		return fmt.Errorf("invalid DrawingAction %s", str)
	}
	*a = action
	return nil
}

type DrawingUnitAction int8

const (
	DrawingUnitActionDraw DrawingUnitAction = 0
	DrawingUnitActionUndo DrawingUnitAction = 1
)

func (r DrawingUnitAction) String() string {
	switch r {
	case DrawingUnitActionDraw:
		return "draw"
	case DrawingUnitActionUndo:
		return "undo"
	default:
		return "unknown"
	}
}

func (r DrawingUnitAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (a *DrawingUnitAction) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("data should be a string, got %s", data)
	}

	var action DrawingUnitAction
	switch str {
	case "draw":
		action = DrawingUnitActionDraw
	case "undo":
		action = DrawingUnitActionUndo
	default:
		return fmt.Errorf("invalid DrawingUnitAction %s", str)
	}
	*a = action
	return nil
}

type EmbeddingAction int8

const (
	EmbeddingActionShow EmbeddingAction = 0
	EmbeddingActionHide EmbeddingAction = 1
)

func (r EmbeddingAction) String() string {
	switch r {
	case EmbeddingActionShow:
		return "show"
	case EmbeddingActionHide:
		return "hide"
	default:
		return "unknown"
	}
}

func (r EmbeddingAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (a *EmbeddingAction) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("data should be a string, got %s", data)
	}

	var action EmbeddingAction
	switch str {
	case "show":
		action = EmbeddingActionShow
	case "hide":
		action = EmbeddingActionHide
	default:
		return fmt.Errorf("invalid EmbeddingAction %s", str)
	}
	*a = action
	return nil
}

type GraphicAction int8

const (
	GraphicActionShow GraphicAction = 0
	GraphicActionHide GraphicAction = 1
)

func (r GraphicAction) String() string {
	switch r {
	case GraphicActionShow:
		return "show"
	case GraphicActionHide:
		return "hide"
	default:
		return "unknown"
	}
}

func (r GraphicAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (a *GraphicAction) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("data should be a string, got %s", data)
	}

	var action GraphicAction
	switch str {
	case "show":
		action = GraphicActionShow
	case "hide":
		action = GraphicActionHide
	default:
		return fmt.Errorf("invalid GraphicAction %s", str)
	}
	*a = action
	return nil
}

type MusicAction int8

const (
	MusicActionStart MusicAction = 0
	MusicActionStop  MusicAction = 1
)

func (r MusicAction) String() string {
	switch r {
	case MusicActionStart:
		return "start"
	case MusicActionStop:
		return "stop"
	default:
		return "unknown"
	}
}

func (r MusicAction) MarshalJSON() ([]byte, error) {
	return json.Marshal(r.String())
}

func (a *MusicAction) UnmarshalJSON(data []byte) error {
	var str string
	if err := json.Unmarshal(data, &str); err != nil {
		return fmt.Errorf("data should be a string, got %s", data)
	}

	var action MusicAction
	switch str {
	case "start":
		action = MusicActionStart
	case "stop":
		action = MusicActionStop
	default:
		return fmt.Errorf("invalid MusicAction %s", str)
	}
	*a = action
	return nil
}
