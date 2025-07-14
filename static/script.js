// Загружаем задачи при загрузке страницы
document.addEventListener('DOMContentLoaded', loadTasks);

async function loadTasks() {
    const response = await fetch('/api/tasks');
    const tasks = await response.json();
    renderTasks(tasks);
}

function renderTasks(tasks) {
    const tasksList = document.getElementById('tasks-list');
    tasksList.innerHTML = '';
    
    tasks.forEach(task => {
        const taskElement = document.createElement('div');
        taskElement.className = `task ${task.isCompleted ? 'completed' : ''}`;
        taskElement.innerHTML = `
            <div>
                <h3>${task.title}</h3>
                <p>${task.description}</p>
            </div>
            <div class="task-actions">
                <button onclick="toggleTask(${task.id}, ${task.isCompleted})">
                    ${task.isCompleted ? '❌ Отменить' : '✓ Выполнить'}
                </button>
                <button onclick="deleteTask(${task.id})">🗑 Удалить</button>
            </div>
        `;
        tasksList.appendChild(taskElement);
    });
}

async function addTask() {
    console.log("Функция addTask вызвана")
    const title = document.getElementById('task-title').value;
    const description = document.getElementById('task-desc').value;
    
    console.log("Отправляемые данные:", {  // ← Добавьте эту строку
        title,
        description,
        isCompleted: false
    });

    if (!title) {
        alert('Введите название задачи');
        return;
    }
    
    const response = await fetch('/api/tasks', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            title,
            description,
            isCompleted: false
        })
    });
    
    if (response.ok) {
        document.getElementById('task-title').value = '';
        document.getElementById('task-desc').value = '';
        loadTasks();
    }
}

async function toggleTask(id, isCompleted) {
    await fetch(`/api/tasks/${id}`, {
        method: 'PUT',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({
            isCompleted: !isCompleted
        })
    });
    loadTasks();
}

async function deleteTask(id) {
    if (confirm('Удалить задачу?')) {
        await fetch(`/api/tasks/${id}`, {
            method: 'DELETE'
        });
        loadTasks();
    }
}