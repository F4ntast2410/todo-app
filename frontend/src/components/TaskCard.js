
// Функция для обычного списка
export function createActiveTaskCard(task, onComplete, onDelete, onEditDescription) {
    const card = document.createElement('div');
    card.className = 'task-card';

    // Если бэкенд говорит, что задача выполнена — докидываем специальный класс
    if (task.done) {
        card.classList.add('task-card--completed');
    }
    card.innerHTML = `
    <div class="task-main">
        <div class="task-content">
            <span class="task-title">[#${task.user_task_id}] ${task.title}</span>
            ${task.description ? `<p class="task-desc">${task.description}</p>` : ''}
        </div>
        <div class="task-actions">
            ${task.done ? `<button class="icon-btn icon-btn--status">❌</button>` : `<button class="icon-btn icon-btn--status">✅</button>`}
            <button class="icon-btn icon-btn--edit">📝</button>
            <button class="icon-btn icon-btn--delete">🗑️</button>
        </div>
    </div>
    <div class="task-edit-form hidden">
        <textarea class="task-edit-input" rows="2">${task.description || ''}</textarea>
        <div class="task-edit-actions">
            <button class="btn-save-desc">Сохранить</button>
            <button class="btn-cancel-desc">Отмена</button>
        </div>
    </div>
    `;
    // Находим кнопку внутри свежесозданной карточки
    const deleteBtn = card.querySelector('.icon-btn.icon-btn--delete');

    // Вешаем событие клика
    deleteBtn.addEventListener('click', () => {
        onDelete(task.user_task_id); // Передаем ID задачи в колбэк, который прилетел из главного скрипта
    });

    const statusBtn = card.querySelector('.icon-btn.icon-btn--status');
    statusBtn.addEventListener('click', () => {
        onComplete(task.user_task_id, task.done);
    });

    const editForm = card.querySelector('.task-edit-form');
    const editBtn = card.querySelector('.icon-btn.icon-btn--edit');
    const editInput = card.querySelector('.task-edit-input');
    const descEl = card.querySelector('.task-desc');

    const openEdit = () => {
        editForm.classList.remove('hidden');
        editBtn.classList.add('hidden');
        descEl.classList.add('hidden');
        editInput.focus();
    };
    const closeEdit = () => {
        editInput.value = task.description || '';
        editForm.classList.add('hidden');
        editBtn.classList.remove('hidden');
        descEl.classList.remove('hidden');
    };

    editBtn.addEventListener('click', openEdit);
    card.querySelector('.btn-cancel-desc').addEventListener('click', closeEdit);

    card.querySelector('.btn-save-desc').addEventListener('click', () => {
        const newDescription = editInput.value.trim();
        onEditDescription(task.user_task_id, newDescription);
    });
    return card;

}

// Функция для корзины
export function createTrashTaskCard(task, onRestore, onDeleteForever) {
    const card = document.createElement('div');
    card.className = 'task-card';

    // Если бэкенд говорит, что задача выполнена — докидываем специальный класс
    if (task.done) {
        card.classList.add('task-card--completed');
    }
    card.innerHTML = `
    <div class="task-main">
        <div class="task-content">
            <span class="task-title">[#${task.user_task_id}] ${task.title}</span>
            ${task.description ? `<p class="task-desc">${task.description}</p>` : ''}
        </div>
        <div class="task-actions">
            <button class="icon-btn icon-btn--restore">🔄</button>
            <button class="icon-btn icon-btn--delete-forever">🗑️</button>
        </div>
    </div>`;
    // Находим кнопку внутри свежесозданной карточки
    const restoreBtn = card.querySelector('.icon-btn.icon-btn--restore');

    // Вешаем событие клика
    restoreBtn.addEventListener('click', () => {
        onRestore(task.user_task_id); // Передаем ID задачи в колбэк, который прилетел из главного скрипта
    });

    const deleteForverBtn = card.querySelector('.icon-btn.icon-btn--delete-forever');
    deleteForverBtn.addEventListener('click', () => {
        onDeleteForever(task.user_task_id);
    });
    return card;
}