package service

// 매치 승리 시 고정 포인트
const WIN_POINT = 5

// 획득한 포인트를 저장할 구조체
type PointAccumulator struct {
	CurrentPoint     int // 현재 보유 포인트
	AccumulatedPoint int // 획득 예정 포인트
}

// 새로운 PointAccumulator 생성
func NewPointAccumulator(currentPoint int) *PointAccumulator {
	return &PointAccumulator{
		CurrentPoint:     currentPoint,
		AccumulatedPoint: 0,
	}
}

// 승리 포인트 추가
func (p *PointAccumulator) AddWinPoint() {
	p.AccumulatedPoint += WIN_POINT // 매치 승리 시 5점
}

// 베팅 포인트 처리
func (p *PointAccumulator) AddBetPoint(bettingPoint int, winOrLose bool) {
	if winOrLose {
		earnedPoints := (bettingPoint * 2) + 1 // 베팅 포인트 2배 + 1점
		p.AccumulatedPoint += earnedPoints
	}
}

// 최종 포인트 계산 및 적용
func (p *PointAccumulator) FinalizePoints(isDouble bool, isTriple bool) int {
	// 누적된 포인트(매치 승리 포인트 + 베팅 성공 포인트)에 찬스 적용
	finalPoint := p.AccumulatedPoint

	// 더블 찬스 적용
	if isDouble {
		finalPoint = DoublePoint(finalPoint)
	}

	// 트리플 찬스 적용
	if isTriple {
		finalPoint = TriplePoint(finalPoint)
	}

	// 최종 포인트 적용
	p.CurrentPoint += finalPoint
	// 누적 포인트 초기화
	p.AccumulatedPoint = 0

	return p.CurrentPoint
}

func WinPoint(addpoint int, point int) int {
	return point + WIN_POINT
}

func minusPoint(point int, minusPoint int) int {
	return point - minusPoint
}

func DoublePoint(point int) int {
	return point * 2
}

func TriplePoint(point int) int {
	return point * 3
}
