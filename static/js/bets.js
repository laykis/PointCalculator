// 전역 변수와 함수 선언
window.selectedTeam = null;

// 페이지 로드 시 실행
document.addEventListener('DOMContentLoaded', function() {
    console.log('DOMContentLoaded 이벤트 발생');
    
    // 진행중인 매치의 베팅 목록 로드
    const containers = document.querySelectorAll('.bet-list-container');
    console.log('찾은 베팅 목록 컨테이너:', containers.length);
    
    containers.forEach(container => {
        const matchId = container.id.replace('betList-', '');
        console.log('매치 ID:', matchId);
        loadBetList(matchId);
    });
});

// 매치별 베팅 목록 토글
window.toggleBetList = function(matchId, event) {
    console.log('toggleBetList 호출됨:', matchId);
    // 베팅하기 버튼 클릭 시 이벤트 전파 중단
    if (event && event.target.tagName === 'BUTTON') {
        console.log('버튼 클릭으로 인한 이벤트 중단');
        event.stopPropagation();
        return;
    }
    
    const betListContainer = document.getElementById(`betList-${matchId}`);
    console.log('betListContainer:', betListContainer);
    
    // 현재 display 상태 확인
    const isHidden = betListContainer.style.display === 'none' || !betListContainer.style.display;
    
    if (isHidden) {
        console.log('베팅 목록 표시');
        loadBetList(matchId);
        betListContainer.style.display = 'block';
    } else {
        console.log('베팅 목록 숨김');
        betListContainer.style.display = 'none';
    }
}

// 베팅 목록 로드
async function loadBetList(matchId) {
    try {
        console.log(`${matchId} 매치의 베팅 목록 로드 시작`);
        
        const betListBody = document.getElementById(`betListBody-${matchId}`);
        if (!betListBody) {
            console.error(`betListBody-${matchId} 엘리먼트를 찾을 수 없음`);
            return;
        }
        
        const response = await fetch(`/api/matches/${matchId}/bets`);
        console.log('API 응답 상태:', response.status);
        
        if (!response.ok) {
            throw new Error('베팅 목록을 불러오는데 실패했습니다.');
        }

        const bets = await response.json();
        console.log(`${matchId} 매치의 베팅 목록:`, bets);

        betListBody.innerHTML = '';

        if (bets.length === 0) {
            const row = document.createElement('tr');
            row.innerHTML = '<td colspan="5" class="text-center">베팅 내역이 없습니다.</td>';
            betListBody.appendChild(row);
            return;
        }

        bets.forEach(bet => {
            console.log('처리중인 베팅:', bet);
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${bet.team_name || '알 수 없음'}</td>
                <td>${bet.target_team_name || '알 수 없음'}</td>
                <td>${bet.status === 'W' ? '승리' : '패배'}</td>
                <td>${bet.betting_point}</td>
                <td>
                    ${getBetStatusBadge(bet.status)}
                    ${bet.status === 'P' ? `
                        <button class="btn btn-sm btn-outline-danger ms-2" onclick="window.deleteBet(${bet.id})">삭제</button>
                    ` : ''}
                </td>
            `;
            betListBody.appendChild(row);
        });
    } catch (error) {
        console.error(`${matchId} 매치의 베팅 목록 로드 에러:`, error);
    }
}

// 베팅 상태 뱃지 생성
function getBetStatusBadge(status) {
    switch (status) {
        case 'P':
            return '<span class="badge bg-warning">진행중</span>';
        case 'C':
            return '<span class="badge bg-success">완료</span>';
        default:
            return '<span class="badge bg-secondary">알 수 없음</span>';
    }
}

// 베팅 수정 모달
const editBetModal = new bootstrap.Modal(document.getElementById('editBetModal'));
let currentBetId = null;

// 베팅 수정 모달 열기
window.editBet = async (betId) => {
    try {
        const response = await fetch(`/api/bets/${betId}`);
        if (!response.ok) {
            throw new Error('베팅 정보를 불러오는데 실패했습니다.');
        }
        const bet = await response.json();
        
        // 모달에 기존 베팅 정보 설정
        document.getElementById('editBetId').value = bet.ID;
        document.getElementById('editBetType').value = bet.BetType;
        document.getElementById('editBettingPoint').value = bet.BettingPoint;
        document.getElementById('editIsDouble').checked = bet.IsDouble;
        document.getElementById('editIsTriple').checked = bet.IsTriple;
        
        currentBetId = betId;
        editBetModal.show();
    } catch (error) {
        alert(error.message);
    }
};

// 베팅 수정 제출
window.submitEditBet = async () => {
    if (!currentBetId) return;

    const formData = {
        betType: document.getElementById('editBetType').value,
        bettingPoint: parseInt(document.getElementById('editBettingPoint').value),
        isDouble: document.getElementById('editIsDouble').checked,
        isTriple: document.getElementById('editIsTriple').checked
    };

    try {
        const response = await fetch(`/api/bets/${currentBetId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(formData)
        });

        if (!response.ok) {
            throw new Error('베팅 수정에 실패했습니다.');
        }

        editBetModal.hide();
        window.location.reload();
    } catch (error) {
        alert(error.message);
    }
};

// 베팅 삭제
window.deleteBet = async (betId) => {
    if (!confirm('정말로 이 베팅을 삭제하시겠습니까?')) {
        return;
    }

    try {
        const response = await fetch(`/api/bets/${betId}`, {
            method: 'DELETE'
        });

        if (!response.ok) {
            throw new Error('베팅 삭제에 실패했습니다.');
        }

        window.location.reload();
    } catch (error) {
        alert(error.message);
    }
}; 