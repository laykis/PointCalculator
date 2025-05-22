package service

func WinPoint(addpoint int, point int) int {

	return point + addpoint
}

func minusPoint(point int, minusPoint int) int {

	return point - minusPoint
}

func BetPoint(bettingPoint int, point int, winOrLose bool, isDouble bool, isTriple bool) int {

	//추가 베팅 포인트 차감
	point = minusPoint(point, bettingPoint)

	//예측 성공 시
	if winOrLose {

		bettingPoint = DoublePoint(bettingPoint) // 베팅성공 점수
		bettingPoint += 1                        // 기본점수

		//더블 베팅 시
		if isDouble {
			bettingPoint = DoublePoint(bettingPoint)
		}

		//트리플 베팅 시
		if isTriple {
			bettingPoint = TriplePoint(bettingPoint)
		}

		//최종 획득 점수 합산
		point += bettingPoint
	}

	return point
}

func DoublePoint(point int) int {

	return point * 2
}

func TriplePoint(point int) int {

	return point * 3
}
