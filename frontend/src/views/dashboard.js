import { getProfile, logout } from '../api/auth.js';
import { getTasks, createTask, deleteTask, updateStatus, restoreTask, getTrashTasks } from '../api/tasks.js';
import { createActiveTaskCard, createTrashTaskCard } from '../components/TaskCard.js';

async function initDashboard() {
    try {
        // Вызываем проверку профиля при старте страницы
        const response = await getProfile();

        if (!response.ok) {
            // Если бэкенд вернул ошибку, отличную от 401 (например, 500)
            // на всякий случай тоже отправляем на авторизацию
            window.location.href = '/views/auth.html';
            return;
        }

        // Если всё ок, вытаскиваем данные пользователя (например, имя)
        const user = await response.json();
        console.log('Успешная авторизация. Добро пожаловать,', user.username);

        // Здесь запускаем инициализацию остального интерфейса дашборда
        setupDashboardUI(user);
        await loadUserTasks();
        setupFormListeners();

        setupLogoutListener();

    } catch (error) {
        console.error('Ошибка при проверке авторизации:', error);
        window.location.href = '/views/auth.html';
    }
}

function setupDashboardUI(user) {
    const usernameDisplay = document.getElementById('username-display');
    if (usernameDisplay) {
        usernameDisplay.textContent = user.username;
    }
}
async function loadUserTasks() {
    const container = document.getElementById('tasks-container');
    if (!container) return;
    try {
        const response = await getTasks();
        if (!response.ok) throw new Error('Не удалось загрузить задачи');

        const tasks = await response.json();
        container.innerHTML = ''; // Очищаем текст загрузки

        if (!tasks || tasks.length === 0) {
            container.innerHTML = '<p class="empty-text">У вас пока нет задач. Самое время создать первую!</p>';
            return;
        }

        // Рендерим каждую задачу (подстрой поля под структуру твоей DTO/модели в Go)
        tasks.forEach(task => {
            const taskElement = createActiveTaskCard(task, onComplete, onDelete);
            container.appendChild(taskElement);
        });

    } catch (error) {
        console.error(error);
        container.innerHTML = '<p class="error-text">Ошибка при загрузке задач с сервера.</p>';
    }
}
async function loadTrashTasks() {
    const trashContainer = document.getElementById('trash-container');
    if (!trashContainer) return;
    try {
        const response = await getTrashTasks();
        if (!response.ok) throw new Error('Не удалось загрузить задачи из корзины');

        const tasks = await response.json();
        trashContainer.innerHTML = ''; // Очищаем текст загрузки

        if (!tasks || tasks.length === 0) {
            trashContainer.innerHTML = '<p class="empty-text">Корзина пуста.</p>';
            return;
        }

        // Рендерим каждую задачу (подстрой поля под структуру твоей DTO/модели в Go)
        tasks.forEach(task => {
            const taskElement = createTrashTaskCard(task, onRestore, onDeletePermanently);
            trashContainer.appendChild(taskElement);
        });

    } catch (error) {
        console.error(error);
        trashContainer.innerHTML = '<p class="error-text">Ошибка при загрузке корзины с сервера.</p>';
    }
}

// Обработка отправки формы
function setupFormListeners() {
    const form = document.getElementById('add-task-form');
    const titleInput = document.getElementById('task-title-input');
    const descInput = document.getElementById('task-desc-input'); // Добавляем ссылку на описание

    form.addEventListener('submit', async (e) => {
        e.preventDefault();
        
        const title = titleInput.value.trim();
        const description = descInput.value.trim(); // Забираем значение
        
        if (!title) return;

        try {
            // Передаем и title, и description. 
            // Ключ 'description' должен совпадать с json-тегом в твоей Go-структуре!
            const response = await createTask({ 
                title: title, 
                description: description 
            });
            
            if (response.ok) {
                titleInput.value = ''; 
                descInput.value = ''; // Очищаем поле описания после успешного создания
                await loadUserTasks(); 
            } else {
                alert('Не удалось создать задачу');
            }
        } catch (error) {
            console.error('Ошибка при создании задачи:', error);
        }
    });
}
function setupLogoutListener() {
    const logoutBtn = document.getElementById('logout-btn');
    if (!logoutBtn) return;

    logoutBtn.addEventListener('click', async () => {
        try {
            const response = await logout();
            
            if (response.ok) {
                // Если бэкенд успешно удалил куку, отправляем пользователя на авторизацию
                window.location.href = '/views/auth.html';
            } else {
                alert('Не удалось выйти из системы. Попробуйте еще раз.');
            }
        } catch (error) {
            console.error('Ошибка при выходе:', error);
        }
    });
}

async function onDelete(user_task_id) {
    try {
        const response = await deleteTask(user_task_id);

        if (response.ok) {
            await loadUserTasks(); 
        } else {
            alert('Не удалось удалить задачу');
        }
    } catch (error) {
        console.error('Ошибка при удалении задачи:', error);
    }
}
async function onRestore(user_task_id) {
    try {
        const response = await restoreTask(user_task_id);

        if (response.ok) {
            await loadUserTasks(); 
            await loadTrashTasks();
        } else {
            alert('Не удалось восстановить задачу');
        }
    } catch (error) {
        console.error('Ошибка при восстановлении задачи:', error);
    }
}
async function onComplete(user_task_id, done) {
    try {
        const response = await updateStatus(user_task_id, {done: !done});

        if (response.ok) {
            await loadUserTasks(); 
        } else {
            alert('Не удалось обновить задачу');
        }
    } catch (error) {
        console.error('Ошибка при обновлении задачи:', error);
    }
}
async function onDeletePermanently(user_task_id) {
    if (!confirm('Удалить эту задачу навсегда?')) return;
    try {
        const response = await deleteTask(user_task_id);
        if (response.ok) {
            await loadTrashTasks();
        }
    } catch (error) {
        console.error('Ошибка при полном удалении:', error);
    }
}
function setupTabListeners() {
    const btnActive = document.getElementById('btn-tab-active');
    const btnTrash = document.getElementById('btn-tab-trash');
    const activeContainer = document.getElementById('tasks-container');
    const trashContainer = document.getElementById('trash-container');
    const taskForm = document.querySelector('.task-form'); // твоя форма добавления задач

    btnActive.addEventListener('click', () => {
        // Меняем активный класс у кнопок
        btnActive.classList.add('active');
        btnTrash.classList.remove('active');
        
        // Показываем активные задачи и форму создания, скрываем корзину
        activeContainer.classList.remove('hidden');
        if (taskForm) taskForm.classList.remove('hidden');
        trashContainer.classList.add('hidden');
        
        // Обновляем данные
        loadUserTasks();
    });

    btnTrash.addEventListener('click', () => {
        // Меняем активный класс у кнопок
        btnTrash.classList.add('active');
        btnActive.classList.remove('active');
        
        // Скрываем активные и форму, показываем корзину
        activeContainer.classList.add('hidden');
        if (taskForm) taskForm.classList.add('hidden'); // В корзине форма создания не нужна
        trashContainer.classList.remove('hidden');
        
        // Загружаем удалённые задачи
        loadTrashTasks();
    });
}
// Запускаем проверку немедленно при загрузке файла скрипта
setupTabListeners();
initDashboard();