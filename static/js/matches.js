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
        gameId: parsedGameId,
        playerTeamId: parsedTeam1Id,
        opponentTeamId: parsedTeam2Id
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
            throw new Error(errorData.error || '매치 추가에 실패했습니다.');
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

// 승리 팀 설정
async function setWinner(matchId, winnerTeamId) {
    if (!confirm('승리 팀을 확정하시겠습니까?')) {
        return;
    }

    try {
        const response = await fetch(`/api/matches/${matchId}/result`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                winnerTeamId: winnerTeamId
            })
        });

        if (!response.ok) {
            throw new Error('승리 처리 중 오류가 발생했습니다.');
        }

        const result = await response.json();
        console.log('승리 팀 설정 응답:', result);
        
        if (result.success) {
            alert('승리 처리가 완료되었습니다.');
            window.location.reload();
        } else {
            alert(result.message || '승리 처리 중 오류가 발생했습니다.');
        }
    } catch (error) {
        console.error('승리 팀 설정 에러:', error);
        alert(error.message);
    }
}

// 페이지 로드 시 이벤트 리스너 등록
document.addEventListener('DOMContentLoaded', function() {
    // 베팅 버튼 클릭 이벤트 (이벤트 위임 사용)
    document.addEventListener('click', function(event) {
        const button = event.target.closest('.bet-button');
        if (button) {
            event.preventDefault();
            const matchId = button.dataset.matchId;
            showBetModal(matchId);
        }
    });

    // 베팅 제출 버튼 클릭 이벤트
    const submitBetButton = document.getElementById('submitBetButton');
    if (submitBetButton) {
        submitBetButton.addEventListener('click', submitBet);
    }

    // 베팅 폼 제출 이벤트
    const betForm = document.getElementById('betForm');
    if (betForm) {
        betForm.addEventListener('submit', function(event) {
            event.preventDefault();
            submitBet();
        });
    }

    // 모달 관련 설정
    const betModal = document.getElementById('betModal');
    if (betModal) {
        const modalInstance = new bootstrap.Modal(betModal);
        
        // 모달이 닫힐 때 초기화
        betModal.addEventListener('hidden.bs.modal', function () {
            document.body.classList.remove('modal-open');
            const modalBackdrop = document.querySelector('.modal-backdrop');
            if (modalBackdrop) {
                modalBackdrop.remove();
            }
            // 폼 초기화
            betForm.reset();
            const betTeamSelect = document.getElementById('betTeam');
            const targetTeamSelect = document.getElementById('targetTeam');
            if (betTeamSelect) betTeamSelect.innerHTML = '<option value="">팀을 선택하세요</option>';
            if (targetTeamSelect) targetTeamSelect.innerHTML = '<option value="">팀을 선택하세요</option>';
        });
    }
});

// 베팅 모달 표시
async function showBetModal(matchId) {
    try {
        // 매치 정보 가져오기
        const matchResponse = await fetch(`/api/matches/${matchId}`).then(res => {
            if (!res.ok) throw new Error('매치 정보를 가져오는데 실패했습니다.');
            return res.json();
        });
        
        console.log('매치 API 응답:', matchResponse);
        
        // 폼 초기화 및 매치 ID 설정
        const matchIdInput = document.getElementById('matchId');
        if (matchIdInput) matchIdInput.value = matchId;
        
        // 팀 선택 드롭다운 초기화
        const betTeamSelect = document.getElementById('betTeam');
        const targetTeamSelect = document.getElementById('targetTeam');
        
        if (betTeamSelect) betTeamSelect.innerHTML = '<option value="">팀을 선택하세요</option>';
        if (targetTeamSelect) targetTeamSelect.innerHTML = '<option value="">팀을 선택하세요</option>';
        
        // 전체 팀 목록을 베팅하는 팀 선택에 추가
        const teamSelects = document.querySelectorAll('#team1Id option:not(:first-child)');
        teamSelects.forEach(option => {
            if (betTeamSelect) {
                const newOption = option.cloneNode(true);
                betTeamSelect.appendChild(newOption);
            }
        });

        // 매치의 양 팀만 베팅 대상 팀 선택에 추가
        if (matchResponse.playerTeamId && targetTeamSelect) {
            const targetOption1 = document.createElement('option');
            targetOption1.value = matchResponse.playerTeamId;
            targetOption1.textContent = matchResponse.team1Name;
            targetTeamSelect.appendChild(targetOption1);
        }
        
        if (matchResponse.opponentTeamId && targetTeamSelect) {
            const targetOption2 = document.createElement('option');
            targetOption2.value = matchResponse.opponentTeamId;
            targetOption2.textContent = matchResponse.team2Name;
            targetTeamSelect.appendChild(targetOption2);
        }

        // 모달 표시
        const betModal = document.getElementById('betModal');
        if (betModal) {
            const modalInstance = bootstrap.Modal.getInstance(betModal) || new bootstrap.Modal(betModal);
            modalInstance.show();
        }
    } catch (error) {
        console.error('데이터 로딩 에러:', error);
        alert(error.message || '데이터를 불러오는데 실패했습니다.');
    }
}

// 베팅 제출
async function submitBet() {
    try {
        const form = document.getElementById('betForm');
        if (!form) throw new Error('베팅 폼을 찾을 수 없습니다.');

        // 폼 데이터 가져오기
        const formData = new FormData(form);
        const data = Object.fromEntries(formData.entries());
        
        // 데이터 유효성 검사
        if (!data.matchId || !data.teamId || !data.targetTeamId || !data.bet_type || !data.bettingPoint) {
            throw new Error('모든 필드를 입력해주세요.');
        }

        // 포인트 유효성 검사
        const point = parseInt(data.bettingPoint);
        if (isNaN(point) || point <= 0) {
            throw new Error('유효한 포인트를 입력해주세요.');
        }

        // 더블/트리플 찬스 체크 상태 가져오기
        const isDouble = document.getElementById('isDouble').checked;
        const isTriple = document.getElementById('isTriple').checked;

        // API 요청 데이터 구성
        const requestData = {
            match_id: parseInt(data.matchId),
            team_id: parseInt(data.teamId),
            target_team_id: parseInt(data.targetTeamId),
            bet_type: data.bet_type,
            betting_point: point,
            is_double: isDouble,
            is_triple: isTriple
        };

        // API 호출
        const response = await fetch('/api/bets', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestData),
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || '베팅 생성에 실패했습니다.');
        }

        const result = await response.json();
        console.log('베팅 생성 응답:', result);
        
        // 모달 닫기
        const betModal = document.getElementById('betModal');
        if (betModal) {
            const modalInstance = bootstrap.Modal.getInstance(betModal);
            if (modalInstance) modalInstance.hide();
        }

        // 성공 메시지 표시
        alert('베팅이 생성되었습니다.');
        
        // 페이지 새로고침
        window.location.reload();
    } catch (error) {
        console.error('베팅 생성 에러:', error);
        alert(error.message);
    }
}

// 랜덤 매치 생성
document.getElementById('createRandomMatchButton').addEventListener('click', async function() {
    const gameId = document.getElementById('randomGameId').value;
    
    if (!gameId) {
        alert('게임을 선택해주세요.');
        return;
    }

    try {
        const response = await fetch('/api/matches/random', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                gameId: parseInt(gameId)
            }),
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || '랜덤 매치 생성에 실패했습니다.');
        }

        // 성공 시 페이지 새로고침
        window.location.reload();
    } catch (error) {
        console.error('랜덤 매치 생성 에러:', error);
        alert(error.message);
    }
}); 