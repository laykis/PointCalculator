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

// 팀 수정
async function editTeam(id) {
    const newName = prompt('새로운 팀 이름을 입력하세요:');
    if (!newName) return;

    try {
        const response = await fetch(`/api/teams/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name: newName }),
        });

        if (!response.ok) {
            throw new Error('팀 수정에 실패했습니다.');
        }

        window.location.reload();
    } catch (error) {
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
        alert(`${result.name} 팀이 삭제되었습니다.`);
        window.location.reload();
    } catch (error) {
        alert(error.message);
    }
} 