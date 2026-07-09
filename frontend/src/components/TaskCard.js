
// Функция для обычного списка
export function createActiveTaskCard(task, onComplete, onDelete) {
    const card = document.createElement('div');
    card.className = 'task-card';

    // Если бэкенд говорит, что задача выполнена — докидываем специальный класс
    if (task.done) {
        card.classList.add('task-card--completed');
    }
    card.innerHTML = `
    <div class="task-content">
        <span class="task-title">[#${task.user_task_id}] ${task.title}</span>
        ${task.description ? `<p class="task-desc">${task.description}</p>` : ''}
    </div>
    <div class="task-actions">
        ${task.done ? `<button class="btn-status">❌</button>` : `<button class="btn-status">✅</button>`}
        <button class="btn-delete">🗑️</button>
    </div>`;
    // Находим кнопку внутри свежесозданной карточки
    const deleteBtn = card.querySelector('.btn-delete');

    // Вешаем событие клика
    deleteBtn.addEventListener('click', () => {
        onDelete(task.user_task_id); // Передаем ID задачи в колбэк, который прилетел из главного скрипта
    });

    const statusBtn = card.querySelector('.btn-status');
    statusBtn.addEventListener('click', () => {
        onComplete(task.user_task_id, task.done);
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
    <div class="task-content">
        <span class="task-title">[#${task.user_task_id}] ${task.title}</span>
        ${task.description ? `<p class="task-desc">${task.description}</p>` : ''}
    </div>
    <div class="task-actions">
        <button class="btn-restore">🔄</button>
        <button class="btn-delete-forever">🗑️</button>
    </div>`;
    // Находим кнопку внутри свежесозданной карточки
    const restoreBtn = card.querySelector('.btn-restore');

    // Вешаем событие клика
    restoreBtn.addEventListener('click', () => {
        onRestore(task.user_task_id); // Передаем ID задачи в колбэк, который прилетел из главного скрипта
    });

    const deleteForverBtn = card.querySelector('.btn-delete-forever');
    deleteForverBtn.addEventListener('click', () => {
        onDeleteForever(task.user_task_id);
    });
    return card;
}