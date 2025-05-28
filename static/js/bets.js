let selectedTeam = null;

function showBetModal(matchId, team1Name, team2Name) {
    document.getElementById('matchId').value = matchId;
    document.getElementById('betForm').reset();
    
    // 팀 선택 드롭다운 업데이트
    const betTeamSelect = document.getElementById('betTeam');
    const targetTeamSelect = document.getElementById('targetTeam');
    betTeamSelect.innerHTML = '<option value="">팀을 선택하세요</option>';
    targetTeamSelect.innerHTML = '<option value="">팀을 선택하세요</option>';
    
    // 전체 팀 목록 조회 후 베팅하는 팀 옵션 추가
    fetch('/api/teams')
        .then(response => response.json())
        .then(teams => {
            teams.forEach(team => {
                const option = document.createElement('option');
                option.value = team.id;
                option.textContent = team.name;
                betTeamSelect.appendChild(option);
            });
        });
    
    // 매치 정보 가져오기 (베팅 대상 팀 옵션 추가)
    fetch(`/api/matches/${matchId}`)
        .then(response => response.json())
        .then(match => {
            console.log('매치 API 응답:', match); // 실제 응답 확인용
            // 구조체의 json 태그(camelCase)에 맞춰 접근
            const targetOption1 = document.createElement('option');
            targetOption1.value = match.playerTeamId;
            targetOption1.textContent = team1Name;
            targetTeamSelect.appendChild(targetOption1);
            
            const targetOption2 = document.createElement('option');
            targetOption2.value = match.opponentTeamId;
            targetOption2.textContent = team2Name;
            targetTeamSelect.appendChild(targetOption2);
            
            new bootstrap.Modal(document.getElementById('betModal')).show();
        })
        .catch(error => {
            console.error('매치 정보 조회 에러:', error);
            alert('매치 정보를 불러오는데 실패했습니다.');
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

function submitBet() {
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

// 베팅 삭제
async function deleteBet(id) {
    if (!confirm('정말로 이 베팅을 삭제하시겠습니까?')) return;

    try {
        const response = await fetch(`/api/bets/${id}`, {
            method: 'DELETE',
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || '베팅 삭제에 실패했습니다.');
        }

        const result = await response.json();
        console.log('베팅 삭제 응답:', result);
        alert(result.message || '베팅이 삭제되었습니다.');
        window.location.reload();
    } catch (error) {
        console.error('베팅 삭제 에러:', error);
        alert(error.message);
    }
} 