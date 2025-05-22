// 매치 추가
async function submitMatch() {
    const gameId = document.getElementById('gameId').value;
    const team1Id = document.getElementById('team1Id').value;
    const team2Id = document.getElementById('team2Id').value;

    if (!gameId || !team1Id || !team2Id) {
        alert('모든 필드를 입력해주세요.');
        return;
    }

    if (team1Id === team2Id) {
        alert('서로 다른 팀을 선택해주세요.');
        return;
    }

    const parsedGameId = parseInt(gameId);
    const parsedTeam1Id = parseInt(team1Id);
    const parsedTeam2Id = parseInt(team2Id);

    if (isNaN(parsedGameId) || isNaN(parsedTeam1Id) || isNaN(parsedTeam2Id)) {
        alert('유효하지 않은 ID 값입니다.');
        return;
    }

    const requestData = {
        game_id: parsedGameId,
        player_team_id: parsedTeam1Id,
        opponent_team_id: parsedTeam2Id
    };

    console.log('매치 생성 요청 데이터:', requestData);

    try {
        const response = await fetch('/api/matches', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestData),
        });

        if (!response.ok) {
            const errorData = await response.json();
            console.error('매치 생성 실패:', errorData);
            if (errorData.error === 'match already exists') {
                throw new Error('중복된 매치입니다.');
            }
            throw new Error('매치 추가에 실패했습니다.');
        }

        const result = await response.json();
        console.log('매치 생성 응답:', result);
        
        // 모달 닫기
        const modal = bootstrap.Modal.getInstance(document.getElementById('addMatchModal'));
        modal.hide();

        // 성공 메시지 표시
        alert('매치가 생성되었습니다.');
        
        // 페이지 새로고침
        window.location.reload();
    } catch (error) {
        console.error('매치 생성 에러:', error);
        alert(error.message);
    }
}

// 매치 수정 모달 열기
async function editMatch(id) {
    try {
        const response = await fetch(`/api/matches/${id}`);
        if (!response.ok) {
            throw new Error('매치 정보를 불러오는데 실패했습니다.');
        }

        const match = await response.json();
        console.log('수정할 매치 정보:', match);

        // 모달에 현재 값 설정
        document.getElementById('editMatchId').value = match.id;
        document.getElementById('editGameId').value = match.game_id;
        document.getElementById('editTeam1Id').value = match.player_team_id;
        document.getElementById('editTeam2Id').value = match.opponent_team_id;

        // 모달 표시
        const modal = new bootstrap.Modal(document.getElementById('editMatchModal'));
        modal.show();
    } catch (error) {
        console.error('매치 정보 조회 에러:', error);
        alert(error.message);
    }
}

// 매치 수정 제출
async function submitEditMatch() {
    const matchId = document.getElementById('editMatchId').value;
    const gameId = document.getElementById('editGameId').value;
    const team1Id = document.getElementById('editTeam1Id').value;
    const team2Id = document.getElementById('editTeam2Id').value;

    if (!gameId || !team1Id || !team2Id) {
        alert('모든 필드를 입력해주세요.');
        return;
    }

    if (team1Id === team2Id) {
        alert('서로 다른 팀을 선택해주세요.');
        return;
    }

    const parsedGameId = parseInt(gameId);
    const parsedTeam1Id = parseInt(team1Id);
    const parsedTeam2Id = parseInt(team2Id);

    if (isNaN(parsedGameId) || isNaN(parsedTeam1Id) || isNaN(parsedTeam2Id)) {
        alert('유효하지 않은 ID 값입니다.');
        return;
    }

    const requestData = {
        game_id: parsedGameId,
        player_team_id: parsedTeam1Id,
        opponent_team_id: parsedTeam2Id
    };

    console.log('매치 수정 요청 데이터:', requestData);

    try {
        const response = await fetch(`/api/matches/${matchId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestData),
        });

        if (!response.ok) {
            const errorData = await response.json();
            console.error('매치 수정 실패:', errorData);
            throw new Error(errorData.error || '매치 수정에 실패했습니다.');
        }

        const result = await response.json();
        console.log('매치 수정 응답:', result);
        
        // 모달 닫기
        const modal = bootstrap.Modal.getInstance(document.getElementById('editMatchModal'));
        modal.hide();

        // 성공 메시지 표시
        alert('매치가 수정되었습니다.');
        
        // 페이지 새로고침
        window.location.reload();
    } catch (error) {
        console.error('매치 수정 에러:', error);
        alert(error.message);
    }
}

// 매치 삭제
async function deleteMatch(id) {
    if (!confirm('정말로 이 매치를 삭제하시겠습니까?')) return;

    try {
        const response = await fetch(`/api/matches/${id}`, {
            method: 'DELETE',
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || '매치 삭제에 실패했습니다.');
        }

        const result = await response.json();
        console.log('매치 삭제 응답:', result);
        alert(result.message || '매치가 삭제되었습니다.');
        window.location.reload();
    } catch (error) {
        console.error('매치 삭제 에러:', error);
        alert(error.message);
    }
} 