// 전역 변수와 함수 선언
window.selectedTeam = null;

// 페이지 로드 시 실행
document.addEventListener('DOMContentLoaded', function() {
    // 진행중인 매치의 베팅 목록 로드
    document.querySelectorAll('.bet-list-container').forEach(container => {
        const matchId = container.id.replace('betList-', '');
        loadBetList(matchId);
    });

    // 베팅하기 버튼 이벤트 리스너
    document.querySelectorAll('.bet-button').forEach(button => {
        button.addEventListener('click', function(event) {
            const matchId = this.dataset.matchId;
            const team1Name = this.dataset.team1Name;
            const team2Name = this.dataset.team2Name;
            showBetModal(matchId, team1Name, team2Name);
        });
    });
});

// 베팅 모달 표시
window.showBetModal = function(matchId, team1Name, team2Name) {
    document.getElementById('matchId').value = matchId;
    document.getElementById('betForm').reset();
    
    // 팀 선택 드롭다운 업데이트
    const betTeamSelect = document.getElementById('betTeam');
    const targetTeamSelect = document.getElementById('targetTeam');
    betTeamSelect.innerHTML = '<option value="">팀을 선택하세요</option>';
    targetTeamSelect.innerHTML = '<option value="">팀을 선택하세요</option>';
    
    // 매치 정보 먼저 가져오기
    fetch(`/api/matches/${matchId}`)
        .then(response => response.json())
        .then(match => {
            console.log('매치 API 응답:', match);
            
            // 팀 정보 가져오기
            return fetch('/api/teams')
                .then(response => response.json())
                .then(teams => {
                    // 전체 팀 목록을 베팅 팀 선택에 추가
                    teams.forEach(team => {
                        const option = document.createElement('option');
                        option.value = team.id;
                        option.textContent = team.name;
                        betTeamSelect.appendChild(option);
                    });

                    // 매치의 양 팀을 대상 팀 선택에 추가
                    const team1 = teams.find(t => t.id === match.playerTeamId);
                    const team2 = teams.find(t => t.id === match.opponentTeamId);

                    if (team1) {
                        const targetOption1 = document.createElement('option');
                        targetOption1.value = match.playerTeamId;
                        targetOption1.textContent = team1.name;
                        targetTeamSelect.appendChild(targetOption1);
                    }
                    
                    if (team2) {
                        const targetOption2 = document.createElement('option');
                        targetOption2.value = match.opponentTeamId;
                        targetOption2.textContent = team2.name;
                        targetTeamSelect.appendChild(targetOption2);
                    }

                    new bootstrap.Modal(document.getElementById('betModal')).show();
                });
        })
        .catch(error => {
            console.error('데이터 로딩 에러:', error);
            alert('데이터를 불러오는데 실패했습니다.');
        });
}

function selectTeam(teamNumber) {
    selectedTeam = teamNumber;
    const team1Btn = document.querySelector('button[onclick="selectTeam(1)"]');
    const team2Btn = document.querySelector('button[onclick="selectTeam(2)"]');
    
    team1Btn.classList.remove('btn-primary');
    team2Btn.classList.remove('btn-primary');
    team1Btn.classList.add('btn-outline-primary');
    team2Btn.classList.add('btn-outline-primary');
    
    if (teamNumber === 1) {
        team1Btn.classList.remove('btn-outline-primary');
        team1Btn.classList.add('btn-primary');
    } else {
        team2Btn.classList.remove('btn-outline-primary');
        team2Btn.classList.add('btn-primary');
    }
}

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
    console.log('loadBetList 시작:', matchId);
    try {
        const response = await fetch(`/api/matches/${matchId}/bets`);
        console.log('API 응답:', response);
        if (!response.ok) {
            throw new Error('베팅 목록을 불러오는데 실패했습니다.');
        }

        const bets = await response.json();
        console.log('받아온 베팅 목록:', bets);
        const betListBody = document.getElementById(`betListBody-${matchId}`);
        console.log('betListBody 엘리먼트:', betListBody);
        betListBody.innerHTML = '';

        if (bets.length === 0) {
            const row = document.createElement('tr');
            row.innerHTML = '<td colspan="5" class="text-center">베팅 내역이 없습니다.</td>';
            betListBody.appendChild(row);
            return;
        }

        bets.forEach(bet => {
            const row = document.createElement('tr');
            row.innerHTML = `
                <td>${bet.team_name || '알 수 없음'}</td>
                <td>${bet.target_team_name || '알 수 없음'}</td>
                <td>${bet.status === 'W' ? '승리' : '패배'}</td>
                <td>${bet.betting_point}</td>
                <td>
                    ${getBetStatusBadge(bet.status)}
                    ${bet.status === 'P' ? `
                        <button class="btn btn-sm btn-outline-danger ms-2" onclick="deleteBet(${bet.id})">삭제</button>
                    ` : ''}
                </td>
            `;
            betListBody.appendChild(row);
        });
    } catch (error) {
        console.error('베팅 목록 로드 에러:', error);
        alert(error.message);
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

// 베팅 제출
window.submitBet = function() {
    const matchId = document.getElementById('matchId').value;
    const betTeam = document.getElementById('betTeam').value;
    const targetTeam = document.getElementById('targetTeam').value;
    const betType = document.getElementById('betType').value;
    const point = document.getElementById('point').value;
    
    if (!betTeam) {
        alert('베팅하는 팀을 선택해주세요.');
        return;
    }
    
    if (!targetTeam) {
        alert('베팅 대상 팀을 선택해주세요.');
        return;
    }
    
    if (!betType) {
        alert('베팅 내용을 선택해주세요.');
        return;
    }
    
    if (!point || point < 0) {
        alert('베팅 포인트를 입력해주세요.');
        return;
    }

    if (betTeam === targetTeam) {
        alert('베팅하는 팀과 베팅 대상 팀이 같을 수 없습니다.');
        return;
    }
    
    fetch('/api/bets', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            match_id: parseInt(matchId),
            team_id: parseInt(betTeam),
            target_team_id: parseInt(targetTeam),
            betting_point: parseInt(point),
            status: betType
        })
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert(data.error);
        } else {
            alert('베팅이 완료되었습니다.');
            window.location.reload();
        }
    })
    .catch(error => {
        console.error('Error:', error);
        alert('베팅 중 오류가 발생했습니다.');
    });
}

