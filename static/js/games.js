// 게임 추가
async function submitGame() {
    const gameName = document.getElementById('gameName').value;
    if (!gameName) {
        alert('게임 이름을 입력해주세요.');
        return;
    }

    try {
        const response = await fetch('/api/games', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name: gameName }),
        });

        if (!response.ok) {
            const errorData = await response.json();
            if (errorData.error === 'game already exists') {
                throw new Error('중복된 게임 이름입니다.');
            }
            throw new Error('게임 추가에 실패했습니다.');
        }

        const result = await response.json();
        
        // 모달 닫기
        const modal = bootstrap.Modal.getInstance(document.getElementById('addGameModal'));
        modal.hide();

        // 성공 메시지 표시
        alert(`${gameName} 게임이 생성되었습니다.`);
        
        // 페이지 새로고침
        window.location.reload();
    } catch (error) {
        alert(error.message);
    }
}

// 게임 수정
async function editGame(id) {
    const newName = prompt('새로운 게임 이름을 입력하세요:');
    if (!newName) return;

    try {
        const response = await fetch(`/api/games/${id}`, {
            method: 'PUT',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ name: newName }),
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || '게임 수정에 실패했습니다.');
        }

        const result = await response.json();
        alert(`${result.name} 게임이 수정되었습니다.`);
        window.location.reload();
    } catch (error) {
        alert(error.message);
    }
}

// 게임 삭제
async function deleteGame(id) {
    if (!confirm('정말로 이 게임을 삭제하시겠습니까?')) return;

    try {
        const response = await fetch(`/api/games/${id}`, {
            method: 'DELETE',
        });

        if (!response.ok) {
            const errorData = await response.json();
            throw new Error(errorData.error || '게임 삭제에 실패했습니다.');
        }

        const result = await response.json();
        alert('게임이 삭제되었습니다.');
        window.location.reload();
    } catch (error) {
        alert(error.message);
    }
} 