// 팀 추가
async function submitTeam() {
    const teamName = document.getElementById('teamName').value;
    if (!teamName) {
        alert('팀 이름을 입력해주세요.');
        return;
    }

    try {
        const response = await fetch('/api/teams', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name: teamName }),
        });

        if (!response.ok) {
            const errorData = await response.json();
            if (errorData.error === 'team already exists') {
                throw new Error('중복된 팀 이름입니다.');
            }
            throw new Error('팀 추가에 실패했습니다.');
        }

        const result = await response.json();
        
        // 모달 닫기
        const modal = bootstrap.Modal.getInstance(document.getElementById('addTeamModal'));
        modal.hide();

        // 성공 메시지 표시
        alert(`${teamName} 팀이 생성되었습니다.`);
        
        // 페이지 새로고침
        window.location.reload();
    } catch (error) {
        alert(error.message);
    }
}

// 팀 수정 모달 열기
async function editTeam(id) {
    try {
        const response = await fetch(`/api/teams/${id}`);
        if (!response.ok) {
            throw new Error('팀 정보를 불러오는데 실패했습니다.');
        }

        const team = await response.json();
        console.log('수정할 팀 정보:', team);

        // 모달에 현재 값 설정
        document.getElementById('editTeamId').value = team.id;
        document.getElementById('editTeamName').value = team.name;
        document.getElementById('editTeamScore').value = team.point || 0;  // point가 없으면 0으로 설정

        // 모달 표시
        const modal = new bootstrap.Modal(document.getElementById('editTeamModal'));
        modal.show();
    } catch (error) {
        console.error('팀 정보 조회 에러:', error);
        alert(error.message);
    }
}

// 팀 수정 제출
async function submitEditTeam() {
    const teamId = document.getElementById('editTeamId').value;
    const teamName = document.getElementById('editTeamName').value;
    const teamScore = document.getElementById('editTeamScore').value;

    if (!teamName) {
        alert('팀 이름을 입력해주세요.');
        return;
    }

    if (teamScore === '' || isNaN(parseInt(teamScore)) || parseInt(teamScore) < 0) {
        alert('유효한 점수를 입력해주세요.');
        return;
    }

    const requestData = {
        name: teamName,
        point: parseInt(teamScore)
    };

    console.log('팀 수정 요청 데이터:', requestData);

    try {
        const response = await fetch(`/api/teams/${teamId}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(requestData),
        });

        const responseData = await response.json();
        console.log('서버 응답:', responseData);

        if (!response.ok) {
            console.error('팀 수정 실패:', responseData);
            if (responseData.error === 'team already exists') {
                throw new Error('중복된 팀 이름입니다.');
            }
            throw new Error(responseData.error || '팀 수정에 실패했습니다.');
        }
        
        // 모달 닫기
        const modal = bootstrap.Modal.getInstance(document.getElementById('editTeamModal'));
        modal.hide();

        // 성공 메시지 표시
        alert('팀이 수정되었습니다.');
        
        // 페이지 새로고침
        window.location.reload();
    } catch (error) {
        console.error('팀 수정 에러:', error);
        alert(error.message);
    }
}

// 팀 삭제
async function deleteTeam(id) {
    if (!confirm('정말로 이 팀을 삭제하시겠습니까?')) return;

    try {
        const response = await fetch(`/api/teams/${id}`, {
            method: 'DELETE',
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || '팀 삭제에 실패했습니다.');
        }

        const result = await response.json();
        alert('팀이 삭제되었습니다.');
        window.location.reload();
    } catch (error) {
        alert(error.message);
    }
} 