// 베팅 삭제
window.deleteBet = function(id) {
    if (!confirm('정말로 이 베팅을 삭제하시겠습니까?')) return;

    fetch(`/api/bets/${id}`, {
        method: 'DELETE',
    })
    .then(response => {
        if (!response.ok) {
            throw new Error('베팅 삭제에 실패했습니다.');
        }
        return response.json();
    })
    .then(result => {
        alert(result.message || '베팅이 삭제되었습니다.');
        window.location.reload();
    })
    .catch(error => {
        console.error('베팅 삭제 에러:', error);
        alert(error.message);
    });
}

// 매치 선택 시 해당 매치의 팀 목록 업데이트
function updateTeamList(matchId, targetSelectId) {
    const select = document.getElementById(targetSelectId);
    select.innerHTML = '<option value="">팀을 선택하세요</option>';

    if (!matchId) return;

    // 선택된 매치의 팀 정보 가져오기
    fetch(`/api/matches/${matchId}`)
        .then(response => response.json())
        .then(match => {
            // 팀1 추가
            const option1 = document.createElement('option');
            option1.value = match.player_team_id;
            option1.textContent = match.team1_name;
            select.appendChild(option1);

            // 팀2 추가
            const option2 = document.createElement('option');
            option2.value = match.opponent_team_id;
            option2.textContent = match.team2_name;
            select.appendChild(option2);
        })
        .catch(error => {
            console.error('매치 정보 조회 에러:', error);
            alert('매치 정보를 불러오는데 실패했습니다.');
        });
}

// 베팅 수정 모달 열기
async function editBet(id) {
    try {
        const response = await fetch(`/api/bets/${id}`);
        if (!response.ok) {
            throw new Error('베팅 정보를 불러오는데 실패했습니다.');
        }

        const bet = await response.json();
        console.log('수정할 베팅 정보:', bet);

        // 모달에 현재 값 설정
        document.getElementById('editBetId').value = bet.id;
        document.getElementById('editMatchId').value = bet.match_id;
        updateTeamList(bet.match_id, 'editTeamId');
        document.getElementById('editTeamId').value = bet.team_id;
        document.getElementById('editPoint').value = bet.point;

        // 모달 표시
        const modal = new bootstrap.Modal(document.getElementById('editBetModal'));
        modal.show();
    } catch (error) {
        console.error('베팅 정보 조회 에러:', error);
        alert(error.message);
    }
}

// 베팅 수정 제출
async function submitEditBet() {
    const betId = document.getElementById('editBetId').value;
    const matchId = document.getElementById('editMatchId').value;
    const teamId = document.getElementById('editTeamId').value;
    const point = document.getElementById('editPoint').value;

    if (!matchId || !teamId || !point) {
        alert('모든 필드를 입력해주세요.');
        return;
    }

    const parsedMatchId = parseInt(matchId);
    const parsedTeamId = parseInt(teamId);
    const parsedPoint = parseInt(point);

    if (isNaN(parsedMatchId) || isNaN(parsedTeamId) || isNaN(parsedPoint)) {
        alert('유효하지 않은 값입니다.');
        return;
    }

    const requestData = {
        match_id: parsedMatchId,
        team_id: parsedTeamId,
        point: parsedPoint
    };

    try {
        const response = await fetch(`/api/bets/${betId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestData),
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || '베팅 수정에 실패했습니다.');
        }

        const result = await response.json();
        console.log('베팅 수정 응답:', result);
        
        // 모달 닫기
        const modal = bootstrap.Modal.getInstance(document.getElementById('editBetModal'));
        modal.hide();

        // 성공 메시지 표시
        alert('베팅이 수정되었습니다.');
        
        // 페이지 새로고침
        window.location.reload();
    } catch (error) {
        console.error('베팅 수정 에러:', error);
        alert(error.message);
    }
} 